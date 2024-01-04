import React from 'react';
import { useRequest } from 'ahooks';
import { ProList } from '@ant-design/pro-components';
import { Avatar, Flex, Tag, theme } from 'antd';
import { getUserLikes } from '@/services';
import Link from 'next/link';
import { ClockCircleOutlined, EyeOutlined, LikeTwoTone } from '@ant-design/icons';
import PrettyTime from '@/components/PrettyTime';
import { mergeSearchParams, prettyNumber } from '@/utils';
import { useSearchParams } from 'next/navigation';

const Likes: React.FC = () => {
  const searchParams = useSearchParams();

  const current = Number(searchParams.get('page') ?? 1);
  const pageSize = Number(searchParams.get('pageSize') ?? 10);

  const { loading, data, run, params } = useRequest(params => getUserLikes({
    page: current,
    pageSize,
    ...params,
    include: 'Article.Tags,Article.User',
  }), { debounceWait: 300 });

  const { token } = theme.useToken();

  return (
    <>
      <style jsx>{`
        .list {
          :global(.ant-pro-card-body) {
            padding: 0;
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
        <ProList<API.Like>
          itemLayout="vertical"
          rowKey="id"
          dataSource={data?.data}
          loading={loading}
          metas={{
            title: {
              render: (_, like) => (
                <Link className="title" href={`/articles/${like.articleId}`}>
                  {like.article?.title}
                </Link>
              ),
            },
            description: {
              render: (_, { article }) => article?.tags?.map(tag => (
                <Tag key={tag.id}>{tag.name}</Tag>
              )),
            },
            actions: {
              render: (_, { article, createdAt }) => (
                <Flex wrap="wrap" gap={12}>
                  <div className="action">
                    <Avatar src={article?.user?.avatar} size={24} /> {article?.user?.name}
                  </div>
                  <div className="action">
                    <ClockCircleOutlined /><PrettyTime time={createdAt} />
                  </div>
                  <div className="action">
                    <EyeOutlined />{prettyNumber(article?.viewCount ?? 0)}
                  </div>
                  <div className="action">
                    <LikeTwoTone twoToneColor={token.colorPrimary} />{prettyNumber(article?.likeCount ?? 0)}
                  </div>
                </Flex>
              ),
            },
            extra: {
              render: (_: any, { article }: API.Like) => article?.preview && (
                <img src={article.preview} width={320} height={180} style={{ objectFit: 'contain' }} alt="预览" />
              ),
            },
            content: {
              render: (_, { article }) => article?.textContent?.substr(0, 300),
            },
          }}
          pagination={{
            total: data?.total,
            hideOnSinglePage: true,
            showLessItems: true,
            responsive: true,
            current,
            pageSize,
            onChange(page, pageSize) {
              run({ ...params[0], page, pageSize });
            },
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
          }}
        />
      </div>
    </>
  );
};

export default Likes;