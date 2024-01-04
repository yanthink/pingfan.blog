declare type RequestResponse<T = unknown, M = Record<string, any>> = {
  code: number;
  message: string;
  success: boolean;
  data: T;
  total?: number;
  prev_cursor?: string;
  next_cursor?: string;

  [key: string]: any;
} & M;

declare type RequestPageParams = {
  current?: number;
  pageSize?: number;
};
