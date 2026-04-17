export interface Comment {
  id: number;
  content: string;
  authorId: number;
  likeCount: number;
  replies?: Comment[];
  createdAt: number;
}

export interface CommentsResponse {
  list: Comment[];
  pagination: {
    total: number;
    page: number;
    pageSize: number;
    totalPages: number;
  };
}

export interface CreateCommentPayload {
  content: string;
  parentId?: number;
}
