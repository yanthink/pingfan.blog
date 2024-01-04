import Request from './request';
import { useTokenStore } from '@/stores';

const request = new Request;

request.setConfig({
  baseUrl: import.meta.env.VITE_BASE_URL,
});

request.interceptor.request = async options => {
  options.header.Authorization = `Bearer ${useTokenStore().token}`;
};

request.interceptor.response = async (res, options) => {
  if (res.statusCode >= 200 && res.statusCode < 300 && (res.data.success ?? true)) {
    return { ...res.data, response: res };
  }

  if (typeof res.data === 'string') {
    try {
      res.data = JSON.parse(res.data);
    } catch (e) {
      //
    }
  }

  switch (res.data.code) {
    case 401_001:
      uni.$emit('logout');
      break;
    default:
      if (!options.skipErrorHandler) {
        uni.showToast({ icon: 'none', title: res.data.message || '未知错误' });
      }
      break;
  }

  return Promise.reject({ response: res, options });
};

export default request;
