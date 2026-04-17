import { useQuery } from "@tanstack/react-query";
import apiClient, { ApiError } from "@/lib/api-client";
import type { SearchResponse } from "@/types/search";

interface SearchParams {
  keyword: string;
  page?: number;
  pageSize?: number;
}

export function useSearch(params: SearchParams) {
  return useQuery<SearchResponse, ApiError>({
    queryKey: ["search", params.keyword, params.page, params.pageSize],
    queryFn: () =>
      apiClient
        .get("search", {
          searchParams: {
            q: params.keyword,
            page: params.page ?? 1,
            pageSize: params.pageSize ?? 20,
          },
        })
        .json<SearchResponse>(),
    enabled: params.keyword.trim().length > 0,
  });
}
