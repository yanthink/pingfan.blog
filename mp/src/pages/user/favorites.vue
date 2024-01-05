<template>
  <view class="page" data-weui-theme="light">
    <view class="page__bd">
      <view class="weui-panel weui-panel_access">
        <view class="weui-panel__bd">
          <navigator
              v-for="favorite in favorites" :key="favorite.id"
              class="weui-media-box weui-media-box_appmsg"
              :url="`/pages/article/details?id=${favorite.article?.id}`"
          >
            <view class="weui-media-box__bd">
              <text class="weui-media-box__title title">{{ favorite.article?.title }}</text>
              <text class="weui-media-box__desc">{{ favorite.article?.textContent }}</text>
              <view class="weui-media-box__info">
                <text class="weui-media-box__info__meta">{{ prettyTime(favorite.article?.createdAt) }}</text>
                <text class="weui-media-box__info__meta">{{ prettyNumber(favorite.article?.viewCount ?? 0) }} 阅读</text>
                <text class="weui-media-box__info__meta">{{ prettyNumber(favorite.article?.commentCount ?? 0) }} 评论</text>
              </view>
            </view>
            <view v-if="favorite.article?.preview" class="weui-media-box__hd">
              <image class="weui-media-box__thumb" mode="aspectFit" :src="favorite.article?.preview" />
            </view>
          </navigator>
        </view>
      </view>

      <loadmore :status="loadMoreStatus" />
    </view>
  </view>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue';
import { useLockFn, useLoginSwitch } from '@/hooks';
import { getUserFavorites } from '@/services';
import { prettyTime, prettyNumber } from '@/utils';
import { onReachBottom } from '@dcloudio/uni-app';
import Loadmore from '@/components/loadmore/loadmore.vue';

const loading = ref(false);
const favorites = ref<API.Favorite[]>([]);
const pageSize = 10;
let page = 1;
let total = 0;

const loadMoreStatus = computed(() => {
  if (loading.value) return 'loading';
  if (page >= Math.ceil(total / pageSize)) {
    return 'noMore';
  }
  return 'more';
});

const queryFavorites = useLockFn(async (params?: Record<string, any>) => {
  loading.value = true;
  try {
    const { data, ...rest } = await getUserFavorites({
      ...params,
      pageSize,
      page,
      include: 'Article.User',
    });

    total = rest.total ?? 0;

    if (page === 1) {
      favorites.value = [];
    }

    favorites.value.push(...data);
  } finally {
    loading.value = false;
  }
});

function loadMore() {
  if (loadMoreStatus.value === 'noMore') {
    return;
  }

  page += 1;
  queryFavorites({});
}

useLoginSwitch(() => queryFavorites({}));
onReachBottom(loadMore);
</script>

<style lang="scss" scoped>
.weui-panel {
  margin-top: 0;
}

.weui-media-box_appmsg {
  align-items: flex-start;
}

.weui-media-box__desc {
  line-height: 24px;
}

.weui-media-box__info {
  margin-top: 8px;
  padding-bottom: 0;
}

.weui-media-box__hd {
  width: 100px;
  height: 80px;
  margin-left: 16px;
  margin-right: 0;
  background: #f5f5f5;

  image {
    width: 100%;
    height: 100%;
  }
}
</style>