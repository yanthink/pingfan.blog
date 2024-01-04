import { useLockFn, useUserPromise } from "@/hooks";

export function useAuthFn<P extends any[] = any[], V extends any = any>(fn: (...args: P) => Promise<V>) {
  const userPromise = useUserPromise();
  const lockFn = useLockFn(fn);

  return async (...args: P) => {
    const user = await userPromise;
    if (!user.value.name?.length || !user.value.avatar) {
      const { confirm } = await uni.showModal({
        title: '设置用户名和头像',
        content: '请先设置用户名和头像。',
        confirmText: '前去设置',
      });

      if (confirm) {
        uni.navigateTo({ url: '/pages/settings/index' });
      }

      return Promise.reject();
    }

    return lockFn(...args);
  }
}