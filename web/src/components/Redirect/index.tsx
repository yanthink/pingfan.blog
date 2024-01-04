'use client';

import React from 'react';
import { useRouter } from '@/hooks';
import { useMount } from 'ahooks';

interface RedirectProps {
  to?: string;
}

const Redirect: React.FC<RedirectProps> = ({ to }) => {
  const router = useRouter();

  useMount(() => {
    router.replace(to ?? `/login?${new URLSearchParams({ redirect: window.location.pathname + window.location.search })}`);
  });

  return null;
};

export default Redirect;