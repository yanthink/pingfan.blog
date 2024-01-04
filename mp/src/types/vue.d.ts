import request from '@/request';

declare module '@vue/runtime-core' {
  export interface ComponentCustomProperties {
    $http: typeof request;
    appName: string;
    imgBaseUrl: string;

    goBack(): void;

    navigateTo: typeof Uni.navigateTo;
  }
}

interface ImportMetaEnv {
  VITE_BASE_URL: string;
  VITE_SOCKET_URL: string;
}