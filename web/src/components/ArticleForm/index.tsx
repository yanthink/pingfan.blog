'use client';

import React, { useEffect, useRef, useState } from 'react';
import {
  ProForm,
  ProFormInstance,
  ProFormSelect,
  ProFormText,
  ProFormUploadButton,
} from '@ant-design/pro-components';
import { Card, SelectProps, Spin, Image, Row, Col, Space, message } from 'antd';
import type { RcFile, UploadChangeParam, UploadProps } from 'antd/es/upload';
import { useDebounceFn } from 'ahooks';
import Link from 'next/link';
import { articleStrategy } from '@/markdown';
import Editor from '@/components/Editor';
import { getTags } from '@/services';
import { useToken } from '@/hooks';

interface ArticleFormProps {
  article?: API.Article;

  onSubmit?(values: Record<string, any>): void;
}

const formLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 2 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 20 },
  },
};

const getBase64 = (file: RcFile): Promise<string> => new Promise((resolve, reject) => {
  const reader = new FileReader();
  reader.readAsDataURL(file as any);
  reader.onload = () => resolve(reader.result as string);
  reader.onerror = (error) => reject(error);
});

const ArticleForm: React.FC<ArticleFormProps> = ({ article, onSubmit }) => {
  const formRef = useRef<ProFormInstance>();

  const [token] = useToken();

  const [tagOptions, setTagOptions] = useState<SelectProps['options']>([]);
  const [tagsFetching, setTagsFetching] = useState(false);
  const [tagsRemoteSearch, setTagsRemoteSearch] = useState(false);
  const [previewUrl, setPreviewUrl] = useState('');
  const [previewVisible, setPreviewVisible] = useState(false);

  const fetchTagsRef = useRef(0);

  const { run: debounceQueryTags } = useDebounceFn(async (q = '') => {
    fetchTagsRef.current += 1;

    const fetchId = fetchTagsRef.current;
    setTagOptions([]);
    setTagsFetching(true);

    const { data, total = 0 } = await getTags({ q, pageSize: 100 });

    if (fetchId !== fetchTagsRef.current) {
      return;
    }

    setTagOptions(data.map((tag) => ({ key: tag.id, label: tag.name, value: tag.id })));
    setTagsFetching(false);

    if (!q) {
      setTagsRemoteSearch(data.length < total);
    }
  }, { wait: 800 });

  useEffect(() => {
    debounceQueryTags()
  }, [debounceQueryTags]);

  useEffect(() => {
    console.log(article);
    if (article?.preview) {
      formRef.current?.setFieldValue(
        ['previews'] as any,
        [{ uid: '-1', name: 'preview.png', status: 'done', url: article.preview }],
      );
    }
  }, [article]);

  return (
    <Card title={article?.id ? '编辑文章' : '发布文章'}>
      <ProForm<API.Article & {
        previews: UploadProps['fileList'];
        tagIds: number[];
      }>
        formRef={formRef}
        layout="horizontal"
        {...formLayout}
        initialValues={{
          ...article,
          tagIds: article?.tags?.map(tag => tag.id),
        }}
        submitter={{
          render: (props, doms) => {
            return (
              <Row>
                <Col xs={24} sm={{ span: 20, offset: 2 }}>
                  <Space>{doms}</Space>
                </Col>
              </Row>
            );
          },
        }}
        onFinish={async (values) => {
          const { previews, ...data } = values;

          data.preview = (previews as any)?.filter((file: any) => file.status === 'done')[0]?.response?.data?.url
            ?? previews?.[0]?.url
            ?? '';

          await onSubmit?.(data);

          return true;
        }}
      >
        <ProFormText
          label="标题" name="title" placeholder="请输入文章标题"
          rules={[{ required: true, message: '请输入标题' }]}
        />
        <ProForm.Item label="标签">
          <ProFormSelect
            name="tagIds"
            mode="multiple"
            placeholder="文章标签"
            fieldProps={{
              options: tagOptions,
              loading: tagsFetching,
              filterOption: !tagsRemoteSearch,
              showSearch: tagsRemoteSearch || undefined,
              onSearch(q) {
                if (!/^\s*$/.test(q)) {
                  debounceQueryTags(q);
                }
              },
              notFoundContent: tagsFetching ? <Spin size="small" /> : null,
              maxTagCount: 5,
            }}
            noStyle
          />
          <Link href="/">添加标签</Link>
        </ProForm.Item>
        <ProFormUploadButton
          label="预览图" name="previews"
          action={`${process.env.NEXT_PUBLIC_BASE_URL}/api/resources/upload`}
          accept="image/*"
          max={1}
          listType="picture-card"
          fieldProps={{
            name: 'file',
            headers: { Authorization: `Bearer ${token}` },
            data: { type: 'articleImage' },
            onChange({ file }: UploadChangeParam) {
              if (file.status === 'error') {
                message.error(file.response?.message ?? '上传失败');
              }
            },
            async onPreview(file) {
              if (file.originFileObj) {
                const base64 = await getBase64(file.originFileObj as RcFile);
                setPreviewUrl(base64);
                setPreviewVisible(true);

                return;
              }

              if (file.url) {
                setPreviewUrl(file.url);
                setPreviewVisible(true);
              }
            },
          }}
          rules={[{
            validator(rule, value) {
              if (value?.some((file: any) => file.status === 'error')) {
                return Promise.reject(value?.[0]?.response?.message ?? '上传失败');
              }

              return Promise.resolve();
            },
          }]}
        />
        <ProForm.Item label="内容" name="content" rules={[{ required: true, message: '请输入文章内容' }]} required>
          <Editor resourceType="articleImage" options={articleStrategy} placeholder="请输入文章内容" />
        </ProForm.Item>
      </ProForm>
      <Image
        alt="预览"
        preview={{
          visible: previewVisible,
          src: previewUrl,
          onVisibleChange: setPreviewVisible,
        }}
        style={{ display: 'none' }}
      />
    </Card>
  );
};

export default ArticleForm;