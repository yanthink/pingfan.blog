'use client';

import { Anchor } from 'antd';
import type { AnchorLinkItemProps } from 'antd/es/anchor/Anchor';
import type { HtmlElementNode } from '@jsdevtools/rehype-toc';

interface TocProps {
  items?: AnchorLinkItemProps[];
}

const Toc: React.FC<TocProps> = ({ items }) => {
  return (
    <Anchor affix={false} showInkInFixed items={items} targetOffset={60} />
  );
};

export default Toc;