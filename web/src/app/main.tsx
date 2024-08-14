'use client'

import React, { useEffect, useState } from 'react'
import { Avatar, Flex, Tag, Tabs, theme, Form } from 'antd'
import { ProList, LightFilter, ProFormSelect, ProFormDateRangePicker } from '@ant-design/pro-components'
import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import { useRouter } from '@/hooks'
import { mergeSearchParams, prettyNumber } from '@/utils'
import {
  ClockCircleOutlined,
  EyeOutlined,
  HeartOutlined,
  LikeOutlined,
  MessageOutlined,
} from '@ant-design/icons'
import PrettyTime from '@/components/PrettyTime'
import { getTags } from '@/services'
import { useRequest } from 'ahooks'
import dayjs from 'dayjs'

interface MainProps {
  articles: API.Article[];
  total: number;
}

const Main: React.FC<MainProps> = ({ articles, total }) => {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { token } = theme.useToken()
  const [form] = Form.useForm()

  const { data } = useRequest(() => getTags({ pageSize: 100 }), { cacheKey: 'tags' })

  const [tabKey, setTabKey] = useState(searchParams.get('order') ?? 'hot')

  useEffect(() => {
    setTabKey(searchParams.get('order') ?? 'hot')
  }, [searchParams])

  useEffect(() => {
    form.setFieldsValue({
      tag: searchParams.has('tag') ? JSON.parse(searchParams.get('tag')!) : undefined,
      date: searchParams.has('startDate') && searchParams.has('endDate')
          ? [dayjs(searchParams.get('startDate')), dayjs(searchParams.get('endDate'))]
          : undefined,
    })
  }, [form, searchParams])

  const current = Number(searchParams.get('page') ?? 1)
  const pageSize = Number(searchParams.get('pageSize') ?? 10)

  return (
      <>
        <style jsx>{`
            .list {
                background: #fff;
                border-radius: 6px;

                :global(.ant-tabs-nav-list) {
                    margin: 0 24px;
                }

                :global(.ant-pro-list-row-content) {
                    margin-left: 0;
                    margin-right: 0;
                    line-height: 22px;
                    display: -webkit-box;
                    -webkit-line-clamp: 3;
                    -webkit-box-orient: vertical;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    height: 66px;
                }

                :global(.ant-pro-list-row-title) {
                    display: -webkit-box;
                    -webkit-line-clamp: 2;
                    -webkit-box-orient: vertical;
                    overflow: hidden;
                    text-overflow: ellipsis;
                }

                :global(.title) {
                    color: ${token.colorTextHeading};

                    &:hover {
                        color: ${token.colorPrimary};
                    }
                }

                .action {
                    display: flex;
                    line-height: 24px;
                    align-items: center;
                    gap: 4px;
                }

                :global(.ant-list-item-main) {
                    max-width: 100%;
                    overflow: hidden;
                }
            }

            @media screen and (max-width: 768px) {
                .list :global(.ant-list-item-extra) {
                    display: none;
                }

                .form {
                    display: none;
                }
            }
        `}</style>
        <div className="list">
          <Tabs
              activeKey={tabKey}
              items={[
                { label: '热门', key: 'hot' },
                { label: '最新', key: 'latest' },
                { label: '点赞最多', key: 'like' },
                { label: '评论最多', key: 'comment' },
              ]}
              onChange={key => router.push(`?${mergeSearchParams(searchParams, { order: key, page: 1 })}`)}
              tabBarExtraContent={{
                right: (
                    <div className="form">
                      <LightFilter
                          style={{ marginRight: 24 }}
                          form={form}
                          onValuesChange={(_, values) => {
                            const params = new URLSearchParams(searchParams)
                            params.delete('tag')
                            params.delete('startDate')
                            params.delete('endDate')

                            router.push(`?${mergeSearchParams(params, {
                              ...(values.tag && { tag: JSON.stringify(values.tag) }),
                              ...(values.date && { startDate: values.date[0], endDate: values.date[1] }),
                            })}`)
                          }}
                      >
                        <ProFormDateRangePicker
                            name="date"
                            label="时间"
                            fieldProps={{
                              placement: 'bottomRight',
                              presets: [
                                {
                                  label: '1个月内',
                                  value: [dayjs().add(-1, 'month').startOf('d'), dayjs().endOf('d')],
                                },
                                {
                                  label: '3个月内',
                                  value: [dayjs().add(-3, 'month').startOf('d'), dayjs().endOf('d')],
                                },
                                { label: '半年内', value: [dayjs().add(-6, 'month').startOf('d'), dayjs().endOf('d')] },
                                { label: '1年内', value: [dayjs().add(-1, 'year').startOf('d'), dayjs().endOf('d')] },
                                { label: '2年内', value: [dayjs().add(-2, 'year').startOf('d'), dayjs().endOf('d')] },
                                { label: '3年内', value: [dayjs().add(-3, 'year').startOf('d'), dayjs().endOf('d')] },
                              ],
                            }}
                        />
                        <ProFormSelect
                            label="标签"
                            name={['tag', 'id']}
                            options={data?.data?.map(tag => ({
                              label: tag.name,
                              value: tag.id,
                            }))}
                            fieldProps={{
                              placement: 'bottomRight',
                            }}
                            showSearch
                        />
                      </LightFilter>
                    </div>
                ),
              }}
          />
          <ProList<API.Article>
              itemLayout="vertical"
              rowKey="id"
              dataSource={articles}
              metas={{
                title: {
                  render: (_, article) => (
                      <Link className="title" href={`/articles/${article.id}`}>
                        {article.title}
                      </Link>
                  ),
                },
                description: {
                  render: (_, article) => article.tags?.map(tag => (
                      <Link
                          key={tag.id}
                          href={`?${mergeSearchParams(searchParams, {
                            tag: JSON.stringify({ id: tag.id }),
                            page: 1,
                          })}`}
                      >
                        <Tag>{tag.name}</Tag>
                      </Link>
                  )),
                },
                actions: {
                  render: (_, article) => (
                      <Flex wrap="wrap" gap={12}>
                        <div className="action">
                          <Avatar src={article.user?.avatar} size={24} /> {article.user?.name}
                        </div>
                        <div className="action">
                          <ClockCircleOutlined /><PrettyTime time={article.createdAt} />
                        </div>
                        <div className="action">
                          <EyeOutlined />{prettyNumber(article.viewCount ?? 0)}
                        </div>
                        <div className="action">
                          <LikeOutlined />{prettyNumber(article.likeCount ?? 0)}
                        </div>
                        <div className="action">
                          <MessageOutlined />{prettyNumber(article.commentCount ?? 0)}
                        </div>
                        <div className="action">
                          <HeartOutlined /> {prettyNumber(article.favoriteCount ?? 0)}
                        </div>
                      </Flex>
                  ),
                },
                extra: {
                  render: (_: any, article: API.Article) => article.preview && (
                      <img src={article.preview} width={320} height={180} style={{ objectFit: 'contain' }} alt="预览" />
                  ),
                },
                content: {
                  render: (_, article) => (
                      article.textContent?.substring(0, 300)
                  ),
                },
              }}
              pagination={{
                total,
                hideOnSinglePage: true,
                showLessItems: true,
                responsive: true,
                showSizeChanger: true,
                current,
                pageSize,
                itemRender(page, type, defaultDom: any) {
                  if (page < 1 || String(page) === (searchParams.get('page') ?? '1')) {
                    return defaultDom
                  }

                  const element = React.cloneElement(defaultDom)

                  return (
                      <Link href={`?${mergeSearchParams(searchParams, { page })}`}>
                        {element.type === 'a' ? element.props.children : defaultDom}
                      </Link>
                  )
                },
                onShowSizeChange(current, size) {
                  router.push(`?${mergeSearchParams(searchParams, { pageSize: size, page: 1 })}`)
                },
              }}
          />
        </div>
      </>
  )
}

export default Main