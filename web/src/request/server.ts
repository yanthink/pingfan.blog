import { cookies, headers } from 'next/headers';
import baseRequest from './base';

export default async function serverRequest<T = unknown, M = Record<string, any>>(url: string, init: RequestInit = {}) {
  const { ...requestInit } = init as Record<string, any>;

  requestInit.headers ??= {};

  const token = cookies().get(process.env.NEXT_PUBLIC_TOKEN_KEY ?? 'token')?.value ?? '';
  Object.assign(requestInit.headers, {
    'Authorization': `Bearer ${token}`,
    // 如果服务端通过nginx等中间件代理，服务端需要将代理服务ip和Next.js服务ip设置为信任ip。
    // location / {
    //   proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    //   proxy_pass http://localhost:3000;
    // }
    'X-Forwarded-For': headers().get("X-Forwarded-For"),
  })

  return baseRequest<T, M>(url, requestInit as RequestInit);
}