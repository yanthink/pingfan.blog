import React from 'react';
import { useRequest } from 'ahooks';
import { ProList } from '@ant-design/pro-components';
import { theme, Flex } from 'antd';
import { getUserUpvotes } from '@/services';
import Link from 'next/link';
import Viewer from '@/components/Viewer';
import { ClockCircleOutlined, LikeTwoTone, MessageOutlined } from '@ant-design/icons';
import PrettyTime from '@/components/PrettyTime';
import { stripTagsStrategy } from '@/markdown';
import { mergeSearchParams, prettyNumber } from '@/utils';
import { useRouter } from '@/hooks';
import { useSearchParams } from 'next/navigation';

const Upvotes: React.FC = () => {
  const searchParams = useSearchParams();

  const current = Number(searchParams.get('page') ?? 1);
  const pageSize = Number(searchParams.get('pageSize') ?? 10);

  const { loading, data, run, params } = useRequest(params => getUserUpvotes({
    page: current,
    pageSize,
    ...params,
    include: 'Comment.User,Comment.Article',
  }), { debounceWait: 300 });

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

          :global(.ant-pro-list-row-content) {
            margin-left: 0;
            margin-right: 0;
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

          :global(.markdown-body) {
            width: 100%;
            font-size: 14px;
            background-color: transparent;
          }
        }

        .comment-preview {
          padding: 6px 12px;
          background-color: ${token.colorBgLayout};
          display: -webkit-box;
          -webkit-line-clamp: 3;
          -webkit-box-orient: vertical;
          overflow: hidden;
          text-overflow: ellipsis;
          border-radius: ${token.borderRadiusLG}px;
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

        :global(.ant-list-item:hover) .comment-preview {
          background-color: ${token.colorInfoBg};
        }

        @media screen and (max-width: 768px) {
          .list :global(.ant-list-item-extra) {
            display: none;
          }
        }
      `}</style>
      <div className="list">
        <ProList<API.Upvote>
          itemLayout="vertical"
          rowKey="id"
          loading={loading}
          dataSource={data?.data}
          metas={{
            title: {
              render: (_, upvote) => (
                <Link className="title" href={`/articles/${upvote.comment?.articleId}`}>
                  {upvote.comment?.article?.title}
                </Link>
              ),
            },
            actions: {
              render: (_, { comment, createdAt }) => (
                <Flex wrap="wrap" gap={12}>
                  <div className="action">
                    <ClockCircleOutlined /><PrettyTime time={createdAt} />
                  </div>
                  <div className="action">
                    <LikeTwoTone twoToneColor={token.colorPrimary} />{prettyNumber(comment?.upvoteCount ?? 0)}
                  </div>
                  <div className="action">
                    <MessageOutlined />{prettyNumber(comment?.replyCount ?? 0)}
                  </div>
                </Flex>
              ),
            },
            content: {
              render: (_, { comment }) => (
                <div
                  className="comment-preview"
                  onClick={() => router.push(`/articles/${comment?.articleId}?pinnedId=${comment?.id}#评论`)}
                >
                  <a>{comment?.user?.name}</a>：
                  <Viewer options={stripTagsStrategy}>
                    {
                      (comment?.content ?? '')
                        .replace(/^```.*?\n$/g, '')
                        .split('\n')
                        .slice(0, 5)
                        .join('\n')
                    }
                  </Viewer>
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

export default Upvotes;