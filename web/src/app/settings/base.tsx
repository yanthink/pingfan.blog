import { ProForm, ProFormText, ProFormCaptcha } from '@ant-design/pro-components';
import type { ProFormInstance } from '@ant-design/pro-components';
import { AutoComplete, Avatar, Flex, Upload, Button, message, theme } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import React, { useRef, useState } from 'react';
import { useToken, useUser } from '@/hooks';
import ImgCrop from 'antd-img-crop';
import type { UploadChangeParam } from 'antd/es/upload/interface';
import { sendEmailCaptcha, updateUserProfile } from '@/services';

const AvatarPreview: React.FC<{ value?: string }> = ({ value }) => {
  const { token: themeToken } = theme.useToken();

  return <Avatar src={value} size={140} style={{ boxShadow: themeToken.boxShadow }} />;
};

const Base: React.FC = () => {
  const formRef = useRef<ProFormInstance>();
  const [emails, setEmails] = useState<{ value: string }[]>([]);
  const { user, setUser } = useUser();

  const [token] = useToken();

  return (
    <>
      <style jsx>{`
        .container {
          display: flex;
          padding-top: 12px;

          .left {
            max-width: 448px;
          }

          .right {
            flex: 1;
            padding-left: 104px;
          }
        }


        @media screen and (max-width: 1200px) {
          .container {
            flex-direction: column-reverse;

            .right {
              display: flex;
              flex-direction: column;
              align-items: center;
              max-width: 448px;
              padding: 20px;

              .avatar_title {
                display: none;
              }
            }
          }
        }

        @media screen and (min-width: 1200px) {
          .container {
            .left {
              width: 300px;
            }
          }
        }
      `}</style>
      <ProForm<API.User>
        formRef={formRef}
        initialValues={user}
        onFinish={async values => {
          await updateUserProfile(values);
          setUser({ ...user, ...values });

          message.success('修改成功');
        }}
        submitter={{
          render: (_, [resetBtn, submitBtn]) => submitBtn,
        }}
      >
        <div className="container">
          <div className="left">
            <ProFormText
              label="用户名"
              name="name"
              rules={[
                { required: true, message: '请输入您的用户名!' },
                {
                  pattern: /^(?!_)(?!.*?_$)[a-zA-Z0-9_\u4e00-\u9fa5]{2,10}$/,
                  message: '用户名格式不正确!',
                },
              ]}
              extra="2-10位字符，可包含中文，英文，数字和下划线，不能以下划线开头和结尾。"
            />
            <ProForm.Item
              label="邮箱"
              name="email"
              rules={[{ type: 'email', message: '邮箱格式不正确!' }]}
            >
              <AutoComplete
                options={emails}
                placeholder="Email"
                onSearch={value => setEmails(!value || value.indexOf('@') >= 0 ? [] : [
                  { value: `${value}@qq.com` }, { value: `${value}@163.com` }, { value: `${value}@gmail.com` },
                ])}
                style={{ width: '100%' }}
              />
            </ProForm.Item>
            {formRef.current?.getFieldValue('email') && formRef.current?.getFieldValue('email') !== user.email && (
              <>
                <ProFormCaptcha
                  name="emailCode"
                  phoneName="email"
                  rules={[{ required: true, message: '请输入验证码' }]}
                  placeholder="请输入验证码"
                  onGetCaptcha={async (email) => {
                    const { data } = await sendEmailCaptcha({ email, checkUnique: true });
                    formRef.current?.setFieldValue('emailCodeKey', data.key);
                  }}
                />
                <ProForm.Item name="emailCodeKey" hidden />
              </>
            )}
          </div>
          <div className="right">
            <Flex align="center" vertical style={{ display: 'inline-flex' }}>
              <ProForm.Item label="头像" name="avatar">
                <AvatarPreview />
              </ProForm.Item>
              <ImgCrop quality={.8} rotationSlider showReset onModalCancel={resolve => resolve(false)}>
                <Upload
                  name="file"
                  action={`${process.env.NEXT_PUBLIC_BASE_URL}/api/resources/upload`}
                  accept="image/*"
                  showUploadList={false}
                  headers={{ Authorization: `Bearer ${token}` }}
                  data={{ type: 'avatar' }}
                  onChange={({ file }: UploadChangeParam) => {
                    if (file.status === 'error') {
                      message.error(file.response?.message ?? '上传失败');
                    }

                    if (file.status === 'done') {
                      formRef.current?.setFieldValue('avatar', file.response.data.url);
                    }
                  }}
                  maxCount={1}
                >
                  <Button icon={<UploadOutlined />}>更换头像</Button>
                </Upload>
              </ImgCrop>
            </Flex>
          </div>
        </div>
      </ProForm>
    </>
  );
};

export default Base;