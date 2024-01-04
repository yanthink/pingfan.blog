import { getCurrentInstance } from 'vue';

export function useNextTick() {
  const currentInstance = getCurrentInstance();

  // vue3 api 小程序不支持 nextTick
  return async function nextTick(fn?: any) {
    return currentInstance?.proxy?.$nextTick.call(currentInstance?.proxy, fn);
  }
}