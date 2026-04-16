# 08 - OpenAPI 接口规范模板

> 本文档定义前后端协作的接口契约规范，所有 REST API 必须遵循本模板。

---

## 1. 通用规范

### 1.1 基础信息

- **协议**：HTTPS
- **数据格式**：JSON
- **字符编码**：UTF-8
- **Content-Type**：`application/json`
- **版本前缀**：`/api/v1`

### 1.2 认证方式

- 除公开接口外，其余接口需在请求头携带：
  - `Authorization: Bearer <AccessToken>`

### 1.3 分页规范

请求参数：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `page` | int | 否 | 页码，从 1 开始，默认 1 |
| `pageSize` | int | 否 | 每页条数，默认 20，最大 100 |

响应结构：

```json
{
  "code": "OK",
  "message": "success",
  "data": {
    "list": [],
    "pagination": {
      "total": 100,
      "page": 1,
      "pageSize": 20,
      "totalPages": 5
    }
  },
  "timestamp": 1713273600
}
```

---

## 2. 统一响应结构

### 2.1 成功响应

```json
{
  "code": "OK",
  "message": "success",
  "data": { ... },
  "timestamp": 1713273600
}
```

### 2.2 错误响应

```json
{
  "code": "E400001",
  "message": "请求参数错误",
  "timestamp": 1713273600
}
```

---

## 3. 接口定义示例

### 3.1 用户认证

#### POST /api/v1/auth/register
邮箱注册

**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "username": "johndoe"
}
```

**Response (201):**
```json
{
  "code": "OK",
  "message": "success",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "dGhpcyBpcyBhIHJlZnJlc2g...",
    "expiresIn": 7200,
    "user": {
      "id": 1000001,
      "email": "user@example.com",
      "username": "johndoe",
      "avatarUrl": ""
    }
  },
  "timestamp": 1713273600
}
```

---

#### POST /api/v1/auth/login
邮箱登录

**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response (200):**
同注册响应

---

#### POST /api/v1/auth/refresh
刷新 Token

**Request:**
```json
{
  "refreshToken": "dGhpcyBpcyBhIHJlZnJlc2g..."
}
```

**Response (200):**
```json
{
  "code": "OK",
  "data": {
    "accessToken": "eyJhbGci...",
    "refreshToken": "newRefresh...",
    "expiresIn": 7200
  }
}
```

---

#### POST /api/v1/auth/logout
登出

**Request:**
```json
{
  "refreshToken": "dGhpcyBpcyBhIHJlZnJlc2g..."
}
```

**Response (200):**
```json
{ "code": "OK", "message": "success" }
```

---

#### GET /api/v1/auth/oauth/github
获取 GitHub OAuth 授权 URL

**Response (200):**
```json
{
  "code": "OK",
  "data": {
    "authUrl": "https://github.com/login/oauth/authorize?client_id=xxx&state=xxx&..."
  }
}
```

---

#### GET /api/v1/auth/oauth/callback?code=xxx&state=xxx&provider=github
OAuth 回调

**Response (200):**
同注册响应，或返回 `NEED_BIND` 状态（需前端处理绑定流程）

---

### 3.2 文章管理

#### GET /api/v1/posts
获取文章列表

**Query Parameters:**
- `page` (int)
- `pageSize` (int)
- `sort` (string): `hot` | `new` | `recommend`
- `tag` (string, optional)

**Response (200):**
```json
{
  "code": "OK",
  "data": {
    "list": [
      {
        "id": "post_001",
        "title": "深入浅出 Go 并发模式",
        "summary": "本文介绍 Go 语言中的 goroutine 与 channel...",
        "author": {
          "id": 1000001,
          "username": "johndoe",
          "avatarUrl": ""
        },
        "coverImage": "",
        "tags": ["Go", "并发"],
        "likeCount": 128,
        "commentCount": 32,
        "readCount": 1024,
        "publishedAt": "2024-04-15T10:00:00Z"
      }
    ],
    "pagination": {
      "total": 100,
      "page": 1,
      "pageSize": 20,
      "totalPages": 5
    }
  }
}
```

---

#### GET /api/v1/posts/:id
获取文章详情

**Response (200):**
```json
{
  "code": "OK",
  "data": {
    "id": "post_001",
    "title": "深入浅出 Go 并发模式",
    "content": "# 前言\n\nGo 语言...",
    "contentType": "markdown",
    "author": { ... },
    "tags": ["Go", "并发"],
    "likeCount": 128,
    "commentCount": 32,
    "readCount": 1024,
    "isLiked": false,
    "isCollected": false,
    "publishedAt": "2024-04-15T10:00:00Z",
    "updatedAt": "2024-04-15T12:00:00Z"
  }
}
```

---

#### POST /api/v1/posts
发布文章（需登录）

**Request:**
```json
{
  "title": "深入浅出 Go 并发模式",
  "content": "# 前言\n\nGo 语言...",
  "contentType": "markdown",
  "tags": ["Go", "并发"],
  "coverImage": ""
}
```

**Response (201):**
```json
{
  "code": "OK",
  "data": {
    "id": "post_001",
    "status": "pending"
  }
}
```

---

#### PUT /api/v1/posts/:id
更新文章（需作者本人或管理员）

**Request:**
同发布文章

**Response (200):**
同文章详情

---

#### DELETE /api/v1/posts/:id
删除文章（需作者本人或管理员）

**Response (200):**
```json
{ "code": "OK", "message": "success" }
```

---

### 3.3 评论管理

#### GET /api/v1/posts/:id/comments
获取文章评论

**Response (200):**
```json
{
  "code": "OK",
  "data": {
    "list": [
      {
        "id": "cmt_001",
        "content": "写得很好！",
        "author": { ... },
        "likeCount": 5,
        "replies": [
          {
            "id": "cmt_002",
            "content": "感谢支持！",
            "author": { ... },
            "parentId": "cmt_001",
            "likeCount": 2
          }
        ],
        "createdAt": "2024-04-15T11:00:00Z"
      }
    ]
  }
}
```

---

#### POST /api/v1/posts/:id/comments
发表评论

**Request:**
```json
{
  "content": "写得很好！",
  "parentId": "" // 为空表示一级评论
}
```

---

### 3.4 互动

#### POST /api/v1/posts/:id/like
点赞 / 取消点赞

**Response (200):**
```json
{
  "code": "OK",
  "data": {
    "isLiked": true,
    "likeCount": 129
  }
}
```

---

#### POST /api/v1/posts/:id/collect
收藏 / 取消收藏

**Response (200):**
```json
{
  "code": "OK",
  "data": {
    "isCollected": true,
    "collectCount": 45
  }
}
```

---

### 3.5 用户

#### GET /api/v1/users/:id
获取用户主页信息

**Response (200):**
```json
{
  "code": "OK",
  "data": {
    "id": 1000001,
    "username": "johndoe",
    "bio": "Go 后端开发者",
    "avatarUrl": "",
    "followersCount": 128,
    "followingCount": 56,
    "postsCount": 12,
    "isFollowing": false
  }
}
```

---

#### POST /api/v1/users/:id/follow
关注 / 取消关注

**Response (200):**
```json
{
  "code": "OK",
  "data": {
    "isFollowing": true
  }
}
```

---

## 4. Swagger 注解规范（Go）

```go
// @Summary      用户注册
// @Description  使用邮箱和密码注册新用户
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  RegisterRequest  true  "注册信息"
// @Success      201  {object}  response.Response{data=AuthResponse}
// @Failure      400  {object}  response.Response
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) { ... }
```

生成命令：

```bash
swag init -g cmd/gateway/main.go -o docs/swagger
```
