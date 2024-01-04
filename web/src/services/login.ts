import request from '@/request';

export async function login(data: Record<string, any>) {
  return request<API.User & { token: string }>('/api/login', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function getWxQRCode() {
  return request<{ token: string; img: string; expiresIn: number }>('/api/login/wx_qrcode', { method: 'POST' });
}