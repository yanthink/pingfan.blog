<template>
  <view class="page" data-weui-theme="light">
    <view class="page__bd">
      <view class="profile">
        <view class="avatar" :class="{ 'avatar-default': defaultAvatar }">
          <image :src="avatar" @error="avatarLoadError = true" />
        </view>
        <view class="username">{{ user.name || '暂未设置' }}</view>
      </view>

      <view class="weui-cells">
        <navigator class="weui-cell weui-cell_access" url="favorites">
          <text class="weui-cell__bd">我的收藏</text>
          <text class="weui-cell__ft"></text>
        </navigator>
        <!-- #ifndef MP-WEIXIN -->
        <navigator class="weui-cell weui-cell_access" url="comments">
          <text class="weui-cell__bd">我的评论</text>
          <text class="weui-cell__ft"></text>
        </navigator>
        <!-- #endif -->
        <navigator class="weui-cell weui-cell_access" url="likes">
          <text class="weui-cell__bd">我的点赞</text>
          <text class="weui-cell__ft"></text>
        </navigator>
      </view>

      <view class="weui-cells">
        <navigator class="weui-cell weui-cell_access" url="/pages/notifications/index">
          <view class="weui-cell__bd">
            <text>消息通知</text>
          </view>
          <view class="weui-cell__ft" style="display: flex">
            <text v-if="notification.unreadCount > 0" class="weui-badge"
            >{{ notification.unreadCount > 99 ? '99+' : notification.unreadCount }}
            </text>
          </view>
        </navigator>
      </view>

      <view class="weui-cells">
        <navigator class="weui-cell weui-cell_access" url="/pages/settings/index">
          <text class="weui-cell__bd">个人设置</text>
          <text class="weui-cell__ft"></text>
        </navigator>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { useNotification, useUser } from '@/hooks';
import { ref, computed } from 'vue';
import { onShow } from '@dcloudio/uni-app';

const { user } = useUser();
const { notification } = useNotification();

onShow(() => {
  if (notification.value.unreadCount > 0) {
    uni.setTabBarBadge({
      index: 1,
      text: String(notification.value.unreadCount),
    });
  } else {
    uni.removeTabBarBadge({ index: 1 });
  }
});

const avatarLoadError = ref(false);
const defaultAvatar = computed(() => !user.value.avatar || avatarLoadError.value);
const avatar = computed(() => defaultAvatar.value ? '/static/images/user.png' : user.value.avatar);
</script>

<style lang="scss" scoped>
.page__bd {
  display: flex;
  flex-direction: column;

  :deep(.uni-badge--x) {
    display: flex;
  }
}

.profile {
  background: #fff;
  position: relative;
  height: 220px;
  display: flex;
  flex-direction: column;
  align-items: center;
  border-bottom: 1px solid #f0f0f0;
  text-align: center;
}

.avatar {
  display: flex;
  width: 120px;
  height: 120px;
  margin: 35px 10px 16px 10px;
  border-radius: 50%;
  border: 2px solid #fff;
  box-shadow: 3px 3px 10px rgba(0, 0, 0, 0.2);
  overflow: hidden;
  box-sizing: border-box;
  padding: 0;

  image {
    width: 100%;
    height: 100%;
  }

  &-default {
    background: rgba(0, 0, 0, 0.25);
    padding: 24px;

    image {
      width: 72px;
      height: 72px;
    }
  }
}

.username {
  font-size: 18px;
  font-weight: 500;
}
</style>
