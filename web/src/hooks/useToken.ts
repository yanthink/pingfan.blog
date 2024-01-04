import { useSnapshot } from 'valtio';
import { tokenStore, setToken } from '@/stores';

export function useToken(): [string, typeof setToken] {
  const snapshot = useSnapshot(tokenStore);
  
  return [snapshot.token, setToken];
}