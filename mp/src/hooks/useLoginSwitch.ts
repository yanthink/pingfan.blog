import { useUser } from '@/hooks';
import { watch } from 'vue';

export function useLoginSwitch(inFn: () => void, outFn?: () => void) {
  const { isLogin } = useUser();
  watch(isLogin, logged => logged ? inFn() : outFn?.(), { immediate: true });
}