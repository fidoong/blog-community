import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient, { ApiError } from "@/lib/api-client";
import { handleApiError } from "@/lib/error-handler";

export interface FollowStats {
  followersCount: number;
  followingCount: number;
  isFollowing?: boolean;
}

interface FollowListResponse {
  list: { id: number; followerId: number; followingId: number; createdAt: number }[];
  pagination: {
    total: number;
    page: number;
    pageSize: number;
    totalPages: number;
  };
}

const followKeys = {
  stats: (userId: string | number) => ["follow-stats", userId] as const,
  followers: (userId: string | number, page: number) =>
    ["followers", userId, page] as const,
  following: (userId: string | number, page: number) =>
    ["following", userId, page] as const,
};

export function useFollowStats(userId: string | number) {
  return useQuery<FollowStats, ApiError>({
    queryKey: followKeys.stats(userId),
    queryFn: () => apiClient.get(`users/${userId}/follow-stats`).json<FollowStats>(),
    enabled: !!userId,
  });
}

export function useFollowers(userId: string | number, page = 1) {
  return useQuery<FollowListResponse, ApiError>({
    queryKey: followKeys.followers(userId, page),
    queryFn: () =>
      apiClient
        .get(`users/${userId}/followers`, { searchParams: { page, pageSize: 20 } })
        .json<FollowListResponse>(),
    enabled: !!userId,
  });
}

export function useFollowing(userId: string | number, page = 1) {
  return useQuery<FollowListResponse, ApiError>({
    queryKey: followKeys.following(userId, page),
    queryFn: () =>
      apiClient
        .get(`users/${userId}/following`, { searchParams: { page, pageSize: 20 } })
        .json<FollowListResponse>(),
    enabled: !!userId,
  });
}

export function useFollow() {
  const queryClient = useQueryClient();

  return useMutation<void, ApiError, string | number>({
    mutationFn: (userId) => apiClient.post(`users/${userId}/follow`).json<void>(),
    onSuccess: (_, userId) => {
      queryClient.invalidateQueries({ queryKey: followKeys.stats(userId) });
      queryClient.invalidateQueries({ queryKey: ["followers", userId] });
    },
    onError: (error) => handleApiError(error),
  });
}

export function useUnfollow() {
  const queryClient = useQueryClient();

  return useMutation<void, ApiError, string | number>({
    mutationFn: (userId) => apiClient.delete(`users/${userId}/follow`).json<void>(),
    onSuccess: (_, userId) => {
      queryClient.invalidateQueries({ queryKey: followKeys.stats(userId) });
      queryClient.invalidateQueries({ queryKey: ["followers", userId] });
    },
    onError: (error) => handleApiError(error),
  });
}
