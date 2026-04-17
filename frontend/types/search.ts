export interface SearchHit {
  id: number;
  title: string;
  summary: string;
  content?: string;
  authorId: number;
  authorName: string;
  tags: string[];
  viewCount: number;
  likeCount: number;
  commentCount: number;
  createdAt: number;
  score: number;
  highlight?: Record<string, string[]>;
}

export interface SearchResponse {
  list: SearchHit[];
  pagination: {
    total: number;
    page: number;
    pageSize: number;
    totalPages: number;
  };
  took: number;
}
