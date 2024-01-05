<template>
  <view class="page" data-weui-theme="light">
    <view class="page__bd">
      <view class="weui-panel weui-panel_access">
        <view class="weui-panel__bd">
          <navigator
              v-for="notification in notifications" :key="notification.id"
              class="weui-media-box weui-media-box_appmsg"
              :url="notification.url"
          >
            <view class="weui-media-box__bd">
              <view class="weui-media-box__title">
                <viewer :node="notification.subjectNode" />
              </view>
              <view class="content">
                <viewer :node="notification.messageNode" />
              </view>
              <view class="weui-media-box__info">
                <text class="weui-media-box__info__meta">{{ prettyTime(notification.createdAt) }}</text>
              </view>
            </view>
          </navigator>
        </view>
      </view>

      <loadmore :status="loadMoreStatus" />
    </view>
  </view>
</template>

<script lang="ts" setup>
import { useLockFn, useLoginSwitch } from '@/hooks';
import { computed, ref } from 'vue';
import { getUserNotifications } from '@/services';
import { onReachBottom } from '@dcloudio/uni-app';
import type { Nodes } from 'hast';
import { commentStrategy, toHast } from '@/markdown';
import { prettyTime } from '@/utils';
import Loadmore from '@/components/loadmore/loadmore.vue';

interface Notification extends API.Notification {
  subjectNode: Nodes;
  messageNode: Nodes;
  url?: string;
}

const loading = ref(false);
const notifications = ref<Notification[]>([]);
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

const queryUserNotifications = useLockFn(async (params?: Record<string, any>) => {
  loading.value = true;
  try {
    const { data, ...rest } = await getUserNotifications({
      ...params,
      pageSize,
      page,
      include: 'FromUser',
    });

    total = rest.total ?? 0;

    if (page === 1) {
      notifications.value = [];
    }

    notifications.value.push(...data.map(item => {
      const notification = {
        ...item,
        subjectNode: toHast(item.subject!, commentStrategy),
        messageNode: toHast(item.message!, commentStrategy),
      } as Notification;

      switch (notification.type) {
        case 'ArticleComment':
          notification.url = `/pages/article/comments?commentId=${notification.data?.id}`;
          break;
        case 'ArticleHasNewReply':
        case 'CommentHasNewReply':
        case 'CommentReply':
        case 'CommentUpvote':
          notification.url = `/pages/article/comments?commentId=${notification.data?.comment_id}&pinnedId=${notification.data?.id}`;
          break;
        case 'ArticleLike':
          notification.url = `/pages/article/details?id=${notification.data?.article_id}`;
          break;
      }

      return notification;
    }));
  } finally {
    loading.value = false;
  }
});

function loadMore() {
  if (loadMoreStatus.value === 'noMore') {
    return;
  }

  page += 1;
  queryUserNotifications({});
}

useLoginSwitch(() => queryUserNotifications({}));
onReachBottom(loadMore);
</script>

<style lang="scss" scoped>
.weui-media-box__title {
  font-weight: 500;

  :deep(.markdown-body navigator) {
    color: var(--weui-FG-0);
    text-decoration: underline;
  }
}

.content {
  :deep {
    .markdown-body {
      color: #333;

      navigator {
        color: var(--weui-FG-0);
        text-decoration: underline;
      }
    }

    .reply-quote {
      padding: 4px 12px;
      background-color: #f2f2f2;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      text-overflow: ellipsis;
      border-radius: 6px;
      margin: -8px 0 0;
      line-height: 24px;
    }
  }
}
</style>