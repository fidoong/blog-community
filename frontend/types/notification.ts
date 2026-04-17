export type NotificationType =
  | "comment"
  | "reply"
  | "like_post"
  | "like_comment"
  | "follow"
  | "system";

export interface NotificationItem {
  id: number;
  type: NotificationType;
  title: string;
  content: string;
  actorId?: number;
  targetId?: number;
  targetType?: string;
  isRead: boolean;
  createdAt: number;
}

export interface NotificationsResponse {
  list: NotificationItem[];
  pagination: {
    total: number;
    page: number;
    pageSize: number;
    totalPages: number;
  };
}

export interface UnreadCountResponse {
  count: number;
}
