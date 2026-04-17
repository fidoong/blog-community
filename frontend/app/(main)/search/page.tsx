"use client";

import { Suspense } from "react";
import { Search as SearchIcon } from "lucide-react";
import { SearchContent } from "./search-content";

export default function SearchPage() {
  return (
    <Suspense
      fallback={
        <div className="container mx-auto max-w-3xl py-12">
          <div className="h-14 w-full rounded-xl border bg-muted/50 mb-8" />
          <div className="py-16 text-center">
            <SearchIcon className="mx-auto h-12 w-12 text-muted-foreground/50 mb-4" />
            <p className="text-muted-foreground">加载中...</p>
          </div>
        </div>
      }
    >
      <SearchContent />
    </Suspense>
  );
}
