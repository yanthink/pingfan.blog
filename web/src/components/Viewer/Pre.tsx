import { useRef, useState } from 'react';
import { theme } from 'antd';
import { CopyOutlined, CheckOutlined } from '@ant-design/icons';

interface PreProps extends React.DetailedHTMLProps<React.HTMLAttributes<HTMLPreElement>, HTMLPreElement> {
}

const Pre: React.FC<PreProps> = (props) => {
  const { token } = theme.useToken();

  const textInput = useRef<HTMLDivElement>(null);
  const [copied, setCopied] = useState(false);

  const onCopy = () => {
    setCopied(true);
    navigator.clipboard.writeText(textInput.current?.textContent ?? '');

    setTimeout(() => {
      setCopied(false);
    }, 2000);
  };

  return (
    <>
      <style jsx>{`
        .pre-wrap {
          position: relative;
        }

        .copy {
          position: absolute;
          right: 16px;
          top: 16px;
          color: #fff;
          display: none;
        }

        .pre-wrap:hover .copy {
          display: inline-block;
        }
      `}</style>
      <div ref={textInput} className="pre-wrap">
        <a className="copy" onClick={onCopy}>
          {
            copied ?
              <CheckOutlined style={{ color: token.colorSuccess }} /> :
              <CopyOutlined style={{ color: '#cdd6e0' }} />
          }
        </a>
        <pre>{props.children}</pre>
      </div>
    </>
  );
};

export default Pre;