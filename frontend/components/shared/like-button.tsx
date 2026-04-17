"use client";

import { ThumbsUp } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useLikeStatus, useToggleLike } from "@/hooks/queries/use-interaction";
import { useAuthStore } from "@/stores/auth-store";

interface LikeButtonProps {
  targetType: string;
  targetId: string | number;
  initialCount?: number;
}

export function LikeButton({ targetType, targetId, initialCount = 0 }: LikeButtonProps) {
  const { token } = useAuthStore();
  const { data } = useLikeStatus(targetType, targetId);
  const toggleMutation = useToggleLike();

  const isLiked = data?.isLiked ?? false;
  const count = data?.likeCount ?? initialCount;

  const handleClick = () => {
    if (!token) {
      window.location.href = "/login";
      return;
    }
    toggleMutation.mutate({ targetType, targetId });
  };

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={handleClick}
      disabled={toggleMutation.isPending}
      className={isLiked ? "text-primary" : ""}
    >
      <ThumbsUp className="mr-1 h-4 w-4" />
      {count}
    </Button>
  );
}
