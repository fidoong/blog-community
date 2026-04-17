"use client";

import { useState } from "react";
import { useComments, useCreateComment } from "@/hooks/queries/use-comments";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Skeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import type { Comment } from "@/types/comment";

function CommentItem({ comment, onReply }: { comment: Comment; postId: string | number; onReply: (parentId: number) => void }) {
  return (
    <div className="py-4">
      <div className="flex items-start gap-3">
        <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-muted text-xs font-semibold">
          U{comment.authorId}
        </div>
        <div className="flex-1 min-w-0">
          <div className="text-sm font-medium mb-1">用户 {comment.authorId}</div>
          <p className="text-sm text-foreground leading-relaxed">{comment.content}</p>
          <div className="mt-2 flex items-center gap-4 text-xs text-muted-foreground">
            <span>{new Date(comment.createdAt * 1000).toLocaleString()}</span>
            <button className="hover:text-foreground transition-colors" onClick={() => onReply(comment.id)}>
              回复
            </button>
            <span className="flex items-center gap-1">
              <span>👍</span>
              <span>{comment.likeCount}</span>
            </span>
          </div>

          {comment.replies && comment.replies.length > 0 && (
            <div className="mt-4 space-y-3 border-l-2 border-border pl-4">
              {comment.replies.map((reply) => (
                <div key={reply.id}>
                  <div className="text-sm font-medium mb-1">用户 {reply.authorId}</div>
                  <p className="text-sm text-foreground leading-relaxed">{reply.content}</p>
                  <div className="mt-1 flex items-center gap-4 text-xs text-muted-foreground">
                    <span>{new Date(reply.createdAt * 1000).toLocaleString()}</span>
                    <span className="flex items-center gap-1">
                      <span>👍</span>
                      <span>{reply.likeCount}</span>
                    </span>
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
    <div className="border-t pt-8">
      <h3 className="text-xl font-semibold mb-6">评论 ({data?.pagination.total ?? 0})</h3>

      {token ? (
        <div className="mb-8 space-y-3">
          {replyTo && (
            <div className="flex items-center justify-between text-sm">
              <span className="text-muted-foreground">回复评论 #{replyTo}</span>
              <button 
                className="text-foreground/80 hover:text-foreground transition-colors" 
                onClick={() => setReplyTo(null)}
              >
                取消回复
              </button>
            </div>
          )}
          <Textarea
            placeholder="写下你的评论..."
            value={content}
            onChange={(e) => setContent(e.target.value)}
            rows={4}
            className="resize-none text-[15px] leading-relaxed"
          />
          <Button
            onClick={handleSubmit}
            disabled={createMutation.isPending || !content.trim()}
          >
            发表评论
          </Button>
        </div>
      ) : (
        <div className="mb-8 rounded-lg border border-dashed bg-muted/30 p-6 text-center text-sm text-muted-foreground">
          登录后可发表评论
        </div>
      )}

      <div className="divide-y">
        {isLoading && (
          <>
            <div className="py-4">
              <Skeleton className="h-20 w-full" />
            </div>
            <div className="py-4">
              <Skeleton className="h-20 w-full" />
            </div>
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
          <div className="py-12 text-center text-sm text-muted-foreground">
            暂无评论，快来抢沙发吧！
          </div>
        )}
      </div>
    </div>
  );
}
