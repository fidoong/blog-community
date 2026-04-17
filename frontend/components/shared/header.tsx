"use client";

import Link from "next/link";
import { useState } from "react";
import { motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "@/components/theme-toggle";
import { useAuthStore } from "@/stores/auth-store";
import { useLogout } from "@/hooks/queries/use-auth";
import { useUnreadCount } from "@/hooks/queries/use-notifications";
import { PenSquare, Search, Bell } from "lucide-react";

export function Header() {
  const { user, refreshToken, logout } = useAuthStore();
  const logoutMutation = useLogout();
  const { data: unreadData } = useUnreadCount();
  const [searchQuery, setSearchQuery] = useState("");

  const unreadCount = unreadData?.count ?? 0;

  const handleLogout = () => {
    if (refreshToken) {
      logoutMutation.mutate({ refreshToken });
    } else {
      logout();
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = searchQuery.trim();
    if (trimmed) {
      window.location.href = `/search?q=${encodeURIComponent(trimmed)}`;
    }
  };

  return (
    <motion.header
      className="sticky top-0 z-50 w-full border-b bg-background/80 backdrop-blur-md"
      initial={{ opacity: 0, y: -10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, ease: "easeOut" }}
    >
      <div className="container mx-auto flex h-16 items-center justify-between px-4 md:px-6 lg:px-8">
        <div className="flex items-center gap-6 lg:gap-8">
          <Link
            href="/"
            className="text-xl font-bold tracking-tight hover:text-foreground/80 transition-colors shrink-0"
          >
            BlogCommunity
          </Link>

          <nav className="hidden gap-6 text-sm font-medium lg:flex">
            <Link
              href="/"
              className="whitespace-nowrap text-muted-foreground hover:text-foreground transition-colors"
            >
              首页
            </Link>
          </nav>
        </div>

        <div className="flex items-center gap-3">
          {/* 搜索框 */}
          <form
            onSubmit={handleSearch}
            className="hidden md:flex relative max-w-[240px] w-full"
          >
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="搜索文章..."
              className="h-9 w-full rounded-lg border bg-muted/50 pl-9 pr-4 text-sm transition-colors focus:border-foreground focus:bg-background focus:outline-none focus:ring-1 focus:ring-foreground"
            />
          </form>

          {user && (
            <Link href="/posts/new">
              <Button size="sm" variant="ghost" className="gap-1.5">
                <PenSquare className="h-4 w-4" />
                <span>写文章</span>
              </Button>
            </Link>
          )}
          <ThemeToggle />
          {user ? (
            <>
              <Link
                href="/notifications"
                className="relative inline-flex items-center justify-center p-2 text-muted-foreground hover:text-foreground transition-colors"
              >
                <Bell className="h-5 w-5" />
                {unreadCount > 0 && (
                  <span className="absolute top-1 right-1 flex h-4 min-w-4 items-center justify-center rounded-full bg-red-500 px-1 text-[10px] font-medium text-white">
                    {unreadCount > 99 ? "99+" : unreadCount}
                  </span>
                )}
              </Link>
              <Link
                href={`/user/${user.id}`}
                className="hidden text-sm font-medium text-muted-foreground hover:text-foreground transition-colors sm:inline"
              >
                {user.username}
              </Link>
              <Button
                size="sm"
                onClick={handleLogout}
                disabled={logoutMutation.isPending}
                className="bg-black text-white hover:bg-black/90 dark:bg-white dark:text-black dark:hover:bg-white/90"
              >
                退出
              </Button>
            </>
          ) : (
            <Button size="sm" onClick={() => window.location.href = "/login"}>
              登录
            </Button>
          )}
        </div>
      </div>
    </motion.header>
  );
}
