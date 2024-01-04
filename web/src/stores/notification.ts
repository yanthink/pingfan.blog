import { proxy } from 'valtio';
import { getUserUnreadNotificationCount } from '@/services';
import { check } from '@/utils';

export const notificationStore = proxy<{ unreadCount: number }>({
  unreadCount: 0,
});

export async function fetchNotificationUnreadCount() {
  const { data } = await getUserUnreadNotificationCount();
  notificationStore.unreadCount = data;

  return data;
}

export function setUnreadCount(unreadCount: number) {
  notificationStore.unreadCount = unreadCount;
}

if (typeof window !== 'undefined' && check()) {
  fetchNotificationUnreadCount();
}