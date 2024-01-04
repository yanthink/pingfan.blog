import React from 'react';
import { ProForm, ProFormText } from '@ant-design/pro-components';
import { useUser } from '@/hooks';
import { updateUserPassword } from '@/services';
import { message } from 'antd';

const Security: React.FC = () => {
  const { user, setUser } = useUser();

  return (
    <>
      <ProForm
        onFinish={async values => {
          await updateUserPassword(values);
          setUser({ ...user, hasPassword: true });

          message.success('修改成功');
        }}
        style={{ paddingTop: 12 }}
      >
        <ProFormText.Password
          label="旧密码"
          name="oldPassword"
          rules={[{ required: user.hasPassword, message: '请输入你的旧密码。' }]}
          required={user.hasPassword}
          hidden={!user.hasPassword}
        />

        <ProFormText.Password
          label="新密码"
          name="password"
          rules={[
            { required: true, message: '请输入你的密码。' },
            { min: 6, message: '密码长度不能小于6位。' },
          ]}
          required
        />

        <ProFormText.Password
          label="确认密码"
          name="passwordConfirmation"
          dependencies={['password']}
          rules={[
            { required: true, message: '请输入你的确认密码。' },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('password') === value) {
                  return Promise.resolve();
                }

                return Promise.reject(new Error('你输入的新密码和确认密码不匹配。'));
              },
            }),
          ]}
          required
        />
      </ProForm>
    </>
  );
};

export default Security;