'use client';

import { useRequest } from 'ahooks';
import { Typography, theme, Flex } from 'antd';
import Link from 'next/link';
import { ProList } from '@ant-design/pro-components';
import { ClockCircleOutlined, LikeOutlined, MessageOutlined } from '@ant-design/icons';
import { getUserComments } from '@/services';
import Viewer from '@/components/Viewer';
import PrettyTime from '@/components/PrettyTime';
import { stripTagsStrategy } from '@/markdown';
import { mergeSearchParams, prettyNumber } from '@/utils';
import React from 'react';
import { useRouter } from '@/hooks';
import { useSearchParams } from 'next/navigation';

interface CommentsProps {
  type?: number;
}

const Comments: React.FC<CommentsProps> = ({ type = 0 }) => {
  const searchParams = useSearchParams();

  const current = Number(searchParams.get('page') ?? 1);
  const pageSize = Number(searchParams.get('pageSize') ?? 10);

  const { loading, data, run, params } = useRequest(params => getUserComments(mergeSearchParams(searchParams, {
    type,
    page: current,
    pageSize,
    ...params,
    include: type === 1 ? 'Article,Parent.User' : 'Article',
  })), { debounceWait: 300 });

  const router = useRouter();
  const { token } = theme.useToken();

  return (
    <>
      <style jsx>{`
        .list {
          :global(.ant-pro-card-body) {
            padding: 0;
          }

          :global(.ant-list-item-meta) {
            margin-bottom: 0;
          }

          :global(.title a) {
            margin-top: 0;
            margin-bottom: 12px;
            margin-right: 16px;
            word-break: break-all;
            color: rgba(0, 0, 0, 0.88);
            font-size: 16px;
            font-weight: bold;
            line-height: 1.5;
            transition: color 0.3s;
          }

          :global(.title a:hover) {
            color: ${token.colorPrimary};
          }

          :global(.title .ant-typography) {
            white-space: nowrap;
          }

          :global(.ant-pro-list-row-content) {
            margin-left: 0;
            margin-right: 0;
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

          :global(.markdown-body) {
            width: 100%;
            font-size: 14px;
            background-color: transparent;
          }
        }

        .reply-quote {
          padding: 6px 12px;
          background-color: ${token.colorBgLayout};
          display: -webkit-box;
          -webkit-line-clamp: 3;
          -webkit-box-orient: vertical;
          overflow: hidden;
          text-overflow: ellipsis;
          border-radius: ${token.borderRadiusLG}px;
          margin: 8px 0 0;
          cursor: pointer;

          :global(.markdown-body) {
            display: inline;
            margin-top: 0;
          }

          :global(.markdown-body:before) {
            display: none;
          }

          :global(.markdown-body:after) {
            display: none;
          }
        }

        :global(.ant-list-item:hover) .reply-quote {
          background-color: ${token.colorInfoBg};
        }
      `}</style>
      <div className="list">
        <ProList<API.Comment>
          itemLayout="vertical"
          rowKey="id"
          dataSource={data?.data}
          loading={loading}
          metas={{
            actions: {
              render: (_, comment) => (
                <Flex key="actions" wrap="wrap" gap={12}>
                  <div className="action">
                    <LikeOutlined style={{ marginRight: 4 }} />{prettyNumber(comment.upvoteCount ?? 0)}
                  </div>
                  <div className="action">
                    <MessageOutlined style={{ marginRight: 4 }} />{prettyNumber(comment.replyCount ?? 0)}
                  </div>
                </Flex>
              ),
            },
            content: {
              render: (_, comment) => (
                <div className="content">
                  <Flex className="title" justify="space-between">
                    <Link href={`/articles/${comment.articleId}?pinnedId=${comment.id}#评论`}>
                      {comment.article?.title}
                    </Link>
                    <Typography.Text type="secondary">
                      <ClockCircleOutlined style={{ marginRight: 4 }} /><PrettyTime time={comment.createdAt} />
                    </Typography.Text>
                  </Flex>
                  <Viewer>{comment.content}</Viewer>
                  {comment.parent && (
                    <div className="reply-quote"
                      onClick={() => router.push(`/articles/${comment.articleId}?parentId=${comment.parentId}#评论`)}
                    >
                      <a>{comment.parent?.user?.name}</a>：
                      <Viewer options={stripTagsStrategy}>
                        {
                          (comment.parent?.content ?? '')
                          .replace(/^```.*?\n$/g, '')
                          .split('\n')
                          .slice(0, 5)
                          .join('\n')
                        }
                      </Viewer>
                    </div>
                  )}
                </div>
              ),
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

export default Comments;