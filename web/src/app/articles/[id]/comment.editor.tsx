'use client';

import React, { useRef } from 'react';
import { ProForm, type ProFormInstance } from '@ant-design/pro-components';
import Editor from '@/components/Editor';
import type { EditorProps } from '@/components/Editor';

interface CommentEditorProps extends Omit<EditorProps, 'resourceType'> {
  submitText?: string;

  onSubmit?(values: Record<string, any>): void;
}

const CommentEditor: React.FC<CommentEditorProps> = ({ submitText = '评论', onSubmit, ...rest }) => {
  const formRef = useRef<ProFormInstance>();

  return (
    <div style={{ background: '#fff', paddingBottom: 24 }}>
      <ProForm
        formRef={formRef}
        submitter={{
          searchConfig: {
            submitText,
          },
        }}
        onFinish={async values => {
          await onSubmit?.(values)
          formRef.current?.resetFields();
        }}
      >
        <ProForm.Item name="content" rules={[{ required: true, message: '请输入评论内容' }]} required>
          <Editor {...rest} mode="tab" maxLength={512} resourceType="commentImage" />
        </ProForm.Item>
      </ProForm>
    </div>
  );
};

export default CommentEditor;