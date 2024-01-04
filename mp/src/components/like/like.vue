<template>
  <view
      class="like"
      :class="{ 'like-large': size === 'large', 'like-small': size === 'small', liked }"
      @click="$emit('click')"
  >
    <image v-if="!liked" src="/static/images/like_81x81.png" />
    <image v-else src="/static/images/liked_81x81.png" />
    <slot name="suffix" />
  </view>
</template>

<script lang="ts" setup>
interface LikeProps {
  liked?: boolean;
  size?: 'default' | 'large' | 'small';
}

const props = defineProps<LikeProps>();

defineEmits(['click'])
</script>

<style lang="scss" scoped>
.like {
  image {
    width: 24px;
    height: 24px;
  }

  display: flex;
  align-items: flex-end;
  grid-gap: 4px;
}

.like-large {
  image {
    width: 40px;
    height: 40px;
  }
}

.like-small {
  image {
    width: 16px;
    height: 16px;
  }
}
</style>