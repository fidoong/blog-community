"use client"

import { useEffect, useRef } from "react";
import { PostCard, PostCardSkeleton } from "@/components/features/post-card";
import type { Post } from "@/types/post";

interface InfinitePostListProps {
  posts: Post[];
  isLoading: boolean;
  isFetchingNextPage: boolean;
  hasNextPage: boolean;
  onLoadMore: () => void;
}

export function InfinitePostList({
  posts,
  isLoading,
  isFetchingNextPage,
  hasNextPage,
  onLoadMore,
}: InfinitePostListProps) {
  const observerRef = useRef<IntersectionObserver | null>(null);
  const loadMoreRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (isLoading || isFetchingNextPage || !hasNextPage) return;

    observerRef.current = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          onLoadMore();
        }
      },
      { threshold: 0.1 }
    );

    if (loadMoreRef.current) {
      observerRef.current.observe(loadMoreRef.current);
    }

    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, [isLoading, isFetchingNextPage, hasNextPage, onLoadMore]);

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
    <div className="h-full overflow-y-auto">
      <div className="divide-y">
        {posts.map((post, index) => (
          <div key={post.id} className="py-4 first:pt-0">
            <PostCard post={post} index={index} />
          </div>
        ))}
      </div>

      {/* 加载触发器 */}
      {hasNextPage && (
        <div ref={loadMoreRef} className="py-4 text-center">
          {isFetchingNextPage ? (
            <div className="inline-flex items-center gap-2 text-sm text-muted-foreground">
              <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
              <span>加载中...</span>
            </div>
          ) : (
            <div className="h-4" />
          )}
        </div>
      )}

      {!hasNextPage && posts.length > 0 && (
        <div className="py-8 text-center text-sm text-muted-foreground">
          已经到底了
        </div>
      )}
    </div>
  );
}
