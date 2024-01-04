'use client';

import React from 'react';
import { Typography, Space, Tag, theme, Avatar, Flex } from 'antd';
import { ClockCircleOutlined, EyeOutlined, TagsFilled, EditOutlined } from '@ant-design/icons';
import Link from 'next/link';
import PrettyTime from '@/components/PrettyTime';
import Viewer, { type ViewerProps } from '@/components/Viewer';
import { prettyNumber } from '@/utils';
import Like from '@/components/Like';
import Favorite from '@/components/Favorite';
import { useAuthFn, useUser } from '@/hooks';
import { likeArticle, favoriteArticle } from '@/services';

interface ContentProps extends ViewerProps {
  article: API.Article;
  onChange?(article: API.Article): void;
}

const Content: React.FC<ContentProps> = ({ article, onChange, ...viewerProps }) => {
  const { token } = theme.useToken();

  const { user } = useUser();

  const like = useAuthFn(async (id: string | number) => {
    const { data } = await likeArticle(id);
    onChange?.({ ...article, likeCount: data.likeCount, hasLiked: data.hasLiked });
  });

  const favorite = useAuthFn(async (id: string | number) => {
    const { data } = await favoriteArticle(id);
    onChange?.({ ...article, favoriteCount: data.favoriteCount, hasFavorited: data.hasFavorited });
  });

  return (
    <>
      <style jsx>{`
        .header {
          margin-right: -24px;
          margin-left: -24px;
          padding: 0 24px 12px 24px;
          border-bottom: 1px solid ${token.colorBorder};
          margin-bottom: 24px;
        }

        .extra {
          color: ${token.colorTextSecondary};
          font-size: 14px;
          margin-top: 8px;
          overflow: hidden;
        }

        .extra :global(.ant-typography) {
          display: flex;
          align-items: center;
        }

        .extra :global(.anticon), .extra :global(.ant-avatar) {
          margin-right: 6px;
        }

        // 代码块
        .content :global(.markdown-body .pre-wrap) {
          margin: 0 -24px 16px;
        }

        .content :global(.markdown-body pre) {
          border-radius: 0;
        }

        .content :global(.markdown-body *) {
          scroll-margin-top: 60px;
        }
      `}</style>
      <div className="header">
        <Typography.Title level={3}>{article.title}</Typography.Title>
        <div className="extra">
          <Space align="center" split={<Typography.Text type="secondary">⋅</Typography.Text>} wrap>
            <Typography.Text type="secondary">
              <Avatar src={article.user?.avatar} size={24} /> {article.user?.name}
            </Typography.Text>
            <Typography.Text type="secondary">
              <ClockCircleOutlined /> <PrettyTime time={article.createdAt} />
            </Typography.Text>
            <Typography.Text type="secondary">
              <EyeOutlined /> {prettyNumber(article.viewCount ?? 0)} 阅读
            </Typography.Text>
            {article.tags?.length && (
              <Typography.Text type="secondary">
                <TagsFilled style={{ marginRight: 12, fontSize: 16 }} />
                {article.tags?.map(tag => (
                  <Link key={tag.id} href={`/?tagId=${tag.id}`}>
                    <Tag>{tag.name}</Tag>
                  </Link>
                ))}
              </Typography.Text>
            )}
            {user?.id == article.userId && (
              <Typography.Text type="secondary">
                <Link href={`/articles/${article.id}/edit`} style={{ color: 'inherit' }}>
                  <EditOutlined />编辑
                </Link>
              </Typography.Text>
            )}
          </Space>
        </div>
      </div>
      <div className="content">
        <Viewer {...viewerProps} />
      </div>
      <Flex justify="center" gap={24} style={{ marginTop: 40 }}>
        <Like size="large" liked={article.hasLiked} suffix={article.likeCount} onClick={() => like(article.id!)} />
        <Favorite
          size="large" favorited={article.hasFavorited} suffix={article.favoriteCount}
          onClick={() => favorite(article.id!)}
        />
      </Flex>
    </>
  );
};

export default Content;