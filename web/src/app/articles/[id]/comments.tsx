'use client';

import React, { useState } from 'react';
import { Flex, List, message, Space, Tag, theme } from 'antd';
import { Comment } from '@ant-design/compatible';
import { MessageOutlined, DoubleRightOutlined } from '@ant-design/icons';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { useAuthFn, useRouter, useUser } from '@/hooks';
import { mergeSearchParams } from '@/utils';
import PrettyTime from '@/components/PrettyTime';
import Viewer from '@/components/Viewer';
import Like from '@/components/Like';
import CommentEditor from './comment.editor';
import { createComment, upvoteComment } from '@/services';
import { useUpdate } from 'ahooks';
import classnames from 'classnames';

interface CommentsProps {
  article: API.Article;
  comments: API.Comment[];
  total: number;
}

const Comments: React.FC<CommentsProps> = ({ article, comments, total }) => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { token } = theme.useToken();

  const current = Number(searchParams.get('page') ?? 1);
  const pageSize = Number(searchParams.get('pageSize') ?? 10);
  const pinnedId = Number(searchParams.get('pinnedId') ?? 0);

  const [replyCommentId, setReplyCommentId] = useState<number>();
  const update = useUpdate();

  const upvote = useAuthFn(async (comment: API.Comment) => {
    const { data } = await upvoteComment(comment.id!);
    Object.assign(comment, { upvoteCount: data.upvoteCount, hasUpvoted: data.hasUpvoted });
    update();
  });

  const { user } = useUser();

  const create = useAuthFn(async data => {
    const response = await createComment(data);
    response.data.user = user;

    return response;
  });

  return (
    <div className="comment-container">
      <style jsx>{`
        .comment-container {
          :global(.ant-list-item) {
            padding: 0;
          }

          :global(.bytemd) {
            height: 200px;
          }

          :global(.ant-comment) {
            width: 100%;
            background-color: transparent;
          }

          :global(.markdown-body) {
            font-size: 14px;
            background-color: transparent;
            margin-top: 8px;
          }

          :global(.markdown-body p) {
            margin-bottom: 16px;
          }

          :global(.markdown-body pre:last-child) {
            margin-bottom: 0;
          }

          .nested-container {
            background-color: ${token.colorBgLayout};
            border-radius: ${token.borderRadiusLG}px;
          }

          :global(.ant-comment-nested .ant-comment) {
            padding: 0 12px;
          }

          :global(.ant-tag) {
            margin: 0;
            padding: 3px 6px;
            line-height: 11px;
          }

          :global(.ant-comment-actions li) {
            margin-right: 12px;
            cursor: pointer;
          }

          .reply, :global(.like) {
            display: flex;
          }

          .reply :global(span), :global(.like span) {
            margin-right: 2px;
          }

          :global(.more-replies) {
            padding: 0 12px 12px 56px;
          }

          :global(.nested-comment-content) {
            display: flex;
            flex-wrap: wrap;
            overflow: hidden;
          }

          :global(.nested-comment-content > a) {
            line-height: 1.6;
            margin-top: 8px;
            margin-right: 8px;
          }

          :global(.nested-comment-content .markdown-body) {
            max-width: 100%;
          }

          :global(.pinned) {
            background-color: ${token.colorInfoBg};
            padding: 0 16px;
          }
        }
      `}</style>
      <CommentEditor
        onSubmit={async values => {
          const { data } = await create({ ...values, articleId: article.id }) ?? {};
          message.success('评论成功！');
          comments.unshift(data!);
          update();
        }}
      />
      <List
        dataSource={comments}
        rowKey={comment => comment.id === pinnedId && comment === comments[0] ? 'pinned' : comment.id!}
        header={
          <Flex justify="space-between" align="center">
            <span>{total} 评论</span>
            {(searchParams.has('cid') || searchParams.has('parentId') || searchParams.has('commentId')) && (
              <Link href="?#评论">全部评论</Link>
            )}
          </Flex>
        }
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
              <Link href={`?${mergeSearchParams(searchParams, { page })}#评论`}>
                {element.type === 'a' ? element.props.children : defaultDom}
              </Link>
            );
          },
          onShowSizeChange(current, size) {
            router.push(`?${mergeSearchParams(searchParams, { pageSize: size, page: 1 })}#评论`);
          },
        }}
        renderItem={(comment, index) => (
          <List.Item
            id={comment.id === pinnedId && index === 0 ? `pinned-comment-${comment.id}` : `comment-${comment.id}`}
            className={classnames({ pinned: comment.id === pinnedId && index === 0 })}
          >
            <Comment
              author={
                <Space size="small">
                  {comment.user?.name}
                  {comment.userId === article.userId && <Tag>博主</Tag>}
                </Space>
              }
              avatar={comment.user?.avatar}
              content={<Viewer>{comment.content}</Viewer>}
              datetime={<PrettyTime time={comment.createdAt} />}
              actions={[
                <Like key="upvote" liked={comment.hasUpvoted} suffix={comment.upvoteCount}
                  onClick={() => upvote(comment)}
                />,
                <div
                  key="reply" className="reply"
                  onClick={() => setReplyCommentId(replyCommentId === comment.id ? 0 : comment.id)}
                >
                  <MessageOutlined /> 回复
                </div>,
              ]}
            >
              <div className="nested-container" style={{ marginBottom: comment.replies?.length ? 16 : 0 }}>
                {comment.id === replyCommentId && (
                  <CommentEditor
                    submitText="回复" placeholder={`@${comment.user?.name}`}
                    onSubmit={async values => {
                      const { data } = await createComment({
                        ...values,
                        articleId: article.id,
                        parentId: comment.id,
                      }) ?? {};
                      message.success('回复成功！');
                      comment.replies ??= [];
                      comment.replies.unshift(data!);
                      setReplyCommentId(0);
                      update();
                    }}
                  />
                )}
                {comment.replies?.map((reply, i) => (
                  <div key={reply.id}>
                    <Comment
                      author={
                        <Space size="small">
                          {reply.user?.name}
                          {reply.userId === article.userId && <Tag color={token.colorPrimary}>博主</Tag>}
                          {reply.userId === comment?.userId && <Tag>楼主</Tag>}
                        </Space>
                      }
                      avatar={reply.user?.avatar}
                      content={
                        <div className="nested-comment-content">
                          {reply.parentId != reply.commentId && reply.parent?.user && (
                            <Link href={`?${mergeSearchParams({}, { parentId: reply.parentId })}#评论`}>
                              @{reply.parent?.user.name}
                            </Link>
                          )}
                          <Viewer>{reply.content}</Viewer>
                        </div>
                      }
                      datetime={<PrettyTime time={reply.createdAt} />}
                      actions={[
                        <Like key="upvote" liked={reply.hasUpvoted} suffix={reply.upvoteCount}
                          onClick={() => upvote(reply)}
                        />,
                        <div
                          key="reply" className="reply"
                          onClick={() => setReplyCommentId(replyCommentId === reply.id ? 0 : reply.id)}
                        >
                          <MessageOutlined /> 回复
                        </div>,
                      ]}
                    />
                    {reply.id === replyCommentId && (
                      <CommentEditor submitText="回复" placeholder={`@${reply.user?.name}`}
                        onSubmit={async values => {
                          const { data } = await createComment({
                            ...values,
                            articleId: article.id,
                            parentId: reply.id,
                          }) ?? {};
                          message.success('回复成功！');
                          if (i < comment.replies!.length - 1) {
                            comment.replies!.splice(i + 1, 0, data!);
                          } else {
                            comment.replies!.push(data!);
                          }
                          setReplyCommentId(0);
                          update();
                        }}
                      />
                    )}
                  </div>
                ))}
                {!searchParams.get('commentId') && !searchParams.get('parentId') && comment.replyCount! > (comment.replies?.length ?? 0)! && (
                  <div className="more-replies">
                    <Link href={`?${mergeSearchParams({}, { commentId: comment.id })}#评论`}>
                      更多{comment.replyCount! - (comment.replies?.length ?? 0)}条回复 <DoubleRightOutlined />
                    </Link>
                  </div>
                )}
              </div>
            </Comment>
          </List.Item>
        )}
      />
    </div>
  );
};

export default Comments;