import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient, { ApiError } from "@/lib/api-client";
import { handleApiError } from "@/lib/error-handler";
import type {
  NotificationItem,
  NotificationsResponse,
  UnreadCountResponse,
} from "@/types/notification";

export function useNotifications(page = 1, pageSize = 20, unreadOnly = false) {
  return useQuery<NotificationsResponse, ApiError>({
    queryKey: ["notifications", page, pageSize, unreadOnly],
    queryFn: () =>
      apiClient
        .get("notifications", {
          searchParams: { page, pageSize, unread: unreadOnly ? "1" : "0" },
        })
        .json<NotificationsResponse>(),
    enabled: typeof window !== "undefined",
  });
}

export function useUnreadCount() {
  return useQuery<UnreadCountResponse, ApiError>({
    queryKey: ["notifications", "unread-count"],
    queryFn: () =>
      apiClient
        .get("notifications/unread-count")
        .json<UnreadCountResponse>(),
    refetchInterval: 30000, // 每 30 秒刷新一次
    enabled: typeof window !== "undefined",
  });
}

export function useMarkRead() {
  const queryClient = useQueryClient();

  return useMutation<void, ApiError, number>({
    mutationFn: (id) =>
      apiClient.put(`notifications/${id}/read`).json<void>(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
    },
    onError: (error) => handleApiError(error),
  });
}

export function useMarkAllRead() {
  const queryClient = useQueryClient();

  return useMutation<void, ApiError>({
    mutationFn: () =>
      apiClient.put("notifications/read-all").json<void>(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
    },
    onError: (error) => handleApiError(error),
  });
}
