import { defineStore } from 'pinia';
import { ref, watch } from 'vue';
import { getUserUnreadNotificationCount } from '@/services';
import { useLoginSwitch } from '@/hooks';

export const useNotificationStore = defineStore('notification', () => {
  const notification = ref({ unreadCount: 0 });

  async function fetchNotificationUnreadCount() {
    const { data } = await getUserUnreadNotificationCount();
    notification.value.unreadCount = data;

    return data;
  }

  watch(() => notification.value.unreadCount, count => {
    const route = getCurrentPages().at(-1)?.route ?? '';
    if (!['pages/index/index', 'pages/user/index'].includes(route)) {
      return
    }

    if (count > 0) {
      uni.setTabBarBadge({
        index: 1,
        text: String(count),
      });
    } else {
      uni.removeTabBarBadge({ index: 1 });
    }
  });

  useLoginSwitch(fetchNotificationUnreadCount);

  return { notification, fetchNotificationUnreadCount };
});
