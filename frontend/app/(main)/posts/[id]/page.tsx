"use client";

import { useParams } from "next/navigation";
import Link from "next/link";
import { usePost } from "@/hooks/queries/use-posts";
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
import { ArrowLeft } from "lucide-react";
import { CommentSection } from "@/components/shared/comment-section";
import { LikeButton } from "@/components/shared/like-button";
import { CollectButton } from "@/components/shared/collect-button";

export default function PostDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const { data: post, isLoading } = usePost(id);
  const { user: me } = useAuthStore();

  const authorId = post?.authorId;
  const { data: author } = useUserProfile(authorId ?? "");
  const { data: authorStats } = useFollowStats(authorId ?? "");

  const followMutation = useFollow();
  const unfollowMutation = useUnfollow();

  const isMe = me?.id === authorId;
  const isFollowing = authorStats?.isFollowing ?? false;

  const handleFollowToggle = () => {
    if (!authorId) return;
    if (isFollowing) {
      unfollowMutation.mutate(authorId);
    } else {
      followMutation.mutate(authorId);
    }
  };

  if (isLoading) {
    return (
      <div className="container mx-auto max-w-3xl px-4 py-8">
        <Skeleton className="mb-4 h-8 w-3/4" />
        <Skeleton className="mb-2 h-4 w-full" />
        <Skeleton className="mb-2 h-4 w-full" />
        <Skeleton className="h-4 w-2/3" />
      </div>
    );
  }

  if (!post) {
    return (
      <div className="container mx-auto max-w-3xl px-4 py-16 text-center">
        <p className="text-muted-foreground">文章不存在或已被删除</p>
        <Link href="/" className="mt-4 inline-block">
          <Button variant="outline">返回首页</Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="container mx-auto max-w-3xl px-4 py-8">
      <Link href="/" className="mb-4 inline-flex items-center text-sm text-muted-foreground hover:text-foreground">
        <ArrowLeft className="mr-1 h-4 w-4" />
        返回首页
      </Link>

      <Card>
        <CardContent className="p-6 md:p-8">
          <h1 className="mb-4 text-2xl font-bold md:text-3xl">{post.title}</h1>

          {/* Author bar */}
          <div className="mb-6 flex items-center justify-between rounded-lg bg-muted/50 p-3">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-full bg-muted text-sm font-bold">
                {author?.username?.[0]?.toUpperCase() ?? "U"}
              </div>
              <div>
                <Link
                  href={`/user/${post.authorId}`}
                  className="text-sm font-medium hover:text-primary"
                >
                  {author?.username ?? `用户 ${post.authorId}`}
                </Link>
                <div className="text-xs text-muted-foreground">
                  {authorStats?.followersCount ?? 0} 粉丝
                </div>
              </div>
            </div>
            {me && !isMe && (
              <Button
                variant={isFollowing ? "outline" : "default"}
                size="sm"
                onClick={handleFollowToggle}
                disabled={followMutation.isPending || unfollowMutation.isPending}
              >
                {isFollowing ? "已关注" : "关注"}
              </Button>
            )}
          </div>

          <div className="mb-6 flex flex-wrap items-center gap-3 text-sm text-muted-foreground">
            {post.tags.map((tag) => (
              <span key={tag} className="rounded bg-muted px-2 py-0.5">
                {tag}
              </span>
            ))}
            <span className="ml-auto flex items-center gap-4">
              <span>👁 {post.viewCount}</span>
              <span>💬 {post.commentCount}</span>
            </span>
          </div>

          {post.contentType === "markdown" ? (
            <article className="prose prose-zinc dark:prose-invert max-w-none">
              <pre className="whitespace-pre-wrap font-sans text-foreground bg-transparent p-0">
                {post.content}
              </pre>
            </article>
          ) : (
            <article
              className="prose prose-zinc dark:prose-invert max-w-none"
              dangerouslySetInnerHTML={{ __html: post.content || "" }}
            />
          )}

          <div className="mt-8 flex items-center gap-4 border-t pt-6">
            <LikeButton targetType="post" targetId={post.id} initialCount={post.likeCount} />
            <CollectButton targetType="post" targetId={post.id} initialCount={post.collectCount} />
          </div>
        </CardContent>
      </Card>

      <CommentSection postId={id} />
    </div>
  );
}
