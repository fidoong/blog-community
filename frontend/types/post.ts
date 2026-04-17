export interface Post {
  id: number;
  title: string;
  content?: string;
  summary: string;
  contentType: "markdown" | "rich_text";
  coverImage: string;
  authorId: number;
  authorName?: string;
  status: "draft" | "pending" | "published" | "rejected";
  viewCount: number;
  likeCount: number;
  commentCount: number;
  collectCount: number;
  tags: string[];
  publishedAt?: number;
  createdAt: number;
  updatedAt: number;
}

export interface Pagination {
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface PostsResponse {
  list: Post[];
  pagination: Pagination;
}

export interface CreatePostPayload {
  title: string;
  content: string;
  contentType?: "markdown" | "rich_text";
  coverImage?: string;
  tags?: string[];
}
