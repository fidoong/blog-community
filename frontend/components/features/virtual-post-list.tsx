"use client"

import { PostCard, PostCardSkeleton } from "@/components/features/post-card";
import { useVirtualScroll } from "@/hooks/use-virtual-scroll";
import type { Post } from "@/types/post";

interface VirtualPostListProps {
  posts: Post[];
  isLoading: boolean;
  isFetchingNextPage: boolean;
  hasNextPage: boolean;
  onLoadMore: () => void;
}

export function VirtualPostList({
  posts,
  isLoading,
  isFetchingNextPage,
  hasNextPage,
  onLoadMore,
}: VirtualPostListProps) {
  const { parentRef, virtualizer, items } = useVirtualScroll({
    count: posts.length,
    estimateSize: 120,
    overscan: 3,
    onLoadMore,
    hasNextPage,
    isFetchingNextPage,
    threshold: 3,
  });

  if (isLoading) {
    return (
      <div className="divide-y">
        {[...Array(5)].map((_, i) => (
          <div key={i} className="py-4 first:pt-0">
            <PostCardSkeleton />
          </div>
        ))}
      </div>
    );
  }

  if (posts.length === 0) {
    return (
      <div className="py-12 text-center">
        <p className="text-sm text-muted-foreground">暂无文章，快来发布第一篇吧！</p>
      </div>
    );
  }

  return (
    <div
      ref={parentRef}
      className="h-full overflow-auto scrollbar-thin"
    >
      <div
        style={{
          height: `${virtualizer.getTotalSize()}px`,
          width: "100%",
          position: "relative",
        }}
      >
        <div className="divide-y">
          {items.map((virtualItem) => {
            const post = posts[virtualItem.index];
            if (!post) return null;
            
            return (
              <div
                key={virtualItem.key}
                data-index={virtualItem.index}
                ref={virtualizer.measureElement}
                style={{
                  position: "absolute",
                  top: 0,
                  left: 0,
                  width: "100%",
                  transform: `translateY(${virtualItem.start}px)`,
                }}
              >
                <div className="py-4 first:pt-0 bg-background">
                  <PostCard post={post} index={virtualItem.index} />
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* 加载更多指示器 */}
      {isFetchingNextPage && (
        <div className="py-4 text-center bg-background">
          <div className="inline-flex items-center gap-2 text-sm text-muted-foreground">
            <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
            <span>加载中...</span>
          </div>
        </div>
      )}

      {!hasNextPage && posts.length > 0 && (
        <div className="py-8 text-center text-sm text-muted-foreground bg-background">
          已经到底了
        </div>
      )}
    </div>
  );
}

