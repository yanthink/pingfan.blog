<template>
  <view class="page" data-weui-theme="light">
    <view class="page__bd">
      <view v-if="comment.id" class="topic">
        <comment :comment="comment" hide-reply @upvote="upvote">
          <template #suffix>
            <view
                class="quote"
                style="margin-top: -8px"
                @click="navigateTo({ url: `details?id=${comment.articleId}` })"
            >
              {{ comment.article?.title }}
            </view>
          </template>
        </comment>
      </view>
      <view v-if="pinnedComment.id" class="pinned" @click="handleCommentClick(comment, -1)">
        <comment :comment="pinnedComment" hide-reply-count @upvote="upvote" @reply="handleCommentClick(comment, -1)" />
      </view>
      <view class="comments">
        <view
            v-for="(comment, index) in comments"
            :key="comment.rowKey || comment.id"
            :id="`comment-${comment.rowKey || comment.id}`"
            :class="{ 'new-comment': comment.id === newCommentId }"
            @click="handleCommentClick(comment, index)"
        >
          <comment :comment="comment" hide-reply-count @upvote="upvote" @reply="handleCommentClick(comment, index)">
            <template #prefix>
              <view
                  v-if="!!comment.replyQuoteContent"
                  class="quote"
                  style="margin-top: 8px; margin-bottom: 8px"
              >
                {{ comment.replyQuoteContent }}
              </view>
            </template>
          </comment>
        </view>
      </view>
    </view>

    <loadmore :status="loadMoreStatus" />
    <view style="height: 40px" />
    <safearea />

    <view class="comment-form">
      <comment-editor
          ref="editorRef"
          :comment-count="comment.replyCount"
          :placeholder="editorPlaceholder"
          hide-favorite
          hide-like
          @blur="handleBlur"
          @submit="submitComment"
      />
      <safearea />
    </view>
  </view>
</template>

<script lang="ts" setup>
import { onLoad, onReachBottom } from '@dcloudio/uni-app';
import { computed, ref } from "vue";
import { useLockFn, useAuthFn, useNextTick, useUser } from '@/hooks';
import { getComment, getComments, upvoteComment, createComment } from "@/services";
import Comment from '@/components/comment/comment.vue';
import Loadmore from '@/components/loadmore/loadmore.vue';
import CommentEditor from "@/components/comment-editor/comment-editor.vue";
import type {  CommentEditorInstance } from "@/components/comment-editor/comment-editor.vue";
import Safearea from "@/components/safearea/safearea.vue";
import { getProcessor, stripTagsStrategy } from "@/markdown";

let queryParams: Record<string, any> = {};

const nextTick = useNextTick();
const { user } = useUser();

const comment = ref<API.Comment>({});
const pinnedComment = ref<API.Comment>({});
const loading = ref(false);
const comments = ref<(API.Comment & { rowKey?: string; replyQuoteContent?: string })[]>([]);
const nextCursor = ref<string>();
const editorPlaceholder = ref<string>();
const editorRef = ref<CommentEditorInstance>();
const newCommentId = ref(0);

let replyCommentId = 0;
let replayCommentIndex = -1;

const loadMoreStatus = computed(() => {
  if (loading.value) return 'loading';

  return nextCursor.value ? 'more' : 'noMore';
});

const queryComments = useLockFn(async (params?: Record<string, any>) => {
  loading.value = true;
  try {
    const { data, cursor } = await getComments({
      ...queryParams,
      ...params,
      include: 'User,Parent.User',
    });

    if (!params?.cursor) {
      comments.value = [];
    }
    comments.value.push(...data.map(comment => {
      if (comment.parent && String(comment.parentId) !== queryParams.commentId) {
        const content = getProcessor(stripTagsStrategy)
            .processSync(
                comment.parent.content!.replace(/^```.*?\n$/g, '')
                    .split('\n')
                    .slice(0, 5)
                    .join('\n')
            )
            .toString()
            .trim()
            .substring(0, 100)
            .replaceAll('\n', ' ')
            .trim()

        return {
          ...comment,
          replyQuoteContent: `${comment.parent.user?.name}： ${content}`
        }
      }

      return comment;
    }));
    nextCursor.value = cursor?.after;
  } finally {
    loading.value = false;
  }
});

function handleCommentClick(comment: API.Comment, index: number) {
  replyCommentId = comment.id!;
  replayCommentIndex = index;

  editorPlaceholder.value = `回复 ${comment.user?.name}：`;
  editorRef.value?.setContent('');
  editorRef.value?.setFocus(true);
}

let blurTimer: ReturnType<typeof setTimeout>

function handleBlur() {
  clearTimeout(blurTimer);

  blurTimer = setTimeout(() => {
    if (editorRef.value?.getFocus()) {
      return;
    }
    replyCommentId = 0;
    replayCommentIndex = -1;
    editorPlaceholder.value = '';
    editorRef.value?.setContent('');
  }, 600);
}

onLoad(query => {
  queryParams = { ...query };

  if (!queryParams.commentId) {
    return;
  }

  getComment(queryParams.commentId, { include: 'User,Article' }).then(({ data }) => {
    comment.value = data;

    if (queryParams.reply) {
      handleCommentClick(comment.value, -1);
    }
  });
  if (queryParams.pinnedId) {
    getComment(queryParams.pinnedId, { include: 'User' }).then(({ data }) => pinnedComment.value = data);
  }
  queryComments({});
});

const upvote = useAuthFn(async (comment: API.Comment) => {
  const { data } = await upvoteComment(comment.id!);
  Object.assign(comment, { upvoteCount: data.upvoteCount, hasUpvoted: data.hasUpvoted });
});

const submitComment = useAuthFn(async (content: string) => {
  const index = replayCommentIndex + 1;
  const parentId = replyCommentId || comment.value.id;

  if (!comment.value.id) {
    return;
  }

  const { data } = await createComment({
    articleId: comment.value.articleId,
    content,
    parentId,
  });

  const newComment = {
    ...data,
    rowKey: `new_${data.id}`,
    user: user.value,
  };

  comments.value.splice(index, 0, newComment)
  newCommentId.value = newComment.id!;
  editorRef.value?.setContent('');
  comment.value.replyCount! += 1;

  await nextTick();
  await uni.pageScrollTo({ selector: `#comment-${newComment.rowKey}` });
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
  background: #fff;
}

.pinned, .new-comment {
  animation: .6s linear pinned forwards;
}

.quote {
  padding: 6px 12px;
  background: #f2f2f2;
  border-radius: 4px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 16px;
  font-size: 14px;
}

.comments {
  border-top: 1px solid #f0f0f0;
}

@keyframes pinned {
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