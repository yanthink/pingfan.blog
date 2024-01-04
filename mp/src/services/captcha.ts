import request from '@/request';

export async function sendEmailCaptcha(data: Record<string, any>) {
  return request.post<{ key: string }>('/api/captcha/email', data);
}