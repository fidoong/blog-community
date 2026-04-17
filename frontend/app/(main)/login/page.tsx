"use client";

import { useState, useEffect, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { useLogin, useRegister } from "@/hooks/queries/use-auth";
import { useGitHubOAuth, useGoogleOAuth } from "@/hooks/queries/use-oauth";

function OAuthCallbackHandler() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const token = searchParams.get("token");
    const refreshToken = searchParams.get("refreshToken");
    if (token) {
      import("@/stores/auth-store").then(({ useAuthStore }) => {
        useAuthStore.getState().setToken(token);
        if (refreshToken) {
          useAuthStore.getState().setRefreshToken(refreshToken);
        }
        router.replace("/");
      });
    }
  }, [searchParams, router]);

  return null;
}

export default function LoginPage() {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const router = useRouter();

  const loginMutation = useLogin();
  const registerMutation = useRegister();
  const githubOAuth = useGitHubOAuth();
  const googleOAuth = useGoogleOAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isLogin) {
      await loginMutation.mutateAsync({ email, password });
    } else {
      await registerMutation.mutateAsync({ email, username, password });
    }
    router.push("/");
  };

  return (
    <div className="container mx-auto flex min-h-[calc(100vh-4rem)] max-w-md flex-col justify-center px-4 py-12">
      <Suspense fallback={null}>
        <OAuthCallbackHandler />
      </Suspense>
      <Card className="border-border/50">
        <CardHeader className="text-center space-y-2 pb-6">
          <CardTitle className="text-2xl font-bold">
            {isLogin ? "登录" : "注册"}
          </CardTitle>
          <CardDescription className="text-muted-foreground">
            {isLogin ? "欢迎回来，请登录您的账号" : "创建新账号，开始分享技术文章"}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email" className="text-sm font-medium">
                邮箱
              </Label>
              <Input
                id="email"
                type="email"
                placeholder="name@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="h-10"
              />
            </div>
            {!isLogin && (
              <div className="space-y-2">
                <Label htmlFor="username" className="text-sm font-medium">
                  用户名
                </Label>
                <Input
                  id="username"
                  placeholder="johndoe"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                  className="h-10"
                />
              </div>
            )}
            <div className="space-y-2">
              <Label htmlFor="password" className="text-sm font-medium">
                密码
              </Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="h-10"
              />
            </div>
            <Button
              type="submit"
              className="w-full h-10"
              disabled={loginMutation.isPending || registerMutation.isPending}
            >
              {isLogin ? "登录" : "注册"}
            </Button>
          </form>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-background px-2 text-muted-foreground">
                或通过以下方式继续
              </span>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <Button
              variant="outline"
              onClick={() => githubOAuth.mutate()}
              disabled={githubOAuth.isPending}
              className="h-10"
            >
              GitHub
            </Button>
            <Button
              variant="outline"
              onClick={() => googleOAuth.mutate()}
              disabled={googleOAuth.isPending}
              className="h-10"
            >
              Google
            </Button>
          </div>

          <div className="text-center text-sm">
            <button
              type="button"
              onClick={() => setIsLogin(!isLogin)}
              className="text-foreground/80 hover:text-foreground transition-colors"
            >
              {isLogin ? "还没有账号？立即注册" : "已有账号？立即登录"}
            </button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
