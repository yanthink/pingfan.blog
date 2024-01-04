'use client';

import React, { useContext, useEffect, useState } from 'react';
import classnames from 'classnames';
import { PageContainer } from '@ant-design/pro-components';
import { Tabs, Checkbox, Spin, message, notification } from 'antd';
import { LoginFormPage, ProFormText } from '@ant-design/pro-components';
import { LockOutlined, UserOutlined, LoadingOutlined } from '@ant-design/icons';
import Link from 'next/link';
import { useCountDown, useLockFn, useRequest } from 'ahooks';
import { getWxQRCode, login } from '@/services';
import { useSearchParams } from 'next/navigation';
import { useToken, useUser, useNotification, useRouter } from '@/hooks';
import { WebSocketContext, ReadyState } from '@/app/websocket';

interface MainProps {
}

const Main: React.FC<MainProps> = () => {
  const { readyState, sendMessage, latestMessage } = useContext(WebSocketContext);

  const [activeKey, setActiveKey] = useState('scan');
  const [leftTime, setLeftTime] = useState<number>();

  const router = useRouter();
  const searchParams = useSearchParams();

  const [, setToken] = useToken();
  const { setUser } = useUser();
  const { refresh: refreshNotification } = useNotification();

  const { data: wxQRCode, loading, run: refreshWxQRCode } = useRequest(getWxQRCode);
  useCountDown({ leftTime, onEnd: () => refreshWxQRCode() });

  useEffect(() => {
    if (readyState === ReadyState.Open && wxQRCode?.data.token) {
      sendMessage(JSON.stringify({ event: 'login', data: { token: wxQRCode?.data.token } }));
    }
  }, [readyState, sendMessage, wxQRCode]);

  useEffect(() => {
    if (wxQRCode?.data.expiresIn) {
      setLeftTime(wxQRCode.data.expiresIn);
    }
  }, [wxQRCode]);

  useEffect(() => {
    if (latestMessage) {
      const ret = JSON.parse(latestMessage.data) as { event: string; data?: Record<string, any> };

      if (ret.event === 'WxScanLoginSuccess') {
        setToken(ret.data!.token!);
        message.success('登录成功！');
        setUser(ret.data!.user);
        refreshNotification();

        const redirect = searchParams.get('redirect');
        if (!redirect && !ret.data!.user.name) {
          router.replace( '/settings');
        } else {
          router.replace(redirect || '/user');
        }
      }
    }
  }, [latestMessage]);

  const submit = useLockFn(async (values: Record<string, any>) => {
    const { data: { token, ...user } } = await login(values);
    setToken(token);
    message.success('登录成功！');

    setUser(user);
    refreshNotification();

    router.replace(searchParams.get('redirect') || '/');
  });

  function renderQRCode() {
    if (loading || !wxQRCode?.data.img) {
      return (
        <Spin size="large" indicator={<LoadingOutlined />} tip="正在加载...">
          <div style={{ width: 328, height: 328 }} />
        </Spin>
      );
    }

    return (
      <img
        className="qrcode"
        src={`data:image/png;base64,${wxQRCode.data.img}`}
        alt="小程序码"
        style={{ width: '100%' }}
      />
    );
  }

  return (
    <PageContainer>
      <div className={classnames('container', { 'active-scan': activeKey === 'scan' })}>
        <style jsx>{`
          .background-image {
            position: fixed;
            z-index: 0;
            background: url("https://gw.alipayobjects.com/zos/rmsportal/TVYTbAXWheQpRcWDaDMu.svg");
            inset: 24px;
          }

          .container :global(.ant-pro-form-login-page-container) {
            height: 480px;
          }

          .active-scan :global(button) {
            display: none;
          }
        `}</style>
        <div className="background-image" />
        <LoginFormPage onFinish={submit}>
          <Tabs
            size="large"
            items={[
              { key: 'scan', label: '微信扫码登录' },
              { key: 'account', label: '账户密码登录' },
            ]}
            activeKey={activeKey}
            onChange={setActiveKey}
            centered
          />
          <div
            style={{
              display: activeKey === 'scan' ? 'flex' : 'none',
              justifyContent: 'center',
              alignItems: 'center',
            }}
          >
            {renderQRCode()}
          </div>
          <div style={{ display: activeKey === 'account' ? 'block' : 'none' }}>
            <ProFormText
              name="name"
              fieldProps={{ size: 'large', prefix: <UserOutlined />, style: { marginTop: 12 } }}
              placeholder="用户名"
              rules={[{ required: true, message: '请输入用户名!' }]}
            />
            <ProFormText.Password
              name="password"
              fieldProps={{ size: 'large', prefix: <LockOutlined /> }}
              placeholder="密码"
              rules={[{ required: true, message: '请输入密码！' }]}
            />
            <div style={{ marginBlockEnd: 12, display: activeKey === 'account' ? 'block' : 'none' }}>
              <Checkbox>自动登录</Checkbox>
              <Link style={{ float: 'right' }} href="#">忘记密码</Link>
            </div>
          </div>
        </LoginFormPage>
      </div>
    </PageContainer>
  );
};

export default Main;