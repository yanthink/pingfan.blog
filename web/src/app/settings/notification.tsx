import React from 'react';
import { CheckCard, ProForm } from '@ant-design/pro-components';
import { useUser } from '@/hooks';
import { updateUserMeta } from '@/services';
import { message } from 'antd';

const Notification: React.FC = () => {
  const { user, setUser } = useUser();

  return (
    <>
      <ProForm
        initialValues={user}
        onFinish={async values => {
          await updateUserMeta(values);
          setUser({ ...user, ...values });

          message.success('修改成功');
        }}
        submitter={{
          render: (_, doms) => doms.pop(),
        }}
        style={{ paddingTop: 12 }}
      >
        <ProForm.Item label="邮件通知" name={['meta', 'emailNotify']}>
          <CheckCard.Group size="small">
            <CheckCard
              title="关闭通知"
              description="系统将不发送任何邮件通知"
              value={0}
              style={{ width: 250 }}
            />
            <CheckCard
              title="开启通知"
              description="系统将发送所有邮件通知"
              value={1}
              style={{ width: 250 }}
            />
            <CheckCard
              title="离线通知"
              description="系统在你离线时将以邮件的形式通知"
              value={2}
              style={{ width: 250 }}
            />
          </CheckCard.Group>
        </ProForm.Item>
      </ProForm>
    </>
  );
};

export default Notification;