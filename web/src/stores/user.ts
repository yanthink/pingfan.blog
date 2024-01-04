'use client';

import { proxy } from 'valtio';
import { getAuthUser } from '@/services';

export const userStore = proxy<{ user: API.User, loading: boolean }>({
  user: {},
  loading: false,
});

export async function fetchAuthUser() {
  userStore.loading = true;

  try {
    const { data } = await getAuthUser();
    userStore.user = data;

    return data;
  } finally {
    userStore.loading = false;
  }
}

export function setUser(user: API.User) {
  userStore.user = user;
}
