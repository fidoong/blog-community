"use client";

import { useState, useEffect, useRef } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { Search, X, TrendingUp } from "lucide-react";
import { usePosts } from "@/hooks/queries/use-posts";
import { useHotKeywords } from "@/hooks/queries/use-hot-keywords";
import { Container } from "@/components/ui/container";
import { PostCard, PostCardSkeleton } from "@/components/features/post-card";
import { EmptyState } from "@/components/ui/empty-state";

export function SearchContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const initialQ = searchParams.get("q") ?? "";

  const [query, setQuery] = useState(initialQ);
  const [searchTerm, setSearchTerm] = useState(initialQ);
  const inputRef = useRef<HTMLInputElement>(null);

  const { data, isLoading } = usePosts({
    keyword: searchTerm,
    pageSize: 20,
  });
  const { data: hotData } = useHotKeywords(10);

  const posts = data?.list ?? [];
  const hotKeywords = hotData?.list ?? [];

  // 当 URL 参数变化时同步
  useEffect(() => {
    const q = searchParams.get("q") ?? "";
    setQuery(q);
    setSearchTerm(q);
  }, [searchParams]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = query.trim();
    if (trimmed) {
      router.push(`/search?q=${encodeURIComponent(trimmed)}`);
    }
  };

  const handleClear = () => {
    setQuery("");
    inputRef.current?.focus();
  };

  return (
    <Container size="md" className="py-12">
      {/* 搜索框 */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
      >
        <form onSubmit={handleSubmit} className="relative mb-8">
          <Search className="absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
          <input
            ref={inputRef}
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="搜索文章标题、摘要..."
            className="h-14 w-full rounded-xl border bg-background pl-12 pr-12 text-base shadow-sm transition-colors focus:border-foreground focus:outline-none focus:ring-1 focus:ring-foreground"
          />
          {query && (
            <button
              type="button"
              onClick={handleClear}
              className="absolute right-4 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
            >
              <X className="h-5 w-5" />
            </button>
          )}
        </form>
      </motion.div>

      {/* 搜索结果 */}
      {searchTerm && (
        <div className="mb-4 text-sm text-muted-foreground">
          {isLoading ? (
            <span>搜索中...</span>
          ) : (
            <span>
              共找到 <strong className="text-foreground">{data?.pagination.total ?? 0}</strong> 篇与「
              <strong className="text-foreground">{searchTerm}</strong>」相关的文章
            </span>
          )}
        </div>
      )}

      {isLoading && (
        <div className="divide-y">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="py-4 first:pt-0">
              <PostCardSkeleton />
            </div>
          ))}
        </div>
      )}

      {!isLoading && searchTerm && posts.length === 0 && (
        <EmptyState
          title="未找到相关文章"
          description={`没有找到与「${searchTerm}」相关的文章，换个关键词试试吧`}
        />
      )}

      {!isLoading && posts.length > 0 && (
        <div className="divide-y">
          {posts.map((post, index) => (
            <div key={post.id} className="py-4 first:pt-0">
              <PostCard post={post} index={index} />
            </div>
          ))}
        </div>
      )}

      {!searchTerm && (
        <>
          {hotKeywords.length > 0 && (
            <motion.div
              className="mb-8"
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.3, delay: 0.1 }}
            >
              <div className="flex items-center gap-2 mb-3">
                <TrendingUp className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm font-medium text-muted-foreground">热搜词</span>
              </div>
              <div className="flex flex-wrap gap-2">
                {hotKeywords.map((keyword, index) => (
                  <button
                    key={keyword}
                    onClick={() => {
                      setQuery(keyword);
                      router.push(`/search?q=${encodeURIComponent(keyword)}`);
                    }}
                    className="inline-flex items-center gap-1.5 rounded-lg bg-muted px-3 py-1.5 text-sm hover:bg-muted/80 transition-colors"
                  >
                    <span className={`text-xs font-semibold ${
                      index < 3 ? "text-foreground" : "text-muted-foreground"
                    }`}>
                      {index + 1}
                    </span>
                    <span>{keyword}</span>
                  </button>
                ))}
              </div>
            </motion.div>
          )}
          <div className="py-16 text-center">
            <Search className="mx-auto h-12 w-12 text-muted-foreground/50 mb-4" />
            <p className="text-muted-foreground">输入关键词开始搜索文章</p>
          </div>
        </>
      )}
    </Container>
  );
}
