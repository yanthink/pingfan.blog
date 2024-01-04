'use client';

import React, { useEffect } from 'react';
import { Form, Avatar, Card, Flex, Tag } from 'antd';
import {
  ProList,
  LightFilter,
  ProFormSwitch,
  ProFormSelect,
} from '@ant-design/pro-components';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { useRouter } from '@/hooks';
import { mergeSearchParams, prettyNumber } from '@/utils';
import {
  ClockCircleOutlined,
  EyeOutlined, FireOutlined,
  HeartOutlined,
  LikeOutlined, MessageOutlined,
} from '@ant-design/icons';
import PrettyTime from '@/components/PrettyTime';

interface MainProps {
  articles: API.Article[];
  total: number;
}

const Main: React.FC<MainProps> = ({ articles, total }) => {
  const [form] = Form.useForm();
  const router = useRouter();
  const searchParams = useSearchParams();

  const current = Number(searchParams.get('page') ?? 1);
  const pageSize = Number(searchParams.get('pageSize') ?? 10);

  useEffect(() => {
    form.setFieldsValue({
      sortFields: searchParams.get('sortFields') ?? undefined,
      filters: searchParams.get('filters') ?? undefined,
      onlyQueryTitle: searchParams.get('queryFields') === 'title',
    });
  }, [form, searchParams]);

  return (
    <>
      <style jsx>{`
        .list {
          :global(.ant-pro-list .ant-pro-card-body) {
            padding-inline: 0;
          }

          :global(.ant-pro-list-row-content) {
            margin-left: 0;
            margin-right: 0;
            line-height: 22px;
            display: -webkit-box;
            -webkit-line-clamp: 3;
            -webkit-box-orient: vertical;
            overflow: hidden;
            text-overflow: ellipsis;
            height: 66px;
          }

          :global(.ant-pro-list-row-title) {
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
            overflow: hidden;
            text-overflow: ellipsis;
          }

          :global(.title) {
            color: inherit;
          }

          .action {
            display: flex;
            line-height: 24px;
            align-items: center;
            gap: 4px;
          }

          :global(.ant-list-item-main) {
            max-width: 100%;
            overflow: hidden;
          }
        }

        @media screen and (max-width: 768px) {
          .list :global(.ant-list-item-extra) {
            display: none;
          }
        }
      `}</style>
      <div className="list">
        <Card
          title={<span>找到 {total} 条关于 <Tag style={{ margin: 0 }}>{searchParams.get('q')}</Tag> 的内容</span>}
          extra={
            <LightFilter
              form={form}
              onValuesChange={(_, values) => {
                const params = new URLSearchParams(searchParams)
                params.delete('sortFields');
                params.delete('filters');
                params.delete('queryFields');

                router.push(`?${mergeSearchParams(params, {
                  ...(values.sortFields && { sortFields: values.sortFields }),
                  ...(values.filters && { filters: values.filters }),
                  ...(values.onlyQueryTitle && { queryFields: 'title' }),
                })}`);
              }}
            >
              <ProFormSelect
                label="排序"
                name="sortFields"
                options={[
                  { label: <><LikeOutlined /> 点赞最多</>, value: '-like_count' },
                  { label: <><LikeOutlined /> 点赞最少</>, value: 'like_count' },
                  { label: <><MessageOutlined /> 评论最多</>, value: '-comment_count' },
                  { label: <><MessageOutlined /> 评论最少</>, value: 'comment_count' },
                  { label: <><FireOutlined /> 热度最高</>, value: '-hotness' },
                  { label: <><FireOutlined /> 热度最低</>, value: 'hotness' },
                  { label: <><ClockCircleOutlined /> 最新创建</>, value: '-created_at' },
                  { label: <><ClockCircleOutlined /> 最老创建</>, value: 'created_at' },
                ]}
              />
              <ProFormSelect
                label="时间"
                name="filters"
                options={[
                  { label: '1个月内', value: '1month' },
                  { label: '3个月内', value: '3months' },
                  { label: '半年内', value: '6months' },
                  { label: '1年内', value: '1year' },
                  { label: '2年内', value: '2years' },
                  { label: '3年内', value: '3years' },
                ]}
              />
              <ProFormSwitch label="只搜标题" name="onlyQueryTitle" />
            </LightFilter>
          }
        >
          <ProList<API.Article>
            itemLayout="vertical"
            rowKey="id"
            dataSource={articles}
            metas={{
              title: {
                render: (_, article) => (
                  <Link className="title" href={`/articles/${article.id}`}>
                    <span dangerouslySetInnerHTML={{ __html: article.highlights?.title?.[0] ?? article.title! }} />
                  </Link>
                ),
              },
              description: {
                render: (_, article) => article.tags?.map(tag => (
                  <Link key={tag.id} href={`?${mergeSearchParams(searchParams, { tagId: tag.id, page: 1 })}`}>
                    <Tag>{tag.name}</Tag>
                  </Link>
                )),
              },
              actions: {
                render: (_, article) => (
                  <Flex wrap="wrap" gap={12}>
                    <div className="action">
                      <Avatar src={article.user?.avatar} size={24} /> {article.user?.name}
                    </div>
                    <div className="action">
                      <ClockCircleOutlined /><PrettyTime time={article.createdAt} />
                    </div>
                    <div className="action">
                      <EyeOutlined />{prettyNumber(article.viewCount ?? 0)}
                    </div>
                    <div className="action">
                      <LikeOutlined />{prettyNumber(article.likeCount ?? 0)}
                    </div>
                    <div className="action">
                      <MessageOutlined />{prettyNumber(article.commentCount ?? 0)}
                    </div>
                    <div className="action">
                      <HeartOutlined /> {prettyNumber(article.favoriteCount ?? 0)}
                    </div>
                  </Flex>
                ),
              },
              extra: {
                render: (_: any, article: API.Article) => article.preview && (
                  <img src={article.preview} width={320} height={180} style={{ objectFit: 'contain' }} alt="预览" />
                ),
              },
              content: {
                render: (_, article) => (
                  <div
                    dangerouslySetInnerHTML={{
                      __html: article.highlights?.content?.[0] ??
                        article.textContent?.substring(0, 300),
                    }}
                  />
                ),
              },
            }}
            pagination={{
              total,
              hideOnSinglePage: true,
              showLessItems: true,
              responsive: true,
              showSizeChanger: true,
              current,
              pageSize,
              itemRender(page, type, defaultDom: any) {
                if (page < 1 || String(page) === (searchParams.get('page') ?? '1')) {
                  return defaultDom;
                }

                const element = React.cloneElement(defaultDom);

                return (
                  <Link href={`?${mergeSearchParams(searchParams, { page })}`}>
                    {element.type === 'a' ? element.props.children : defaultDom}
                  </Link>
                );
              },
              onShowSizeChange(current, size) {
                router.push(`?${mergeSearchParams(searchParams, { pageSize: size, page: 1 })}`);
              },
            }}
          />
        </Card>
      </div>
    </>
  );
};

export default Main;