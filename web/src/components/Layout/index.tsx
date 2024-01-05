'use client';

import React, { useEffect, useMemo, useRef, useState } from 'react';
import { HappyProvider } from '@ant-design/happy-work-theme';
import { ProLayout } from '@ant-design/pro-components';
import { Layout as AntdLayout, Menu, FloatButton, Space, Button, Drawer, theme, type MenuProps } from 'antd';
import {
  MessageOutlined,
  ReadOutlined,
  SettingOutlined,
  TagsOutlined,
  UserOutlined,
  GithubOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined,
} from '@ant-design/icons';
import { usePathname } from 'next/navigation';
import Link from 'next/link';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';
import { useUser } from '@/hooks';
import { check, logout } from '@/utils';
import HeaderAvatar from '@/components/HeaderAvatar';
import HeaderNoticeIcon from '@/components/HeaderNoticeIcon';
import HeaderSearch from '@/components/HeaderSearch';

dayjs.locale('zh-cn');
dayjs.extend(relativeTime);

const { Header, Content, Footer } = AntdLayout;

type Menus = {
  key: string;
  label: React.ReactNode;
  icon?: React.ReactElement;
  access?: string;
  children?: Menus;
}[]

const menus: Menus = [
  { label: '文章列表', icon: <ReadOutlined />, key: '/' },
  {
    label: '系统管理',
    icon: <SettingOutlined />,
    key: 'system',
    access: 'admin',
    children: [
      { label: '用户列表', icon: <UserOutlined />, key: '/users' },
      { label: '评论列表', icon: <MessageOutlined />, key: '/comments' },
      { label: '标签列表', icon: <TagsOutlined />, key: '/tags' },
    ],
  },
];

interface LayoutProps extends React.PropsWithChildren {
  pathname?: string | null;
  user?: API.User;
}

const Layout: React.FC<LayoutProps> = (props) => {
  const pathname = usePathname() ?? props.pathname;
  const serverSideSetUser = useRef(false);
  
  const { token: themeToken } = theme.useToken();
  const { user, setUser } = useUser();

  if (!serverSideSetUser.current) {
    if (props.user) {
      setUser(props.user); 
    } else if (check()) {
      logout();
    }
    serverSideSetUser.current = true;
  }

  const menuItems = useMemo((): MenuProps['items'] => {
    function getMenus(menus: Menus): MenuProps['items'] {
      return menus.filter(item => !item.access || item.access === 'admin' && user.role === 1).map(item => {
        let label = item.label;

        if (!item.children && pathname !== item.key) {
          label = <Link href={item.key}>{item.label}</Link>;
        }

        return { ...item, label, children: item.children && getMenus(item.children) };
      });
    }

    return getMenus(menus);
  }, [pathname, user.role]);

  const [isMobile, setIsMobile] = useState(false);
  const [collapsed, setCollapsed] = useState(true);

  useEffect(() => {
    setCollapsed(true);
  }, [pathname]);

  function resize() {
    requestAnimationFrame(() => setIsMobile(window.innerWidth < 768));
  }

  useEffect(() => {
    window.addEventListener('resize', resize);
    resize();
    return () => window.removeEventListener('resize', resize);
  }, []);

  return (
    <>
      <style jsx>{`
        .header-main {
          display: flex;
          align-items: center;
          height: 56px;
          max-width: 1200px;
          margin: 0 auto;

          :global(.ant-menu) {
            border-bottom: 0;
          }

          :global(.ant-menu-submenu.ant-menu-submenu-horizontal:hover) {
            //border-radius: 6px;
            //background-color: rgba(0, 0, 0, 0.03);
          }
        }

        .content-main {
          max-width: 1200px;
          margin: 24px auto;
          min-height: calc(100vh - 172px);

          :global(a) {
            &:hover {
              --ant-color-link-hover: ${themeToken.colorPrimary};
            }
          }
        }

        .footer-main {
          max-width: 1200px;
          text-align: center;
          margin: auto;
          a {
            color: ${themeToken.colorText};
            &:hover {
              text-decoration: underline;
            }
          }
        }

        .logo {
          display: flex;
          height: 100%;
          align-items: center;
          overflow: hidden;

          h1 {
            display: inline-block;
            line-height: 24px;
            margin-left: 6px;
            font-weight: 600;
            font-size: 16px;
            color: rgba(0, 0, 0, 0.88);
            vertical-align: top;
          }
        }

        .dropdown-avatar {
          height: 40px;
          line-height: 40px;
          padding: 8px;
          border-radius: 6px;
          color: rgba(0, 0, 0, 0.45);
          cursor: pointer;
          display: flex;
          align-items: center;

          &:hover {
            background-color: rgba(0, 0, 0, 0.03);
          }
        }

        .actions {
          display: flex;
          align-items: center;

          > * {
            display: flex;
            padding: 8px;
            border-radius: 6px;
            font-size: 16px;

            &:hover {
              background-color: rgba(0, 0, 0, 0.03);
            }
          }
        }

        @media screen and (max-width: 1200px) {
          .content-main {
            margin: 12px;
          }
        }
      `}</style>
      <HappyProvider>
        {/* todo 这里如果不用 ProLayout 不知为何Pro组件的主色无法生效 */}
        <ProLayout colorPrimary={themeToken.colorPrimary} pure>
          <AntdLayout>
            <Header
              style={{
                position: 'fixed',
                width: '100%',
                zIndex: 999,
                backdropFilter: 'blur(8px)',
                borderBottom: '1px solid rgba(5, 5, 5, 0.06)',
              }}
            >
              <div className="header-main">
                {isMobile ? (
                  <Space style={{ flex: 1 }}>
                    <Button
                      type="text"
                      icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                      onClick={() => setCollapsed(!collapsed)}
                      style={{
                        fontSize: '16px',
                      }}
                    />
                    <div className="logo">
                      <Link href="/" style={{ display: 'flex', alignItems: 'center' }}>
                        <img alt="logo" loading="lazy" width="32" height="32" src="/logo.svg" />
                      </Link>
                    </div>
                  </Space>
                ) : (
                  <>
                    <div className="logo">
                      <Link href="/" style={{ display: 'flex', alignItems: 'center' }}>
                        <img alt="logo" loading="lazy" width="32" height="32" src="/logo.svg" />
                        <h1 className="jsx-3694805390">平凡的博客</h1>
                      </Link>
                    </div>
                    <Menu
                      mode="horizontal"
                      selectedKeys={[pathname]}
                      items={menuItems}
                      style={{ padding: '0 12px', flex: 1 }}
                    />
                  </>
                )}
                <Space>
                  {!isMobile && <HeaderSearch />}
                  <div className="actions">
                    <a
                      href="https://github.com/yanthink/pingfan.blog"
                      target="_blank"
                      style={{ color: 'inherit' }}
                    >
                      <GithubOutlined />
                    </a>
                    <div><HeaderNoticeIcon key="notice" /></div>
                  </div>
                  <div className="dropdown-avatar">
                    <HeaderAvatar />
                  </div>
                </Space>
              </div>
            </Header>
            <Content style={{ paddingTop: 56 }}>
              <div className="content-main">
                {props.children}
              </div>
            </Content>
            <Footer>
              <div className="footer-main">©2023 平凡的博客 <a href="https://beian.miit.gov.cn" target="_blank">粤ICP备18080782号-1</a></div>
            </Footer>
          </AntdLayout>
          <Drawer
            placement="left"
            open={!collapsed}
            width={256}
            closable={false}
            onClose={() => setCollapsed(true)}
            styles={{ body: { padding: 0 } }}
          >
            <Menu
              mode="inline"
              selectedKeys={[pathname]}
              items={menuItems}
            />
          </Drawer>
          <FloatButton.BackTop />
        </ProLayout>
      </HappyProvider>
    </>
  );
};

export default Layout;