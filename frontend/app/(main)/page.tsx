"use client";

import Link from "next/link";
import { useState, useMemo, useEffect, useRef } from "react";
import { motion } from "framer-motion";
import { useInfiniteFeed } from "@/hooks/queries/use-infinite-feed";
import { Button } from "@/components/ui/button";
import { Container } from "@/components/ui/container";
import { SidebarSection } from "@/components/ui/sidebar";
import { PostCard, PostCardSkeleton } from "@/components/features/post-card";
import { LeftNav } from "@/components/shared/left-nav";
import { useAuthStore } from "@/stores/auth-store";

type FeedTab = "following" | "hot" | "latest";

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
  const loadMoreRef = useRef<HTMLDivElement>(null);

  const feedType = activeTab === "hot" ? "hot" : activeTab === "following" ? "following" : "latest";
  const { 
    data, 
    isLoading, 
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage 
  } = useInfiniteFeed({
    type: feedType,
    pageSize: 20,
    ...(activeTab === "hot" ? { period: "7d" } : {}),
  });

  // 扁平化所有页面的数据
  const allPosts = useMemo(() => {
    return data?.pages.flatMap((page) => page.list) ?? [];
  }, [data]);

  // 自动加载更多
  useEffect(() => {
    if (!loadMoreRef.current || isLoading || isFetchingNextPage || !hasNextPage) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          fetchNextPage();
        }
      },
      { threshold: 0.1 }
    );

    observer.observe(loadMoreRef.current);

    return () => {
      observer.disconnect();
    };
  }, [isLoading, isFetchingNextPage, hasNextPage, fetchNextPage]);

  return (
    <div className="h-full overflow-hidden">
      <Container size="full" className="h-full py-6">
        <div className="h-full grid grid-cols-1 lg:grid-cols-[200px_1fr_300px] gap-6">
          {/* 左侧导航 */}
          <aside className="hidden lg:block overflow-y-auto">
            <LeftNav />
          </aside>

          {/* 主内容区 */}
          <main className="min-w-0 flex flex-col overflow-hidden">
            {/* 顶部分类导航 - 固定不滚动 */}
            <div className="flex-shrink-0 mb-4 pb-3 border-b">
              <div className="flex gap-6 overflow-x-auto">
                {[
                  { key: "following", label: "关注" },
                  { key: "latest", label: "最新" },
                  { key: "hot", label: "热榜" },
                ].map((tab) => (
                  <button
                    key={tab.key}
                    onClick={() => setActiveTab(tab.key as FeedTab)}
                    className={`whitespace-nowrap pb-3 text-sm font-medium transition-colors relative ${
                      activeTab === tab.key
                        ? "text-foreground"
                        : "text-muted-foreground hover:text-foreground"
                    }`}
                  >
                    {tab.label}
                    {activeTab === tab.key && (
                      <motion.div
                        className="absolute bottom-0 left-0 right-0 h-0.5 bg-foreground"
                        layoutId="activeTab"
                        transition={{ duration: 0.2, ease: "easeOut" }}
                      />
                    )}
                  </button>
                ))}
              </div>
            </div>

          {/* 文章列表 - 独立滚动区域 */}
          <div className="flex-1 overflow-y-auto">
            {activeTab === "following" && !user && (
              <div className="rounded-lg border border-dashed bg-muted/30 p-12 text-center">
                <p className="text-sm text-muted-foreground mb-4">登录后查看关注作者的最新文章</p>
                <Button size="sm" onClick={() => window.location.href = "/login"}>去登录</Button>
              </div>
            )}

            {(activeTab !== "following" || user) && (
              <div className="divide-y">
                {isLoading && (
                  <>
                    {[...Array(5)].map((_, i) => (
                      <div key={i} className="py-4 first:pt-0">
                        <PostCardSkeleton />
                      </div>
                    ))}
                  </>
                )}
                
                {!isLoading && allPosts.length === 0 && (
                  <div className="py-12 text-center">
                    <p className="text-sm text-muted-foreground">暂无文章，快来发布第一篇吧！</p>
                  </div>
                )}

                {allPosts.map((post, index) => (
                  <div key={`${activeTab}-${post.id}-${index}`} className="py-4 first:pt-0">
                    <PostCard 
                      post={post} 
                      index={index} 
                      disableAnimation={index >= 20}
                    />
                  </div>
                ))}

                {/* 自动加载触发器 */}
                {hasNextPage && !isLoading && (
                  <div ref={loadMoreRef} className="py-8 text-center">
                    {isFetchingNextPage ? (
                      <div className="inline-flex items-center gap-2 text-sm text-muted-foreground">
                        <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                        <span>加载中...</span>
                      </div>
                    ) : (
                      <div className="h-4" />
                    )}
                  </div>
                )}

                {!hasNextPage && allPosts.length > 0 && (
                  <div className="py-8 text-center text-sm text-muted-foreground">
                    已经到底了
                  </div>
                )}
              </div>
            )}
          </div>
        </main>

        {/* 右侧信息栏 */}
        <aside className="hidden lg:block overflow-y-auto">
          <div className="space-y-4">
            {/* 作者榜 */}
            <SidebarSection title="作者榜">
              <div className="space-y-4">
                {MOCK_AUTHORS.map((author, index) => (
                  <motion.div 
                    key={author.id} 
                    className="flex items-start gap-3"
                    initial={{ opacity: 0, x: 10 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ 
                      duration: 0.2, 
                      delay: index * 0.05,
                      ease: "easeOut" 
                    }}
                  >
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
                  </motion.div>
                ))}
              </div>
            </SidebarSection>

            {/* 热门标签 */}
            <SidebarSection title="热门标签">
              <div className="flex flex-wrap gap-2">
                {MOCK_TAGS.slice(0, 6).map((tag, index) => (
                  <motion.button
                    key={tag.name}
                    className="inline-flex items-center gap-1.5 rounded-md bg-muted px-2.5 py-1.5 text-xs hover:bg-muted/80 transition-colors"
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ 
                      duration: 0.2, 
                      delay: index * 0.05,
                      ease: "easeOut" 
                    }}
                  >
                    <span className="font-medium">{tag.name}</span>
                  </motion.button>
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
    </div>
  );
}
