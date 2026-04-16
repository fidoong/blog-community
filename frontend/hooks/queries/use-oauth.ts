import { useMutation } from "@tanstack/react-query";
import apiClient from "@/lib/api-client";

interface OAuthURLResponse {
  authUrl: string;
}

async function getGitHubAuthURL(): Promise<OAuthURLResponse> {
  return await apiClient.get("auth/oauth/github").json<OAuthURLResponse>();
}

async function getGoogleAuthURL(): Promise<OAuthURLResponse> {
  return await apiClient.get("auth/oauth/google").json<OAuthURLResponse>();
}

export function useGitHubOAuth() {
  return useMutation({
    mutationFn: getGitHubAuthURL,
    onSuccess: (data) => {
      window.location.href = data.authUrl;
    },
  });
}

export function useGoogleOAuth() {
  return useMutation({
    mutationFn: getGoogleAuthURL,
    onSuccess: (data) => {
      window.location.href = data.authUrl;
    },
  });
}
