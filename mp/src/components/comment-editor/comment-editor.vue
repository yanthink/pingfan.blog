<template>
  <view class="comment-editor" :class="{ focus }">
    <view class="actions">
      <view class="trigger" @click="focus = true">{{ triggerPlaceholder || '写评论...' }}</view>
      <view style="display: flex; align-items: center; grid-gap: 24px">
        <view class="message" @click="$emit('message')">
          <text v-if="commentCount" class="weui-badge">{{ prettyNumber(commentCount) }}</text>
          <image src="/static/images/message_81x81.png" />
        </view>
        <view v-if="!hideFavorite" class="favorite-btn" @click="$emit('favorite')">
          <favorite :favorited="favorited">
            <template v-if="favoriteCount" #suffix>
              <text style="font-size: 13px; line-height: 1">{{ prettyNumber(favoriteCount) }}</text>
            </template>
          </favorite>
        </view>
        <view v-if="!hideLike" class="like-btn" @click="$emit('like')">
          <like :liked="liked">
            <template v-if="likeCount" #suffix>
              <text style="font-size: 13px; line-height: 1">{{ prettyNumber(likeCount) }}</text>
            </template>
          </like>
        </view>
      </view>
    </view>

    <view v-show="focus" class="comment-editor-mask" @click.stop.prevent="focus = false" />

    <view v-show="focus" class="comment-editor-form" :style="{ bottom: `${bottom}px` }">
      <view class="comment-editor-form__bd">
        <textarea
            :placeholder="placeholder || '写评论...'"
            v-model="content"
            :maxlength="512"
            :focus="focus"
            :show-confirm-bar="false"
            :adjust-position="false"
            disable-default-padding
            hold-keyboard
            fixed
            @focus="focus = true"
            @blur="focus = false"
            @keyboardheightchange="keyboardheightchange"
        />
        <view v-if="content.length" class="submit-btn" @click.prevent="$emit('submit', content)">发布</view>
      </view>
      <view class="comment-editor-form__ft">
        <view class="upload-image" @click.prevent="chooseImage">
          <image src="/static/images/image_81x81.png" />
        </view>
      </view>
    </view>
  </view>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import { useToken } from "@/hooks";
import Favorite from "@/components/favorite/favorite.vue";
import Like from "@/components/like/like.vue";
import { prettyNumber } from '@/utils';

interface CommentEditorProps {
  commentCount?: number;
  favorited?: boolean;
  favoriteCount?: number;
  liked?: boolean;
  likeCount?: number;
  hideFavorite?: boolean;
  hideLike?: boolean;
  triggerPlaceholder?: string;
  placeholder?: string;
}

defineProps<CommentEditorProps>();
const emit = defineEmits<{
  (e: 'message'): void;
  (e: 'like'): void;
  (e: 'favorite'): void;
  (e: 'blur'): void;
  (e: 'focus'): void;
  (e: 'submit', value: string): void;
}>();

const token = useToken();

const content = ref('');
const focus = ref(false);
const bottom = ref(0);

watch(focus, focus => focus ? emit('focus') : emit('blur'));

function keyboardheightchange(e: any) {
  bottom.value = e.detail.height;
}

async function chooseImage() {
  const oldContent = content.value;

  let filePath: string;

  try {
    // #ifdef MP-WEIXIN
    await new Promise((resolve, reject) => {
      uni.chooseMedia({
        count: 1,
        mediaType: ['image'],
        async success(res) {
          filePath = res.tempFiles[0].tempFilePath;
          resolve(filePath);
        },
        fail: reject,
      });
    })
    // #endif
    // #ifndef MP-WEIXIN
    const { tempFilePaths } = await uni.chooseImage({ count: 1 });
    filePath = tempFilePaths[0];
    // #endif

    focus.value = true;
    content.value = oldContent;

    const { data } = await uni.uploadFile({
      url: `${import.meta.env.VITE_BASE_URL}/api/resources/upload`,
      header: {
        Authorization: `Bearer ${token.value}`,
      },
      formData: { type: 'commentImage' },
      name: 'file',
      filePath,
    });

    const result = JSON.parse(data);
    if (content.value.split('\n').at(-1).trim() !== '') {
      content.value += '\n';
    }
    content.value += `![file](${result.data.url})`;
  } catch (e:any) {
    focus.value = true;
    content.value = oldContent;
  }
}

function setContent(value: string) {
  content.value = value;
}

function getFocus() {
  return focus.value;
}

function setFocus(value: boolean) {
  focus.value = value;
}

export interface CommentEditorInstance {
  setContent: typeof setContent;
  getFocus: typeof getFocus;
  setFocus: typeof setFocus;
}

defineExpose({ setContent, getFocus, setFocus });
</script>

<style lang="scss" scoped>
.comment-editor {
  padding: 8px 16px;
  background: #fff;
  border-top: 1px solid #f0f0f0;
}

.actions {
  width: 100%;
  height: 36px;
  display: flex;
  align-items: center;
  justify-items: flex-end;
  grid-gap: 16px;

  .trigger {
    flex: 1;
    padding: 8px;
    background: #f2f2f2;
    border-radius: 4px;
    color: #666666;
  }

  image {
    width: 24px;
    height: 24px;
  }

  .message {
    position: relative;
    height: 24px;

    .weui-badge {
      position: absolute;
      top: -10px;
      right: -10px;
    }
  }
}

.comment-editor-mask {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 998;
}

.comment-editor-form {
  padding: 8px 16px;
  background: #fff;
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 999;
  border-top: 1px solid #f0f0f0;
  //transition: bottom .2s linear;
}

.comment-editor-form__bd {
  display: flex;
  align-items: center;
  grid-gap: 16px;
}

textarea {
  flex: 1;
  background: #f2f2f2;
  border-radius: 6px;
  padding: 8px;
  line-height: 20px;
  height: 76px;
}

.submit-btn {
  color: #13C2C2;
  font-size: 16px;
  font-weight: bold;
  line-height: 76px;
}

.comment-editor-form__ft {
  display: flex;
  align-items: center;
  grid-gap: 24px;
  height: 24px;
  margin-top: 8px;
}

.upload-image {
  height: 24px;

  image {
    width: 24px;
    height: 24px;
  }
}
</style>