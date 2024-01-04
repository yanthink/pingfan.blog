import request from '@/request';

export async function sendEmailCaptcha(data: Record<string, any>) {
  return request<{ key: string }>('/api/captcha/email', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}