'use client';

import { useRequest } from 'ahooks';
import { ProList } from '@ant-design/pro-components';
import { Typography, theme, Flex } from 'antd';
import { ClockCircleOutlined } from '@ant-design/icons';
import { getUserNotifications, userNotificationsMarkAsRead } from '@/services';
import Viewer from '@/components/Viewer';
import PrettyTime from '@/components/PrettyTime';
import { useNotification } from '@/hooks';
import Link from 'next/link';
import React from 'react';
import { mergeSearchParams } from '@/utils';
import { useSearchParams } from 'next/navigation';

const Main: React.FC = () => {
  const searchParams = useSearchParams();

  const current = Number(searchParams.get('page') ?? 1);
  const pageSize = Number(searchParams.get('pageSize') ?? 10);

  const { unreadCount, setUnreadCount } = useNotification();

  const { loading, data, run, params } = useRequest(params => getUserNotifications({
    page: current,
    pageSize,
    ...params,
    include: 'FromUser',
  }), {
    debounceWait: 100,
    onSuccess() {
      if (unreadCount > 0) {
        userNotificationsMarkAsRead();
        setUnreadCount(0);
      }
    },
  });

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

          :global(.ant-pro-list-row-header-container, .ant-pro-list-row-title) {
            width: 100%;
          }

          :global(.title) {
            width: 100%;
            font-size: 16px;
            font-weight: bold;
            margin-bottom: 8px;
          }

          :global(.title a) {
            word-break: break-all;
            color: rgba(0, 0, 0, 0.88);
            line-height: 1.5;
            transition: color 0.3s;
          }

          :global(.title a:hover) {
            color: ${token.colorPrimary};
          }

          :global(.title .ant-typography) {
            white-space: nowrap;
          }

          :global(.unread) {
            background-color: ${token.colorPrimaryBg};
          }

          :global(.ant-list-item-meta-title a) {
            color: inherit;
            text-decoration: underline;
          }

          :global(.markdown-body) {
            width: 100%;
            font-size: 14px;
            background-color: transparent;
          }

          :global(.markdown-body .reply-quote) {
            padding: 6px 12px;
            background-color: ${token.colorBgLayout};
            display: -webkit-box;
            -webkit-line-clamp: 3;
            -webkit-box-orient: vertical;
            overflow: hidden;
            text-overflow: ellipsis;
            border-radius: ${token.borderRadiusLG}px;
            cursor: pointer;
            text-decoration: none;
            color: ${token.colorTextSecondary};
          }

          :global(.ant-list-item:hover .markdown-body .reply-quote) {
            background-color: ${token.colorInfoBg};
          }
        }
      `}</style>
      <div className="list">
        <ProList<API.Notification>
          rowKey="id"
          dataSource={data?.data}
          loading={loading}
          metas={{
            avatar: {
              dataIndex: ['fromUser', 'avatar'],
            },
            title: {
              render: (_, notification) => (
                <Flex className="title" justify="space-between">
                  <Viewer>{notification.subject}</Viewer>
                  <Typography.Text type="secondary">
                    <ClockCircleOutlined style={{ marginRight: 4 }} /><PrettyTime time={notification.createdAt} />
                  </Typography.Text>
                </Flex>
              ),
            },
            description: {
              render: (_, notification) => (
                <Viewer>{notification.message}</Viewer>
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

export default Main;