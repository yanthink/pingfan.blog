import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import { login as doLogin, getUser } from '@/services';
import { useToken } from '@/hooks';

export const useUserStore = defineStore('user', () => {
  const token = useToken();
  const user = ref<API.User>({});
  const isLogin = computed(() => !!token.value);

  const launchOptions = uni.getLaunchOptionsSync();

  async function fetchUser() {
    const { data } = await getUser({ ...launchOptions.query });
    user.value = data;

    return user;
  }

  async function login() {
    try {
      const { code } = await uni.login({});

      const { data } = await doLogin(code, launchOptions.query);

      token.value = data.token!;
      delete data.token;
      user.value = data;

      uni.$emit('login');
    } catch (e: any) {
      uni.showToast({
        icon: 'none',
        title: '登录授权失败',
        duration: 1500,
      });
    }
  }

  let retries = 0;

  function logout() {
    token.value = '';
    user.value = {};

    if (retries === 0) {
      retries++;
      login();
    }
  }

  isLogin.value ? fetchUser() : login();

  uni.$on('logout', logout);

  return {
    user,
    fetchUser,
    login,
    logout,
    isLogin,
  };
});