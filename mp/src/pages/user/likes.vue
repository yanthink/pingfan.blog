<template>
  <view class="page" data-weui-theme="light">
    <view class="page__bd">
      <view class="tabs">
        <tabs ref="tabsRef" :items="items" v-model:current="current" disable-scroll />
      </view>

      <view class="weui-panel weui-panel_access">
        <view class="weui-panel__bd">
          <template v-if="current === 0">
            <navigator
                v-for="like in likes" :key="like.id"
                class="weui-media-box weui-media-box_appmsg"
                :url="`/pages/article/details?id=${like.articleId}`"
            >
              <view class="weui-media-box__bd">
                <text class="weui-media-box__title title">{{ like.article?.title }}</text>
                <text class="weui-media-box__desc">{{ like.article?.textContent }}</text>
                <view class="weui-media-box__info">
                  <text class="weui-media-box__info__meta">{{ prettyTime(like.article?.createdAt) }}</text>
                  <text class="weui-media-box__info__meta">{{ prettyNumber(like.article?.likeCount ?? 0) }} 点赞</text>
                  <text class="weui-media-box__info__meta">
                    {{ prettyNumber(like.article?.commentCount ?? 0) }} 评论
                  </text>
                </view>
              </view>
              <view v-if="like.article?.preview" class="weui-media-box__hd">
                <image class="weui-media-box__thumb" mode="aspectFit" :src="like.article?.preview" />
              </view>
            </navigator>
          </template>
          <template v-else>
            <navigator
                v-for="upvote in upvotes" :key="upvote.id"
                class="weui-media-box weui-media-box_appmsg"
                :url="`/pages/article/details?id=${upvote.comment?.articleId}`"
            >
              <view class="weui-media-box__bd">
                <view class="weui-media-box__title title">{{ upvote.comment?.article?.title }}</view>
                <view class="content">
                  <view class="comment-preview">
                    <text>{{ upvote.comment?.user?.name }}：{{ upvote.comment?.textContent }}</text>
                  </view>
                </view>
                <view class="weui-media-box__info">
                  <text class="weui-media-box__info__meta">{{ prettyTime(upvote.createdAt) }}</text>
                  <text class="weui-media-box__info__meta">
                    {{ prettyNumber(upvote.comment?.upvoteCount ?? 0) }} 点赞
                  </text>
                  <text class="weui-media-box__info__meta">
                    {{ prettyNumber(upvote.comment?.replyCount ?? 0) }} 回复
                  </text>
                </view>
              </view>
            </navigator>
          </template>
        </view>
      </view>
    </view>
  </view>
</template>

<script lang="ts" setup>
import { useLockFn, useLoginSwitch, useToken } from '@/hooks';
import { computed, ref, watch } from 'vue';
import { getUserLikes, getUserUpvotes } from '@/services';
import { prettyTime, prettyNumber } from '@/utils';
import { getProcessor, stripTagsStrategy } from '@/markdown';

const items = ['文章', '评论'];

const token = useToken();

const current = ref(0);
const loading = ref(false);
const likes = ref<API.Like[]>([]);
const upvotes = ref<API.Upvote[]>([]);
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

const queryLikes = useLockFn(async (params?: Record<string, any>) => {
  loading.value = true;
  try {
    const { data, ...rest } = await getUserLikes({
      ...params,
      pageSize,
      page,
      include: 'Article.User',
    });

    total = rest.total ?? 0;

    if (page === 1) {
      likes.value = [];
    }

    likes.value.push(...data);
  } finally {
    loading.value = false;
  }
});

const queryUpvotes = useLockFn(async (params?: Record<string, any>) => {
  loading.value = true;
  try {
    const { data, ...rest } = await getUserUpvotes({
      ...params,
      pageSize,
      page,
      include: 'Comment.User,Comment.Article',
    });

    total = rest.total ?? 0;

    if (page === 1) {
      upvotes.value = [];
    }

    upvotes.value.push(...data.map(item => ({
      ...item,
      comment: {
        ...item.comment,
        textContent: getProcessor(stripTagsStrategy)
            .processSync(
                item.comment?.content!.replace(/^```.*?\n$/g, '')
                    .split('\n')
                    .slice(0, 5)
                    .join('\n'),
            )
            .toString(),
      },
    })));
  } finally {
    loading.value = false;
  }
});

function loadMore() {
  if (loadMoreStatus.value === 'noMore') {
    return;
  }

  page += 1;
  current.value === 0 ? queryLikes({}) : queryUpvotes({});
}

useLoginSwitch(() => current.value === 0 ? queryLikes({}) : queryUpvotes({}));

watch(current, current => {
  page = 1;
  if (token) {
    likes.value = [];
    upvotes.value = [];
    total = 0;
    current === 0 ? queryLikes({}) : queryUpvotes({});
  }
});
</script>

<style lang="scss" scoped>
.tabs {
  background: #fff;
  position: sticky;
  top: -1px;
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

.comment-preview {
  padding: 4px 12px;
  background-color: #f2f2f2;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  border-radius: 6px;
  margin: 16px 0;
  line-height: 24px;
}
</style>