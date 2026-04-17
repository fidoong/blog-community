"use client";

import { useParams } from "next/navigation";
import Link from "next/link";
import { motion } from "framer-motion";
import { usePost } from "@/hooks/queries/use-posts";
import { useRelatedPosts } from "@/hooks/queries/use-related-posts";
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
import { ArrowLeft } from "lucide-react";
import { CommentSection } from "@/components/shared/comment-section";
import { LikeButton } from "@/components/shared/like-button";
import { CollectButton } from "@/components/shared/collect-button";

export default function PostDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const { data: post, isLoading } = usePost(id);
  const { data: relatedData } = useRelatedPosts(id, 5);
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

  // 生成用户头像显示文本
  const getAvatarText = (authorId: number) => {
    const idStr = String(authorId);
    return idStr.length > 3 ? idStr.slice(-3) : idStr;
  };

  if (isLoading) {
    return (
      <Container size="lg" className="py-6">
        <Skeleton className="mb-4 h-8 w-3/4" />
        <Skeleton className="mb-3 h-4 w-full" />
        <Skeleton className="mb-3 h-4 w-full" />
        <Skeleton className="h-4 w-2/3" />
      </Container>
    );
  }

  if (!post) {
    return (
      <Container size="lg" className="py-12">
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
    <Container size="lg" className="py-6">
      <div className="grid grid-cols-1 lg:grid-cols-[1fr_240px] gap-8">
        {/* 主内容区 */}
        <motion.article 
          className="min-w-0"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4, ease: "easeOut" }}
        >
          {/* 返回按钮 */}
          <Link href="/" className="mb-4 inline-flex items-center text-sm text-muted-foreground hover:text-foreground transition-colors">
            <ArrowLeft className="mr-1.5 h-4 w-4" />
            返回首页
          </Link>

          {/* 标题 */}
          <h1 className="mb-4 text-3xl font-bold tracking-tight leading-tight">{post.title}</h1>
          
          {/* 作者信息 */}
          <div className="mb-4 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-muted text-[10px] font-semibold">
                {getAvatarText(post.authorId)}
              </div>
              <div>
                <Link href={`/user/${post.authorId}`} className="text-sm font-medium hover:text-foreground/80 transition-colors truncate max-w-[200px]">
                  {post.authorName || author?.username || `用户 ${post.authorId}`}
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

          {/* 标签和统计 */}
          <div className="mb-6 flex flex-wrap items-center gap-2 text-xs text-muted-foreground pb-6 border-b">
            {post.tags.map((tag) => (
              <span key={tag} className="rounded-md bg-muted px-2 py-1 font-medium">
                {tag}
              </span>
            ))}
            <div className="ml-auto flex items-center gap-3">
              <span className="flex items-center gap-1">
                <span>👁</span>
                <span>{post.viewCount}</span>
              </span>
              <span className="flex items-center gap-1">
                <span>💬</span>
                <span>{post.commentCount}</span>
              </span>
            </div>
          </div>

          {/* 文章内容 */}
          <div className="mb-6">
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
          <div className="flex items-center gap-3 py-4 border-y">
            <LikeButton targetType="post" targetId={post.id} initialCount={post.likeCount} />
            <CollectButton targetType="post" targetId={post.id} initialCount={post.collectCount} />
          </div>

          {/* 评论区 */}
          <div className="mt-6">
            <CommentSection postId={id} />
          </div>

          {/* 相关推荐 */}
          {relatedData && relatedData.list.length > 0 && (
            <motion.div
              className="mt-10 pt-8 border-t"
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.3, delay: 0.2 }}
            >
              <h3 className="text-lg font-semibold mb-4">相关推荐</h3>
              <div className="space-y-4">
                {relatedData.list.map((relatedPost, index) => (
                  <motion.div
                    key={relatedPost.id}
                    initial={{ opacity: 0, x: -8 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ duration: 0.2, delay: index * 0.05 }}
                  >
                    <Link href={`/posts/${relatedPost.id}`} className="group block">
                      <h4 className="text-sm font-medium group-hover:text-foreground/80 transition-colors line-clamp-2">
                        {relatedPost.title}
                      </h4>
                      <div className="mt-1 flex items-center gap-2 text-xs text-muted-foreground">
                        <span className="truncate max-w-[120px]">
                          {relatedPost.authorName || `用户 ${relatedPost.authorId}`}
                        </span>
                        <span>·</span>
                        <span>👍 {relatedPost.likeCount ?? 0}</span>
                        <span>👁 {relatedPost.viewCount ?? 0}</span>
                      </div>
                    </Link>
                  </motion.div>
                ))}
              </div>
            </motion.div>
          )}
        </motion.article>

        {/* 右侧栏 */}
        <aside className="hidden lg:block">
          <motion.div 
            className="sticky top-[5rem] space-y-4"
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.4, delay: 0.1, ease: "easeOut" }}
          >
            {/* 作者信息卡片 */}
            <div className="rounded-lg border bg-card p-4">
              <div className="mb-3 text-sm font-semibold">关于作者</div>
              <div className="flex items-center gap-3 mb-3">
                <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-muted text-xs font-semibold">
                  {getAvatarText(post.authorId)}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="text-sm font-medium truncate">
                    {post.authorName || author?.username || `用户 ${post.authorId}`}
                  </div>
                  <div className="text-xs text-muted-foreground">
                    {authorStats?.followersCount ?? 0} 粉丝
                  </div>
                </div>
              </div>
              {me && !isMe && (
                <Button
                  variant={isFollowing ? "outline" : "default"}
                  size="sm"
                  className="w-full"
                  onClick={handleFollowToggle}
                  disabled={followMutation.isPending || unfollowMutation.isPending}
                >
                  {isFollowing ? "已关注" : "关注"}
                </Button>
              )}
            </div>

            {/* 目录 */}
            <div className="rounded-lg border bg-card p-4">
              <div className="text-sm font-semibold">目录</div>
              <div className="mt-3 text-xs text-muted-foreground">
                暂无目录
              </div>
            </div>
          </motion.div>
        </aside>
      </div>
    </Container>
  );
}
