<template>
  <view
      class="wx_loading_view"
      id="wx_loading_view"
      :class="classnames([{ wx_loading_view__animated: animated, wx_loading_view__hide: !show, }, extClass])"
      :style="[styles]"
  >
    <view v-if="type === 'dot-white'" class="loading wx_dot_loading wx_dot_loading_white" :style="{ fontSize: `${fontSize ?? 14}px` }" />
    <view v-else-if="type==='circle'" class="weui-loadmore">
      <view class="weui-loading" :style="{ fontSize: `${fontSize ?? 14}px` }" />
      <view class="weui-loadmore__tips">{{ tips }}</view>
    </view>
    <view v-else class="loading wx_dot_loading" :style="{ fontSize: `${fontSize ?? 14}px` }"/>
  </view>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import classnames from 'classnames';

interface LoadingProps {
  extClass?: string;
  show?: boolean;
  animated?: boolean;
  duration?: number;
  type?: 'dot-gray' | 'dot-white' | 'circle';
  tips?: string;
  fontSize?: number;
}

const props = defineProps<LoadingProps>();

const styles = computed(() => {
  if (props.animated) {
    return {
      transition: `height ${props.duration ?? 350}ms ease`,
      fontSize: `${props.fontSize ?? 14}px`,
    };
  }

  return {
    fontSize: `${props.fontSize ?? 14}px`,
  };
});
</script>

