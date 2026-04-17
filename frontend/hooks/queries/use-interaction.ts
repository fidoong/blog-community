import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient, { ApiError } from "@/lib/api-client";
import { handleApiError } from "@/lib/error-handler";

interface LikeStatusResponse {
  isLiked: boolean;
  likeCount: number;
}

interface CollectStatusResponse {
  isCollected: boolean;
  collectCount: number;
}

export function useLikeStatus(targetType: string, targetId: string | number) {
  return useQuery<LikeStatusResponse, ApiError>({
    queryKey: ["likeStatus", targetType, targetId],
    queryFn: () =>
      apiClient.get(`likes/${targetType}/${targetId}`).json<LikeStatusResponse>(),
    enabled: !!targetId,
  });
}

export function useToggleLike() {
  const queryClient = useQueryClient();

  return useMutation<LikeStatusResponse, ApiError, { targetType: string; targetId: string | number }>({
    mutationFn: ({ targetType, targetId }) =>
      apiClient.post(`likes/${targetType}/${targetId}`).json<LikeStatusResponse>(),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["likeStatus", variables.targetType, variables.targetId],
      });
    },
    onError: (error) => handleApiError(error),
  });
}

export function useCollectStatus(targetType: string, targetId: string | number) {
  return useQuery<CollectStatusResponse, ApiError>({
    queryKey: ["collectStatus", targetType, targetId],
    queryFn: () =>
      apiClient.get(`collects/${targetType}/${targetId}`).json<CollectStatusResponse>(),
    enabled: !!targetId,
  });
}

export function useToggleCollect() {
  const queryClient = useQueryClient();

  return useMutation<CollectStatusResponse, ApiError, { targetType: string; targetId: string | number }>({
    mutationFn: ({ targetType, targetId }) =>
      apiClient.post(`collects/${targetType}/${targetId}`).json<CollectStatusResponse>(),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["collectStatus", variables.targetType, variables.targetId],
      });
    },
    onError: (error) => handleApiError(error),
  });
}
