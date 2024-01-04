<template>
  <view
      class="weui-loadmore"
      :class="{ 'weui-loadmore_line': status === 'noMore' }"
  >
    <template v-if="status === 'loading'">
      <view class="weui-primary-loading">
        <text class="weui-primary-loading__dot" />
      </view>
      <text class="weui-loadmore__tips">正在加载</text>
    </template>
    <text v-else-if="status === 'noMore'" class="weui-loadmore__tips">已无更多数据</text>
    <text v-else class="weui-hidden_abs">加载更多</text>
  </view>
</template>

<script lang="ts" setup>
interface LoadmoreProps {
  status?: 'loading' | 'more' | 'noMore';
}

defineProps<LoadmoreProps>();
</script>