import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

export const TOKEN_KEY = 'token';

export const useTokenStore = defineStore('token', () => {
  const token = ref<string>(uni.getStorageSync(TOKEN_KEY));

  watch(token, newToken => {
    uni.setStorageSync(TOKEN_KEY, newToken);
  });

  return { token };
});
