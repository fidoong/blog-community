"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { usePost, useUpdatePost } from "@/hooks/queries/use-posts";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Container } from "@/components/ui/container";
import { PageHeader } from "@/components/ui/page-header";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

export default function EditPostPage() {
  const params = useParams();
  const router = useRouter();
  const postId = params.id as string;

  const { data: post, isLoading } = usePost(postId);
  const updateMutation = useUpdatePost();

  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [tags, setTags] = useState("");

  useEffect(() => {
    if (post) {
      setTitle(post.title);
      setContent(post.content ?? "");
      setTags(post.tags?.join(", ") ?? "");
    }
  }, [post]);

  const handleUpdate = async () => {
    if (!title.trim() || !content.trim()) return;
    await updateMutation.mutateAsync({
      id: postId,
      payload: {
        title: title.trim(),
        content: content.trim(),
        tags: tags.split(",").map((t) => t.trim()).filter(Boolean),
      },
    });
    toast.success("文章更新成功");
    router.push(`/posts/${postId}`);
  };

  if (isLoading) {
    return (
      <Container size="md" className="py-12">
        <Skeleton className="h-8 w-48 mb-8" />
        <div className="space-y-6">
          <div className="space-y-2">
            <Skeleton className="h-4 w-16" />
            <Skeleton className="h-11 w-full" />
          </div>
          <div className="space-y-2">
            <Skeleton className="h-4 w-16" />
            <Skeleton className="h-64 w-full" />
          </div>
          <Skeleton className="h-10 w-32" />
        </div>
      </Container>
    );
  }

  if (!post) {
    return (
      <Container size="md" className="py-16 text-center">
        <p className="text-muted-foreground">文章不存在或暂无权限编辑</p>
      </Container>
    );
  }

  return (
    <Container size="md" className="py-12">
      <PageHeader title="编辑文章" description="修改你的文章内容" />

      <div className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="title" className="text-sm font-medium">标题</Label>
          <Input
            id="title"
            placeholder="请输入文章标题"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="h-11 text-base"
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="content" className="text-sm font-medium">内容 (Markdown)</Label>
          <Textarea
            id="content"
            rows={20}
            placeholder="开始写作..."
            value={content}
            onChange={(e) => setContent(e.target.value)}
            className="resize-none text-[15px] leading-relaxed font-mono"
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="tags" className="text-sm font-medium">标签</Label>
          <Input
            id="tags"
            placeholder="用英文逗号分隔，如：Go, 后端, 微服务"
            value={tags}
            onChange={(e) => setTags(e.target.value)}
            className="h-10"
          />
        </div>

        <div className="flex gap-3 pt-4">
          <Button
            variant="outline"
            onClick={() => router.push(`/posts/${postId}`)}
          >
            取消
          </Button>
          <Button
            onClick={handleUpdate}
            disabled={updateMutation.isPending || !title || !content}
          >
            {updateMutation.isPending ? "保存中..." : "保存修改"}
          </Button>
        </div>
      </div>
    </Container>
  );
}
