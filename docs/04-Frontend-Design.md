# 04 - 前端企业级设计方案

> 基于 Next.js 15 (App Router) 的企业级前端工程实践。

---

## 1. 项目目录结构

```
frontend/
├── app/                        # App Router (Next.js 15)
│   ├── (main)/                 # 用户端主站 (Route Group)
│   │   ├── page.tsx            # 首页 Feed
│   │   ├── post/[id]/page.tsx  # 文章详情 (SSR)
│   │   ├── user/[id]/page.tsx  # 个人主页
│   │   ├── search/page.tsx     # 搜索结果
│   │   ├── login/page.tsx      # 登录页
│   │   ├── settings/page.tsx   # 用户设置
│   │   └── layout.tsx          # 主站公共布局 (Header + Footer)
│   ├── admin/                  # 管理后台 (Route Group)
│   │   ├── layout.tsx          # Admin 布局 (Sidebar + Topbar)
│   │   ├── dashboard/page.tsx
│   │   ├── posts/page.tsx
│   │   └── users/page.tsx
│   ├── api/                    # Next.js API Routes (BFF)
│   │   ├── auth/
│   │   │   ├── login/route.ts
│   │   │   ├── refresh/route.ts
│   │   │   └── oauth/
│   │   │       ├── github/route.ts
│   │   │       └── callback/route.ts
│   │   └── upload/route.ts
│   ├── layout.tsx              # 根布局 (Providers, Fonts, Theme)
│   └── globals.css             # 全局样式 + Tailwind + CSS Variables
├── components/
│   ├── ui/                     # 原子/基础组件 (shadcn/ui 风格)
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   ├── avatar.tsx
│   │   ├── badge.tsx
│   │   ├── dialog.tsx
│   │   ├── dropdown-menu.tsx
│   │   ├── input.tsx
│   │   ├── skeleton.tsx
│   │   ├── toast.tsx
│   │   └── ...
│   ├── shared/                 # 业务公共组件
│   │   ├── header.tsx          # 顶部导航
│   │   ├── footer.tsx          # 页脚
│   │   ├── post-card.tsx       # 文章卡片
│   │   ├── comment-tree.tsx    # 嵌套评论树
│   │   ├── user-profile.tsx    # 用户信息卡片
│   │   ├── theme-toggle.tsx    # 明暗主题切换
│   │   └── pagination.tsx      # 分页器
│   └── features/               # 功能域组件 (按领域拆分)
│       ├── auth/
│       │   ├── login-form.tsx
│       │   └── oauth-buttons.tsx
│       ├── editor/
│       │   ├── rich-editor.tsx
│   │   │   ├── markdown-editor.tsx
│   │   │   └── publish-modal.tsx
│       ├── feed/
│       │   ├── feed-tabs.tsx
│   │   │   ├── post-list.tsx
│   │   │   └── trending-sidebar.tsx
│       └── profile/
│           ├── stats-chart.tsx
│           └── post-grid.tsx
├── hooks/                      # 自定义 Hooks
│   ├── use-auth.ts
│   ├── use-theme.ts
│   ├── use-media-query.ts
│   ├── use-debounce.ts
│   ├── use-local-storage.ts
│   └── use-infinite-scroll.ts
├── lib/
│   ├── api-client.ts           # HTTP 客户端封装
│   ├── error-handler.ts        # 前端统一错误处理
│   ├── utils.ts                # 通用工具 (cn, formatDate)
│   └── constants.ts            # 常量定义
├── stores/                     # Zustand 状态
│   ├── auth-store.ts
│   ├── theme-store.ts
│   └── ui-store.ts
├── types/                      # 全局类型定义
│   ├── api.ts
│   ├── user.ts
│   ├── post.ts
│   └── comment.ts
├── styles/
│   └── globals.css
├── public/
│   └── images/
└── tests/                      # 单元/集成测试
    ├── unit/
    └── integration/
```

---

## 2. UI 设计体系与配色（参考 Next.js / Vercel）

### 2.1 设计原则

- **极简主义**：大量留白，信息层次分明
- **高对比度**：文字与背景对比度符合 WCAG AA 标准
- **精致圆角**：卡片、按钮统一圆角 (`--radius: 0.5rem`)
- **细腻动效**：Hover 状态使用 `transition-all duration-200`
- **深色模式**：通过 CSS Variables 一键切换

### 2.2 CSS Variables 配色方案

```css
:root {
  /* 背景 */
  --background: #ffffff;
  --foreground: #000000;
  --muted: #f5f5f5;
  --muted-foreground: #737373;
  
  /* 边框与卡片 */
  --border: #e5e5e5;
  --card: #ffffff;
  --card-foreground: #000000;
  
  /* 主题色（Next.js 蓝） */
  --primary: #0070f3;
  --primary-foreground: #ffffff;
  --secondary: #f3f3f3;
  --secondary-foreground: #000000;
  
  /* 状态色 */
  --accent: #fafafa;
  --accent-foreground: #000000;
  --destructive: #ef4444;
  --destructive-foreground: #ffffff;
  
  /* 圆角 */
  --radius: 0.5rem;
}

.dark {
  --background: #000000;
  --foreground: #ffffff;
  --muted: #171717;
  --muted-foreground: #a3a3a3;
  --border: #262626;
  --card: #0a0a0a;
  --card-foreground: #ffffff;
  --primary: #0070f3;
  --secondary: #1a1a1a;
  --accent: #1a1a1a;
}
```

### 2.3 Tailwind 扩展类（v4 写法）

```css
/* globals.css */
@import "tailwindcss";

@theme {
  --font-sans: 'Geist', 'Inter', system-ui, sans-serif;
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --color-primary: var(--primary);
  --color-muted: var(--muted);
  --color-border: var(--border);
  --radius-default: var(--radius);
}
```

---

## 3. HTTP 客户端封装

文件：`lib/api-client.ts`

```typescript
import ky from 'ky';
import { useAuthStore } from '@/stores/auth-store';

const apiClient = ky.create({
  prefixUrl: process.env.NEXT_PUBLIC_API_BASE_URL,
  headers: { 'Content-Type': 'application/json' },
  timeout: 10000,
  hooks: {
    beforeRequest: [
      (request) => {
        const token = useAuthStore.getState().token;
        if (token) {
          request.headers.set('Authorization', `Bearer ${token}`);
        }
      },
    ],
    afterResponse: [
      async (request, options, response) => {
        if (!response.ok) {
          const body = await response.json<{ message?: string; code?: string }>().catch(() => ({}));
          throw new ApiError(body.message || '请求失败', response.status, body.code);
        }
        return response;
      },
    ],
  },
});

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public code?: string
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

export default apiClient;
```

---

## 4. TanStack Query 封装规范

文件：`hooks/queries/use-posts.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import apiClient, { ApiError } from '@/lib/api-client';
import { toast } from 'sonner';
import { handleApiError } from '@/lib/error-handler';

const postKeys = {
  all: ['posts'] as const,
  lists: () => [...postKeys.all, 'list'] as const,
  list: (filters: PostFilters) => [...postKeys.lists(), filters] as const,
  details: () => [...postKeys.all, 'detail'] as const,
  detail: (id: string) => [...postKeys.details(), id] as const,
};

export function usePostList(filters: PostFilters) {
  return useQuery({
    queryKey: postKeys.list(filters),
    queryFn: async () => {
      return apiClient.get('posts', { searchParams: filters }).json<PostListResp>();
    },
    staleTime: 60 * 1000,
  });
}

export function useCreatePost() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreatePostPayload) =>
      apiClient.post('posts', { json: payload }).json<Post>(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: postKeys.lists() });
      toast.success('发布成功');
    },
    onError: (error) => handleApiError(error),
  });
}
```

---

## 5. 统一错误处理

文件：`lib/error-handler.ts`

```typescript
import { ApiError } from './api-client';
import { toast } from 'sonner';

export function handleApiError(error: unknown) {
  if (error instanceof ApiError) {
    if (error.status === 401) {
      window.location.href = '/login';
      return;
    }
    if (error.status === 403) {
      toast.error('权限不足');
      return;
    }
    if (error.status >= 500) {
      toast.error('服务器繁忙，请稍后重试');
      return;
    }
    toast.error(error.message);
  } else {
    toast.error('网络异常，请检查连接');
  }
}
```

---

## 6. 通用类型定义

文件：`types/api.ts`

```typescript
export interface ApiResponse<T> {
  code: string;
  message: string;
  data: T;
  timestamp: number;
}

export interface PaginationParams {
  page: number;
  pageSize: number;
}

export interface PaginationMeta {
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface PaginatedResponse<T> {
  list: T[];
  pagination: PaginationMeta;
}
```

---

## 7. 页面渲染策略

| 页面 | 策略 | 原因 |
|------|------|------|
| 文章详情 `/post/[id]` | **SSR (RSC)** | SEO 核心，需要首屏 HTML 与结构化数据 |
| 首页 Feed | **Streaming SSR** | 快速首字节，推荐内容流式加载 |
| 用户主页 `/user/[id]` | **SSR** | SEO 与社交分享 |
| 搜索结果 `/search` | **SSR** | SEO，动态搜索参数 |
| 创作者后台 `/admin/*` | **CSR** | 强交互，无需 SEO |
| 文章发布 `/editor` | **CSR** | 富文本编辑器较重，客户端渲染更流畅 |
| 登录/注册 | **CSR** | 表单交互为主 |

---

## 8. 状态管理策略

| 状态类型 | 管理方式 | 示例 |
|----------|----------|------|
| 服务端状态 | TanStack Query | 文章列表、用户详情、评论数据 |
| 全局客户端状态 | Zustand | 当前登录用户、主题模式、全局 Toast |
| 局部 UI 状态 | useState / useReducer | Modal 开关、表单输入、Tab 索引 |
| URL 状态 | Next.js SearchParams | 搜索关键词、筛选条件、分页页码 |
