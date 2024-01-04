import { storeToRefs } from 'pinia';
import { useTokenStore } from '@/stores';

export function useToken() {
  const { token } = storeToRefs(useTokenStore());

  return token;
}