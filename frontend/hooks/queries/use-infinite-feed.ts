import { useInfiniteQuery } from "@tanstack/react-query";
import apiClient from "@/lib/api-client";
import type { Post } from "@/types/post";

interface FeedParams {
  type: "latest" | "hot" | "following";
  pageSize?: number;
  period?: "24h" | "7d" | "30d";
}

interface FeedResponse {
  list: Post[];
  pagination: {
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
  };
}

export function useInfiniteFeed({ type, pageSize = 20, period }: FeedParams) {
  return useInfiniteQuery({
    queryKey: ["feed", type, period],
    queryFn: async ({ pageParam = 1 }) => {
      const params: Record<string, string | number> = {
        type,
        page: pageParam,
        pageSize,
      };

      if (type === "hot" && period) {
        params.period = period;
      }

      const response = await apiClient
        .get("feed", { searchParams: params })
        .json<FeedResponse>();

      return response;
    },
    getNextPageParam: (lastPage) => {
      const { page, totalPages } = lastPage.pagination;
      return page < totalPages ? page + 1 : undefined;
    },
    initialPageParam: 1,
    staleTime: 1000 * 60 * 5, // 5分钟
  });
}
