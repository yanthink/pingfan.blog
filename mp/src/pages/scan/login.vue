<template>
  <view class="page" data-weui-theme="light">
    <view class="page__hd">
      <view style="text-align: center">
        <image src="/static/images/desktop_200x200.png" />
        <view style="text-align: center">
          <view style="font-size: 16px">WEB端登录确认</view>
        </view>
      </view>

      <view style="text-align: center">
        <button
            style="width: 200px"
            class="weui-btn weui-btn_primary"
            @click="login"
        >
          登录
        </button>
        <view style="font-size: 16px; margin-top: 20px; color: rgba(0, 0, 0, 0.5)">
          <navigator open-type="exit" target="miniProgram">取消登录</navigator>
        </view>
      </view>
    </view>
  </view>
</template>

<script lang="ts" setup>
import { onLoad } from '@dcloudio/uni-app';
import { ref } from 'vue';
import { useLockFn, useUserPromise } from '@/hooks';
import { scanLogin } from '@/services';

const uuid = ref('');
const userPromise = useUserPromise();

onLoad(query => {
  uuid.value = decodeURIComponent(query?.scene ?? '');
});

const login = useLockFn(async () => {
  uni.showLoading();
  try {
    const user = await userPromise;
    await scanLogin(uuid.value);
    uni.showToast({ title: '登录成功' });
    if (!user.value.name) {
      uni.reLaunch({ url: '/pages/settings/index' });
    } else {
      uni.reLaunch({ url: '/pages/user/index' });
    }
  } finally {
    uni.hideLoading();
  }
});
</script>

<style lang="scss" scoped>
.page__hd {
  box-sizing: border-box;
  padding: 60px 0 100px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: space-between;

  image {
    width: 200px;
    height: 200px;
  }
}
</style>