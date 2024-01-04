'use client';

import { tokenStore, setToken, setUser, setUnreadCount } from '@/stores';

export function check(): boolean {
  return !!tokenStore.token;
}

export async function logout() {
  setToken('');
  setUser({});
  setUnreadCount(0);
}