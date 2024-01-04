'use client';

import React, { useEffect, useState } from 'react';
import { ProCard } from '@ant-design/pro-components';
import { Affix, Button, Card, Col, Avatar, Typography, Divider, Row, Flex, Space, theme } from 'antd';
import { MailOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons';
import { useSearchParams } from 'next/navigation';
import { useRouter, useUser } from '@/hooks';
import Favorites from './favorites';
import Comments from './comments';
import Likes from './likes';
import Upvotes from './upvotes';
import { mergeSearchParams } from '@/utils';
import Authorized from '@/components/Authorized';
import Redirect from '@/components/Redirect';

const tabKeys = ['favorites', 'comments', 'replies', 'likes', 'upvotes'];

const Main: React.FC = () => {
  const { user, loading } = useUser();
  const router = useRouter();

  const { token } = theme.useToken();

  const searchParams = useSearchParams();
  const tab = searchParams.get('tab') ?? 'favorites';

  const [tabKey, setTabKey] = useState<string>(tabKeys.includes(tab) ? tab : 'favorites');

  useEffect(() => {
    setTabKey(tab);
  }, [tab]);

  return (
    <>
      <style jsx>{`
        @media (max-width: 768px) {
          .main :global(.left) {
            width: 100%;
            margin-bottom: 12px;
            display: none;
          }

          .main :global(.affix) {
            > div:first-child[aria-hidden='true'] {
              display: none;
            }

            :global(.ant-affix) {
              position: static !important;
            }
          }
        }
      `}</style>
      <Authorized noMatch={<Redirect />}>
        <div className="main">
          <Row gutter={24}>
            <Col lg={6} md={24} className="left">
              <Affix className="affix" offsetTop={56}>
                <Card bordered={false} loading={loading}>
                  <Flex align="center" vertical>
                    <Avatar alt="头像" src={user.avatar} size={120}
                      style={{ boxShadow: token.boxShadow, marginBottom: 16 }}
                    />
                    <Typography.Title level={4} style={{ marginBottom: 4 }}>{user.name}</Typography.Title>
                    <span>暂无个人描述</span>
                  </Flex>
                  <Space size={4} style={{ marginTop: 24, paddingLeft: 24 }}><MailOutlined /> {user.email || '暂无~'}
                  </Space>
                  <Divider />
                  <Button
                    block
                    size="large"
                    icon={<EditOutlined />}
                    onClick={() => router.push('/settings')}
                  >
                    编辑个人资料
                  </Button>
                  {user.role === 1 && (
                    <Button
                      size="large"
                      type="primary"
                      icon={<PlusOutlined />}
                      block
                      onClick={() => router.push('/articles/create')}
                      style={{ marginTop: 12 }}
                    >
                      发布文章
                    </Button>
                  )}
                </Card>
              </Affix>
            </Col>
            <Col lg={18} md={24} className="right">
              <ProCard
                tabs={{
                  activeKey: tabKey,
                  items: [
                    { label: '收藏', key: 'favorites', children: <Favorites /> },
                    { label: '评论', key: 'comments', children: <Comments /> },
                    { label: '回复', key: 'replies', children: <Comments type={1} /> },
                    { label: '文章点赞', key: 'likes', children: <Likes /> },
                    { label: '评论点赞', key: 'upvotes', children: <Upvotes /> },
                  ],
                  onChange: key => router.push(`?${mergeSearchParams(searchParams, { tab: key, page: 1 })}`),
                }}
              />
            </Col>
          </Row>
        </div>
      </Authorized>
    </>
  );
};

export default Main;