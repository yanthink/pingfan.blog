<template>
  <view class="comment" @click="$emit('click', comment)">
    <view class="avatar" :class="{ 'avatar-default': defaultAvatar }">
      <image mode="aspectFit" :src="avatar" />
    </view>
    <view class="main">
      <view class="author-upvote">
        <text class="author">{{ comment.user?.name }}</text>
        <Like size="small" :liked="comment.hasUpvoted" @click="$emit('upvote', comment)">
          <template #suffix>
            <text class="upvote-count">{{ comment.upvoteCount }}</text>
          </template>
        </Like>
      </view>
      <view class="content">
        <slot name="prefix" />
        <viewer :node="node" />
        <slot name="suffix" />
      </view>
      <view class="meta">
        <view v-if="!hideReply" class="comment-replay-btn" @click.stop.prevent="$emit('reply', comment)">
          {{ (!hideReplyCount && !comment.commentId && comment.replyCount > 0) ? `${comment.replyCount} ` : '' }}回复
        </view>
        <text>{{ prettyTime(comment.createdAt) }}</text>
      </view>
    </view>
  </view>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue';
import type { Root } from 'hast';
import Viewer from "@/components/viewer/viewer.vue";
import { commentStrategy, toHast } from "@/markdown";
import { prettyTime } from '@/utils';
import Like from "@/components/like/like.vue";

interface CommentProps {
  comment: API.Comment;
  hideReply?: boolean;
  hideReplyCount?: boolean;
}

const props = defineProps<CommentProps>();
defineEmits<{
  (e: 'upvote', comment: API.Comment): void;
  (e: 'click', comment: API.Comment): void;
  (e: 'reply', comment: API.Comment): void;
}>();

const node = ref<Root>(toHast(props.comment.content!, commentStrategy));

const avatarLoadError = ref(false);
const defaultAvatar = computed(() => !props.comment.user?.avatar || avatarLoadError.value);
const avatar = computed(() => defaultAvatar.value ? '/static/images/user.png' : props.comment.user?.avatar);
</script>

<style lang="scss" scoped>
.comment {
  padding: 16px;
  display: flex;

  &:active {
    background: var(--weui-BG-COLOR-ACTIVE);
  }
}

.avatar {
  width: 32px;

  image {
    width: 32px;
    height: 32px;
    border-radius: 50%;
  }

  &-default {
    image {
      box-sizing: border-box;
      background: rgba(0, 0, 0, 0.25);
      padding: 6px;
    }
  }
}

.main {
  margin-left: 16px;
  flex: 1;
  overflow: hidden;
}

.author {
  font-size: 14px;
  font-weight: 500;
}

.author-upvote {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.upvote-count {
  font-size: 13px;
  line-height: 13px;
}

.content {
  margin: 4px 0;

  .markdown-body {
    font-size: 15px !important;
  }
}

.meta {
  display: flex;
  margin-top: 0;
  align-items: center;
  grid-gap: 12px;
  color: var(--weui-FG-2);
  font-size: 13px;
}

.comment-replay-btn {
  background-color: #f2f2f2;
  line-height: 1;
  padding: 6px 18px 6px 12px;
  border-radius: 15px;
  color: #333;
  position: relative;

  &:after {
    content: " ";
    width: 8px;
    height: 16px;
    mask-position: 0 0;
    mask-repeat: no-repeat;
    mask-size: 100%;
    background-color: currentColor;
    mask-image: url(data:image/svg+xml,%3Csvg%20width%3D%2212%22%20height%3D%2224%22%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%3E%3Cpath%20d%3D%22M2.454%206.58l1.06-1.06%205.78%205.779a.996.996%200%20010%201.413l-5.78%205.779-1.06-1.061%205.425-5.425-5.425-5.424z%22%20fill%3D%22%23B2B2B2%22%20fill-rule%3D%22evenodd%22%2F%3E%3C%2Fsvg%3E);
    position: absolute;
    top: 50%;
    right: 8px;
    margin-top: -8px;
  }
}
</style>