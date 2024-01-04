'use client';

import React, { useEffect } from 'react';
import { theme } from 'antd';
import { usePathname, useSearchParams } from 'next/navigation';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

const NProgressBar: React.FC = () => {
  const { token } = theme.useToken();

  const pathname = usePathname();
  const searchParams = useSearchParams();

  useEffect(() => {
    NProgress.done();
  }, [pathname, searchParams]);

  useEffect(() => {
    const handleAnchorClick = (event: MouseEvent) => {
      const anchorElement = event.currentTarget as HTMLAnchorElement;

      if (event.metaKey || event.ctrlKey) return;

      const targetUrl = new URL(anchorElement.href);

      const currentUrl = new URL(location.href);

      if (
        targetUrl.host === currentUrl.host &&
        targetUrl.port === currentUrl.port &&
        targetUrl.pathname === currentUrl.pathname &&
        targetUrl.search === currentUrl.search
      ) {
        return;
      }

      NProgress.start();
    };

    const handleMutation: MutationCallback = () => {
      const anchorElements = document.querySelectorAll('a');
      Array.from(anchorElements).forEach(anchor => {
        if (anchor.href && anchor.target !== '_blank' && anchor.href[0] !== '#') {
          anchor.addEventListener('click', handleAnchorClick);
        }
      });
    };

    const mutationObserver = new MutationObserver(handleMutation);
    mutationObserver.observe(document, { childList: true, subtree: true });

    window.history.pushState = new Proxy(window.history.pushState, {
      apply: (target, thisArg, argArray: any) => {
        NProgress.done();
        return target.apply(thisArg, argArray);
      },
    });
  }, []);


  return (
    <>
      <style jsx>{`
        :global(#nprogress .bar) {
          background: ${token.colorPrimary};
        }

        :global(#nprogress .peg) {
          box-shadow: 0 0 10px ${token.colorPrimary}, 0 0 5px ${token.colorPrimary};
        }

        :global(#nprogress .spinner-icon) {
          border-top-color: ${token.colorPrimary};
          border-left-color: ${token.colorPrimary};
        }
      `}</style>
    </>
  );
};

export default NProgressBar;