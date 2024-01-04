import React from 'react';
import Link from 'next/link';
import { Badge } from 'antd';
import { BellOutlined } from '@ant-design/icons';
import { useNotification, useUser } from '@/hooks';

const HeaderNoticeIcon: React.FC = () => {
  const { user } = useUser();
  const { unreadCount } = useNotification();

  if (user.id) {
    return (
      <>
        <style jsx>{`
          .action {
            display: flex;
          }
        `}</style>
        <Link href="/notifications" style={{ display: 'inline-flex' }}>
          <Badge count={unreadCount} styles={{ root: { fontSize: 16 } }}>
            <BellOutlined type="bell" style={{ color: 'rgba(0, 0, 0, 0.45)' }} />
          </Badge>
        </Link>
      </>
    );
  }

  return null;
};

export default HeaderNoticeIcon;