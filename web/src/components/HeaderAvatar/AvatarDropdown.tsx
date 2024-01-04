'use client';

import React from 'react';
import { Avatar, Dropdown, Space } from 'antd';
import { UserOutlined, SettingOutlined, LogoutOutlined } from '@ant-design/icons';
import { useRouter, useLogout } from '@/hooks';

interface AvatarDropdownProps {
  user: API.User;
}

const AvatarDropdown: React.FC<AvatarDropdownProps> = ({ user }) => {
  const router = useRouter();
  const logout = useLogout();

  return (
    <Dropdown
      placement="bottomRight"
      menu={{
        onClick(e) {
          switch (e.key) {
            case 'logout':
              logout();
              break;
            case 'user':
              router.push('/user');
              break;
            case 'settings':
              router.push('/settings');
              break;
          }
        },
        items: [
          {
            key: 'user',
            icon: <UserOutlined />,
            label: '个人中心',
          },
          {
            key: 'settings',
            icon: <SettingOutlined />,
            label: '个人设置',
          },
          {
            type: 'divider',
          },
          {
            key: 'logout',
            icon: <LogoutOutlined />,
            label: '退出登录',
          },
        ],
      }}

    >
      <Space>
        <Avatar src={user.avatar} alt="头像" />
        <span>{user.name}</span>
      </Space>
    </Dropdown>
  );
};

export default AvatarDropdown;