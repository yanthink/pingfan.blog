import React from 'react';
import { Nodes } from '@/markdown';
import { Fragment, jsx, jsxs } from 'react/jsx-runtime';
import { toJsxRuntime } from 'hast-util-to-jsx-runtime';
import Link from 'next/link';
import Pre from './Pre';

interface BodyProps {
  hast: Nodes;
}

const Body: React.FC<BodyProps> = ({ hast }) => {
  return (
    toJsxRuntime(hast as any, {
      Fragment,
      jsx,
      jsxs,
      ignoreInvalidStyle: true,
      passKeys: true,
      passNode: true,
      components: {
        a({ node, href = '', target, ...props }) {
          if (href && !target && !/^(https)?:\/\//i.test(href) && href[0] !== '#') {
            return <Link href={href} {...props as any} />;
          }

          if (href && !target && href[0] !== '#') {
            target = '_blank';
          }

          return <a {...props} href={href} target={target} />;
        },
        pre({ node, ...props }) {
          return <Pre {...props} />;
        },
      },
    })
  );
};

export default React.memo(Body, (prevProps, nextProps) => {
  return prevProps.hast === nextProps.hast;
});