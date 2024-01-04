'use client';

import React from 'react';
import { Card, Col, Row, Affix } from 'antd';
import { useCreation, useUpdate } from 'ahooks';
import { parse, articleStrategy, hastToTocAnchor } from '@/markdown';
import Content from './content';
import Toc from './toc';
import Comments from './comments';

interface MainProps {
  article: API.Article;
  comments: API.Comment[];
  commentTotal: number;
}

const Main: React.FC<MainProps> = ({ article, comments, commentTotal }) => {
  const { anchorItems, ...viewerProps } = useCreation(() => {
    const viewerProps = parse(article.content!, articleStrategy);
    const anchorItems = hastToTocAnchor(viewerProps.hast, [{
      type: 'element',
      tagName: 'h1',
      properties: { id: '评论' },
      children: [],
    }]);

    return { ...viewerProps, anchorItems };
  }, []);

  const update = useUpdate();
  const setArticle = (newArticle: API.Article) => {
    Object.assign(article, newArticle);
    update();
  };

  return (
    <Row gutter={24}>
      <Col xl={18} lg={24} md={24} sm={24} xs={24}>
        <Card bordered={false}>
          <Content article={article} onChange={setArticle} {...viewerProps} />
        </Card>
        <Card
          id={anchorItems[anchorItems.length - 1]?.key as string ?? '评论'}
          bordered={false} style={{ marginTop: 24 }}
        >
          <Comments article={article} comments={comments} total={commentTotal} />
        </Card>
      </Col>
      <Col xl={6} lg={0} md={0} sm={0} xs={0}>
        {anchorItems.length > 0 && (
          <Affix offsetTop={56}>
            <Card bordered={false} style={{ maxHeight: 'calc(100vh - 56px)', overflowY: 'auto' }}>
              <Toc items={anchorItems} />
            </Card>
          </Affix>
        )}
      </Col>
    </Row>
  );
};

export default Main;