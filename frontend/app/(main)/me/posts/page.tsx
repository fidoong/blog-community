"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { useAuthStore } from "@/stores/auth-store";
import { usePosts, useDeletePost, usePublishPost } from "@/hooks/queries/use-posts";
import { Button } from "@/components/ui/button";
import { Container } from "@/components/ui/container";
import { Skeleton } from "@/components/ui/skeleton";
import { PageHeader } from "@/components/ui/page-header";
import { EmptyState } from "@/components/ui/empty-state";
import { PenSquare, Trash2, Send, Eye } from "lucide-react";
import { toast } from "sonner";
import type { Post } from "@/types/post";

type PostStatus = "all" | "draft" | "pending" | "published" | "rejected";

const STATUS_MAP: Record<string, { label: string; color: string }> = {
  draft: { label: "草稿", color: "bg-muted text-muted-foreground" },
  pending: { label: "审核中", color: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200" },
  published: { label: "已发布", color: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200" },
  rejected: { label: "被驳回", color: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200" },
};

const TABS: { key: PostStatus; label: string }[] = [
  { key: "all", label: "全部" },
  { key: "published", label: "已发布" },
  { key: "draft", label: "草稿" },
  { key: "pending", label: "审核中" },
  { key: "rejected", label: "被驳回" },
];

function PostRow({ post, index }: { post: Post; index: number }) {
  const router = useRouter();
  const deleteMutation = useDeletePost();
  const publishMutation = usePublishPost();
  const status = STATUS_MAP[post.status] ?? { label: post.status, color: "bg-muted" };

  const handleDelete = async () => {
    if (!confirm(`确定要删除文章《${post.title}》吗？此操作不可恢复。`)) return;
    await deleteMutation.mutateAsync(post.id);
    toast.success("文章已删除");
  };

  const handlePublish = async () => {
    await publishMutation.mutateAsync(post.id);
    toast.success("文章已发布");
  };

  return (
    <motion.div
      className="flex items-center gap-4 py-4 border-b last:border-b-0"
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.2, delay: Math.min(index * 0.03, 0.3), ease: "easeOut" }}
    >
      <div className="flex-1 min-w-0">
        <Link href={`/posts/${post.id}`} className="group">
          <h3 className="text-base font-semibold truncate group-hover:text-foreground/80 transition-colors">
            {post.title}
          </h3>
        </Link>
        <div className="mt-1.5 flex items-center gap-3 text-xs text-muted-foreground">
          <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-medium ${status.color}`}>
            {status.label}
          </span>
          <span>{new Date(post.createdAt * 1000).toLocaleDateString('zh-CN')}</span>
          <span className="flex items-center gap-1">
            <Eye className="h-3 w-3" />
            {post.viewCount ?? 0}
          </span>
          <span>👍 {post.likeCount ?? 0}</span>
          <span>💬 {post.commentCount ?? 0}</span>
        </div>
      </div>

      <div className="flex items-center gap-2 shrink-0">
        {post.status === "draft" && (
          <Button
            size="sm"
            variant="outline"
            onClick={handlePublish}
            disabled={publishMutation.isPending}
            className="h-8 gap-1"
          >
            <Send className="h-3.5 w-3.5" />
            发布
          </Button>
        )}
        <Link href={`/posts/${post.id}/edit`}>
          <Button size="sm" variant="ghost" className="h-8 gap-1">
            <PenSquare className="h-3.5 w-3.5" />
            编辑
          </Button>
        </Link>
        <Button
          size="sm"
          variant="ghost"
          onClick={handleDelete}
          disabled={deleteMutation.isPending}
          className="h-8 gap-1 text-muted-foreground hover:text-destructive"
        >
          <Trash2 className="h-3.5 w-3.5" />
          删除
        </Button>
      </div>
    </motion.div>
  );
}

export default function MyPostsPage() {
  const { user } = useAuthStore();
  const [activeTab, setActiveTab] = useState<PostStatus>("all");

  const { data, isLoading } = usePosts({
    authorId: user?.id,
    status: activeTab === "all" ? undefined : activeTab,
    pageSize: 50,
  });

  const posts = data?.list ?? [];

  return (
    <Container size="md" className="py-12">
      <PageHeader
        title="我的文章"
        description="管理你发布的所有文章"
      />

      {/* 状态筛选 */}
      <div className="flex items-center gap-1 mb-6 border-b">
        {TABS.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={`px-4 py-2.5 text-sm font-medium transition-colors relative ${
              activeTab === tab.key
                ? "text-foreground"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            {tab.label}
            {activeTab === tab.key && (
              <motion.div
                className="absolute bottom-0 left-0 right-0 h-0.5 bg-foreground"
                layoutId="myPostsTab"
                transition={{ duration: 0.2, ease: "easeOut" }}
              />
            )}
          </button>
        ))}
      </div>

      {/* 文章列表 */}
      {isLoading && (
        <div className="space-y-4">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="flex items-center gap-4 py-4 border-b">
              <div className="flex-1 space-y-2">
                <Skeleton className="h-5 w-3/4" />
                <Skeleton className="h-4 w-1/2" />
              </div>
              <Skeleton className="h-8 w-20" />
            </div>
          ))}
        </div>
      )}

      {!isLoading && posts.length === 0 && (
        <EmptyState
          title="暂无文章"
          description={
            activeTab === "all"
              ? "你还没有发布过文章，快去写一篇吧！"
              : `暂无${TABS.find((t) => t.key === activeTab)?.label}状态的文章`
          }
        />
      )}

      {!isLoading && posts.length > 0 && (
        <div>
          {posts.map((post, index) => (
            <PostRow key={post.id} post={post} index={index} />
          ))}
        </div>
      )}
    </Container>
  );
}
