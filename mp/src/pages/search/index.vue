<template>
  <view class="page" data-weui-theme="light">
    <view class="page__bd">
      <searchbar v-model="queryParams.q" @search="refresh" />

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
              <view class="weui-media-box__title title">
                <viewer :node="article.titleNode" />
              </view>
              <view class="weui-media-box__desc">
                <viewer :node="article.contentNode" />
              </view>
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
import { useLockFn } from '@/hooks';
import { searchArticles } from '@/services';
import { onLoad, onReachBottom } from '@dcloudio/uni-app';
import type { Nodes } from 'hast';
import { toHast, articleSearchStrategy } from '@/markdown';

interface Article extends API.Article {
  titleNode: Nodes;
  contentNode: Nodes;
}

let queryParams: Record<string, any> = {};

const items = ['热门文章', '最新文章', '点赞最多', '评论最多'];

const current = ref(0);
const loading = ref(false);
const articles = ref<Article[]>([]);
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

const queryArticles = useLockFn(async (params?: Record<string, any>) => {
  loading.value = true;
  try {
    const { data, ...rest } = await searchArticles({
      ...queryParams,
      ...params,
      sortFields: ['-hotness', '-created_at', '-like_count', '-comment_count'][current.value],
    });

    total = rest.total ?? 0;

    if (page === 1) {
      articles.value = [];
    }

    articles.value.push(...data.map(item => ({
      ...item,
      titleNode: toHast(item.highlights?.title?.[0] ?? item.title, articleSearchStrategy),
      contentNode: toHast(item.highlights?.content?.[0] ?? item.textContent?.substring(0, 300), articleSearchStrategy),
    } as Article)));
  } finally {
    loading.value = false;
  }
});

async function refresh() {
  await queryArticles();
  uni.stopPullDownRefresh();
}

function loadMore() {
  if (loadMoreStatus.value === 'noMore') {
    return;
  }

  page += 1;
  queryArticles({});
}

watch(current, () => {
  articles.value = [];
  refresh();
});

onLoad(query => {
  queryParams = { ...query };
  refresh();
})

onReachBottom(loadMore);
</script>

<style lang="scss" scoped>
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
