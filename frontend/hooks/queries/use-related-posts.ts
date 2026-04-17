import { useQuery } from "@tanstack/react-query";
import apiClient from "@/lib/api-client";
import type { Post } from "@/types/post";

interface RelatedPostsResponse {
  list: Post[];
}

export function useRelatedPosts(postId: string | number, limit = 5) {
  return useQuery<RelatedPostsResponse>({
    queryKey: ["related-posts", postId, limit],
    queryFn: () =>
      apiClient
        .get(`posts/${postId}/related`, { searchParams: { limit } })
        .json<RelatedPostsResponse>(),
    enabled: !!postId,
  });
}
