import { logout } from '@/utils/auth';
import { useCallback } from 'react';
import { useRouter } from '@/hooks';

export function useLogout() {
  const router = useRouter();

  return useCallback(() => {
    logout();
    router.replace(`/login?${new URLSearchParams({ redirect: window.location.pathname + window.location.search })}`);
  }, [router]);
}