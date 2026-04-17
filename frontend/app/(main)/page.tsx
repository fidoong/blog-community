"use client";

import Link from "next/link";
import { useState } from "react";
import { useFeed } from "@/hooks/queries/use-feed";
import { Button } from "@/components/ui/button";
import { Container } from "@/components/ui/container";
import { Tabs } from "@/components/ui/tabs";
import { EmptyState } from "@/components/ui/empty-state";
import { PostCard, PostCardSkeleton } from "@/components/features/post-card";
import { useAuthStore } from "@/stores/auth-store";

type FeedTab = "recommend" | "hot" | "latest";

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

  const tabs = [
    { key: "recommend", label: "推荐" },
    { key: "hot", label: "热榜" },
    { key: "latest", label: "最新" },
  ];

  return (
    <Container size="md" className="py-12">
      <div className="mb-8 flex items-center justify-between gap-4">
        <Tabs 
          tabs={tabs} 
          activeTab={activeTab} 
          onTabChange={(key) => setActiveTab(key as FeedTab)} 
        />
        {user && (
          <Link href="/posts/new">
            <Button size="sm">写文章</Button>
          </Link>
        )}
      </div>

      {activeTab === "recommend" && (
        <EmptyState 
          title="推荐流正在开发中" 
          description="敬请期待 👷"
        />
      )}

      {activeTab !== "recommend" && (
        <div className="space-y-4">
          {isLoading && (
            <>
              <PostCardSkeleton />
              <PostCardSkeleton />
              <PostCardSkeleton />
            </>
          )}
          {data?.list.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
          {!isLoading && data?.list.length === 0 && (
            <EmptyState 
              title="暂无文章" 
              description="快来发布第一篇吧！"
            />
          )}
        </div>
      )}
    </Container>
  );
}
