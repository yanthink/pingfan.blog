<template>
  <view class="page" data-weui-theme="light">
    <view class="weui-form weui-bottom-fixed-opr-page">
      <view class="weui-form__bd weui-bottom-fixed-opr-page__content">
        <view class="weui-form__control-area">
          <text class="weui-cells__title">基本设置</text>
          <view class="weui-cells">
            <view
                class="weui-cell weui-cell_access"
                @click="chooseAvatar"
            >
              <text class="weui-cell__bd">头像</text>
              <view class="weui-cell__ft" style="height: 24px">
                <image
                    class="avatar" :class="{ 'avatar-default': !formData.avatar }"
                    :src="avatar"
                />
              </view>
            </view>

            <view class="weui-cell weui-cell_active" :class="{'weui-cell_warn': errors.name}">
              <view class="weui-cell__hd"><label class="weui-label">用户名</label></view>
              <view class="weui-cell__bd">
                <input class="weui-input" v-model="formData.name" placeholder="请输入用户名" />
              </view>
              <view class="weui-cell__ft">
                <view class="weui-btn_reset weui-btn_icon" @click="showNameHelp">
                  <text class="weui-icon-info-circle" />
                </view>
              </view>
            </view>

            <view class="weui-cell weui-cell_active" :class="{'weui-cell_warn': errors.email}">
              <view class="weui-cell__hd"><label class="weui-label">邮箱</label></view>
              <view class="weui-cell__bd">
                <input class="weui-input" type="email" v-model="formData.email" placeholder="请输入邮箱号码" />
              </view>
            </view>

            <view
                v-if="formData.email && formData.email !== user.email"
                class="weui-cell weui-cell_access weui-cell_vcode" :class="{'weui-cell_warn': errors.emailCode}"
            >
              <view class="weui-cell__hd"><label class="weui-label">验证码</label></view>
              <view class="weui-cell__bd weui-flex" style="align-items: center">
                <input
                    class="weui-input weui-cell__control weui-cell__control_flex"
                    v-model="formData.emailCode"
                    placeholder="输入验证码"
                    style="flex: 1"
                />
                <view class="weui-vcode-btn" @click="sendCode">{{ leftSeconds > 0 ? `${leftSeconds}秒后重试` : '获取验证码' }}</view>
              </view>
            </view>
          </view>

          <text class="weui-cells__title">邮件通知</text>
          <view class="weui-cells weui-cells_radio">
            <view class="weui-cell weui-cell_active" @click="formData.meta.emailNotify = 0">
              <text class="weui-cell__bd">关闭通知</text>
              <view class="weui-cell__ft">
                <view class="weui-check" :aria-checked="formData.meta.emailNotify === 0 ? 'true' : 'false'" />
                <text class="weui-icon-checked" />
              </view>
            </view>
            <view class="weui-cell weui-cell_active" @click="formData.meta.emailNotify = 1">
              <text class="weui-cell__bd">开启通知</text>
              <view class="weui-cell__ft">
                <view class="weui-check" :aria-checked="formData.meta.emailNotify === 1 ? 'true' : 'false'" />
                <text class="weui-icon-checked" />
              </view>
            </view>
            <view class="weui-cell weui-cell_active" @click="formData.meta.emailNotify = 2">
              <text class="weui-cell__bd">离线通知</text>
              <view class="weui-cell__ft">
                <view class="weui-check" :aria-checked="formData.meta.emailNotify === 2 ? 'true' : 'false'" />
                <text class="weui-icon-checked" />
              </view>
            </view>
          </view>
        </view>
      </view>
      <view class="weui-form__ft weui-bottom-fixed-opr-page__tool">
        <view class="weui-form__opr-area">
          <button class="weui-btn weui-btn_primary" @click="submit">确定</button>
        </view>
      </view>
    </view>

    <img-crop
        v-if="chooseAvatarUrl"
        :url="chooseAvatarUrl"
        :width="200"
        :height="200"
        @cancel="chooseAvatarUrl = ''"
        @success="updateAvatar"
    />
  </view>
</template>

<script lang="ts" setup>
import { ref, watch, computed } from 'vue';
import { useCountdown, useLockFn, useToken, useUser } from '@/hooks';
import ImgCrop from '@/components/img-crop/img-crop.vue';
import { sendEmailCaptcha, updateUserProfile } from '@/services';

const { user } = useUser();
const token = useToken();
const { leftSeconds, setLeftSeconds, start } = useCountdown();

const formData = ref<API.User & { emailCode?: string; emailCodeKey?: string }>({
  name: '',
  avatar: '',
  email: '',
  meta: {
    emailNotify: 0,
  },
});
const errors = ref<Record<string, boolean>>({});

const avatar = computed(() => formData.value.avatar || '/static/images/user.png');
const chooseAvatarUrl = ref('');

watch(user, user => formData.value = {
  ...formData.value, ...user,
  meta: {
    emailNotify: user.meta?.emailNotify ?? 0,
  },
}, { immediate: true });

function showNameHelp() {
  uni.showModal({
    title: '用户名格式',
    content: '2-10位字符，可包含中、英、数、_，不能以下划线开头和结尾。',
    showCancel: false,
  });
}

async function chooseAvatar() {
  // #ifdef MP-WEIXIN
  uni.chooseMedia({
    count: 1,
    mediaType: ['image'],
    success({ tempFiles }) {
      chooseAvatarUrl.value = tempFiles[0].tempFilePath;
    },
  });
  // #endif
  // #ifndef MP-WEIXIN
  const { tempFilePaths } = await uni.chooseImage({ count: 1 });
  chooseAvatarUrl.value = tempFilePaths[0];
  // #endif
}

async function updateAvatar(filePath: string) {
  chooseAvatarUrl.value = '';

  const { data } = await uni.uploadFile({
    url: `${import.meta.env.VITE_BASE_URL}/api/resources/upload`,
    header: {
      Authorization: `Bearer ${token.value}`,
    },
    formData: { type: 'avatar' },
    name: 'file',
    filePath,
  });

  const res = JSON.parse(data);

  formData.value.avatar = res.data.url;
}

const sendCode = useLockFn(async () => {
  if (leftSeconds.value > 0) {
    return;
  }

  try {
    const { data } = await sendEmailCaptcha({ email: formData.value.email, checkUnique: true });
    formData.value.emailCodeKey = data.key;
    setLeftSeconds(120);
    start();
  } catch (e: any) {
    if (e?.response?.data?.errors) {
      errors.value = Object.keys(e?.response?.data?.errors)
          .reduce((errors, fieldName) => Object.assign(errors, { [fieldName]: true }), {});
    }
  }
});

const submit = useLockFn(async () => {
  const wasChanged = ['avatar', 'name', 'email', 'meta'].some(key => {
    if (key === 'meta') {
      return JSON.stringify(formData.value[key]) !== JSON.stringify(user.value[key] ?? {});
    }

    return formData.value[key as keyof API.User] !== user.value[key as keyof API.User];
  });

  if (!wasChanged) {
    return;
  }

  try {
    const { data } = await updateUserProfile(formData.value);
    user.value = { ...user.value, ...data };
    uni.showToast({ title: '修改成功' });
  } catch (e: any) {
    if (e?.response?.data?.errors) {
      errors.value = Object.keys(e?.response?.data?.errors)
          .reduce((errors, fieldName) => Object.assign(errors, { [fieldName]: true }), {});
    }
  }
})
</script>

<style lang="scss" scoped>
.page {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  overflow-y: auto;
}

.weui-form__bd.weui-bottom-fixed-opr-page__content {
  padding-top: 0;
}

.weui-cells__title {
  //font-size: 16px;
  display: flex;
  margin: 8px 0;
}

.avatar {
  width: 24px;
  height: 24px;
}

.avatar-default {
  background: rgba(0, 0, 0, 0.25);
}

.weui-vcode-btn {
  width: 100px;
  text-align: center;
  &:active {
    background: #e5e5e5;
  }
}
</style>

