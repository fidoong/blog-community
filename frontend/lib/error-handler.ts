import { toast } from "sonner";
import { ApiError } from "./api-client";

export function handleApiError(error: unknown) {
  if (error instanceof ApiError) {
    if (error.status === 401) {
      if (typeof window !== "undefined") {
        window.location.href = "/login";
      }
      return;
    }
    if (error.status === 403) {
      toast.error("权限不足");
      return;
    }
    if (error.status >= 500) {
      toast.error("服务器繁忙，请稍后重试");
      return;
    }
    toast.error(error.message);
  } else {
    toast.error("网络异常，请检查连接");
  }
}
