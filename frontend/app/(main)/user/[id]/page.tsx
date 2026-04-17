"use client";

import { useParams } from "next/navigation";
import { useUserProfile } from "@/hooks/queries/use-user";
import {
  useFollowStats,
  useFollow,
  useUnfollow,
} from "@/hooks/queries/use-follow";
import { useAuthStore } from "@/stores/auth-store";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Container } from "@/components/ui/container";
import { EmptyState } from "@/components/ui/empty-state";

export default function UserProfilePage() {
  const params = useParams();
  const userId = params.id as string;
  const { user: me } = useAuthStore();

  const { data: profile, isLoading: profileLoading } = useUserProfile(userId);
  const { data: stats, isLoading: statsLoading } = useFollowStats(userId);

  const followMutation = useFollow();
  const unfollowMutation = useUnfollow();

  const isMe = me?.id.toString() === userId;
  const isFollowing = stats?.isFollowing ?? false;

  const handleFollowToggle = () => {
    if (isFollowing) {
      unfollowMutation.mutate(userId);
    } else {
      followMutation.mutate(userId);
    }
  };

  if (profileLoading) {
    return (
      <Container size="md" className="py-12">
        <div className="flex items-center gap-4 mb-8">
          <Skeleton className="h-20 w-20 rounded-full" />
          <div className="space-y-2">
            <Skeleton className="h-6 w-32" />
            <Skeleton className="h-4 w-48" />
          </div>
        </div>
      </Container>
    );
  }

  if (!profile) {
    return (
      <Container size="md" className="py-16">
        <EmptyState title="用户不存在" />
      </Container>
    );
  }

  return (
    <Container size="md" className="py-12">
      {/* 用户信息卡片 */}
      <div className="mb-12">
        <div className="flex items-start justify-between gap-6">
          <div className="flex items-center gap-4">
            <div className="flex h-20 w-20 items-center justify-center rounded-full bg-muted text-2xl font-bold">
              {profile.username?.[0]?.toUpperCase() ?? "U"}
            </div>
            <div>
              <h1 className="text-2xl font-bold mb-1">{profile.username}</h1>
              <p className="text-sm text-muted-foreground mb-3">{profile.email}</p>
              <div className="flex gap-6 text-sm">
                <span>
                  <strong className="font-semibold text-foreground">{stats?.followersCount ?? 0}</strong>{" "}
                  <span className="text-muted-foreground">粉丝</span>
                </span>
                <span>
                  <strong className="font-semibold text-foreground">{stats?.followingCount ?? 0}</strong>{" "}
                  <span className="text-muted-foreground">关注</span>
                </span>
              </div>
            </div>
          </div>

          {me && !isMe && (
            <Button
              variant={isFollowing ? "outline" : "default"}
              onClick={handleFollowToggle}
              disabled={followMutation.isPending || unfollowMutation.isPending}
            >
              {isFollowing ? "已关注" : "关注"}
            </Button>
          )}
        </div>
      </div>

      {/* 文章列表 */}
      <div>
        <h2 className="text-lg font-semibold mb-6">发布的文章</h2>
        <EmptyState 
          title="文章列表开发中" 
          description="敬请期待 👷"
        />
      </div>
    </Container>
  );
}
