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
      setAuth(data.user, data.accessToken);
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
      setAuth(data.user, data.accessToken);
    },
    onError: (error) => handleApiError(error),
  });
}
