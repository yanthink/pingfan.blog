import React from 'react';
import { Tooltip, theme } from 'antd';
import { LikeOutlined, LikeTwoTone } from '@ant-design/icons';
import classnames from 'classnames';

interface LikeProps {
  liked?: boolean;
  suffix?: React.ReactNode;
  size?: 'default' | 'large';

  onClick?(): void;
}

const Like: React.FC<LikeProps> = ({ liked, suffix = '', size, onClick }) => {
  const { token } = theme.useToken();

  return (
    <>
      <style jsx>{`
        .like {
          cursor: pointer;
        }

        .like-large :global(.anticon) {
          color: ${token.colorTextSecondary};
          font-size: 40px;
        }
      `}</style>
      <Tooltip title="点赞">
        <div className={classnames('like', { 'like-large': size === 'large' })} onClick={onClick}>
          {liked ? <LikeTwoTone twoToneColor={token.colorPrimary} /> : <LikeOutlined />} {suffix}
        </div>
      </Tooltip>
    </>
  );
};

export default Like;