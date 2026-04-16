export interface ApiResponse<T> {
  code: string;
  message: string;
  data: T;
  timestamp: number;
}

export interface PaginationParams {
  page: number;
  pageSize: number;
}

export interface PaginationMeta {
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface PaginatedResponse<T> {
  list: T[];
  pagination: PaginationMeta;
}
