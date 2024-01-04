'use client';

import Authorized from '@/components/Authorized';
import { ProTable, type ProColumns, ActionType } from '@ant-design/pro-components';
import React, { useRef } from 'react';
import { getUsers, updateUser } from '@/services';
import dayjs from 'dayjs';
import { message, Popconfirm } from 'antd';
import { useLockFn } from 'ahooks';

const Main: React.FC = () => {
  const actionRef = useRef<ActionType>();

  const update = useLockFn(updateUser);

  const columns: ProColumns<API.User>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '头像',
      dataIndex: 'avatar',
      valueType: 'avatar',
      width: 60,
      hideInSearch: true,
    },
    {
      title: '用户名',
      dataIndex: 'name',
      width: 120,
      ellipsis: true,
      copyable: true,
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      width: 160,
      ellipsis: true,
      copyable: true,
    },
    {
      title: 'openid',
      dataIndex: 'openid',
      width: 160,
      ellipsis: true,
      copyable: true,
    },
    {
      title: '角色',
      dataIndex: 'role',
      valueType: 'select',
      valueEnum: new Map([
        [0, '普通用户'],
        [1, '管理员'],
      ]),
      width: 80,
    },
    {
      title: '状态',
      dataIndex: 'status',
      valueType: 'radio',
      valueEnum: new Map([
        [0, { text: '正常', status: 'Success' }],
        [1, { text: '锁定', status: 'Error' }],
      ]),
      width: 60,
    },
    {
      title: '注册时间',
      dataIndex: 'createdAt',
      valueType: 'dateTime',
      hideInSearch: true,
      width: 170,
    },
    {
      title: '注册时间',
      dataIndex: 'dateTimeRange',
      valueType: 'dateTimeRange',
      fieldProps: {
        presets: [
          { label: '今天', value: [dayjs().startOf('d'), dayjs().endOf('d')] },
          {
            label: '昨天',
            value: [dayjs().add(-1, 'd').startOf('d'), dayjs().add(-1, 'd').endOf('d')],
          },
          {
            label: '前天',
            value: [dayjs().add(-2, 'd').startOf('d'), dayjs().add(-2, 'd').endOf('d')],
          },
          { label: '最近一周', value: [dayjs().add(-7, 'd').startOf('d'), dayjs()] },
          { label: '最近一月', value: [dayjs().add(-1, 'month').startOf('d'), dayjs()] },
        ],
        getPopupContainer: () => document.body,
      },
      search: {
        transform: (value) => {
          return {
            startAt: value[0],
            endAt: value[1],
          };
        },
      },
      hideInTable: true,
    },
    {
      title: '更新时间',
      dataIndex: 'updatedAt',
      valueType: 'dateTime',
      hideInSearch: true,
      width: 170,
    },
    {
      title: '操作',
      valueType: 'option',
      fixed: 'right',
      width: 60,
      render(_, record) {
        return (
          <Popconfirm
            title={record.status ? '解锁用户' : '锁定用户'}
            description={record.status ? '确定解锁用户吗？' : '确定锁定该用户吗？'}
            placement="topRight"
            onConfirm={async () => {
              await update(record.id!, { status: 1 - record.status! });
              message.success('操作成功');
              actionRef.current?.reload();
            }}
            onOpenChange={() => console.log('open change')}
          >
            <a>{record.status ? '解锁' : '锁定'}</a>
          </Popconfirm>
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
      <ProTable<API.User>
        actionRef={actionRef}
        rowKey="id"
        columns={columns}
        request={getUsers}
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