import React from 'react';
import Main from './main';
import request from '@/request/server';
import qs from 'querystring';

export default async function Home({ searchParams }: { searchParams?: Record<string, any> }) {
  console.log(`/api/articles?${qs.stringify({
    ...searchParams,
    include: 'User,Tags',
  })}`)

  const res = await request<API.Article[]>(`/api/articles?${qs.stringify({
    ...searchParams,
    include: 'User,Tags',
  })}`, { next: { tags: ['articles'] } });

  return <Main total={res.total ?? 0} articles={res.data} />;
}
