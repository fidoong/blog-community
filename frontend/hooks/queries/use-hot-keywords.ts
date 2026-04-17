import { useQuery } from "@tanstack/react-query";
import apiClient from "@/lib/api-client";

interface HotKeywordsResponse {
  list: string[];
}

export function useHotKeywords(limit = 10) {
  return useQuery<HotKeywordsResponse>({
    queryKey: ["hot-keywords", limit],
    queryFn: () =>
      apiClient
        .get("search/hot", { searchParams: { limit } })
        .json<HotKeywordsResponse>(),
    staleTime: 1000 * 60 * 2, // 2分钟
  });
}
