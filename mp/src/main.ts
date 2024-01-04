import { createSSRApp } from 'vue';
import * as Pinia from 'pinia';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';
import App from './App.vue';

dayjs.locale('zh-cn');
dayjs.extend(relativeTime);

export function createApp() {
  const app = createSSRApp(App);
  app.use(Pinia.createPinia());

  app.config.globalProperties.goBack = () => {
    const pages = getCurrentPages();

    if (pages.length === 1 && pages[0].route !== '/pages/index/index') {
      return uni.redirectTo({ url: '/pages/index/index' });
    }

    uni.navigateBack();
  };
  app.config.globalProperties.navigateTo = uni.navigateTo;

  return {
    app,
    Pinia,
  };
}
