'use client';

import React from 'react';
import { Spin } from 'antd';
import { useUser } from '@/hooks';
import AvatarDropdown from '@/components/HeaderAvatar/AvatarDropdown';
import LoginNav from '@/components/HeaderAvatar/LoginNav';

const HeaderAvatar: React.FC = () => {
  const { user, loading } = useUser();

  if (loading) {
    return <Spin size="small" />;
  }

  if (user.id) {
    return <AvatarDropdown user={user} />;
  }

  return <LoginNav />;
};

export default HeaderAvatar;