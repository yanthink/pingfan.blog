<template>
  <view class="page" data-weui-theme="light">
    <view class="page__bd">
      <view class="tabs">
        <tabs ref="tabsRef" :items="items" v-model:current="current" disable-scroll />
      </view>

      <view class="weui-panel weui-panel_access">
        <view class="weui-panel__bd">
          <navigator
              v-for="comment in comments" :key="comment.id"
              class="weui-media-box weui-media-box_appmsg"
              :url="`/pages/article/comments?commentId=${comment.commentId || comment.id}&pinnedId=${comment.id !== comment.commentId ? comment.id : 0}`"
          >
            <view class="weui-media-box__bd">
              <text class="weui-media-box__title title">{{ comment.article?.title }}</text>
              <view class="content">
                <viewer :node="comment.node" />
                <view v-if="comment.parent" class="reply-quote">
                  <text>{{ comment.parent?.user?.name }}：{{ comment.parent?.textContent }}</text>
                </view>
              </view>
              <view class="weui-media-box__info">
                <text class="weui-media-box__info__meta">{{ prettyTime(comment.createdAt) }}</text>
                <text class="weui-media-box__info__meta">{{ prettyNumber(comment.upvoteCount ?? 0) }} 点赞</text>
                <text class="weui-media-box__info__meta">{{ prettyNumber(comment.replyCount ?? 0) }} 回复</text>
              </view>
            </view>
          </navigator>
        </view>
      </view>
    </view>
  </view>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue';
import { useLockFn, useLoginSwitch, useToken } from '@/hooks';
import { getUserComments } from '@/services';
import { onReachBottom } from '@dcloudio/uni-app';
import type { Nodes } from 'hast';
import { toHast, getProcessor, commentStrategy, stripTagsStrategy } from '@/markdown';
import { prettyNumber, prettyTime } from '@/utils';

const items = ['评论', '回复'];

interface Comment extends Omit<API.Comment, 'parent'> {
  node?: Nodes;
  textContent?: string;
  parent?: Comment;
}

const token = useToken();

const current = ref(0);
const loading = ref(false);
const comments = ref<Comment[]>([]);
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

const queryComments = useLockFn(async (params?: Record<string, any>) => {
  loading.value = true;
  try {
    const { data, ...rest } = await getUserComments({
      ...params,
      pageSize,
      page,
      type: current.value,
      include: current.value === 1 ? 'Article,Parent.User' : 'Article',
    });

    total = rest.total ?? 0;

    if (page === 1) {
      comments.value = [];
    }

    comments.value.push(...data.map(item => {
      const comment = {
        ...item,
        node: toHast(item.content!, commentStrategy),
      };

      if (comment.parent?.content) {
        comment.parent = {
          ...comment.parent,
          textContent: getProcessor(stripTagsStrategy)
              .processSync(
                  comment.parent.content!.replace(/^```.*?\n$/g, '')
                      .split('\n')
                      .slice(0, 5)
                      .join('\n'),
              )
              .toString(),
        } as Comment;
      }

      return comment;
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
  queryComments({});
}

useLoginSwitch(() => queryComments({}));
onReachBottom(loadMore);

watch(current, () => {
  page = 1;
  if (token) {
    comments.value = [];
    queryComments({});
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

.title {
  font-weight: 500;
}

.content {
  margin-top: 8px;
  color: #333;
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
</style>