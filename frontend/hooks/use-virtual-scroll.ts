import { useRef, useEffect } from "react";
import { useVirtualizer } from "@tanstack/react-virtual";

interface UseVirtualScrollOptions {
  count: number;
  estimateSize?: number;
  overscan?: number;
  onLoadMore?: () => void;
  hasNextPage?: boolean;
  isFetchingNextPage?: boolean;
  threshold?: number; // 距离底部多少项时触发加载
}

/**
 * 通用虚拟滚动 Hook
 * 支持无限滚动和预加载
 */
export function useVirtualScroll({
  count,
  estimateSize = 150,
  overscan = 5,
  onLoadMore,
  hasNextPage = false,
  isFetchingNextPage = false,
  threshold = 5,
}: UseVirtualScrollOptions) {
  const parentRef = useRef<HTMLDivElement>(null);

  const virtualizer = useVirtualizer({
    count,
    getScrollElement: () => parentRef.current,
    estimateSize: () => estimateSize,
    overscan,
  });

  const items = virtualizer.getVirtualItems();

  // 监听滚动，触发加载更多
  useEffect(() => {
    if (!onLoadMore) return;

    const [lastItem] = [...items].reverse();
    if (!lastItem) return;

    // 当滚动到倒数第 threshold 项时，触发加载
    if (
      lastItem.index >= count - threshold &&
      hasNextPage &&
      !isFetchingNextPage
    ) {
      onLoadMore();
    }
  }, [items, count, hasNextPage, isFetchingNextPage, onLoadMore, threshold]);

  return {
    parentRef,
    virtualizer,
    items,
  };
}
