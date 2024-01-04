import React from 'react';
import { Tooltip, theme } from 'antd';
import { HeartTwoTone, HeartOutlined } from '@ant-design/icons';
import classnames from 'classnames';

interface FavoriteProps {
  favorited?: boolean;
  suffix?: string | number;
  size?: 'default' | 'large';

  onClick?(): void;
}

const Favorite: React.FC<FavoriteProps> = ({ favorited, suffix = '', size, onClick }) => {
  const { token } = theme.useToken();

  return (
    <>
      <style jsx>{`
        .favorite {
          cursor: pointer;
        }

        .favorite-large :global(.anticon) {
          color: ${token.colorTextSecondary};
          font-size: 40px;
        }
      `}</style>
      <Tooltip title="收藏">
        <div className={classnames('favorite', { 'favorite-large': size === 'large' })} onClick={onClick}>
          {favorited ? <HeartTwoTone twoToneColor="#eb2f96" /> : <HeartOutlined />} {suffix}
        </div>
      </Tooltip>
    </>
  );
};

export default Favorite;