"use client";

import { useState } from "react";
import { useComments, useCreateComment } from "@/hooks/queries/use-comments";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import type { Comment } from "@/types/comment";

function CommentItem({ comment, onReply }: { comment: Comment; postId: string | number; onReply: (parentId: number) => void }) {
  return (
    <div className="py-3">
      <div className="flex items-start gap-3">
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-muted text-xs font-medium">
          U{comment.authorId}
        </div>
        <div className="flex-1">
          <div className="text-sm font-medium">用户 {comment.authorId}</div>
          <p className="mt-1 text-sm text-foreground">{comment.content}</p>
          <div className="mt-2 flex items-center gap-4 text-xs text-muted-foreground">
            <span>{new Date(comment.createdAt * 1000).toLocaleString()}</span>
            <button className="hover:text-foreground" onClick={() => onReply(comment.id)}>
              回复
            </button>
            <span>👍 {comment.likeCount}</span>
          </div>

          {comment.replies && comment.replies.length > 0 && (
            <div className="mt-3 space-y-2 border-l-2 border-muted pl-4">
              {comment.replies.map((reply) => (
                <div key={reply.id} className="py-1">
                  <div className="text-sm font-medium">用户 {reply.authorId}</div>
                  <p className="mt-0.5 text-sm text-foreground">{reply.content}</p>
                  <div className="mt-1 flex items-center gap-4 text-xs text-muted-foreground">
                    <span>{new Date(reply.createdAt * 1000).toLocaleString()}</span>
                    <span>👍 {reply.likeCount}</span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export function CommentSection({ postId }: { postId: string | number }) {
  const { token } = useAuthStore();
  const { data, isLoading } = useComments(postId);
  const createMutation = useCreateComment();
  const [content, setContent] = useState("");
  const [replyTo, setReplyTo] = useState<number | null>(null);

  const handleSubmit = async () => {
    if (!content.trim()) return;
    await createMutation.mutateAsync({
      postId,
      payload: {
        content: content.trim(),
        ...(replyTo ? { parentId: replyTo } : {}),
      },
    });
    setContent("");
    setReplyTo(null);
  };

  return (
    <Card className="mt-6">
      <CardContent className="p-6">
        <h3 className="mb-4 text-lg font-semibold">评论 ({data?.pagination.total ?? 0})</h3>

        {token ? (
          <div className="mb-6 space-y-2">
            {replyTo && (
              <div className="flex items-center justify-between text-sm text-muted-foreground">
                <span>回复评论 #{replyTo}</span>
                <button className="text-primary hover:underline" onClick={() => setReplyTo(null)}>
                  取消回复
                </button>
              </div>
            )}
            <Textarea
              placeholder="写下你的评论..."
              value={content}
              onChange={(e) => setContent(e.target.value)}
              rows={3}
            />
            <Button
              onClick={handleSubmit}
              disabled={createMutation.isPending || !content.trim()}
            >
              发表评论
            </Button>
          </div>
        ) : (
          <div className="mb-6 rounded bg-muted p-4 text-center text-sm text-muted-foreground">
            登录后可发表评论
          </div>
        )}

        <div className="divide-y">
          {isLoading && (
            <>
              <Skeleton className="my-3 h-20 w-full" />
              <Skeleton className="my-3 h-20 w-full" />
            </>
          )}
          {data?.list.map((comment) => (
            <CommentItem
              key={comment.id}
              comment={comment}
              postId={postId}
              onReply={setReplyTo}
            />
          ))}
          {!isLoading && data?.list.length === 0 && (
            <div className="py-8 text-center text-sm text-muted-foreground">
              暂无评论，快来抢沙发吧！
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
