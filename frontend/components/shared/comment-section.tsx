"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import { useComments, useCreateComment } from "@/hooks/queries/use-comments";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Skeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import type { Comment } from "@/types/comment";

function CommentItem({ comment, onReply, index = 0 }: { comment: Comment; postId: string | number; onReply: (parentId: number) => void; index?: number }) {
  // 生成用户头像显示文本（取ID后3位或首字母）
  const getAvatarText = (authorId: number) => {
    const idStr = String(authorId);
    return idStr.length > 3 ? idStr.slice(-3) : idStr;
  };

  return (
    <motion.div 
      className="py-5"
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ 
        duration: 0.3, 
        delay: index * 0.05,
        ease: "easeOut" 
      }}
    >
      <div className="flex items-start gap-3">
        {/* 头像 */}
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-muted text-[10px] font-semibold">
          {getAvatarText(comment.authorId)}
        </div>
        
        {/* 内容区 */}
        <div className="flex-1 min-w-0">
          {/* 用户名 */}
          <div className="mb-1.5 text-sm font-medium text-foreground">
            用户 {comment.authorId}
          </div>
          
          {/* 评论内容 */}
          <p className="mb-2 text-sm text-foreground leading-relaxed break-words">
            {comment.content}
          </p>
          
          {/* 操作栏 */}
          <div className="flex items-center gap-4 text-xs text-muted-foreground">
            <span>{new Date(comment.createdAt * 1000).toLocaleString('zh-CN', { 
              year: 'numeric', 
              month: '2-digit', 
              day: '2-digit',
              hour: '2-digit',
              minute: '2-digit'
            })}</span>
            <button 
              className="transition-colors hover:text-foreground" 
              onClick={() => onReply(comment.id)}
            >
              回复
            </button>
            <button className="flex items-center gap-1 transition-colors hover:text-foreground">
              <span>👍</span>
              <span>{comment.likeCount}</span>
            </button>
          </div>

          {/* 楼中楼回复 */}
          {comment.replies && comment.replies.length > 0 && (
            <div className="mt-4 space-y-4 rounded-lg bg-muted/30 p-4">
              {comment.replies.map((reply) => (
                <div key={reply.id} className="flex items-start gap-3">
                  {/* 回复者头像 */}
                  <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-muted text-[10px] font-semibold">
                    {getAvatarText(reply.authorId)}
                  </div>
                  
                  {/* 回复内容 */}
                  <div className="flex-1 min-w-0">
                    <div className="mb-1 text-sm font-medium text-foreground">
                      用户 {reply.authorId}
                    </div>
                    <p className="mb-1.5 text-sm text-foreground leading-relaxed break-words">
                      {reply.content}
                    </p>
                    <div className="flex items-center gap-4 text-xs text-muted-foreground">
                      <span>{new Date(reply.createdAt * 1000).toLocaleString('zh-CN', { 
                        year: 'numeric', 
                        month: '2-digit', 
                        day: '2-digit',
                        hour: '2-digit',
                        minute: '2-digit'
                      })}</span>
                      <button className="flex items-center gap-1 transition-colors hover:text-foreground">
                        <span>👍</span>
                        <span>{reply.likeCount}</span>
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </motion.div>
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
      <h3 className="mb-6 text-xl font-semibold">评论 ({data?.pagination.total ?? 0})</h3>

      {token ? (
        <div className="mb-8 space-y-3">
          {replyTo && (
            <div className="flex items-center justify-between rounded-lg bg-muted/50 px-3 py-2 text-sm">
              <span className="text-muted-foreground">回复评论 #{replyTo}</span>
              <button 
                className="text-foreground/80 transition-colors hover:text-foreground" 
                onClick={() => setReplyTo(null)}
              >
                取消
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
            <div className="py-5">
              <Skeleton className="h-24 w-full" />
            </div>
            <div className="py-5">
              <Skeleton className="h-24 w-full" />
            </div>
          </>
        )}
        {data?.list.map((comment, index) => (
          <CommentItem
            key={comment.id}
            comment={comment}
            postId={postId}
            onReply={setReplyTo}
            index={index}
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
