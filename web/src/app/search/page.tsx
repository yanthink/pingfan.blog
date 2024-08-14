import React from 'react';
import Main from './main';
import request from '@/request/server';
import qs from 'querystring';

export default async function Home({ searchParams }: { searchParams?: Record<string, any> }) {
  const res = await request<API.Article[]>(`/api/articles/search?${qs.stringify({
    ...searchParams,
    include: 'User,Tags',
  })}`, { cache: 'no-store' });

  return <Main total={res.total ?? 0} articles={res.data} />;
}
