declare interface RequestConfig extends UniPromiseOptions<UniNamespace.RequestOptions> {
  baseUrl?: string;
  skipErrorHandler?: boolean;
}

declare type RequestResponse<T = unknown, M = Record<string, any>> = {
  code: number;
  message: string;
  success: boolean;
  data: T;
  errors?: Record<string, string[]>;
  total?: number;
  cursor?: {
    after?: string;
    before?: string;
  };
  response?: UniRequestResponse<T, M>;
  meta?: Record<string, any>;

  [key: string]: any;
} & M;

declare type UniRequestResponse<T = unknown, M = Record<string, any>> =
  { data: RequestResponse<T, M> }
  & Omit<UniNamespace.RequestSuccessCallbackResult, 'data'>;

declare type RequestInterceptor = {
  request?: (options: RequestConfig) => Promise<void>;
  response?: <T = unknown, M = Record<string, any>>(
    response: UniRequestResponse<T, M>,
    options: RequestConfig,
  ) => Promise<RequestResponse<T, M>>
}

