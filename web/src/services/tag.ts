import request from '@/request';

export async function getTags(params: Record<string, any>) {
  return request<API.Tag[]>(`/api/tags?${new URLSearchParams(params)}`);
}

export async function createTag(data: Record<string, any>) {
  return request<API.Tag>('/api/tags', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function updateTag(id: number | string, data: Record<string, any>) {
  return request<API.Tag>(`/api/tags/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}