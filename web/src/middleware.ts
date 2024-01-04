import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import ip from 'ip';

export function middleware(request: NextRequest) {
  const { headers, nextUrl } = request;

  headers.set('pathname', new URL(request.url).pathname);
  const proxies = process.env.PROXIES?.split(',') ?? []

  if (request.ip) {
    if (proxies.some(proxy => ip.cidrSubnet(proxy).contains(request.ip!))) {
      headers.append('X-Forwarded-For', request.ip!);
    } else {
      headers.set('X-Forwarded-For', request.ip!);
    }
  }

  // nextUrl.searchParams.set('url', request.url);

  return NextResponse.rewrite(nextUrl, {
    request: {
      headers,
    },
  });
}

export const config = {
  matcher: '/:path*',
};