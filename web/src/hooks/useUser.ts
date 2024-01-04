import { useSnapshot } from 'valtio';
import { fetchAuthUser, userStore, setUser } from '@/stores';

function refresh() {
  return fetchAuthUser();
}

export function useUser(): { user: API.User; loading: boolean; setUser: typeof setUser; refresh: typeof refresh } {
  const snapshot = useSnapshot(userStore);

  return {
    user: snapshot.user,
    loading: snapshot.loading,
    setUser,
    refresh,
  };
}