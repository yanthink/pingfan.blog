'use client';

import React, { useRef, useState } from 'react';
import { AutoComplete, Avatar, Space, Spin, Tag, theme, Typography } from 'antd';
import {
  ClockCircleOutlined,
  EyeOutlined,
  HeartOutlined,
  LikeOutlined,
  SearchOutlined,
  TagsFilled,
} from '@ant-design/icons';
import { useRouter } from '@/hooks';
import { useRequest } from 'ahooks';
import { searchArticles } from '@/services';
import PrettyTime from '@/components/PrettyTime';
import { prettyNumber } from '@/utils';

const HeaderSearch: React.FC = () => {
  const { token } = theme.useToken();

  const router = useRouter();

  const selectCalled = useRef(false);

  const [value, setValue] = useState('');

  const { data, loading, run, cancel } = useRequest(async (q = '') => {
    if (!q.trim()) {
      return { data: [], total: 0 };
    }

    return searchArticles({ q });
  }, { manual: true, debounceWait: 800 });

  function toSearch() {
    if (value.trim().length) {
      cancel();
      router.push(`/search?q=${value.trim()}`);
    }
  }

  return (
    <>
      <style jsx>{`
        :global(.search-popup .ant-select-item) {
          border-bottom: 1px ${token.colorBorder} solid;
          border-radius: 0;
        }

        .item {
          display: flex;
          flex-direction: column;
          width: 800px;
          white-space: normal;

          &:hover {
            text-decoration: underline;
          }
        }

        .content {
          display: -webkit-box;
          -webkit-line-clamp: 3;
          -webkit-box-orient: vertical;
          overflow: hidden;
          text-overflow: ellipsis;
          margin: 4px 0 8px 0;
        }

        .view-all {
          display: flex;
          justify-content: center;
          align-items: center;
          padding: 12px;
          cursor: pointer;

          &:hover {
            background: ${token.colorBgTextHover};
          }
        }
      `}</style>
      <div className="search" style={{ display: 'flex', alignItems: 'center' }} aria-hidden>
        <AutoComplete
          value={value}
          options={data?.data?.map(article => ({
            value: article.id,
            label: (
              <div className="item">
                <h4 dangerouslySetInnerHTML={{ __html: article.highlights?.title?.[0] ?? article.title! }} />
                <div className="content" dangerouslySetInnerHTML={{ __html: article.highlights?.content?.[0] ?? '' }} />
                <Space align="center" size={12} wrap>
                  <Typography.Text type="secondary">
                    <Avatar src={article.user?.avatar} size={24} /> {article.user?.name}
                  </Typography.Text>
                  <Typography.Text type="secondary">
                    <ClockCircleOutlined /> <PrettyTime time={article.createdAt} />
                  </Typography.Text>
                  <Typography.Text type="secondary">
                    <EyeOutlined /> {prettyNumber(article.viewCount ?? 0)}
                  </Typography.Text>
                  <Typography.Text type="secondary">
                    <LikeOutlined /> {prettyNumber(article.likeCount ?? 0)}
                  </Typography.Text>
                  <Typography.Text type="secondary">
                    <HeartOutlined /> {prettyNumber(article.favoriteCount ?? 0)}
                  </Typography.Text>
                  {article.tags?.length && (
                    <Typography.Text type="secondary">
                      <TagsFilled style={{ marginRight: 4, fontSize: 16, transform: 'translateY(2px)' }} />
                      {article.tags?.map(tag => <Tag key={tag.id}>{tag.name}</Tag>)}
                    </Typography.Text>
                  )}
                </Space>
              </div>
            ),
          })) ?? []}
          dropdownRender={menu => (
            <>
              {menu}
              {!loading && (
                <div className="view-all" onClick={toSearch}>
                  <SearchOutlined style={{ marginRight: 4 }} />
                  查看全部{data?.total ?? 0}个结果
                </div>
              )}
            </>
          )}
          notFoundContent={loading ? <Spin /> : null}
          popupMatchSelectWidth={false}
          placement="bottomRight"
          suffixIcon={
            <SearchOutlined
              style={{ color: 'rgba(0, 0, 0, 0.15)' }}
              onClick={toSearch}
            />
          }
          bordered={false}
          defaultActiveFirstOption={false}
          onSearch={value => {
            setValue(value);
            run(value);
          }}
          onSelect={(value) => {
            selectCalled.current = true;
            router.push(`/articles/${value}`);
          }}
          onKeyDown={e => {
            if (e.code === 'Enter' && !selectCalled.current) {
              toSearch();
            }
            selectCalled.current = false;
          }}
          showSearch
          listHeight={500}
          popupClassName="search-popup"
          style={{ width: 200, borderRadius: 4, backgroundColor: 'rgba(0,0,0,0.03)' }}
          placeholder="站内搜索"
        />
      </div>
    </>
  );
};

export default HeaderSearch;