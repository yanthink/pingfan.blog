'use client';

import Authorized from '@/components/Authorized';
import { ProList } from '@ant-design/pro-components';
import React, { useEffect, useState } from 'react';
import { getComments, updateComment } from '@/services';
import { useRouter } from '@/hooks';
import PrettyTime from '@/components/PrettyTime';
import { useLockFn } from 'ahooks';

const Main: React.FC = () => {
  const router = useRouter();

  const [itemLayout, setItemLayout] = useState<'vertical' | 'horizontal'>();

  function resize() {
    requestAnimationFrame(() => setItemLayout(window.innerWidth < 768 ? 'vertical' : 'horizontal'));
  }

  useEffect(() => {
    window.addEventListener('resize', resize);
    resize();
    return () => window.removeEventListener('resize', resize);
  }, []);

  const update = useLockFn(updateComment);

  return (
    <>
      <style jsx>{`
        .list {
          :global(.ant-pro-list-row-content) {
            flex: none;
            margin-right: 0;
          }
        }
      `}</style>
      <Authorized authority={user => user.role === 1}>
        <div className="list">
          <ProList<API.Comment>
            rowKey="id"
            itemLayout={itemLayout}
            headerTitle="评论列表"
            request={async params => getComments({ ...params, include: 'Article,User', sort: 'desc' })}
            pagination={{}}
            editable={{
              actionRender: (row, config, dom) => [dom.save, dom.cancel],
              onSave: async (key, record, originRow) => {
                await update(originRow.id!, record);
                return true;
              },
            }}
            metas={{
              title: {
                render: (_, record) => (
                  <div
                    onClick={() => {
                      router.push(`/articles/${record.articleId}?pinnedId=${record.id}&parentId=${record.parentId}#评论`);
                    }}
                  >
                    {record.article?.title}
                  </div>
                ),
                editable: false,
              },
              avatar: {
                dataIndex: ['user', 'avatar'],
                editable: false,
              },
              content: {
                render: (_, record) => <PrettyTime time={record.createdAt} />,
                editable: false,
              },
              description: {
                dataIndex: 'content',
              },
              actions: {
                render: (text, row, index, action) => [
                  <a
                    onClick={() => {
                      action?.startEditable(row.id!);
                    }}
                    key="link"
                  >
                    编辑
                  </a>,
                ],
              },
            }}
          />
        </div>
      </Authorized>
    </>
  );
};

export default Main;