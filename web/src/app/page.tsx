import React from 'react';
import Main from './main';
import request from '@/request/server';

export default async function Home({ searchParams }: { searchParams?: Record<string, any> }) {
  const res = await request<API.Article[]>(`/api/articles?${new URLSearchParams({
    ...searchParams,
    include: 'User,Tags',
  })}`, { next: { tags: ['articles'] } });

  return <Main total={res.total ?? 0} articles={res.data} />;
}
