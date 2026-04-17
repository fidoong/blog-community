"use client";

import { useState } from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import {
  useNotifications,
  useMarkRead,
  useMarkAllRead,
} from "@/hooks/queries/use-notifications";
import { Button } from "@/components/ui/button";
import { Container } from "@/components/ui/container";
import { Skeleton } from "@/components/ui/skeleton";
import { PageHeader } from "@/components/ui/page-header";
import { EmptyState } from "@/components/ui/empty-state";
import {
  MessageSquare,
  Reply,
  Heart,
  UserPlus,
  Info,
  Check,
  CheckCheck,
} from "lucide-react";
import type { NotificationItem, NotificationType } from "@/types/notification";

const TYPE_CONFIG: Record<
  NotificationType,
  { label: string; icon: React.ReactNode; color: string }
> = {
  comment: {
    label: "评论",
    icon: <MessageSquare className="h-4 w-4" />,
    color: "bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300",
  },
  reply: {
    label: "回复",
    icon: <Reply className="h-4 w-4" />,
    color:
      "bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300",
  },
  like_post: {
    label: "点赞",
    icon: <Heart className="h-4 w-4" />,
    color: "bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300",
  },
  like_comment: {
    label: "点赞",
    icon: <Heart className="h-4 w-4" />,
    color: "bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300",
  },
  follow: {
    label: "关注",
    icon: <UserPlus className="h-4 w-4" />,
    color:
      "bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300",
  },
  system: {
    label: "系统",
    icon: <Info className="h-4 w-4" />,
    color: "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300",
  },
};

function getTargetLink(n: NotificationItem): string | undefined {
  if (!n.targetType || !n.targetId) return undefined;
  switch (n.targetType) {
    case "post":
      return `/posts/${n.targetId}`;
    case "comment":
      return `/posts/${n.targetId}`; // approximate
    case "user":
      return `/user/${n.targetId}`;
    default:
      return undefined;
  }
}

function NotificationRow({
  notification: n,
  index,
}: {
  notification: NotificationItem;
  index: number;
}) {
  const markRead = useMarkRead();
  const config = TYPE_CONFIG[n.type] ?? TYPE_CONFIG.system;
  const link = getTargetLink(n);
  const timeStr = new Date(n.createdAt * 1000).toLocaleString("zh-CN", {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });

  const handleClick = () => {
    if (!n.isRead) {
      markRead.mutate(n.id);
    }
  };

  const content = (
    <motion.div
      className={`flex items-start gap-4 py-4 border-b last:border-b-0 cursor-pointer transition-colors ${
        n.isRead ? "opacity-70" : ""
      } hover:bg-muted/30 rounded-lg px-2 -mx-2`}
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{
        duration: 0.2,
        delay: Math.min(index * 0.03, 0.3),
        ease: "easeOut",
      }}
      onClick={handleClick}
    >
      <div
        className={`mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-full ${config.color}`}
      >
        {config.icon}
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium">{n.title}</span>
          {!n.isRead && (
            <span className="h-2 w-2 rounded-full bg-red-500 shrink-0" />
          )}
        </div>
        <p className="mt-0.5 text-sm text-muted-foreground line-clamp-2">
          {n.content}
        </p>
        <span className="mt-1 text-xs text-muted-foreground">{timeStr}</span>
      </div>
      {!n.isRead && (
        <Button
          size="sm"
          variant="ghost"
          className="h-8 gap-1 shrink-0 text-muted-foreground hover:text-foreground"
          onClick={(e) => {
            e.stopPropagation();
            markRead.mutate(n.id);
          }}
          disabled={markRead.isPending}
        >
          <Check className="h-3.5 w-3.5" />
          已读
        </Button>
      )}
    </motion.div>
  );

  if (link) {
    return <Link href={link}>{content}</Link>;
  }
  return content;
}

export default function NotificationsPage() {
  const [page, setPage] = useState(1);
  const [unreadOnly, setUnreadOnly] = useState(false);

  const { data, isLoading } = useNotifications(page, 20, unreadOnly);
  const markAllRead = useMarkAllRead();

  const notifications = data?.list ?? [];

  return (
    <Container size="md" className="py-12">
      <div className="flex items-center justify-between mb-2">
        <PageHeader title="消息通知" description="查看你的所有通知和互动消息" />
        <Button
          size="sm"
          variant="outline"
          className="gap-1.5"
          onClick={() => markAllRead.mutate()}
          disabled={markAllRead.isPending || notifications.length === 0}
        >
          <CheckCheck className="h-4 w-4" />
          全部已读
        </Button>
      </div>

      {/* 筛选 */}
      <div className="flex items-center gap-1 mb-6 border-b">
        {[
          { key: false as const, label: "全部" },
          { key: true as const, label: "未读" },
        ].map((tab) => (
          <button
            key={String(tab.key)}
            onClick={() => {
              setUnreadOnly(tab.key);
              setPage(1);
            }}
            className={`px-4 py-2.5 text-sm font-medium transition-colors relative ${
              unreadOnly === tab.key
                ? "text-foreground"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            {tab.label}
            {unreadOnly === tab.key && (
              <motion.div
                className="absolute bottom-0 left-0 right-0 h-0.5 bg-foreground"
                layoutId="notifTab"
                transition={{ duration: 0.2, ease: "easeOut" }}
              />
            )}
          </button>
        ))}
      </div>

      {/* 列表 */}
      {isLoading && (
        <div className="space-y-4">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="flex items-start gap-4 py-4 border-b">
              <Skeleton className="h-9 w-9 rounded-full" />
              <div className="flex-1 space-y-2">
                <Skeleton className="h-5 w-1/2" />
                <Skeleton className="h-4 w-3/4" />
              </div>
            </div>
          ))}
        </div>
      )}

      {!isLoading && notifications.length === 0 && (
        <EmptyState
          title="暂无通知"
          description={
            unreadOnly ? "没有未读通知，好好休息吧！" : "你还没有收到任何通知"
          }
        />
      )}

      {!isLoading && notifications.length > 0 && (
        <div>
          {notifications.map((n, index) => (
            <NotificationRow key={n.id} notification={n} index={index} />
          ))}
        </div>
      )}

      {/* 分页 */}
      {!isLoading && data && data.pagination.totalPages > 1 && (
        <div className="mt-8 flex items-center justify-center gap-2">
          <Button
            size="sm"
            variant="outline"
            disabled={page <= 1}
            onClick={() => setPage((p) => p - 1)}
          >
            上一页
          </Button>
          <span className="text-sm text-muted-foreground">
            第 {page} / {data.pagination.totalPages} 页
          </span>
          <Button
            size="sm"
            variant="outline"
            disabled={page >= data.pagination.totalPages}
            onClick={() => setPage((p) => p + 1)}
          >
            下一页
          </Button>
        </div>
      )}
    </Container>
  );
}
