"use client";

import { useParams } from "next/navigation";
import Link from "next/link";
import { usePost } from "@/hooks/queries/use-posts";
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
