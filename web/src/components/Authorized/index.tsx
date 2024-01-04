'use client';

import React, { useMemo } from 'react';
import { Result, Button, Spin } from 'antd';
import { useRouter, useUser } from '@/hooks';

interface AuthorizedProps extends React.PropsWithChildren {
  authority?: boolean | ((user: API.User) => boolean);
  noMatch?: React.ReactNode;
}

const Authorized: React.FC<AuthorizedProps> = ({ authority, noMatch, children }) => {
  const { user, loading } = useUser();
  const router = useRouter();

  const authorized = useMemo(() => {
    if (!user.id || authority === false) {
      return false;
    }

    if (typeof authority === 'function') {
      return authority(user);
    }

    return true;
  }, [user, authority]);

  if (loading) {
    return (
      <Spin size="large" />
    )
  }

  if (!authorized) {
    return noMatch ?? (
      <Result
        status="403"
        title="403"
        subTitle="你没有权限访问该页面。"
        extra={<Button type="primary" onClick={() => router.push('/')}>返回首页</Button>}
      />
    );
  }

  return children
};

export default Authorized;