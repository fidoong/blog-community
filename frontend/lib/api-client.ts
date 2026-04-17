import ky from "ky";
import { useAuthStore } from "@/stores/auth-store";

const apiClient = ky.create({
  prefix: process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1",
  headers: { "Content-Type": "application/json" },
  timeout: 10000,
  hooks: {
    beforeRequest: [
      (state) => {
        const token = useAuthStore.getState().token;
        if (token) {
          state.request.headers.set("Authorization", `Bearer ${token}`);
        }
      },
    ],
    afterResponse: [
      async (state) => {
        if (!state.response.ok) {
          const body = await state.response
            .json<{ message?: string; code?: string }>()
            .catch(() => ({ message: "", code: "" }));
          throw new ApiError(body.message || "请求失败", state.response.status, body.code || undefined);
        }
        // Unwrap standard API envelope: { code, message, data, timestamp }
        const clone = state.response.clone();
        const body = await clone.json().catch(() => null);
        if (
          body &&
          typeof body === "object" &&
          "data" in body &&
          "code" in body
        ) {
          return new Response(JSON.stringify(body.data), {
            status: state.response.status,
            statusText: state.response.statusText,
            headers: state.response.headers,
          });
        }
        return state.response;
      },
    ],
  },
});

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public code?: string
  ) {
    super(message);
    this.name = "ApiError";
  }
}

export default apiClient;
