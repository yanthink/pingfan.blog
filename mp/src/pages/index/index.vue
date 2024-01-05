<template>
  <view class="page" data-weui-theme="light">
    <view class="page__bd">
      <view class="github-url">https://github.com/yanthink/pingfan.blog</view>

      <searchbar @search="navigateTo({ url: `/pages/search/index?q=${$event}` })" />

      <view class="tabs">
        <tabs ref="tabsRef" :items="items" v-model:current="current" disable-scroll />
      </view>

      <view class="weui-panel weui-panel_access">
        <view class="weui-panel__bd">
          <navigator
              v-for="article in articles" :key="article.id"
              class="weui-media-box weui-media-box_appmsg"
              :url="`/pages/article/details?id=${article.id}`"
          >
            <view class="weui-media-box__bd">
              <text class="weui-media-box__title title">{{ article.title }}</text>
              <text class="weui-media-box__desc">{{ article.textContent }}</text>
              <view class="weui-media-box__info">
                <text class="weui-media-box__info__meta">{{ prettyTime(article.createdAt) }}</text>
                <text class="weui-media-box__info__meta">{{ prettyNumber(article.viewCount ?? 0) }} 阅读</text>
                <text class="weui-media-box__info__meta">{{ prettyNumber(article.commentCount ?? 0) }} 评论</text>
              </view>
            </view>
            <view v-if="article.preview" class="weui-media-box__hd">
              <image class="weui-media-box__thumb" mode="aspectFit" :src="article.preview" />
            </view>
          </navigator>
        </view>
      </view>

      <loadmore :status="loadMoreStatus" />
    </view>
  </view>
</template>

<script setup lang="ts">
import Searchbar from '@/components/searchbar/searchbar.vue';
import Tabs from '@/components/tabs/tabs.vue';
import Loadmore from '@/components/loadmore/loadmore.vue';
import { prettyTime, prettyNumber } from '@/utils';
import { computed, ref, watch } from 'vue';
import { useLockFn, useNotification } from '@/hooks';
import { getArticles } from '@/services';
import { onPullDownRefresh, onReachBottom, onShow } from '@dcloudio/uni-app';

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

const items = ['热门文章', '最新文章', '点赞最多', '评论最多'];

const current = ref(0);
const loading = ref(false);
const articles = ref<API.Article[]>([]);
const nextCursor = ref<string>();

const loadMoreStatus = computed(() => {
  if (loading.value) return 'loading';

  return nextCursor.value ? 'more' : 'noMore';
});

const queryArticles = useLockFn(async (params?: Record<string, any>) => {
  loading.value = true;
  try {
    const { data, cursor } = await getArticles({
      ...params,
      order: ['hot', 'latest', 'like', 'comment'][current.value],
    });

    if (!params?.cursor) {
      articles.value = [];
    }

    articles.value.push(...data.map(item => ({
      ...item,
      textContent: item.textContent?.trim().substring(0, 100),
    })));
    nextCursor.value = cursor?.after;
  } finally {
    loading.value = false;
  }
});

async function refresh() {
  await queryArticles();
  uni.stopPullDownRefresh();
}

function loadMore() {
  if (!nextCursor.value) {
    return;
  }

  queryArticles({ cursor: nextCursor.value });
}

watch(current, () => {
  articles.value = [];
  refresh();
}, { immediate: true });

onPullDownRefresh(refresh);
onReachBottom(loadMore);
</script>

<style lang="scss" scoped>
.github-url {
  width: 100vw;
  text-align: center;
  color: var(--weui-FG-2);
  position: fixed;
  top: 16px;
}

.tabs {
  position: sticky;
  top: -1px;
  background-color: #fff;
  border-bottom: 1px solid #f0f0f0;
  z-index: 9;
}

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
