import Link from "next/link";
import { Skeleton } from "@/components/ui/skeleton";
import type { Post } from "@/types/post";
import { cn } from "@/lib/utils";

interface PostCardProps {
  post: Post;
  className?: string;
}

export function PostCard({ post, className }: PostCardProps) {
  // 生成用户头像显示文本（取ID后3位）
  const getAvatarText = (authorId: number) => {
    const idStr = String(authorId);
    return idStr.length > 3 ? idStr.slice(-3) : idStr;
  };

  return (
    <article className={cn("group", className)}>
      <div className="flex items-start gap-4">
        {/* 作者头像 */}
        <Link href={`/user/${post.authorId}`} className="shrink-0">
          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-muted text-[10px] font-semibold transition-opacity hover:opacity-80">
            {getAvatarText(post.authorId)}
          </div>
        </Link>

        {/* 内容区 */}
        <div className="flex-1 min-w-0">
          {/* 标题 */}
          <Link href={`/posts/${post.id}`}>
            <h3 className="mb-2 text-base font-semibold leading-tight line-clamp-2 transition-colors group-hover:text-foreground/80">
              {post.title}
            </h3>
          </Link>

          {/* 摘要 */}
          <p className="mb-3 text-sm text-muted-foreground line-clamp-2 leading-relaxed">
            {post.summary || post.content?.slice(0, 120) || "暂无摘要"}
          </p>

          {/* 底部信息 */}
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            <Link 
              href={`/user/${post.authorId}`} 
              className="transition-colors hover:text-foreground"
            >
              用户 {post.authorId}
            </Link>
            
            {post.tags && post.tags.length > 0 && (
              <>
                <span>·</span>
                <div className="flex items-center gap-2">
                  {post.tags.slice(0, 2).map((tag) => (
                    <span key={tag} className="cursor-pointer transition-colors hover:text-foreground">
                      {tag}
                    </span>
                  ))}
                </div>
              </>
            )}

            <div className="ml-auto flex shrink-0 items-center gap-3">
              <span className="flex items-center gap-1">
                <span>👁</span>
                <span>{post.viewCount ?? 0}</span>
              </span>
              <span className="flex items-center gap-1">
                <span>👍</span>
                <span>{post.likeCount ?? 0}</span>
              </span>
              <span className="flex items-center gap-1">
                <span>💬</span>
                <span>{post.commentCount ?? 0}</span>
              </span>
            </div>
          </div>
        </div>
      </div>
    </article>
  );
}

export function PostCardSkeleton() {
  return (
    <div className="flex items-start gap-4">
      <Skeleton className="h-10 w-10 rounded-full shrink-0" />
      <div className="flex-1 space-y-3">
        <Skeleton className="h-5 w-3/4" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-2/3" />
        <div className="flex gap-2">
          <Skeleton className="h-4 w-16" />
          <Skeleton className="h-4 w-12" />
          <Skeleton className="ml-auto h-4 w-24" />
        </div>
      </div>
    </div>
  );
}
