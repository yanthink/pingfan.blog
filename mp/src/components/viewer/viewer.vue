<template>
  <block v-if="node.type === 'root'">
    <view class="markdown-body">
      <viewer v-for="(subNode, index) in node.children" :key="index" :node="subNode" />
    </view>
  </block>
  <block v-else-if="node.type === 'element'">
    <image
        v-if="node.tagName === 'img'" :src="node.properties.src" :alt="node.properties.alt" mode="widthFix"
        @click.stop.prevent="previewImage(node.properties.src)"
    />
    <navigator
        v-else-if="node.tagName === 'a'"
        :url="node.properties.href"
        @click.stop.prevent="copyText(node.properties.href)"
    >
      <viewer v-for="(subNode, index) in node.children ?? []" :key="index" :node="subNode" />
    </navigator>
    <checkbox
        v-else-if="node.tagName === 'input' && node.properties.type === 'checkbox'"
        :checked="node.properties.checked"
        :disabled="node.properties.disabled"
    />
    <view v-else-if="node.tagName === 'table'" class="table-wrapper">
      <view :class="classnames(node.properties.className, node.tagName)">
        <viewer v-for="(subNode, index) in node.children ?? []" :key="index" :node="subNode" />
      </view>
    </view>
    <view
        v-else
        :class="classnames(node.properties.className, node.tagName)"
        :data-line="node.properties.line?.[0]"
    >
      <viewer v-for="(subNode, index) in node.children ?? []" :key="index" :node="subNode" />
    </view>
  </block>
  <block v-else-if="node.type === 'text'">
    <text v-if="node.value.trim().length">{{ node.value }}</text>
    <block v-else>{{ node.value }}</block>
  </block>
</template>

<script lang="ts" setup>
import type { Nodes } from 'hast';
import classnames from 'classnames';

const props = defineProps<{ node: Nodes; }>();

function previewImage(url: string) {
  uni.previewImage({
    urls: [url],
  });
}

function copyText(content: string) {
  uni.setClipboardData({ data: content })
}
</script>

<script lang="ts">
export default {
  options: {
    virtualHost: true,
  },
};
</script>

<style lang="scss" scoped>
@import "./prism-material-dark.css";

.markdown-body {
  line-height: 1.8;
  font-size: 16px;
  max-width: 100%;
  word-break: break-all;
  position: relative;

  // 标题样式
  .h1, .h2, .h3, .h4, .h5, .h6 {
    margin-top: 24px;
    margin-bottom: 16px;
    font-weight: 600;
    line-height: 1.25;
  }

  .h1 {
    font-weight: 600;
    padding-bottom: 0.3em;
    font-size: 1.8em;
    border-bottom: 1px solid hsl(210, 18%, 87%);
  }

  .h2 {
    font-weight: 600;
    padding-bottom: 0.3em;
    font-size: 1.5em;
    border-bottom: 1px solid hsl(210, 18%, 87%);
  }

  .h3 {
    font-weight: 600;
    font-size: 1.25em;
  }

  .h4 {
    font-weight: 600;
    font-size: 1em;
  }

  .h5 {
    font-weight: 600;
    font-size: 0.875em;
  }

  .h6 {
    font-weight: 600;
    font-size: 0.85em;
    color: #656d76;
  }

  .p, .blockquote, .ul, .ol, .dl, .table, .pre, .details {
    margin-top: 0;
    margin-bottom: 16px;
  }

  // 列表样式
  .ul, .ol {
    padding-left: 2em;
  }

  .ol {
    list-style-type: decimal;
  }

  .li {
    display: list-item;
    text-align: -webkit-match-parent;
  }

  .task-list-item {
    list-style-type: none;
    transform: translateX(-1em);
  }

  checkbox {
    transform: scale(.75) translateY(-2px);
  }

  // 内联样式
  .a, .br, .strong, .b, .em, .i, .del, .s, .strike, .code, .span, .samp, .small, .kbd, .sub, .sup, .tt, .var, .q, navigator, .mark {
    display: inline;
  }

  .strong, .b {
    font-weight: bold;
  }

  .em {
    font-style: italic;
  }

  .del, .s, .strike {
    text-decoration: line-through;
  }

  .mark {
    background: #ffe58f;
  }

  // 内联代码
  .p > .code, .inline-code {
    margin: 0 4px;
    color: #476582;
    background-color: rgba(27, 31, 35, .05);
    border-radius: 4px;
    border: 1px solid #e4e4e4;
    padding: 2px 6px;
    font-size: 14px;
    line-height: 14px;
  }

  // 代码块
  .pre, .pre .code {
    font-size: 1em;
    tab-size: 4;
    color: #8ecdff;
    background: #383838;
    display: block;
    padding: 12px;
    width: max-content;
    min-width: 100%;
  }

  .pre, .pre view {
    white-space: pre-wrap;
  }

  .pre {
    width: 100%;
    overflow-x: auto;
    padding: 0;
  }

  // 行号
  .code-line,
  .line-number:before {
    display: block;
    border-left: 4px solid transparent;
  }

  .line-number {
    background: inherit;
    padding-left: 48px;

    &:before {
      display: inline-block;
      content: attr(data-line);
      user-select: none;
      box-sizing: border-box;
      width: 40px;
      padding: 0 6px;
      text-align: right;
      color: #808080;
      position: absolute;
      left: 0;
      background: inherit;
    }
  }

  // 高亮行
  .highlight-line,
  .highlight-line.line-number:before {
    border-color: #3992ff;
    background-color: rgba(51, 65, 85, .5);
    background-color: #353f4d;
  }

  .code-line.inserted {
    background-color: #3a5248;
  }

  .code-line.deleted {
    background-color: #5b3d3a;
  }

  // 引用
  .blockquote {
    color: rgba(0, 0, 0, 0.45);
    background-color: #e6f7ff;
    border-left: 8px solid #91d5ff;
    padding: 12px;

    > view:first-child {
      margin-top: 0;
    }

    > view:last-child {
      margin-bottom: 0;
    }
  }

  // 表格
  .table-wrapper {
    display: block;
    width: max-content;
    max-width: 100%;
    overflow: auto;
  }

  .table {
    display: table;
    border-collapse: collapse;
  }

  .thead {
    display: table-header-group;
    font-weight: bold;
  }

  .tr {
    display: table-row;

    &:nth-child(2n) {
      background-color: #f6f8fa;
    }
  }

  .th {
    white-space: nowrap;
  }

  .th, .td {
    display: table-cell;
    padding: 8px;
    border: 1px solid #ddd;
  }

  .tbody {
    display: table-row-group;
  }

  .hr {
    height: 1px;
    background-color: hsl(210, 18%, 87%);
  }

  // 内置组件
  navigator {
    color: #0969da;
  }

  image {
    max-width: 100%;
  }
}
</style>