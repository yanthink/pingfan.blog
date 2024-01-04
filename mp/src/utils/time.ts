import type { ConfigType } from 'dayjs';
import dayjs from 'dayjs';

export function prettyTime(time: ConfigType): string {
  const now = dayjs();
  const current = dayjs(time);

  if (now.diff(current, 'd') <= 15) {
    return current.fromNow();
  }

  if (now.year() === current.year()) {
    return current.format('MM-DD');
  }

  return current.format('YYYY-MM-DD');
}