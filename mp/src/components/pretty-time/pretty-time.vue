<template>
  <text>{{ formatted }}</text>
</template>

<script lang="ts" setup>
import type { ConfigType } from 'dayjs';
import dayjs from 'dayjs';
import { computed } from 'vue';

const props = defineProps<{ time: ConfigType }>();

const formatted = computed(() => {
  const now = dayjs();
  const current = dayjs(props.time);

  if (now.diff(current, 'd') <= 15) {
    return current.fromNow();
  }

  if (now.year() === current.year()) {
    return current.format('MM-DD');
  }

  return current.format('YYYY-MM-DD');
});
</script>