<template>
  <view class="page" data-weui-theme="light">
    <view class="page__bd">
      <view class="header">
        <text class="title">{{ article.title }}</text>
        <view class="weui-media-box__info meta">
          <text class="weui-media-box__info__meta">{{ prettyNumber(article.viewCount ?? 0) }} 阅读</text>
          <text class="weui-media-box__info__meta">{{ prettyNumber(article.commentCount ?? 0) }} 评论</text>
          <view class="weui-media-box__info__meta">{{ prettyTime(article.createdAt) }}</view>
        </view>
      </view>
      <view class="content" user-select>
        <viewer v-if="node.children.length > 0" :node="node" />
        <view v-else style="display: flex; justify-content: center; padding: 24px">
          <loading type="circle" :font-size="48" />
        </view>
        <!-- #ifdef MP-WEIXIN -->
        <view class="actions">
          <like size="large" :liked="article.hasLiked" @click="like">
            <template #suffix>
              <text>{{ article.likeCount }}</text>
            </template>
          </like>
          <favorite size="large" :favorited="article.hasFavorited" @click="favorite">
            <template #suffix>
              <text>{{ article.favoriteCount }}</text>
            </template>
          </favorite>
        </view>
        <!-- #endif -->
      </view>

      <!-- #ifndef MP-WEIXIN -->
      <view id="comments" class="comments">
        <view class="comments-header">评论 {{ article.commentCount }}</view>
        <view
            v-for="(comment, index) in comments" :key="comment.rowKey || comment.id"
            :class="{ 'new-comment': index === 0 && !!comment.rowKey }"
        >
          <comment
              :comment="comment"
              @upvote="upvote"
              @click="navigateTo({ url: `comments?commentId=${$event.id}` })"
              @reply="navigateTo({ url: `comments?commentId=${$event.id}&reply=1` })"
          />
        </view>
      </view>
      <!-- #endif -->
    </view>

    <!-- #ifndef MP-WEIXIN -->
    <loadmore :status="loadMoreStatus" />
    <view style="height: 40px" />
    <safearea />

    <view class="comment-form">
      <comment-editor
          ref="editorRef"
          :comment-count="article.commentCount"
          :favorited="article.hasFavorited"
          :liked="article.hasLiked"
          :favorite-count="article.favoriteCount"
          :like-count="article.likeCount"
          @message="scrollToComments"
          @like="like"
          @favorite="favorite"
          @submit="submitComment"
      />
      <safearea />
    </view>
    <!-- #endif -->
  </view>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { onLoad, onReachBottom } from '@dcloudio/uni-app';
import { getArticle, likeArticle, favoriteArticle, getComments, upvoteComment, createComment } from '@/services';
import { prettyNumber, prettyTime } from '@/utils';
import { useHast, useLockFn, useAuthFn, useUser } from '@/hooks';
import { articleStrategy, toHast } from '@/markdown';
import Viewer from '@/components/viewer/viewer.vue';
import Loading from '@/components/loading/loading.vue';
import Comment from '@/components/comment/comment.vue';
import Loadmore from '@/components/loadmore/loadmore.vue';
import CommentEditor from '@/components/comment-editor/comment-editor.vue';
import type { CommentEditorInstance } from '@/components/comment-editor/comment-editor.vue';
import Safearea from '@/components/safearea/safearea.vue';
import Like from '@/components/like/like.vue';
import Favorite from '@/components/favorite/favorite.vue';

let articleId = 0;

const { user } = useUser();

const article = ref<API.Article>({});
const { node, setChunkHast } = useHast();

const editorRef = ref<CommentEditorInstance>();
const commentsLoading = ref(false);
const comments = ref<(API.Comment & { rowKey?: string })[]>([]);
const nextCursor = ref<string>();

const loadMoreStatus = computed(() => {
  if (commentsLoading.value) return 'loading';

  return nextCursor.value ? 'more' : 'noMore';
});

const queryComments = useLockFn(async (params?: Record<string, any>) => {
  commentsLoading.value = true;
  try {
    const { data, cursor } = await getComments({
      commentId: 0,
      ...params,
      articleId,
      include: 'User',
    });

    if (!params?.cursor) {
      comments.value = [];
    }
    comments.value.push(...data);
    nextCursor.value = cursor?.after;
  } finally {
    commentsLoading.value = false;
  }
});

onLoad(async query => {
  articleId = parseInt(query?.id ?? 0);

  const { data } = await getArticle(articleId);
  uni.setNavigationBarTitle({ title: data.title! });

  article.value = {
    id: data.id,
    title: data.title,
    viewCount: data.viewCount,
    commentCount: data.commentCount,
    createdAt: data.createdAt,
    hasLiked: data.hasLiked,
    likeCount: data.likeCount,
    hasFavorited: data.hasFavorited,
    favoriteCount: data.favoriteCount,
  };

  const hast = toHast(data.content!, articleStrategy);
  await setChunkHast(hast);

  // #ifndef MP-WEIXIN
  queryComments();
  // #endif
});

async function scrollToComments() {
  return uni.pageScrollTo({ selector: `#comments` });
}

const like = useAuthFn(async () => {
  const { data } = await likeArticle(articleId);
  article.value = { ...article.value, likeCount: data.likeCount, hasLiked: data.hasLiked };
});

const favorite = useAuthFn(async () => {
  const { data } = await favoriteArticle(articleId);
  article.value = {
    ...article.value,
    favoriteCount: data.favoriteCount,
    hasFavorited: data.hasFavorited,
  };
});

const upvote = useAuthFn(async (comment: API.Comment) => {
  const { data } = await upvoteComment(comment.id!);
  Object.assign(comment, { upvoteCount: data.upvoteCount, hasUpvoted: data.hasUpvoted });
});

const submitComment = useAuthFn(async (content: string) => {
  const { data } = await createComment({ articleId, content });
  comments.value.unshift({
    ...data,
    rowKey: `new_${data.id}`,
    user: user.value,
  });
  article.value.commentCount! += 1;
  editorRef.value?.setContent('');

  await scrollToComments();

  editorRef.value?.setFocus(false);

  await new Promise(resolve => setTimeout(resolve, 15_000));
});

function loadMore() {
  if (!nextCursor.value) {
    return;
  }

  queryComments({ cursor: nextCursor.value });
}

onReachBottom(loadMore);
</script>

<style lang="scss" scoped>
.page__bd {
  padding: 16px;
  background: #fff;
}

.header {
  margin: 0 -16px;
  padding: 16px 16px 8px;
  border-bottom: 1px solid #f0f0f0;
}

.content {
  margin: 16px 0;
}

.title {
  font-size: 1.25em;
  font-weight: 600;

  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.meta {
  font-size: 14px;
  margin-top: 8px;
}

.actions {
  padding: 24px;
  display: flex;
  justify-content: center;
  grid-gap: 24px;
}

.comments {
  margin: 24px -16px 0;
}

.comments-header {
  padding: 0 16px 12px;
  font-size: 18px;
  font-weight: 500;
  border-bottom: 1px solid #f0f0f0;
}

.new-comment {
  animation: .6s linear newComment forwards;
}

@keyframes newComment {
  0% {
    background: #13C2C2;
  }

  100% {
    background: transparent;
  }
}

.comment-form {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  background: #fff;
}
</style>
