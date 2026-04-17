"use client";

import Link from "next/link";
import { useState } from "react";
import { useFeed } from "@/hooks/queries/use-feed";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useAuthStore } from "@/stores/auth-store";
import type { Post } from "@/types/post";

type FeedTab = "recommend" | "hot" | "latest";

function PostCard({ post }: { post: Post }) {
  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="pb-2">
        <Link href={`/posts/${post.id}`}>
          <CardTitle className="text-lg leading-snug hover:text-primary cursor-pointer">
            {post.title}
          </CardTitle>
        </Link>
      </CardHeader>
      <CardContent className="pt-0">
        <p className="text-sm text-muted-foreground line-clamp-2">
          {post.summary || post.content?.slice(0, 120) || "暂无摘要"}
        </p>
        <div className="mt-3 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
          <Link href={`/user/${post.authorId}`} className="font-medium hover:text-foreground">
            用户 {post.authorId}
          </Link>
          {post.tags?.map((tag) => (
            <span key={tag} className="rounded bg-muted px-2 py-0.5">
              {tag}
            </span>
          ))}
          <span className="ml-auto flex items-center gap-3">
            <span>👁 {post.viewCount ?? 0}</span>
            <span>👍 {post.likeCount ?? 0}</span>
            <span>💬 {post.commentCount ?? 0}</span>
          </span>
        </div>
      </CardContent>
    </Card>
  );
}

function PostSkeleton() {
  return (
    <Card>
      <CardHeader className="pb-2">
        <Skeleton className="h-5 w-3/4" />
      </CardHeader>
      <CardContent className="pt-0 space-y-2">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-2/3" />
        <div className="mt-3 flex gap-2">
          <Skeleton className="h-4 w-12" />
          <Skeleton className="h-4 w-12" />
          <Skeleton className="ml-auto h-4 w-24" />
        </div>
      </CardContent>
    </Card>
  );
}

export default function HomePage() {
  const [activeTab, setActiveTab] = useState<FeedTab>("latest");
  const { user } = useAuthStore();

  const feedType = activeTab === "hot" ? "hot" : "latest";
  const { data, isLoading } = useFeed({
    type: feedType,
    page: 1,
    pageSize: 20,
    ...(activeTab === "hot" ? { period: "7d" } : {}),
  });

  const tabs: { key: FeedTab; label: string }[] = [
    { key: "recommend", label: "推荐" },
    { key: "hot", label: "热榜" },
    { key: "latest", label: "最新" },
  ];

  return (
    <div className="container mx-auto px-4 py-6">
      <div className="mx-auto max-w-3xl">
        <div className="mb-6 flex items-center justify-between">
          <div className="flex gap-1 rounded-lg bg-muted p-1">
            {tabs.map((tab) => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key)}
                className={`rounded-md px-4 py-1.5 text-sm font-medium transition-colors ${
                  activeTab === tab.key
                    ? "bg-background text-foreground shadow-sm"
                    : "text-muted-foreground hover:text-foreground"
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>
          {user && (
            <Link href="/posts/new">
              <Button>写文章</Button>
            </Link>
          )}
        </div>

        {activeTab === "recommend" && (
          <div className="rounded-lg border border-dashed p-8 text-center text-muted-foreground">
            <p className="text-sm">推荐流正在开发中，敬请期待 👷</p>
          </div>
        )}

        {activeTab !== "recommend" && (
          <div className="space-y-4">
            {isLoading && (
              <>
                <PostSkeleton />
                <PostSkeleton />
                <PostSkeleton />
              </>
            )}
            {data?.list.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}
            {!isLoading && data?.list.length === 0 && (
              <div className="py-12 text-center text-muted-foreground">
                暂无文章，快来发布第一篇吧！
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
