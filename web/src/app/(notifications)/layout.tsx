'use client';

import { Menu, theme } from 'antd';
import { useEffect, useState } from 'react';
import { BellOutlined, MailOutlined, NotificationOutlined } from '@ant-design/icons';
import { usePathname } from 'next/navigation';
import { useRouter } from '@/hooks';
import Redirect from '@/components/Redirect';
import Authorized from '@/components/Authorized';

const items = [
  { label: <span><BellOutlined /> 通知</span>, key: '/notifications' },
  { label: <span><MailOutlined /> 私信</span>, key: '/messages' },
  { label: <span><NotificationOutlined /> 系统</span>, key: '/systems' },
];

interface NotificationLayoutProps extends React.PropsWithChildren {
}

const NotificationLayout: React.FC<NotificationLayoutProps> = ({ children }) => {
  const pathname = usePathname();
  const router = useRouter();

  const [mode, setMode] = useState<'inline' | 'horizontal'>();

  const { token } = theme.useToken();

  function resize() {
    requestAnimationFrame(() => setMode(window.innerWidth < 768 ? 'horizontal' : 'inline'));
  }

  useEffect(() => {
    window.addEventListener('resize', resize);
    resize();
    return () => window.removeEventListener('resize', resize);
  }, []);

  return (
    <>
      <style jsx>{`
        .main {
          display: flex;
          min-height: calc(100vh - 200px);
          padding-top: 16px;
          padding-bottom: 16px;
          background-color: #fff;
        }

        .menu {
          width: 224px;
          border-right: 1px solid ${token.colorSplit};
        }

        .menu :global(.ant-menu-inline) {
          border: none;
        }

        .menu :global(.ant-menu:not(.ant-menu-horizontal) .ant-menu-item-selected) {
          font-weight: bold;
        }

        .content {
          flex: 1;
          padding: 8px 32px;
          overflow: hidden;
        }

        .title {
          margin-bottom: 12px;
          color: ${token.colorTextHeading};
          font-weight: 500;
          font-size: 20px;
          line-height: 28px;
        }

        @media screen and (max-width: 768px) {
          .main {
            flex-direction: column;
            padding-top: 8px;
          }

          .menu {
            width: 100%;
            border: none;

            :global(.ant-menu-item) {
              padding-left: 12px !important;
            }
          }

          .content {
            padding: 24px 12px 0 12px;

            :global(.ant-pro-page-container-children-container) {
              padding: 0;
            }
          }
        }

      `}</style>
      <Authorized noMatch={<Redirect />}>
        <div className="main">
          <div className="menu">
            <Menu
              mode={mode}
              selectedKeys={[pathname]}
              items={items}
              onClick={({ key }) => router.push(key)}
            />
          </div>
          <div className="content">
            <div className="title">{items.find(item => item.key === pathname)?.label}</div>
            {children}
          </div>
        </div>
      </Authorized>
    </>
  );
};

export default NotificationLayout;