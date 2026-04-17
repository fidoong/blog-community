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
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

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
      <div className="container mx-auto max-w-3xl px-4 py-8">
        <Card>
          <CardContent className="flex items-center gap-4 py-6">
            <Skeleton className="h-16 w-16 rounded-full" />
            <div className="space-y-2">
              <Skeleton className="h-5 w-32" />
              <Skeleton className="h-4 w-24" />
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="container mx-auto max-w-3xl px-4 py-12 text-center text-muted-foreground">
        用户不存在
      </div>
    );
  }

  return (
    <div className="container mx-auto max-w-3xl px-4 py-8">
      <Card>
        <CardContent className="py-6">
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-4">
              <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted text-xl font-bold">
                {profile.username?.[0]?.toUpperCase() ?? "U"}
              </div>
              <div>
                <h1 className="text-xl font-bold">{profile.username}</h1>
                <p className="text-sm text-muted-foreground">{profile.email}</p>
                <div className="mt-2 flex gap-4 text-sm">
                  <span>
                    <strong>{stats?.followersCount ?? 0}</strong> 粉丝
                  </span>
                  <span>
                    <strong>{stats?.followingCount ?? 0}</strong> 关注
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
        </CardContent>
      </Card>

      <div className="mt-6">
        <h2 className="mb-4 text-lg font-semibold">发布的文章</h2>
        <div className="rounded-lg border border-dashed p-8 text-center text-sm text-muted-foreground">
          文章列表开发中，敬请期待 👷
        </div>
      </div>
    </div>
  );
}
