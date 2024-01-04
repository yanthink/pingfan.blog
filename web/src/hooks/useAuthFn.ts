import { Modal } from 'antd';
import { useSearchParams } from 'next/navigation';
import { useLockFn } from 'ahooks';
import { useRouter, useUser } from '@/hooks';
import { check } from '@/utils';

const { confirm } = Modal;

export function useAuthFn<P extends any[] = any[], V extends any = any>(fn: (...args: P) => Promise<V>) {
  const searchParams = useSearchParams();
  const router = useRouter();

  const lockFn = useLockFn(fn);
  const { user } = useUser();

  return async (...args: P) => {
    if (!check()) {
      confirm({
        title: '登录确认?',
        content: '您还没有登录，点击【确定】前去登录。',
        okText: '确定',
        cancelText: '取消',
        onOk() {
          const redirect = searchParams.get('redirect') || window.location.href;
          router.replace(`/login?${new URLSearchParams({ redirect })}`);
        },
      });

      return Promise.reject();
    }

    if (!user.name?.length || !user.avatar) {
      confirm({
        title: '设置用户名和头像',
        content: '请先设置用户名和头像。',
        okText: '前去设置',
        cancelText: '取消',
        onOk() {
          router.push('/settings');
        },
      });

      return Promise.reject();
    }

    return lockFn(...args);
  };
}