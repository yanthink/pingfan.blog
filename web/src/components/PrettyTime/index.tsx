import React, { useMemo } from 'react';
import type { ConfigType } from 'dayjs';
import dayjs from 'dayjs';
import { Tooltip } from 'antd';

interface PrettyTimeProps {
  time: ConfigType;
}


const PrettyTime: React.FC<PrettyTimeProps> = ({ time }) => {
  const formatted = useMemo(() => {
    const now = dayjs();
    const current = dayjs(time);

    if (now.diff(current, 'd') <= 15) {
      return current.fromNow();
    }

    if (now.year() === current.year()) {
      return current.format('MM-DD HH:mm');
    }

    return current.format('YYYY-MM-DD HH:mm:ss');
  }, [time]);

  if (formatted !== dayjs(time).format('YYYY-MM-DD HH:mm:ss')) {
    return (
      <Tooltip title={dayjs(time).format('YYYY-MM-DD HH:mm:ss')}>
        {formatted}
      </Tooltip>
    );
  }

  return formatted;
};

export default PrettyTime;