import { defineStore } from 'pinia';
import { ref, watch } from 'vue';
import { useNotification, useToken } from '@/hooks';

export enum ReadyState {
  Connecting = 0,
  Open = 1,
  Closing = 2,
  Closed = 3,
}

const RETRIES_LIMIT = 3;

export const useWebsocketStore = defineStore('websocket', () => {
  let websocket: UniNamespace.SocketTask | undefined = undefined;
  let retries = 0;
  let hasLogged = false;

  let reconnectTimer: ReturnType<typeof setTimeout> | undefined = undefined;
  let pingTimer: ReturnType<typeof setInterval> | undefined = undefined;

  const token = useToken();
  const { refresh: refreshNotification } = useNotification();

  const readyState = ref<ReadyState>(ReadyState.Closed);
  const latestMessage = ref<UniNamespace.OnSocketMessageCallbackResult>();

  function reconnect() {
    if (retries < RETRIES_LIMIT && readyState.value !== ReadyState.Open) {
      clearTimeout(reconnectTimer);
    }

    reconnectTimer = setTimeout(() => {
      connectWs();
      retries++;
    }, 3_000);
  }

  function connectWs() {
    const ws = uni.connectSocket({
      url: import.meta.env.VITE_SOCKET_URL,
      success(res) {
        console.log('socket 连接成功！', res);
      },
      fail(res) {
        console.log('socket 连接失败！', res);
      },
    });
    readyState.value = ReadyState.Connecting;

    ws?.onError(() => {
      if (websocket !== ws) {
        return;
      }

      readyState.value = ReadyState.Closed;
      reconnect();
    });

    ws?.onOpen(() => {
      if (websocket !== ws) {
        return;
      }

      readyState.value = ReadyState.Open;
    });

    ws?.onMessage(message => {
      if (websocket !== ws) {
        return;
      }

      console.log('socket onMessage', message);

      const { event } = JSON.parse(message.data) as { event: string; data?: Record<string, any> };
      if (event.startsWith('notifications.')) {
        refreshNotification();
      }

      latestMessage.value = message;
    });

    ws?.onClose(() => {
      if (websocket !== ws) {
        return;
      }

      if (!websocket || websocket === ws) {
        readyState.value = ReadyState.Closed;
      }
    });

    websocket = ws;
  }

  function connect() {
    retries = 0;
    connectWs();
  }

  function send(options: UniNamespace.SendSocketMessageOptions) {
    if (readyState.value !== ReadyState.Open) {
      throw new Error('WebSocket disconnected');
    }

    websocket?.send(options);
  }

  // ping
  watch(readyState, newReadyState => {
    clearInterval(pingTimer);

    if (newReadyState === ReadyState.Open) {
      pingTimer = setInterval(() => send({
        data: JSON.stringify({ event: 'ping' }),
      }), 50_000);
    }
  });

  // 登录/退出
  watch([readyState, token], ([newReadyState, newToken]) => {
    if (newReadyState === ReadyState.Open && !!newToken) {
      send({ data: JSON.stringify({ event: 'login', data: { token: newToken } }) });
      hasLogged = true;
    } else if (hasLogged) {
      send({ data: JSON.stringify({ event: 'logout' }) });
      hasLogged = false;
    }
  });

  connect();

  return { websocket, readyState, latestMessage, connect, send };
});