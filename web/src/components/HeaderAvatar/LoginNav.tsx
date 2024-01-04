'use client';

import React from 'react';
import { Avatar, Space } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import Link from 'next/link';

interface LoginNavProps {
}

const PageLoginNav: React.FC<LoginNavProps> = () => {
  return (
    <Link href={'/login'}>
      <Space size={8}>
        <Avatar icon={<UserOutlined />} alt="登录" size="small" />
        <span style={{ color: 'rgba(0, 0, 0, .65)' }}>账户中心</span>
      </Space>
    </Link>
  );
};

export default PageLoginNav;