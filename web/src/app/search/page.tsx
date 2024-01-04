import React from 'react';
import Main from './main';
import request from '@/request/server';

export default async function Home({ searchParams }: { searchParams?: Record<string, any> }) {
  const res = await request<API.Article[]>(`/api/articles/search?${new URLSearchParams({
    ...searchParams,
    include: 'User,Tags',
  })}`, { cache: 'no-cache' });

  return <Main total={res.total ?? 0} articles={res.data} />;
}
