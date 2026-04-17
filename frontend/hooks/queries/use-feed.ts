import { useQuery } from "@tanstack/react-query";
import apiClient, { ApiError } from "@/lib/api-client";
import type { Post, PostsResponse } from "@/types/post";

export type FeedType = "latest" | "hot" | "following" | "recommend";

interface FeedParams {
  type?: FeedType;
  page?: number;
  pageSize?: number;
  period?: string; // "24h" | "7d" | "30d" for hot
}

export interface FeedResponse {
  list: Post[];
  pagination: {
    total: number;
    page: number;
    pageSize: number;
    totalPages: number;
  };
}

export function useFeed(params: FeedParams = {}) {
  return useQuery<FeedResponse, ApiError>({
    queryKey: ["feed", params],
    queryFn: () =>
      apiClient
        .get("feed", {
          searchParams: {
            type: params.type ?? "latest",
            page: params.page ?? 1,
            pageSize: params.pageSize ?? 20,
            ...(params.period ? { period: params.period } : {}),
          },
        })
        .json<FeedResponse>(),
  });
}
