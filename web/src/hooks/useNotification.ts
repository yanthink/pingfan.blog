import { useSnapshot } from 'valtio';
import { notificationStore, fetchNotificationUnreadCount, setUnreadCount } from '@/stores';

function refresh() {
  return fetchNotificationUnreadCount();
}

export function useNotification(): { unreadCount: number; setUnreadCount: typeof setUnreadCount; refresh: typeof refresh } {
  const snapshot = useSnapshot(notificationStore);

  return {
    unreadCount: snapshot.unreadCount,
    setUnreadCount,
    refresh,
  };
}