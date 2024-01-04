'use client';

import { Result, Button } from 'antd';
import { useRouter } from '@/hooks';

interface ErrorProps {
  error: Error & { digest?: string; };
  reset: () => void;
}

const Error: React.FC<ErrorProps> = ({ error, reset }) => {
  const router = useRouter();

  return (
    <Result
      status={500}
      title="500"
      subTitle="Sorry, something went wrong."
      extra={[
        <Button key="home" type="primary" onClick={() => router.push('/')}>
          返回首页
        </Button>,
        <Button key="reset" onClick={reset}>再试一次</Button>,
      ]}
    >
      <pre style={{ overflow: 'auto' }}>
        {error.stack}
      </pre>
    </Result>
  );
};

export default Error;