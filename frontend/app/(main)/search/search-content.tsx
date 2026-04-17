"use client";

import { useState, useEffect, useRef } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import Link from "next/link";
import { motion } from "framer-motion";
import { Search, X, TrendingUp, Clock } from "lucide-react";
import { useSearch } from "@/hooks/queries/use-search";
import { useHotKeywords } from "@/hooks/queries/use-hot-keywords";
import { Container } from "@/components/ui/container";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/ui/empty-state";

import type { SearchHit } from "@/types/search";

function getHighlightedTitle(hit: SearchHit): string {
  if (hit.highlight?.title?.length) {
    return hit.highlight.title[0];
  }
  return hit.title;
}

function getHighlightedSummary(hit: SearchHit): string {
  if (hit.highlight?.summary?.length) {
    return hit.highlight.summary[0];
  }
  if (hit.highlight?.content?.length) {
    return hit.highlight.content[0];
  }
  return hit.summary;
}

function SearchResultItem({ hit, index }: { hit: SearchHit; index: number }) {
  const titleHtml = getHighlightedTitle(hit);
  const summaryHtml = getHighlightedSummary(hit);

  return (
    <motion.article
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.25, delay: Math.min(index * 0.04, 0.3) }}
      className="group py-5 border-b last:border-b-0"
    >
      <Link href={`/posts/${hit.id}`} className="block">
        <h3
          className="text-lg font-semibold leading-snug group-hover:text-foreground/80 transition-colors"
          dangerouslySetInnerHTML={{ __html: titleHtml }}
        />
        <p
          className="mt-1.5 text-sm text-muted-foreground line-clamp-2 leading-relaxed"
          dangerouslySetInnerHTML={{ __html: summaryHtml }}
        />
        <div className="mt-2 flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
          <span className="font-medium text-foreground/70">{hit.authorName}</span>
          <span>·</span>
          <span className="flex items-center gap-1">
            <Clock className="h-3 w-3" />
            {new Date(hit.createdAt * 1000).toLocaleDateString("zh-CN")}
          </span>
          {hit.tags.length > 0 && (
            <>
              <span>·</span>
              <div className="flex gap-1">
                {hit.tags.slice(0, 3).map((tag) => (
                  <span key={tag} className="inline-flex items-center rounded-md bg-muted px-1.5 py-0.5 text-[10px] font-medium text-muted-foreground">
                    {tag}
                  </span>
                ))}
              </div>
            </>
          )}
          <span className="ml-auto text-[10px] opacity-50">
            相关度 {hit.score.toFixed(2)}
          </span>
        </div>
      </Link>
    </motion.article>
  );
}

export function SearchContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const initialQ = searchParams.get("q") ?? "";

  const [query, setQuery] = useState(initialQ);
  const [searchTerm, setSearchTerm] = useState(initialQ);
  const inputRef = useRef<HTMLInputElement>(null);

  const { data, isLoading } = useSearch({
    keyword: searchTerm,
    pageSize: 20,
  });
  const { data: hotData } = useHotKeywords(10);

  const hits = data?.list ?? [];
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
            placeholder="搜索文章标题、内容、标签..."
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
        <div className="mb-4 flex items-center justify-between text-sm text-muted-foreground">
          {isLoading ? (
            <span>搜索中...</span>
          ) : (
            <span>
              共找到 <strong className="text-foreground">{data?.pagination.total ?? 0}</strong> 篇与「
              <strong className="text-foreground">{searchTerm}</strong>」相关的文章
            </span>
          )}
          {!isLoading && data && data.took > 0 && (
            <span className="text-xs opacity-50">{data.took}ms</span>
          )}
        </div>
      )}

      {isLoading && (
        <div className="space-y-4">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="py-5 border-b">
              <Skeleton className="h-6 w-3/4 mb-2" />
              <Skeleton className="h-4 w-full mb-1" />
              <Skeleton className="h-4 w-2/3" />
            </div>
          ))}
        </div>
      )}

      {!isLoading && searchTerm && hits.length === 0 && (
        <EmptyState
          title="未找到相关文章"
          description={`没有找到与「${searchTerm}」相关的文章，换个关键词试试吧`}
        />
      )}

      {!isLoading && hits.length > 0 && (
        <div>
          {hits.map((hit, index) => (
            <SearchResultItem key={hit.id} hit={hit} index={index} />
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
