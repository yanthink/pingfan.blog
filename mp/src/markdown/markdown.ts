import { unified, type Processor } from 'unified';
import type { Root as MdastRoot } from 'mdast';
import type { Root as HastRoot } from 'hast';
import { VFile } from 'vfile';
import remarkParse from 'remark-parse';
import remarkGfm from 'remark-gfm';
import remarkBreaks from 'remark-breaks';
import remarkGemoji from 'remark-gemoji';
import { defaultSchema, type Schema } from 'hast-util-sanitize';
import remarkRehype from 'remark-rehype';
import type { Options as RehypeOptions } from 'remark-rehype';
import rehypeRaw from 'rehype-raw';
import rehypeSanitize from 'rehype-sanitize';
import rehypeStringify from 'rehype-stringify';
import { visit } from 'unist-util-visit';
import { rehypePrismCommon } from 'rehype-prism-plus';
import { defaultHandlers } from 'mdast-util-to-hast';

interface Plugin {
  remark?: (p: Processor<MdastRoot, MdastRoot | undefined>) => Processor<MdastRoot, MdastRoot | undefined>;
  rehype?: (p: Processor<MdastRoot, MdastRoot, HastRoot>) => Processor<MdastRoot, MdastRoot, HastRoot>;
}

interface Options {
  rehypeOptions?: RehypeOptions;
  plugins?: Plugin[];
  sanitize?: (schema: Schema) => Schema;
}

export const articleStrategy: Options = {
  rehypeOptions: {
    handlers: {
      code(h, node) {
        const tree = defaultHandlers.code(h, node);

        visit(tree, 'element', node => {
          if (node.tagName === 'code' && typeof (node.data as any)?.meta === 'string') {
            node.properties!.dataMeta = (node.data as any).meta;
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
  plugins: [{
    remark: processor => processor.use(remarkGfm).use(remarkBreaks).use(remarkGemoji),
    rehype: processor => processor
      .use(() => (tree: HastRoot) => {
        visit(tree, 'element', (node, index, parent?: any) => {
          if (!parent || parent.tagName !== 'pre' || node.tagName !== 'code') {
            return;
          }

          visit(node, 'text', textNode => {
            textNode.value = textNode.value?.replace('    ', '  ')
          })

          if (!node.properties?.dataMeta) {
            return;
          }

          // rehype-raw 插件会导致代码块的 meta 数据丢失， 而显示行号，行高亮，diff等功能是通过 rehype-prism-plus 插件读取 meta 实现的
          // 我们通过 remark-rehype 插件的 handler 将代码块的 meta 数据保存到 data-meta 里面，然后在这里恢复
          // 注意 sanitize 需要将 code 的 dataMeta 属性加入白名单，不然会被忽略掉
          node.data = { ...node.data, meta: node.properties.dataMeta };
        });
      })
      .use(rehypePrismCommon, { ignoreMissing: true }),
  }],
};

export const commentStrategy: Options = {
  plugins: [{ remark: processor => processor.use(remarkGfm).use(remarkBreaks).use(remarkGemoji) }],
  sanitize(schema) {
    return {
      ...schema,
      clobberPrefix: '',
      tagNames: schema.tagNames?.filter(tagName => !['h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'hr'].includes(tagName)),
    };
  },
};

export const articleSearchStrategy: Options = {
  sanitize(schema) {
    return {
      ...schema,
      clobberPrefix: '',
      tagNames: ['mark'],
    };
  },
}

export const stripTagsStrategy: Options = {
  sanitize: schema => ({ ...schema, tagNames: ['img'] }),
  plugins: [
    {
      rehype(processor) {
        return processor
          .use(() => (tree: HastRoot) => {
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

export function toHast(md: string, options: Options = commentStrategy): HastRoot {
  const file = new VFile(md);
  const processor = getProcessor(options);

  const mdast = processor.parse(file);

  return processor.runSync(mdast, file);
}

export function getProcessor({ plugins, sanitize, rehypeOptions }: Options) {
  let mdastProcessor = unified().use(remarkParse);

  plugins?.forEach(({ remark }) => {
    if (remark) mdastProcessor = remark(mdastProcessor);
  });

  const schema = {
    ...defaultSchema,
    clobberPrefix: '',
    attributes: {
      ...defaultSchema.attributes,
      '*': [...defaultSchema.attributes?.['*'] ?? [], 'className'],
      code: [...defaultSchema.attributes?.code ?? [], 'className', 'dataMeta'],
      img: [...defaultSchema.attributes?.code ?? [], 'src', 'alt', 'className', 'dataIndex'],
    },
  };

  let hastProcessor = mdastProcessor
    .use(remarkRehype, { allowDangerousHtml: true, ...rehypeOptions })
    .use(rehypeRaw)
    .use(rehypeSanitize, sanitize?.(schema) ?? schema);

  plugins?.forEach(({ rehype }) => {
    if (rehype) hastProcessor = rehype(hastProcessor);
  });

  return hastProcessor.use(rehypeStringify);
}
