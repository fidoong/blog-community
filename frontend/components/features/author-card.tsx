"use client"

import Link from "next/link";
import { motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface AuthorCardProps {
  authorId: string;
  username?: string;
  followersCount?: number;
  isMe?: boolean;
  isFollowing?: boolean;
  onFollowToggle?: () => void;
  isLoading?: boolean;
  className?: string;
}

export function AuthorCard({
  authorId,
  username,
  followersCount = 0,
  isMe = false,
  isFollowing = false,
  onFollowToggle,
  isLoading = false,
  className,
}: AuthorCardProps) {
  return (
    <motion.div 
      className={cn("flex items-center justify-between rounded-lg border bg-card p-4", className)}
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.3, ease: "easeOut" }}
    >
      <div className="flex items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-muted text-sm font-semibold">
          {username?.[0]?.toUpperCase() ?? "U"}
        </div>
        <div>
          <Link
            href={`/user/${authorId}`}
            className="text-sm font-medium hover:text-foreground/80 transition-colors truncate max-w-[140px]"
            title={username || `用户 ${authorId}`}
          >
            {username ?? `用户 ${authorId}`}
          </Link>
          <div className="text-xs text-muted-foreground">
            {followersCount} 粉丝
          </div>
        </div>
      </div>
      {!isMe && onFollowToggle && (
        <Button
          variant={isFollowing ? "outline" : "default"}
          size="sm"
          onClick={onFollowToggle}
          disabled={isLoading}
        >
          {isFollowing ? "已关注" : "关注"}
        </Button>
      )}
    </motion.div>
  );
}
