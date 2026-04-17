"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useCreatePost, usePublishPost } from "@/hooks/queries/use-posts";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "sonner";

export default function NewPostPage() {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [tags, setTags] = useState("");
  const router = useRouter();

  const createMutation = useCreatePost();
  const publishMutation = usePublishPost();

  const handleSaveDraft = async () => {
    const post = await createMutation.mutateAsync({
      title,
      content,
      tags: tags.split(",").map((t) => t.trim()).filter(Boolean),
    });
    toast.success("草稿保存成功");
    router.push(`/posts/${post.id}`);
  };

  const handlePublish = async () => {
    const post = await createMutation.mutateAsync({
      title,
      content,
      tags: tags.split(",").map((t) => t.trim()).filter(Boolean),
    });
    await publishMutation.mutateAsync(post.id);
    toast.success("文章发布成功，等待审核");
    router.push(`/posts/${post.id}`);
  };

  return (
    <div className="container mx-auto max-w-3xl px-4 py-8">
      <Card>
        <CardHeader>
          <CardTitle>写文章</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">标题</Label>
            <Input
              id="title"
              placeholder="请输入文章标题"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="content">内容 (Markdown)</Label>
            <Textarea
              id="content"
              rows={16}
              placeholder="开始写作..."
              value={content}
              onChange={(e) => setContent(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="tags">标签</Label>
            <Input
              id="tags"
              placeholder="用英文逗号分隔，如：Go, 后端, 微服务"
              value={tags}
              onChange={(e) => setTags(e.target.value)}
            />
          </div>

          <div className="flex gap-3 pt-2">
            <Button
              variant="outline"
              onClick={handleSaveDraft}
              disabled={createMutation.isPending || !title || !content}
            >
              保存草稿
            </Button>
            <Button
              onClick={handlePublish}
              disabled={createMutation.isPending || publishMutation.isPending || !title || !content}
            >
              发布文章
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
