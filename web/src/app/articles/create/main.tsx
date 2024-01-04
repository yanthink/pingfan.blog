'use client';

import React from 'react';
import ArticleForm from '@/components/ArticleForm'
import { useLockFn } from "ahooks";
import { createArticle } from "@/services";
import { message } from "antd";
import { useRouter } from "@/hooks";

const Main: React.FC = () => {
  const router = useRouter();

  const create = useLockFn(async (values: Record<string, any>) => {
    const { data } = await createArticle(values)
    message.success("文章新建成功")
    router.replace(`/articles/${data.id}`);
  })

  return (
    <ArticleForm onSubmit={create} />
  )
};

export default Main;