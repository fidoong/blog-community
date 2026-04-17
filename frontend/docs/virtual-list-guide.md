# 虚拟列表和无限滚动使用指南

## 概述

本项目使用 `@tanstack/react-virtual` 实现虚拟列表滚动，配合 `@tanstack/react-query` 的无限查询功能，实现高性能的列表渲染和懒加载。

## 核心组件和 Hook

### 1. useInfiniteFeed Hook

位置：`/hooks/queries/use-infinite-feed.ts`

用于获取无限滚动的文章列表数据。

```typescript
const { 
  data,           // 所有页面的数据
  isLoading,      // 首次加载状态
  isFetchingNextPage,  // 加载下一页状态
  hasNextPage,    // 是否还有下一页
  fetchNextPage   // 加载下一页的函数
} = useInfiniteFeed({
  type: "latest",  // "latest" | "hot"
  pageSize: 20,
  period: "7d"     // 仅 hot 类型需要
});
```

### 2. useVirtualScroll Hook

位置：`/hooks/use-virtual-scroll.ts`

通用的虚拟滚动 Hook，可复用于任何列表场景。

```typescript
const { parentRef, virtualizer, items } = useVirtualScroll({
  count: posts.length,        // 列表项总数
  estimateSize: 150,          // 预估每项高度（px）
  overscan: 5,                // 预渲染项数
  onLoadMore,                 // 加载更多回调
  hasNextPage,                // 是否有下一页
  isFetchingNextPage,         // 是否正在加载
  threshold: 5                // 距离底部多少项时触发加载
});
```

### 3. VirtualPostList 组件

位置：`/components/features/virtual-post-list.tsx`

封装好的文章虚拟列表组件。

```typescript
<VirtualPostList
  posts={allPosts}
  isLoading={isLoading}
  isFetchingNextPage={isFetchingNextPage}
  hasNextPage={hasNextPage}
  onLoadMore={() => fetchNextPage()}
/>
```

## 使用示例

### 首页实现

```typescript
export default function HomePage() {
  const [activeTab, setActiveTab] = useState<FeedTab>("latest");
  
  // 1. 使用无限查询
  const { 
    data, 
    isLoading, 
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage 
  } = useInfiniteFeed({
    type: activeTab === "hot" ? "hot" : "latest",
    pageSize: 20,
    ...(activeTab === "hot" ? { period: "7d" } : {}),
  });

  // 2. 扁平化所有页面的数据
  const allPosts = useMemo(() => {
    return data?.pages.flatMap((page) => page.list) ?? [];
  }, [data]);

  // 3. 使用虚拟列表组件
  return (
    <VirtualPostList
      posts={allPosts}
      isLoading={isLoading}
      isFetchingNextPage={isFetchingNextPage}
      hasNextPage={hasNextPage ?? false}
      onLoadMore={() => fetchNextPage()}
    />
  );
}
```

## 布局优化

### 固定 Tab 栏，只滚动列表

```tsx
<main className="min-w-0 flex flex-col">
  {/* 固定的 Tab 栏 */}
  <div className="flex-shrink-0 mb-4 border-b bg-background">
    {/* Tab 内容 */}
  </div>

  {/* 可滚动的列表区域 */}
  <div className="flex-1 min-h-0">
    <VirtualPostList {...props} />
  </div>
</main>
```

### 容器高度设置

```tsx
<Container size="full" className="py-6">
  <div className="grid grid-cols-1 lg:grid-cols-[200px_1fr_300px] gap-6 h-[calc(100vh-8rem)]">
    {/* 左侧导航 - 可滚动 */}
    <aside className="hidden lg:block overflow-y-auto">
      <LeftNav />
    </aside>

    {/* 主内容区 - flex 布局 */}
    <main className="min-w-0 flex flex-col">
      {/* ... */}
    </main>

    {/* 右侧栏 - 可滚动 */}
    <aside className="hidden lg:block overflow-y-auto">
      {/* ... */}
    </aside>
  </div>
</Container>
```

## 性能优化要点

1. **虚拟滚动**：只渲染可见区域的 DOM，大幅减少渲染开销
2. **预渲染（overscan）**：提前渲染上下各 5 项，避免滚动时白屏
3. **预加载（threshold）**：距离底部 5 项时提前加载下一页
4. **数据缓存**：React Query 自动缓存数据，切换 Tab 时无需重新请求
5. **动画优化**：使用 Framer Motion 的 `layoutId` 实现流畅的 Tab 切换动画

## 扩展使用

### 创建其他类型的虚拟列表

```typescript
// 1. 创建数据查询 Hook
export function useInfiniteComments(postId: string) {
  return useInfiniteQuery({
    queryKey: ["comments", postId],
    queryFn: async ({ pageParam = 1 }) => {
      // 获取评论数据
    },
    getNextPageParam: (lastPage) => {
      // 返回下一页页码
    },
    initialPageParam: 1,
  });
}

// 2. 创建虚拟列表组件
export function VirtualCommentList({ comments, ... }) {
  const { parentRef, virtualizer, items } = useVirtualScroll({
    count: comments.length,
    estimateSize: 100,  // 评论项高度
    // ...
  });

  return (
    <div ref={parentRef} className="h-full overflow-auto">
      {/* 渲染逻辑 */}
    </div>
  );
}
```

## 注意事项

1. **estimateSize**：尽量接近实际高度，避免滚动条跳动
2. **容器高度**：虚拟列表容器必须有固定高度
3. **key 值**：使用 `virtualItem.key` 而非 `index`，避免重复渲染
4. **动画性能**：虚拟列表内的动画要谨慎使用，避免影响滚动性能
