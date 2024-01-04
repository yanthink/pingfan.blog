'use client';

import { message } from 'antd';
import { tokenStore } from '@/stores';
import { logout } from '@/utils/auth';
import baseRequest from '@/request/base';

interface InitOptions extends RequestInit {
  skipErrorHandler?: boolean;
  revalidate?: { paths?: string[]; tags?: string[] };
}

export default async function request<T = unknown, M = Record<string, any>>(url: string, init: InitOptions = {}) {
  const { skipErrorHandler, revalidate, ...requestInit } = init;

  const { token } = tokenStore;
  if (token) {
    requestInit.headers = {
      'Authorization': `Bearer ${token}`,
      ...requestInit.headers,
    } as any;
  }

  try {
    const response = await baseRequest<T, M>(url, requestInit);

    if (revalidate) {
      baseRequest(String(new URL(`/api/revalidate`, location.href)), {
        method: 'POST',
        body: JSON.stringify(revalidate),
      });
    }

    return response;
  } catch (e: any) {
    if (!skipErrorHandler) {
      message.error(e.message || e.statusText || 'Request Error');

      switch (e?.code) {
        case 401_001:
          logout();
          break;
      }
    }

    return Promise.reject(e);
  }
}
