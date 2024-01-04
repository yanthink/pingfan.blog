'use client';

import { Menu, theme } from 'antd';
import React, { useEffect, useState } from 'react';
import { ProfileOutlined, BellOutlined, SafetyOutlined } from '@ant-design/icons';
import { useRouter } from '@/hooks';
import Base from './base';
import Notification from './notification';
import Security from './security';
import Redirect from '@/components/Redirect';
import Authorized from '@/components/Authorized';

const items = [
  { label: <span><ProfileOutlined /> 基本设置</span>, key: 'base' },
  { label: <span><BellOutlined /> 通知设置</span>, key: 'notification' },
  { label: <span><SafetyOutlined /> 修改密码</span>, key: 'security' },
];

const Main: React.FC = () => {
  const router = useRouter();

  const [mode, setMode] = useState<'inline' | 'horizontal'>();
  const [selectKey, setSelectKey] = useState('base');

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
              selectedKeys={[selectKey]}
              items={items}
              onClick={({ key }) => setSelectKey(key)}
            />
          </div>
          <div className="content">
            <div className="title">{items.find(item => item.key === selectKey)?.label}</div>
            {selectKey === 'base' && <Base />}
            {selectKey === 'notification' && <Notification />}
            {selectKey === 'security' && <Security />}
          </div>
        </div>
      </Authorized>
    </>
  );
};

export default Main;