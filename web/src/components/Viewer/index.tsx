'use client';

import React, { useEffect, useMemo, useRef, useState } from 'react';
import { Image as AntdImage } from 'antd';
import type { VFile } from 'vfile';
import { parse, type Options, Nodes } from '@/markdown';
import Body from './Body';

export interface ViewerProps {
  children?: string;
  hast?: Nodes;
  file?: VFile;
  options?: Options;
}

const { PreviewGroup } = AntdImage;

const Viewer: React.FC<ViewerProps> = ({ children = '', ...props }) => {
  const viewerRef = useRef<HTMLDivElement>(null);

  const [previewVisible, setPreviewVisible] = useState(false);
  const [previewCurrent, setPreviewCurrent] = useState(0);
  const [images, setImages] = useState<string[]>([]);

  useEffect(() => {
    const imgElements = viewerRef.current!.querySelectorAll<HTMLImageElement>('img[data-index]');
    setImages(Array.from(imgElements).map((img) => img.src));
  }, [children]);

  const { hast, file, options } = useMemo(() => {
    return props.hast ? props : parse(children, props.options);
  }, [children, props]);

  useEffect(() => {
    if (!viewerRef.current || !file) return;

    const cbs = options?.plugins?.map(({ viewerEffect }) =>
      viewerEffect?.({ markdownBody: viewerRef.current!, file }),
    );

    return () => {
      cbs?.forEach(cb => cb?.());
    };
  }, [options, file]);

  return (
    <>
      <div
        ref={viewerRef}
        onClick={(e) => {
          const $ = e.target as HTMLElement;
          if ($.tagName === 'IMG') {
            const src = $.getAttribute('src');
            const index = $.dataset.index;

            if (src && index !== '' && index !== undefined) {
              setPreviewVisible(true);
              setPreviewCurrent(parseInt(index, 10));
            }
          }
        }}
        className="markdown-body"
      >
        <Body hast={hast!} />
      </div>
      {images.length > 0 && (
        <PreviewGroup
          preview={{
            visible: previewVisible,
            current: previewCurrent,
            onChange(current) {
              setPreviewCurrent(current);
            },
            onVisibleChange(visible) {
              setPreviewVisible(visible);
            },
          }}
          items={images}
        />
      )}
    </>
  );
};

export default React.memo(Viewer, (prevProps, nextProps) => {
  return prevProps.children === nextProps.children;
});