import Main from './main';
import { Metadata } from "next";
import request from "@/request/server";
import qs from 'querystring';

export const metadata: Metadata = {
  title: '编辑文章 - 平凡的博客',
};

export default async function Create({ params, searchParams }: any) {
  if (!params.id || params.id === 'undefined') {
    return null;
  }

  const { data } = await request<API.Article>(`/api/articles/${params.id}?${qs.stringify({
    ...searchParams,
    include: 'Tags'
  })}`);

  return <Main article={data} />;
}