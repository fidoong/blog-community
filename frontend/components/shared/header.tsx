"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "@/components/theme-toggle";
import { useAuthStore } from "@/stores/auth-store";
import { useLogout } from "@/hooks/queries/use-auth";
import { PenSquare } from "lucide-react";

export function Header() {
  const { user, refreshToken, logout } = useAuthStore();
  const logoutMutation = useLogout();

  const handleLogout = () => {
    if (refreshToken) {
      logoutMutation.mutate({ refreshToken });
    } else {
      logout();
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
        <div className="flex items-center gap-8">
          <Link href="/" className="text-xl font-bold tracking-tight hover:text-foreground/80 transition-colors">
            BlogCommunity
          </Link>
          <nav className="hidden gap-6 text-sm font-medium md:flex">
            <Link 
              href="/" 
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              首页
            </Link>
            <Link 
              href="/search" 
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              搜索
            </Link>
          </nav>
        </div>

        <div className="flex items-center gap-3">
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
