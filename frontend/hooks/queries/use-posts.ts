import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient, { ApiError } from "@/lib/api-client";
import { handleApiError } from "@/lib/error-handler";
import type { Post, PostsResponse, CreatePostPayload } from "@/types/post";

interface ListPostsParams {
  page?: number;
  pageSize?: number;
  sort?: "new" | "hot";
  status?: string;
  authorId?: string | number;
}

export function usePosts(params: ListPostsParams = {}) {
  return useQuery<PostsResponse, ApiError>({
    queryKey: ["posts", params],
    queryFn: () =>
      apiClient
        .get("posts", {
          searchParams: {
            page: params.page ?? 1,
            pageSize: params.pageSize ?? 20,
            sort: params.sort ?? "new",
            ...(params.status ? { status: params.status } : {}),
            ...(params.authorId ? { authorId: params.authorId } : {}),
          },
        })
        .json<PostsResponse>(),
  });
}

export function usePost(id: string | number) {
  return useQuery<Post, ApiError>({
    queryKey: ["post", id],
    queryFn: () => apiClient.get(`posts/${id}`).json<Post>(),
    enabled: !!id,
  });
}

export function useCreatePost() {
  const queryClient = useQueryClient();

  return useMutation<Post, ApiError, CreatePostPayload>({
    mutationFn: (payload) =>
      apiClient.post("posts", { json: payload }).json<Post>(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
    },
    onError: (error) => handleApiError(error),
  });
}

export function usePublishPost() {
  const queryClient = useQueryClient();

  return useMutation<Post, ApiError, string | number>({
    mutationFn: (id) =>
      apiClient.post(`posts/${id}/publish`).json<Post>(),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ["post", id] });
      queryClient.invalidateQueries({ queryKey: ["posts"] });
    },
    onError: (error) => handleApiError(error),
  });
}
