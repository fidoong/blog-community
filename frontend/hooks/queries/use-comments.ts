import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient, { ApiError } from "@/lib/api-client";
import { handleApiError } from "@/lib/error-handler";
import type { Comment, CommentsResponse, CreateCommentPayload } from "@/types/comment";

export function useComments(postId: string | number) {
  return useQuery<CommentsResponse, ApiError>({
    queryKey: ["comments", postId],
    queryFn: () =>
      apiClient.get(`posts/${postId}/comments`).json<CommentsResponse>(),
    enabled: !!postId,
  });
}

export function useCreateComment() {
  const queryClient = useQueryClient();

  return useMutation<Comment, ApiError, { postId: string | number; payload: CreateCommentPayload }>({
    mutationFn: ({ postId, payload }) =>
      apiClient.post(`posts/${postId}/comments`, { json: payload }).json<Comment>(),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["comments", variables.postId] });
    },
    onError: (error) => handleApiError(error),
  });
}

export function useDeleteComment() {
  const queryClient = useQueryClient();

  return useMutation<void, ApiError, { commentId: number; postId: string | number }>({
    mutationFn: ({ commentId }) =>
      apiClient.delete(`comments/${commentId}`).json<void>(),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["comments", variables.postId] });
    },
    onError: (error) => handleApiError(error),
  });
}
