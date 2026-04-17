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
import { Skeleton } from "@/components/ui/skeleton";
import { Container } from "@/components/ui/container";
import { EmptyState } from "@/components/ui/empty-state";
import { AuthorCard } from "@/components/features/author-card";
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
      <Container size="md" className="py-12">
        <Skeleton className="mb-6 h-8 w-3/4" />
        <Skeleton className="mb-4 h-4 w-full" />
        <Skeleton className="mb-4 h-4 w-full" />
        <Skeleton className="h-4 w-2/3" />
      </Container>
    );
  }

  if (!post) {
    return (
      <Container size="md" className="py-16">
        <EmptyState 
          title="文章不存在或已被删除"
          action={
            <Link href="/">
              <Button variant="outline">返回首页</Button>
            </Link>
          }
        />
      </Container>
    );
  }

  return (
    <Container size="md" className="py-12">
      <Link href="/" className="mb-8 inline-flex items-center text-sm text-muted-foreground hover:text-foreground transition-colors">
        <ArrowLeft className="mr-1.5 h-4 w-4" />
        返回首页
      </Link>

      <article className="space-y-8">
        {/* 标题区 */}
        <div className="space-y-6">
          <h1 className="text-4xl font-bold tracking-tight leading-tight">{post.title}</h1>
          
          <AuthorCard
            authorId={post.authorId}
            username={author?.username}
            followersCount={authorStats?.followersCount}
            isMe={isMe}
            isFollowing={isFollowing}
            onFollowToggle={me && !isMe ? handleFollowToggle : undefined}
            isLoading={followMutation.isPending || unfollowMutation.isPending}
          />

          <div className="flex flex-wrap items-center gap-2">
            {post.tags.map((tag) => (
              <span key={tag} className="rounded-md bg-muted px-2.5 py-1 text-xs text-muted-foreground font-medium">
                {tag}
              </span>
            ))}
            <div className="ml-auto flex items-center gap-4 text-sm text-muted-foreground">
              <span className="flex items-center gap-1.5">
                <span>👁</span>
                <span>{post.viewCount}</span>
              </span>
              <span className="flex items-center gap-1.5">
                <span>💬</span>
                <span>{post.commentCount}</span>
              </span>
            </div>
          </div>
        </div>

        {/* 内容区 */}
        <div className="border-t pt-8">
          {post.contentType === "markdown" ? (
            <div className="prose prose-zinc dark:prose-invert max-w-none">
              <pre className="whitespace-pre-wrap font-sans text-foreground bg-transparent p-0 text-[15px] leading-relaxed">
                {post.content}
              </pre>
            </div>
          ) : (
            <div
              className="prose prose-zinc dark:prose-invert max-w-none"
              dangerouslySetInnerHTML={{ __html: post.content || "" }}
            />
          )}
        </div>

        {/* 互动区 */}
        <div className="flex items-center gap-3 border-t pt-6">
          <LikeButton targetType="post" targetId={post.id} initialCount={post.likeCount} />
          <CollectButton targetType="post" targetId={post.id} initialCount={post.collectCount} />
        </div>
      </article>

      {/* 评论区 */}
      <div className="mt-12">
        <CommentSection postId={id} />
      </div>
    </Container>
  );
}
