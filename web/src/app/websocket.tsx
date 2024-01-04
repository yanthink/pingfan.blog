'use client';

import React, { useEffect, useRef } from 'react';
import { useWebSocket } from 'ahooks';
import type { Result } from 'ahooks/es/useWebSocket';
import { useNotification, useRouter, useToken } from '@/hooks';
import { notification } from 'antd';

export enum ReadyState {
  Connecting = 0,
  Open = 1,
  Closing = 2,
  Closed = 3,
}

export const WebSocketContext = React.createContext<Result>({
  sendMessage(data) {
  },
  disconnect() {
  },
  connect() {
  },
  readyState: ReadyState.Connecting,
});

interface WebsocketProps extends React.PropsWithChildren {
}

const Websocket: React.FC<WebsocketProps> = ({ children }) => {
  const { refresh: refreshNotification } = useNotification();
  const router = useRouter();

  const result = useWebSocket(process.env.NEXT_PUBLIC_SOCKET_URL!, {
    onMessage(m) {
      console.log('socket onMessage', m);
      const ret = JSON.parse(m.data) as { event: string; data?: Record<string, any> };

      if (ret.event.startsWith('notifications.')) {
        notification.info({
          message: '您有一条新通知',
          description: ret.data?.message,
          onClick() {
            router.push('/notifications');
          },
        });

        refreshNotification();
      }
    },
  });

  const { readyState, sendMessage } = result;

  const hasLogged = useRef(false);

  const [token] = useToken();

  useEffect(() => {
    if (readyState === ReadyState.Open) {
      if (token) {
        sendMessage?.(JSON.stringify({ event: 'login', data: { token } }));
        hasLogged.current = true;
      } else if (hasLogged.current) {
        sendMessage?.(JSON.stringify({ event: 'logout' }));
        hasLogged.current = false;
      }
    }
  }, [readyState, sendMessage, token]);

  // ping
  useEffect(() => {
    if (readyState === ReadyState.Open) {
      const timer = setInterval(() => sendMessage?.(JSON.stringify({ event: 'ping' })), 50_000);
      return () => clearInterval(timer);
    }
  }, [readyState, sendMessage]);

  return (
    <WebSocketContext.Provider value={result}>
      {children}
    </WebSocketContext.Provider>
  );
};

export default Websocket;