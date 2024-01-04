'use client';

import React from 'react';
import { Editor as ByteMDE } from '@bytemd/react';
import type { EditorProps as ByteMDEProps } from '@bytemd/react';
import bytemdZhHans from 'bytemd/locales/zh_Hans.json';
import { commentStrategy, type Options } from '@/markdown';
import { upload } from '@/services/upload';
import 'bytemd/dist/index.css';

export interface EditorProps extends Omit<ByteMDEProps, 'value'> {
  resourceType: string;
  value?: string;
  options?: Options;
}

const Editor: React.FC<EditorProps> = ({ resourceType, options = commentStrategy, ...rest }) => {
  return (
    <ByteMDE
      {...options as any}
      locale={bytemdZhHans as any}
      uploadImages={files => upload(files, resourceType)}
      {...rest}
    />
  );
};

export default Editor;