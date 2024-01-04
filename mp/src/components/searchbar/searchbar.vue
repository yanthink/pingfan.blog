<template>
  <view class="weui-search-bar" :class="{'weui-search-bar_focusing': focus || modelValue}">
    <form id="searchForm" class="weui-search-bar__form">
      <view class="weui-search-bar__box">
        <text class="weui-icon-search"></text>
        <input
            type="search" class="weui-search-bar__input" placeholder="搜索"
            :value="modelValue"
            :focus="focus"
            @input="$emit('update:modelValue', $event.detail.value)"
            @focus="focus = true"
            @blur="focus = false"
            @confirm="$emit('search', $event.detail.value)"
        />
        <view class="weui-icon-clear"></view>
      </view>
      <label for="searchInput" class="weui-search-bar__label" @click="focus = true">
        <text class="weui-icon-search" />
        <text>搜索</text>
      </label>
    </form>
    <view class="weui-search-bar__cancel-btn">取消</view>
  </view>
</template>

<script lang="ts" setup>
import { ref } from 'vue';

const props = defineProps<{ modelValue?: string }>();
defineEmits<{
  (e: 'update:modelValue', v: string): void;
  (e: 'search', value: string): void;
}>();

const focus = ref(false);
</script>

<style lang="scss">
.weui-search-bar {
  background-color: #ffffff;
}

.weui-search-bar__box, .weui-search-bar__label {
  background-color: #F8F8F8;
}
</style>