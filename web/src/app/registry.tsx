'use client';

import React, { useRef, useState } from 'react';
import { createCache, extractStyle, StyleProvider } from '@ant-design/cssinjs';
import type Entity from '@ant-design/cssinjs/es/Cache';
import { useServerInsertedHTML } from 'next/navigation';
import { StyleRegistry, createStyleRegistry } from 'styled-jsx';
import { useCreation } from 'ahooks';

const StyledJsxRegistry: React.FC<React.PropsWithChildren> = ({ children }) => {
  const inserted =  useRef(false);
  const cache = useCreation<Entity>(() => createCache(), []);
  const [jsxStyleRegistry] = useState(() => createStyleRegistry());

  useServerInsertedHTML(() => {
    if (inserted.current) {
      return;
    }

    inserted.current = true;

    const styles = jsxStyleRegistry.styles();
    jsxStyleRegistry.flush();

    return (
      <>
        <style id="antd" dangerouslySetInnerHTML={{ __html: extractStyle(cache, true) }} />
        {styles}
      </>
    );
  });

  return (
    <StyleProvider cache={cache}>
      <StyleRegistry registry={jsxStyleRegistry}>{children}</StyleRegistry>
    </StyleProvider>
  );
};

export default StyledJsxRegistry;