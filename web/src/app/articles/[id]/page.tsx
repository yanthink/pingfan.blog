import request from '@/request/server';
import Main from './main';
import type { Metadata } from 'next';
import qs from 'querystring';

interface ArticleProps {
  params?: { id?: string };
  searchParams?: Record<string, any>;
}

async function getArticle(id: string, params: Record<string, any>) {
  return request<API.Article>(`/api/articles/${id}?${qs.stringify(params)}`);
}

async function getComments(id: string, params: Record<string, any>) {
  if (!params.commentId && !params.parentId) {
    params.commentId = '0';
  }

  if (params.commentId || params.parentId) {
    params.wrapId = params.commentId || params.parentId;
  }

  return request<API.Comment[]>(`/api/comments?${new URLSearchParams(params)}`);
}

export async function generateMetadata({ params, searchParams }: ArticleProps): Promise<Metadata> {
  if (!params?.id || params?.id === 'undefined') {
    return {
      title: '文章详情',
    };
  }

  const { data } = await getArticle(params.id, { include: 'User,Tags' });

  return {
    title: `${data.title} - 平凡的博客`,
  };
}

export default async function Article({ params, searchParams }: ArticleProps) {
  if (!params?.id || params?.id === 'undefined') {
    return null;
  }

  const [article, comments] = await Promise.all([
    getArticle(params.id, { include: 'User,Tags' }),
    getComments(params.id, { ...searchParams, articleId: params.id, include: 'User,Replies.User', withReplyUser: '1' }),
  ]);

  return <Main article={article.data} comments={comments.data} commentTotal={comments.total ?? 0} />;
}