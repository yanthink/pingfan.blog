'use client';

import React from 'react';
import ArticleForm from '@/components/ArticleForm';
import { useLockFn } from 'ahooks';
import { updateArticle } from '@/services';
import { message } from 'antd';
import { useRouter } from '@/hooks';
import Redirect from '@/components/Redirect';
import Authorized from '@/components/Authorized';

interface PageFormProps {
  article: API.Article;
}

const Main: React.FC<PageFormProps> = ({ article }) => {
  const router = useRouter();

  const update = useLockFn(async (values: Record<string, any>) => {
    const { data } = await updateArticle(article.id!, values);
    message.success('文章更新成功');
    router.replace(`/articles/${data.id}`);
  });

  return (
    <Authorized noMatch={<Redirect />}>
      <ArticleForm article={article} onSubmit={update} />
    </Authorized>
  );
};

export default Main;