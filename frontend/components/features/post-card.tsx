import Link from "next/link";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import type { Post } from "@/types/post";
import { cn } from "@/lib/utils";

interface PostCardProps {
  post: Post;
  className?: string;
}

export function PostCard({ post, className }: PostCardProps) {
  return (
    <Card className={cn("group hover:border-foreground/20 transition-all", className)}>
      <CardHeader className="pb-3">
        <Link href={`/posts/${post.id}`}>
          <CardTitle className="text-lg font-semibold leading-tight group-hover:text-foreground/80 transition-colors line-clamp-2">
            {post.title}
          </CardTitle>
        </Link>
      </CardHeader>
      <CardContent className="pt-0 space-y-3">
        <p className="text-sm text-muted-foreground line-clamp-2 leading-relaxed">
          {post.summary || post.content?.slice(0, 120) || "暂无摘要"}
        </p>
        <div className="flex flex-wrap items-center gap-2 text-xs">
          <Link 
            href={`/user/${post.authorId}`} 
            className="font-medium text-muted-foreground hover:text-foreground transition-colors"
          >
            用户 {post.authorId}
          </Link>
          <div className="flex flex-wrap gap-2">
            {post.tags?.slice(0, 3).map((tag) => (
              <span 
                key={tag} 
                className="rounded-md bg-muted px-2 py-0.5 text-muted-foreground"
              >
                {tag}
              </span>
            ))}
          </div>
          <div className="ml-auto flex items-center gap-3 text-muted-foreground shrink-0">
            <span className="flex items-center gap-1">
              <span className="text-xs">👁</span>
              <span>{post.viewCount ?? 0}</span>
            </span>
            <span className="flex items-center gap-1">
              <span className="text-xs">👍</span>
              <span>{post.likeCount ?? 0}</span>
            </span>
            <span className="flex items-center gap-1">
              <span className="text-xs">💬</span>
              <span>{post.commentCount ?? 0}</span>
            </span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export function PostCardSkeleton() {
  return (
    <Card>
      <CardHeader className="pb-3">
        <Skeleton className="h-5 w-3/4" />
      </CardHeader>
      <CardContent className="pt-0 space-y-3">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-2/3" />
        <div className="flex gap-2">
          <Skeleton className="h-4 w-16" />
          <Skeleton className="h-4 w-12" />
          <Skeleton className="ml-auto h-4 w-24" />
        </div>
      </CardContent>
    </Card>
  );
}
