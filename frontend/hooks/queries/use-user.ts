import { useQuery } from "@tanstack/react-query";
import apiClient, { ApiError } from "@/lib/api-client";

export interface UserProfile {
  id: number;
  username: string;
  email: string;
  avatarUrl?: string;
  role: string;
}

export function useUserProfile(id: string | number) {
  return useQuery<UserProfile, ApiError>({
    queryKey: ["user", id],
    queryFn: () => apiClient.get(`users/${id}`).json<UserProfile>(),
    enabled: !!id,
  });
}
