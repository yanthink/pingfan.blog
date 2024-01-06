import React from 'react';
import type { Metadata } from 'next';
import { cookies, headers } from 'next/headers';
import { ConfigProvider } from 'antd';
import locale from 'antd/locale/zh_CN';
import theme from '@/theme';
import StyledRegistry from './registry';
import Websocket from './websocket';
import Layout from '@/components/Layout';
import NProgressBar from '@/components/NProgressBar';
import request from '@/request/server';
import './global.scss';

export const metadata: Metadata = {
  title: '平凡的博客',
  description: '平凡的博客 - 个人博客网站，这是一个由Next.js、Uniapp（微信小程序）和Golang打造的开源个人博客网站。探索技术、生活和创意的空间。分享前沿技术、编程经验和生活见闻。',
};

async function getAuthUser(): Promise<API.User | undefined> {
  const token = cookies().get(process.env.NEXT_PUBLIC_TOKEN_KEY ?? 'token')?.value ?? '';
  if (token) {
    const { data } = await request<API.User>('/api/user', { cache: 'no-store' });

    return data;
  }
}

interface RootLayoutProps extends React.PropsWithChildren {
}

const RootLayout: React.FC<RootLayoutProps> = async ({ children }) => {
  const pathname = headers().get('pathname');
  const user = await getAuthUser();

  return (
    <html lang="zh">
    <body>
    <ConfigProvider locale={locale} theme={theme}>
      <StyledRegistry>
        <Websocket>
          <Layout pathname={pathname} user={user}>
            {children}
          </Layout>
          <NProgressBar />
        </Websocket>
      </StyledRegistry>
    </ConfigProvider>
    </body>
    </html>
  );
};

export default RootLayout;
