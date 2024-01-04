'use client';

import Authorized from '@/components/Authorized';
import { EditableProTable, type ProColumns, ActionType } from '@ant-design/pro-components';
import React, { useRef } from 'react';
import { getTags, createTag, updateTag } from '@/services';
import { useLockFn } from 'ahooks';
import { EditOutlined } from '@ant-design/icons';
import { message } from 'antd';

const Main: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const uidRef = useRef(-1);

  const create = useLockFn(createTag);
  const update = useLockFn(updateTag);

  const columns: ProColumns<API.Tag>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      editable: false,
    },
    {
      title: '名称',
      dataIndex: 'name',
    },
    {
      title: '排序',
      dataIndex: 'sort',
      valueType: 'digit',
      hideInSearch: true,
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      valueType: 'dateTime',
      hideInSearch: true,
      editable: false,
    },
    {
      title: '更新时间',
      dataIndex: 'updatedAt',
      valueType: 'dateTime',
      hideInSearch: true,
      editable: false,
    },
    {
      title: '操作',
      valueType: 'option',
      fixed: 'right',
      width: 100,
      render(_, record, i, action) {
        return (
          <a onClick={() => action?.startEditable?.(record.id!)}>
            <EditOutlined /> 编辑
          </a>
        );
      },
    },
  ];

  const x = columns.reduce(
    (x, { width = 120, hideInTable }) => x + (hideInTable ? 0 : Number(width)),
    0,
  );

  return (
    <Authorized authority={user => user.role === 1}>
      <EditableProTable<API.Tag>
        actionRef={actionRef}
        rowKey="id"
        columns={columns}
        request={getTags}
        pagination={{
          defaultPageSize: 10,
          hideOnSinglePage: true,
          showLessItems: true,
          responsive: true,
          showSizeChanger: true,
        }}
        recordCreatorProps={{
          position: 'top',
          newRecordType: 'cache',
          creatorButtonText: '新建标签',
          record() {
            return {
              id: uidRef.current,
              sort: 0,
            };
          },
        }}
        editable={{
          actionRender: (row, config, dom) => [dom.save, dom.cancel],
          async onSave(key, row, originRow, newLine) {
            if (newLine) {
              await create(row);
              uidRef.current--;
            } else {
              await update(originRow.id!, row);
            }
            actionRef.current?.reload();
            message.success('操作成功！');
          },
        }}
        options={{ fullScreen: true, reload: true, setting: true, density: true }}
        scroll={{ x }}
        sticky={{
          offsetHeader: 56,
        }}
      />
    </Authorized>
  );
};

export default Main;