<template>
  <scroll-view
      class="scroll"
      :show-scrollbar="false" scroll-x scroll-with-animation :scroll-left="scrollLeft"
      @scroll="scroll"
  >
    <view class="items">
      <view
          v-for="(item, index) in items"
          :key="index"
          :style="itemStyles(index)"
          :class="['item', `item-${index}`, { active: index === current }]"
          @click="current !== index && $emit('update:current', index)"
      >
        <slot :item="item" :index="index" :current="current">{{ item }}</slot>
      </view>
    </view>

    <view ref="bar" class="bar" :style="barStyles" />
  </scroll-view>
</template>

<script lang="ts" setup>
import type { PropType, CSSProperties } from 'vue';
import { computed, getCurrentInstance, onMounted, ref, watch } from 'vue';
import { useNextTick } from '@/hooks';
import { getRect } from '@/utils';

const props = defineProps({
  items: {
    type: Array as PropType<string[]>,
    default: [],
  },
  current: {
    type: Number,
    default: 0,
  },
  // 导航栏的高度和行高，单位rpx
  height: {
    type: Number,
    default: 80,
  },
  // 单个tab的左或右内边距（各占一半），单位rpx
  gutter: {
    type: Number,
    default: 40,
  },
  // 未选中字体颜色
  color: {
    type: String,
    default: 'rgba(0, 0, 0, 0.88)',
  },
  // 已选中字体颜色
  activeColor: {
    type: String,
    default: '#13C2C2',
  },
  // 底部Bar宽度，单位rpx
  barWidth: {
    type: Number,
    default: 40,
  },
  // 底部Bar高度，单位rpx
  barHeight: {
    type: Number,
    default: 4,
  },
  barStyles: {
    type: Object as PropType<CSSProperties>,
  },
  disableScroll: Boolean,
});

defineEmits(['update:current']);

const currentInstance = getCurrentInstance();
const nextTick = useNextTick();

type Rect = { left: number; right: number; width: number; height: number; };

let scrollRect: Rect = { left: 0, right: 0, width: 0, height: 0 };
let itemsRect = ref<Rect[]>([]);

const scrollLeft = ref(0);
const scrollBarLeftOffset = ref(0);
const preActiveIndex = ref(-1); // 预渲染index
const barExtendWidthPx = ref(0); // 底部Bar额外宽度

let touchScrollLeft = 0;
let itemsWidth = 0;

onMounted(async () => {
  touchScrollLeft = 0;
  await nextTick();
  init();
});

watch(() => props.current, () => {
  change();
});

watch(() => props.items, () => {
  nextTick(init);
}, { flush: 'post' });

async function init() {
  touchScrollLeft = 0;

  scrollRect = await getRect('.scroll', currentInstance) as Rect;
  itemsRect.value = await Promise.all(props.items.map((_, i) => getRect(`.item-${i}`, currentInstance) as Promise<Rect>));

  if (itemsRect.value[0] && itemsRect.value[0].left !== scrollRect.left) {
    const offset = scrollRect.left - itemsRect.value[0].left;
    itemsRect.value.forEach(item => {
      item.left += offset;
      item.right += offset;
    });
  }

  const lastItemRect = itemsRect.value[itemsRect.value.length - 1];
  itemsWidth = lastItemRect.right - itemsRect.value[0].left;

  setScrollViewToCenter();
}

const itemStyles = computed(() => {
  return (index: number) => {
    const style: CSSProperties = {
      height: `${props.height}rpx`,
      paddingLeft: `${props.gutter / 2}rpx`,
      paddingRight: `${props.gutter / 2}rpx`,
      color: props.color,
    };

    if (props.disableScroll) {
      style.flex = 1;
    }

    const activeIndex = preActiveIndex.value > 0 ? preActiveIndex.value : props.current;

    if (index === activeIndex) {
      style.color = props.activeColor;
    }

    return style;
  };
});

const barWidthPx = computed(() => uni.upx2px(props.barWidth) + barExtendWidthPx.value);

const scrollBarLeft = computed(() => {
  const currentItemRect = itemsRect.value[props.current];

  if (currentItemRect) {
    const center = currentItemRect.left + currentItemRect.width / 2;
    return center - barWidthPx.value / 2 - itemsRect.value[0].left + scrollBarLeftOffset.value;
  }

  return 0;
});

const barStyles = computed(() => Object.assign({
  width: `${barWidthPx.value}px`,
  height: `${props.barHeight}rpx`,
  borderRadius: `${props.barHeight / 2}rpx`,
  left: `${scrollBarLeft.value}px`,
  top: `${props.height - props.barHeight}rpx`,
  backgroundColor: props.activeColor,
}, props.barStyles));

function scroll(e: any) {
  touchScrollLeft = e.detail.scrollLeft;
}

// 把活动tab移动到屏幕中心点
async function setScrollViewToCenter() {
  if (props.disableScroll) return;

  const currentItemRect = itemsRect.value[props.current];
  if (!currentItemRect) return;

  const center = currentItemRect.left + currentItemRect.width / 2;
  let left = center - scrollRect.width / 2 - itemsRect.value[0].left;
  left = Math.min(Math.max(left, 0), itemsWidth - scrollRect.width);

  scrollLeft.value = touchScrollLeft;
  await nextTick();
  scrollLeft.value = left;
}

// 配合swiper组件@transition事件可以实现左右滑动平滑切换tab
function setDx(dx: number, width: number) {
  let nextItemIndex = dx > 0 ? props.current + 1 : props.current - 1;
  nextItemIndex = Math.min(Math.max(nextItemIndex, 0), props.items.length - 1);

  const nextItemRect = itemsRect.value[nextItemIndex];
  const nextItemX = nextItemRect.left + nextItemRect.width / 2;

  const currentItemRect = itemsRect.value[props.current];
  const currentItemX = currentItemRect.left + currentItemRect.width / 2;

  // 两个tab之间的距离，因为下一个tab可能在当前tab的左边或者右边，取绝对值即可
  const distanceX = Math.abs(nextItemX - currentItemX);
  // bar最大加长宽度
  const barMaxExtendWidth = 20;
  // 滚动比例
  const ratio = Math.abs(dx) / (width || scrollRect.width);
  // bar当前滚动长度
  const offset = ratio * distanceX;

  if (ratio < 0.5) {
    // bar加长滚动长度的3/4
    barExtendWidthPx.value = Math.min(offset * 0.75, barMaxExtendWidth);
  } else {
    barExtendWidthPx.value = Math.max(barMaxExtendWidth - (offset - distanceX / 2) * 0.75, 0);
  }

  preActiveIndex.value = ratio > 0.75 ? nextItemIndex : -1;

  let leftOffset = offset - barExtendWidthPx.value / 2;

  if (dx < 0) {
    leftOffset = -leftOffset;
  }

  scrollBarLeftOffset.value = leftOffset;
}

function change() {
  barExtendWidthPx.value = 0;
  scrollBarLeftOffset.value = 0;
  preActiveIndex.value = -1;
  setScrollViewToCenter();
}

defineExpose({ setDx });
</script>

<style lang="scss" scoped>
::-webkit-scrollbar, .scroll :deep ::-webkit-scrollbar {
  display: none;
  width: 0 !important;
  height: 0 !important;
  -webkit-appearance: none;
  background: transparent;
}

.scroll {
  flex-direction: row;
  flex: 1;
  position: relative;
}

.items {
  display: flex;
  flex-direction: row;
  flex-wrap: nowrap;


  .item {
    display: flex;
    justify-content: center;
    align-items: center;
    text-align: center;
    flex-wrap: nowrap;
    white-space: nowrap;
    color: #333;
    font-weight: 600;

    &-active {
      color: #FAAE3E;
    }
  }
}

.bar {
  position: absolute;
  background-color: #FAAE3E;
}
</style>