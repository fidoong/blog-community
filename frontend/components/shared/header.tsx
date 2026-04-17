"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "@/components/theme-toggle";
import { useAuthStore } from "@/stores/auth-store";
import { useLogout } from "@/hooks/queries/use-auth";

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
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto flex h-14 items-center justify-between px-4">
        <div className="flex items-center gap-6">
          <Link href="/" className="text-lg font-bold tracking-tight">
            BlogCommunity
          </Link>
          <nav className="hidden gap-4 text-sm font-medium md:flex">
            <Link href="/" className="text-muted-foreground hover:text-foreground transition-colors">
              首页
            </Link>
            <Link href="/search" className="text-muted-foreground hover:text-foreground transition-colors">
              搜索
            </Link>
          </nav>
        </div>

        <div className="flex items-center gap-2">
          <ThemeToggle />
          {user ? (
            <>
              <Link
                href={`/user/${user.id}`}
                className="hidden text-sm font-medium text-muted-foreground hover:text-foreground transition-colors sm:inline"
              >
                {user.username}
              </Link>
              <Button variant="ghost" size="sm" onClick={handleLogout} disabled={logoutMutation.isPending}>
                登出
              </Button>
            </>
          ) : (
            <Button size="sm" onClick={() => window.location.href = "/login"}>登录</Button>
          )}
        </div>
      </div>
    </header>
  );
}
