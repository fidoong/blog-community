"use client";

import { Bookmark } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useCollectStatus, useToggleCollect } from "@/hooks/queries/use-interaction";
import { useAuthStore } from "@/stores/auth-store";

interface CollectButtonProps {
  targetType: string;
  targetId: string | number;
  initialCount?: number;
}

export function CollectButton({ targetType, targetId, initialCount = 0 }: CollectButtonProps) {
  const { token } = useAuthStore();
  const { data } = useCollectStatus(targetType, targetId);
  const toggleMutation = useToggleCollect();

  const isCollected = data?.isCollected ?? false;
  const count = data?.collectCount ?? initialCount;

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
      className={isCollected ? "text-primary" : ""}
    >
      <Bookmark className="mr-1 h-4 w-4" />
      {count}
    </Button>
  );
}
