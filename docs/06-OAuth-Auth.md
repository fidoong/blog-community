# 06 - OAuth2 + JWT 鉴权方案详细设计

> 本文档描述博客社区的统一认证授权体系，涵盖邮箱登录、GitHub OAuth、Google OAuth、JWT 令牌管理及权限控制。

---

## 1. 认证方式总览

| 登录方式 | 技术方案 | 适用场景 |
|----------|----------|----------|
| 邮箱 + 密码 | bcrypt 哈希 + JWT + RefreshToken | 常规用户注册登录 |
| GitHub OAuth | OAuth2 Authorization Code Flow | 开发者快捷登录 |
| Google OAuth | OAuth2 Authorization Code Flow | 国际用户快捷登录 |

---

## 2. 核心概念定义

| 概念 | 说明 |
|------|------|
| **Access Token** | JWT，有效期 2 小时，携带用户 ID、角色、权限 |
| **Refresh Token** | 随机字符串，有效期 7 天，存储于 Redis / HttpOnly Cookie |
| **OAuth State** | 防止 CSRF 攻击的随机状态值，存储于 Redis，有效期 10 分钟 |
| **Token Blacklist** | 登出时 RefreshToken 加入 Redis 黑名单，AccessToken 等待自然过期 |

---

## 3. 邮箱密码登录流程

```
用户提交邮箱 + 密码
        │
        ▼
Gateway / User Service
        ├── 查询 User 表 (Ent)
        ├── bcrypt.CompareHashAndPassword(密码)
        ├── 生成 AccessToken (JWT)
        ├── 生成 RefreshToken (随机 32 字节，Base64)
        ├── Redis 存储 RefreshToken (key: refresh:{user_id}:{token}, TTL: 7d)
        └── 返回 { access_token, refresh_token, expires_in }
```

---

## 4. GitHub OAuth 登录流程

```
1. 前端点击 "GitHub 登录"
        │
        ▼
2. 前端请求后端 /auth/oauth/github
        │
        ▼
3. 后端生成 state，存入 Redis（TTL 10min）
        │
        ▼
4. 后端返回 GitHub 授权 URL
   https://github.com/login/oauth/authorize?client_id=xxx&state=xxx&...
        │
        ▼
5. 前端跳转 GitHub 授权页
        │
        ▼
6. GitHub 回调到 /auth/oauth/callback?code=xxx&state=xxx
        │
        ▼
7. 后端校验 state，用 code 换 access_token
        │
        ▼
8. 用 access_token 请求 GitHub API 获取用户邮箱 + ID
        │
        ▼
9. 后端查询本地用户绑定关系 (oauth_provider + oauth_id)
        ├── 已存在 → 直接签发 JWT + RefreshToken
        ├── 不存在但邮箱已注册 → 返回 "需要绑定" 状态码，引导用户登录后关联
        └── 完全新用户 → 自动创建账号（生成随机用户名），签发 JWT
```

### Google OAuth 流程

与 GitHub 完全一致，仅替换 OAuth 配置：

- 授权端点：`https://accounts.google.com/o/oauth2/v2/auth`
- Token 端点：`https://oauth2.googleapis.com/token`
- 用户信息端点：`https://www.googleapis.com/oauth2/v2/userinfo`

---

## 5. JWT 规范

### 5.1 Token 载荷结构

```json
{
  "sub": "1000001",
  "email": "user@example.com",
  "username": "johndoe",
  "role": "user",
  "iat": 1713273600,
  "exp": 1713280800,
  "jti": "uuid-token-id"
}
```

### 5.2 密钥管理

- **签名算法**：HS256（初期）或 RS256（多服务验证时推荐）
- **密钥来源**：环境变量 `JWT_SECRET`（HS256）或私钥文件（RS256）
- **密钥轮换**：RS256 支持公钥分发，便于无状态服务校验

### 5.3 Go 代码示例

```go
// pkg/auth/jwt.go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   uint64 `json:"sub"`
    Email    string `json:"email"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateAccessToken(userID uint64, email, username, role, secret string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Email:    email,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   fmt.Sprintf("%d", userID),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func ParseAccessToken(tokenStr, secret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
        return []byte(secret), nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, jwt.ErrSignatureInvalid
}
```

---

## 6. OAuth 用户绑定策略

### 6.1 绑定状态矩阵

| oauth_provider + oauth_id | 邮箱已注册 | 处理方式 |
|---------------------------|-----------|----------|
| 存在 | - | 直接登录，签发 Token |
| 不存在 | 存在 | 返回 `NEED_BIND_OAUTH`（前端引导登录后绑定） |
| 不存在 | 不存在 | 自动创建新用户，随机用户名如 `user_abc123`，签发 Token |

### 6.2 绑定 API

```http
POST /api/v1/auth/bind-oauth
Authorization: Bearer <AccessToken>
Content-Type: application/json

{
  "provider": "github",
  "oauth_id": "12345678",
  "oauth_email": "github@example.com"
}
```

---

## 7. Token 刷新与登出

### 7.1 刷新 AccessToken

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "xxx"
}
```

**后端逻辑**：
1. 校验 Redis 中是否存在 `refresh:{user_id}:{token}`
2. 若不存在或已过期 → 返回 401，要求重新登录
3. 生成新的 AccessToken 与 RefreshToken
4. 删除旧 RefreshToken，存入新 RefreshToken
5. 返回新 Token 对

### 7.2 登出

```http
POST /api/v1/auth/logout
Authorization: Bearer <AccessToken>
Content-Type: application/json

{
  "refresh_token": "xxx"
}
```

**后端逻辑**：
1. 解析 AccessToken 获取 UserID
2. 从 Redis 删除 `refresh:{user_id}:{refresh_token}`
3. （可选）将 AccessToken 的 `jti` 加入黑名单，TTL 设为剩余过期时间
4. 返回成功

---

## 8. RBAC 权限控制（Casbin）

### 8.1 角色定义

| 角色 | 权限 |
|------|------|
| `guest` | 只读访问 |
| `user` | 发布文章、评论、点赞、收藏 |
| `creator` | 进入创作者中心、创建专栏 |
| `reviewer` | 审核文章、处理举报 |
| `admin` | 全站管理权限 |

### 8.2 Casbin Model 示例

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

### 8.3 Policy 示例

```csv
p, admin, *, *
p, reviewer, post, review
p, user, post, create
p, user, comment, create
p, guest, post, read

g, alice, admin
g, bob, user
```

### 8.4 Gin 中间件集成

```go
func Authorize(obj, act string, enforcer *casbin.Enforcer) gin.HandlerFunc {
    return func(c *gin.Context) {
        role := c.GetString("role") // 从 JWT Context 提取
        allowed, _ := enforcer.Enforce(role, obj, act)
        if !allowed {
            c.AbortWithStatusJSON(403, gin.H{"code": "E403001", "message": "权限不足"})
            return
        }
        c.Next()
    }
}
```

---

## 9. 安全加固清单

- [x] OAuth `state` 参数校验，防止 CSRF
- [x] JWT 使用强密钥，定期轮换
- [x] RefreshToken 仅存于 Redis / HttpOnly Cookie，前端不可读取
- [x] 密码使用 bcrypt（cost >= 10）哈希存储
- [x] 登录接口限流（如 5 次/分钟/IP）
- [x] 敏感操作二次验证（预留）
- [x] HTTPS 全站加密
