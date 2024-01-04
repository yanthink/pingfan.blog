import { storeToRefs } from 'pinia';
import { useUserStore } from '@/stores';

export function useUser() {
  const store = useUserStore();
  const { user, isLogin } = storeToRefs(store);

  return { user, isLogin, refresh: store.fetchUser };
}