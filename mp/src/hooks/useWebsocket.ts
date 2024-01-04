import { storeToRefs } from 'pinia';
import { useWebsocketStore } from '@/stores';
import { watch } from 'vue';

export interface Options {
  onMessage?: (message: UniNamespace.OnSocketMessageCallbackResult) => void;
}

export function useWebsocket({ onMessage }: Options = {}) {
  const store = useWebsocketStore();

  const { websocket, readyState, latestMessage } = storeToRefs(store);

  if (onMessage) {
    watch(latestMessage, message => {
      if (message) {
        onMessage(message);
      }
    });
  }

  return { websocket, readyState, latestMessage, connect: store.connect, send: store.send };
}