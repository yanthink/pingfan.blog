import { getProcessor, type ViewerProps } from 'bytemd';
import type { Root, RootContent, Element } from 'hast';
import gfm from '@bytemd/plugin-gfm';
import breaks from '@bytemd/plugin-breaks';
import gemoji from '@bytemd/plugin-gemoji';
import { rehypePrismCommon } from 'rehype-prism-plus';
import mermaid from '@bytemd/plugin-mermaid';
import rehypeSlug from 'rehype-slug';
import autolinkHeadings from 'rehype-autolink-headings';
import { visit } from 'unist-util-visit';
import { h } from 'hastscript';
import type { AnchorLinkItemProps } from 'antd/es/anchor/Anchor';
import { VFile } from 'vfile';
import { defaultHandlers } from 'mdast-util-to-hast';
import remarkMath from 'remark-math';
import type { default as K } from 'katex';

export type Nodes = Root | RootContent

interface AnchorLinkItem extends AnchorLinkItemProps {
  headingNumber: number;
  id: string;
  children: AnchorLinkItem[];
}

export interface Options extends Omit<ViewerProps, 'value'> {
}

let katex: typeof K;

// 文章策略
export const articleStrategy: Options = {
  remarkRehype: {
    handlers: {
      code(h, node) {
        const tree = defaultHandlers.code(h, node);

        visit(tree, 'element', node => {
          if (node.tagName === 'code' && typeof node.data?.meta === 'string') {
            node.properties!.dataMeta = node.data.meta;
          }
        });

        return tree;
      },
      inlineCode(h, node) {
        const tree = defaultHandlers.inlineCode(h, node);
        tree.properties!.className = 'inline-code';

        return tree;
      },
    },
  },
  sanitize(schema) {
    return {
      ...schema,
      clobberPrefix: '',
      attributes: {
        ...schema.attributes,
        code: [...schema.attributes?.code ?? [], 'className', 'dataMeta'],
        img: [...schema.attributes?.code ?? [], 'src', 'alt', 'className', 'dataIndex'],
      },
    };
  },
  plugins: [gfm(), breaks(), gemoji(), mermaid(),
    {
      remark: (processor) => processor.use(remarkMath),
      viewerEffect({ markdownBody }) {
        const renderMath = async (selector: string, displayMode: boolean) => {
          const els = markdownBody.querySelectorAll<HTMLElement>(selector);
          if (els.length === 0) return;

          if (!katex) {
            katex = await import('katex').then(m => m.default);
          }

          els.forEach((el) => {
            // 多次运行 viewerEffect 会导致数学公式错乱
            if (el.querySelector('.katex')) {
              return;
            }

            katex.render(el.innerText, el, {
              output: 'mathml',
              throwOnError: false,
              displayMode,
            });
          });
        };

        renderMath('.math.math-inline', false);
        renderMath('.math.math-display', true);
      },
    },
    {
      rehype(processor) {
        return processor
          .use(() => tree => {
            let imagesCount = 0;
            visit(tree, 'element', (node, index, parent?: any) => {
              if (node.tagName === 'img' && node.properties.src) {
                node.properties.dataIndex = String(imagesCount++);
                return;
              }

              if (!parent || parent.tagName !== 'pre' || node.tagName !== 'code' || !node.properties?.dataMeta) {
                return;
              }

              // rehype-raw 插件会导致代码块的 meta 数据丢失， 而显示行号，行高亮，diff等功能是通过 rehype-prism-plus 插件读取 meta 实现的
              // 我们通过 remark-rehype 插件的 handler 将代码块的 meta 数据保存到 data-meta 里面，然后在这里恢复
              // 注意 sanitize 需要将 code 的 dataMeta 属性加入白名单，不然会被忽略掉
              node.data = { ...node.data, meta: node.properties.dataMeta };
            });
          })
          .use(rehypePrismCommon, { ignoreMissing: true })
          .use(rehypeSlug)
          .use(() => autolinkHeadings({
            behavior: 'append',
            content: () => [h('span', '#')],
          }));
      },
    }],
};

// 评论策略
export const commentStrategy: Options = {
  remarkRehype: {
    handlers: {
      code(h, node) {
        const tree = defaultHandlers.code(h, node);

        visit(tree, 'element', node => {
          if (node.tagName === 'code' && typeof node.data?.meta === 'string') {
            node.properties!.dataMeta = node.data.meta;
          }
        });

        return tree;
      },
      inlineCode(h, node) {
        const tree = defaultHandlers.inlineCode(h, node);
        tree.properties!.className = 'inline-code';

        return tree;
      },
    },
  },
  sanitize(schema) {
    return {
      ...schema,
      clobberPrefix: '',
      tagNames: schema.tagNames?.filter(tagName => !['h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'hr'].includes(tagName)),
      attributes: {
        ...schema.attributes,
        code: [...schema.attributes?.code ?? [], 'className', 'dataMeta'],
        img: [...schema.attributes?.code ?? [], 'src', 'alt', 'className', 'dataIndex'],
      },
    };
  },
  plugins: [gfm(), breaks(), gemoji(), mermaid(),
    {
      remark: (processor) => processor.use(remarkMath),
      viewerEffect({ markdownBody }) {
        const renderMath = async (selector: string, displayMode: boolean) => {
          const els = markdownBody.querySelectorAll<HTMLElement>(selector);
          if (els.length === 0) return;

          if (!katex) {
            katex = await import('katex').then(m => m.default);
          }

          els.forEach((el) => {
            // 多次运行 viewerEffect 会导致数学公式错乱
            if (el.querySelector('.katex')) {
              return;
            }

            katex.render(el.innerText, el, {
              output: 'mathml',
              throwOnError: false,
              displayMode,
            });
          });
        };

        renderMath('.math.math-inline', false);
        renderMath('.math.math-display', true);
      },
    },
    {
      rehype(processor) {
        return processor
          .use(() => tree => {
            let imagesCount = 0;
            visit(tree, 'element', (node, index, parent?: any) => {
              if (node.tagName === 'img' && node.properties.src) {
                node.properties.dataIndex = String(imagesCount++);
                return;
              }

              if (!parent || parent.tagName !== 'pre' || node.tagName !== 'code' || !node.properties?.dataMeta) {
                return;
              }

              // rehype-raw 插件会导致代码块的 meta 数据丢失， 而显示行号，行高亮，diff等功能是通过 rehype-prism-plus 插件读取 meta 实现的
              // 我们通过 remark-rehype 插件的 handler 将代码块的 meta 数据保存到 data-meta 里面，然后在这里恢复
              // 注意 sanitize 需要将 code 的 dataMeta 属性加入白名单，不然会被忽略掉
              node.data = { ...node.data, meta: node.properties.dataMeta };
            });
          })
          .use(rehypePrismCommon, { ignoreMissing: true });
      },
    }],
};

// 去标签策略
export const stripTagsStrategy: Options = {
  sanitize: schema => ({ ...schema, tagNames: ['img'] }),
  plugins: [
    {
      rehype(processor) {
        return processor
          .use(() => (tree: Nodes) => {
            visit(tree, 'element', (node, index, parent?: any) => {
              if (node.tagName === 'img' && parent) {
                parent.children.splice(index, 1, { type: 'text', value: '[图片]' });
              }
            });
          });
      },
    },
  ],
};

export function parse(md: string, options: Options = commentStrategy): { hast: Nodes, file: VFile, options: Options } {
  const file = new VFile(md);
  const processor = getProcessor(options);
  const mdastTree = processor.parse(file);

  const hast = processor.runSync(mdastTree, file);

  return {
    hast,
    file,
    options,
  };
}

export function hastToTocAnchor(hast: Nodes, append: Element[] = []): AnchorLinkItemProps[] {
  const headings: Element[] = [];

  visit(hast, 'element', node => {
    if (['h1', 'h2', 'h3', 'h4', 'h5', 'h6'].includes(node.tagName)) {
      headings.push(node);
    }
  });

  headings.push(...append);

  return createAnchorItems(headings);
}

function createAnchorItems(headings: Element[]): AnchorLinkItemProps[] {
  const items: AnchorLinkItem[] = [];

  for (const heading of headings) {
    if (!heading.properties?.id) {
      continue;
    }

    const headingNumber = parseInt(heading.tagName.slice(-1), 10);

    const title = nodeText(heading);
    const id = heading.properties.id as string;

    const item: AnchorLinkItem = {
      key: id,
      headingNumber,
      id,
      href: `#${id}`,
      title: title !== '' ? title : id,
      children: [],
    };

    let lastItem = items.at(-1);

    if (!lastItem || lastItem.headingNumber >= headingNumber) {
      items.push(item);
      continue;
    }

    for (let i = lastItem.headingNumber + 1; i <= 6; i++) {
      const { children } = lastItem as AnchorLinkItem;

      lastItem = children.at(-1);

      if (!lastItem || lastItem.headingNumber >= headingNumber) {
        children.push(item);
        break;
      }
    }
  }

  return items;
}

function nodeText(node: Nodes): string {
  if (typeof node === 'object' && node.type === 'text') {
    return node.value ?? '';
  }

  if ('children' in node) {
    for (const child of node.children) {
      const text = nodeText(child);
      if (text !== '') {
        return text;
      }
    }
  }

  return '';
}

export { getProcessor };