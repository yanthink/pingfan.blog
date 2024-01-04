import { proxy, subscribe } from 'valtio';
import store from 'store';

const TOKEN_KEY = process.env.NEXT_PUBLIC_TOKEN_KEY ?? 'token';

export const tokenStore = proxy<{ token: string }>({
  token: store.get(TOKEN_KEY, ''),
});

export function setToken(token = '') {
  tokenStore.token = token;
}

subscribe(tokenStore, () => {
  store.set(TOKEN_KEY, tokenStore.token);
  const expires = new Date(Date.now() + 10 * 365 * 24 * 60 * 60 * 1000).toUTCString();
  document.cookie = `${TOKEN_KEY}=${tokenStore.token}; expires=${expires}; path=/;`;
});