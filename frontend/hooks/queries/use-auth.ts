import { useMutation } from "@tanstack/react-query";
import apiClient, { ApiError } from "@/lib/api-client";
import { handleApiError } from "@/lib/error-handler";
import { useAuthStore } from "@/stores/auth-store";

interface LoginPayload {
  email: string;
  password: string;
}

interface RegisterPayload {
  email: string;
  username: string;
  password: string;
}

interface AuthResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  user: {
    id: number;
    email: string;
    username: string;
    avatarUrl?: string;
    role: string;
  };
}

export function useLogin() {
  const setAuth = useAuthStore((s) => s.setAuth);

  return useMutation<AuthResponse, ApiError, LoginPayload>({
    mutationFn: (payload) =>
      apiClient.post("auth/login", { json: payload }).json<AuthResponse>(),
    onSuccess: (data) => {
      setAuth(data.user, data.accessToken, data.refreshToken);
    },
    onError: (error) => handleApiError(error),
  });
}

export function useRegister() {
  const setAuth = useAuthStore((s) => s.setAuth);

  return useMutation<AuthResponse, ApiError, RegisterPayload>({
    mutationFn: (payload) =>
      apiClient.post("auth/register", { json: payload }).json<AuthResponse>(),
    onSuccess: (data) => {
      setAuth(data.user, data.accessToken, data.refreshToken);
    },
    onError: (error) => handleApiError(error),
  });
}

interface RefreshPayload {
  refreshToken: string;
}

interface RefreshResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export function useRefresh() {
  const setToken = useAuthStore((s) => s.setToken);
  const setRefreshToken = useAuthStore((s) => s.setRefreshToken);

  return useMutation<RefreshResponse, ApiError, RefreshPayload>({
    mutationFn: (payload) =>
      apiClient.post("auth/refresh", { json: payload }).json<RefreshResponse>(),
    onSuccess: (data) => {
      setToken(data.accessToken);
      setRefreshToken(data.refreshToken);
    },
    onError: (error) => handleApiError(error),
  });
}

interface LogoutPayload {
  refreshToken: string;
}

export function useLogout() {
  const logoutStore = useAuthStore((s) => s.logout);

  return useMutation<{ message: string }, ApiError, LogoutPayload>({
    mutationFn: (payload) =>
      apiClient.post("auth/logout", { json: payload }).json<{ message: string }>(),
    onSuccess: () => {
      logoutStore();
    },
    onError: () => {
      logoutStore();
    },
  });
}
