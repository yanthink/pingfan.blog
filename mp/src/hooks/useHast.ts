import { ref } from 'vue';
import type { Root, ElementContent, RootContent } from 'hast';
import { useNextTick } from '@/hooks/useNextTick';

type Node = Root | RootContent & { children?: ElementContent[] }

function deepClone(obj: any): any {
  if (obj === null || typeof obj !== 'object') {
    return obj;
  }

  if (Array.isArray(obj)) {
    const newArr = [];
    for (let i = 0; i < obj.length; i++) {
      newArr[i] = deepClone(obj[i]);
    }
    return newArr;
  }

  const newObj: any = {};
  for (let key in obj) {
    if (obj.hasOwnProperty(key)) {
      newObj[key] = deepClone(obj[key]);
    }
  }

  return newObj;
}

export function useHast() {
  const node = ref<Root>({ type: 'root', children: [] });
  const nextTick = useNextTick();

  let result: Root = { type: 'root', children: [] };
  let i = 0;

  function* processNode(target: Node, current: Node, chunkSize: number): Generator<Root> {
    i++;

    const { children = [], ...value } = current;
    target.children ??= [];
    target.children.push(value as ElementContent);

    for (let i = 0; i < children.length; i++) {
      yield* processNode(value as ElementContent, children[i], chunkSize);
    }

    if (i >= chunkSize) {
      yield result;
      i = 0;
    }
  }

  function* chunkHastGenerator(hast: Root, chunkSize: number): Generator<Root> {
    result = { type: 'root', children: [] };
    i = 0;

    for (const item of hast.children) {
      yield* processNode(result, item, chunkSize);
    }

    if (i > 0) {
      yield result;
      i = 0;
    }
  }

  function cloneLastNode(target: Node) {
    if (target.children?.length) {
      target.children[target.children.length - 1] = deepClone(target.children[target.children.length - 1]);
    }
  }

  let chunking = false;

  async function setChunkHast(hast: Root, chunkSize = 500) {
    if (chunking) return;

    chunking = true;

    for (const item of chunkHastGenerator(hast, chunkSize)) {
      node.value = { ...item, children: [...item.children] };
      cloneLastNode(node.value);

      await nextTick();
    }

    chunking = false;
  }

  return { node, setChunkHast };
}