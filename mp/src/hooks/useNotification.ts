import { storeToRefs } from 'pinia';
import { useNotificationStore } from '@/stores';

export function useNotification() {
  const store = useNotificationStore();

  const { notification } = storeToRefs(store);

  return { notification, refresh: store.fetchNotificationUnreadCount };
}