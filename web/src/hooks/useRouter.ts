import { useRouter as useNextRouter, usePathname } from 'next/navigation';
import NProgress from 'nprogress';

export function useRouter() {
  const router = useNextRouter();
  const pathname = usePathname();

  const push: typeof router['push'] = (href, options) => {
    if (href === `${location.pathname}${location.search}` || href === location.search) {
      return;
    }

    NProgress.start();

    return router.push(href, options);
  };

  const replace: typeof router['replace'] = (href, options) => {
    if (href === `${location.pathname}${location.search}` || href === location.search) {
      return;
    }

    NProgress.start();

    return router.replace(href, options);
  };

  const back: typeof router['back'] = () => {
    NProgress.start();

    return router.back();
  };

  return { ...router, push, replace, back };
}