class Request {
  config: RequestConfig = {
    url: '',
    baseUrl: '',
    method: 'GET',
    header: {
      Accept: 'application/json',
      'Content-Type': 'application/json; charset=utf-8',
    },
    timeout: 60000,
    dataType: 'json',
    responseType: 'text',
    sslVerify: true,
    withCredentials: false,
  };

  interceptor: RequestInterceptor = {};

  setConfig(config: Partial<RequestConfig>) {
    this.config = {
      ...this.config,
      ...config,
      header: { ...this.config.header, ...config.header },
    };
  };

  async request<T = unknown, M = Record<string, any>>(options: RequestConfig): Promise<RequestResponse<T, M>> {
    const opts = {
      ...this.config,
      ...options,
      header: { ...this.config.header, ...options.header },
      url: `${options.baseUrl ?? this.config.baseUrl}${options.url}`,
    };

    if (this.interceptor.request && typeof this.interceptor.request === 'function') {
      await this.interceptor.request(opts);
    }

    const res = await uni.request(opts);

    if (this.interceptor.response && typeof this.interceptor.response === 'function') {
      return this.interceptor.response<T, M>(res as UniRequestResponse<T, M>, opts);
    }

    return res.data as RequestResponse<T, M>;
  };

  async get<T = unknown, M = Record<string, any>>(url: string, data: RequestConfig['data'] = {}, header: RequestConfig['header'] = {}) {
    return this.request<T, M>({ method: 'GET', url, data, header });
  };

  async post<T = unknown, M = Record<string, any>>(url: string, data: RequestConfig['data'] = {}, header: RequestConfig['header'] = {}) {
    return this.request<T, M>({ method: 'POST', url, data, header });
  };

  async put<T = unknown, M = Record<string, any>>(url: string, data: RequestConfig['data'] = {}, header: RequestConfig['header'] = {}) {
    return this.request<T, M>({ method: 'PUT', url, data, header });
  };

  async delete<T = unknown, M = Record<string, any>>(url: string, data: RequestConfig['data'] = {}, header: RequestConfig['header'] = {}) {
    return this.request<T, M>({ method: 'DELETE', url, data, header });
  };
}

export default Request;
