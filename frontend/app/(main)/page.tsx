"use client";

import Link from "next/link";
import { useState } from "react";
import { useFeed } from "@/hooks/queries/use-feed";
import { Button } from "@/components/ui/button";
import { Container } from "@/components/ui/container";
import { SidebarSection } from "@/components/ui/sidebar";
import { PostCard, PostCardSkeleton } from "@/components/features/post-card";
import { LeftNav } from "@/components/shared/left-nav";
import { useAuthStore } from "@/stores/auth-store";

type FeedTab = "recommend" | "hot" | "latest";

const MOCK_TAGS = [
  { name: "Go", count: 1234 },
  { name: "React", count: 987 },
  { name: "Next.js", count: 856 },
  { name: "TypeScript", count: 743 },
  { name: "Node.js", count: 621 },
  { name: "微服务", count: 512 },
  { name: "数据库", count: 489 },
  { name: "前端", count: 456 },
];

const MOCK_AUTHORS = [
  { id: 1, name: "技术大牛", followers: 12345, articles: 89 },
  { id: 2, name: "前端小王", followers: 8765, articles: 56 },
  { id: 3, name: "后端老李", followers: 6543, articles: 78 },
  { id: 4, name: "全栈张三", followers: 5432, articles: 45 },
];

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

  return (
    <Container size="full" className="py-6">
      <div className="grid grid-cols-1 lg:grid-cols-[200px_1fr_300px] gap-6">
        {/* 左侧导航 */}
        <aside className="hidden lg:block">
          <LeftNav />
        </aside>

        {/* 主内容区 */}
        <main className="min-w-0">
          {/* 顶部分类导航 */}
          <div className="mb-6 flex items-center justify-between border-b">
            <div className="flex gap-6">
              {[
                { key: "recommend", label: "推荐" },
                { key: "latest", label: "最新" },
                { key: "hot", label: "热榜" },
              ].map((tab) => (
                <button
                  key={tab.key}
                  onClick={() => setActiveTab(tab.key as FeedTab)}
                  className={`pb-3 text-sm font-medium transition-colors relative ${
                    activeTab === tab.key
                      ? "text-foreground"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  {tab.label}
                  {activeTab === tab.key && (
                    <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-foreground" />
                  )}
                </button>
              ))}
            </div>
          </div>

          {/* 文章列表 */}
          {activeTab === "recommend" && (
            <div className="rounded-lg border border-dashed bg-muted/30 p-12 text-center">
              <p className="text-sm text-muted-foreground">推荐流正在开发中，敬请期待 👷</p>
            </div>
          )}

          {activeTab !== "recommend" && (
            <div className="divide-y">
              {isLoading && (
                <>
                  <div className="py-4 first:pt-0">
                    <PostCardSkeleton />
                  </div>
                  <div className="py-4">
                    <PostCardSkeleton />
                  </div>
                  <div className="py-4">
                    <PostCardSkeleton />
                  </div>
                </>
              )}
              {data?.list.map((post) => (
                <div key={post.id} className="py-4 first:pt-0">
                  <PostCard post={post} />
                </div>
              ))}
              {!isLoading && data?.list.length === 0 && (
                <div className="py-12 text-center">
                  <p className="text-sm text-muted-foreground">暂无文章，快来发布第一篇吧！</p>
                </div>
              )}
            </div>
          )}
        </main>

        {/* 右侧信息栏 */}
        <aside className="hidden lg:block">
          <div className="sticky top-[5rem] space-y-4">
            {/* 作者榜 */}
            <SidebarSection title="作者榜">
              <div className="space-y-4">
                {MOCK_AUTHORS.map((author, index) => (
                  <div key={author.id} className="flex items-start gap-3">
                    <div className="flex h-5 w-5 shrink-0 items-center justify-center text-xs font-semibold text-muted-foreground">
                      {index + 1}
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="text-sm font-medium truncate hover:text-foreground cursor-pointer transition-colors">
                        {author.name}
                      </div>
                      <div className="text-xs text-muted-foreground mt-0.5">
                        {author.articles} 篇文章 · {(author.followers / 1000).toFixed(1)}k 关注
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </SidebarSection>

            {/* 热门标签 */}
            <SidebarSection title="热门标签">
              <div className="flex flex-wrap gap-2">
                {MOCK_TAGS.slice(0, 6).map((tag) => (
                  <button
                    key={tag.name}
                    className="inline-flex items-center gap-1.5 rounded-md bg-muted px-2.5 py-1.5 text-xs hover:bg-muted/80 transition-colors"
                  >
                    <span className="font-medium">{tag.name}</span>
                  </button>
                ))}
              </div>
            </SidebarSection>

            {/* 链接 */}
            <SidebarSection title="相关链接">
              <div className="space-y-2 text-xs">
                <a href="#" className="block text-muted-foreground hover:text-foreground transition-colors">
                  关于我们
                </a>
                <a href="#" className="block text-muted-foreground hover:text-foreground transition-colors">
                  用户协议
                </a>
                <a href="#" className="block text-muted-foreground hover:text-foreground transition-colors">
                  隐私政策
                </a>
                <a href="#" className="block text-muted-foreground hover:text-foreground transition-colors">
                  帮助中心
                </a>
              </div>
            </SidebarSection>
          </div>
        </aside>
      </div>
    </Container>
  );
}
