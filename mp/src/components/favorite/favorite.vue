<template>
  <view
      class="favorite"
      :class="{ 'favorite-large': size === 'large', 'favorite-small': size === 'small', favorited }"
      @click="$emit('click')"
  >
    <image v-if="!favorited" src="/static/images/favorite_81x81.png" />
    <image v-else src="/static/images/favorited_81x81.png" />
    <slot name="suffix" />
  </view>
</template>

<script lang="ts" setup>
interface FavoriteProps {
  favorited?: boolean;
  size?: 'default' | 'large' | 'small';
}

const props = defineProps<FavoriteProps>();

defineEmits(['click'])
</script>

<style lang="scss" scoped>
.favorite {
  image {
    width: 24px;
    height: 24px;
  }

  display: flex;
  align-items: flex-end;
  grid-gap: 4px;
}

.favorite-large {
  image {
    width: 40px;
    height: 40px;
  }
}

.favorite-small {
  image {
    width: 16px;
    height: 16px;
  }
}
</style>