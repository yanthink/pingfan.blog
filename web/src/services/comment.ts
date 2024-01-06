import request from '@/request';

export async function getComments(params: Record<string, any>) {
  return request<API.Comment[]>(`/api/comments?${new URLSearchParams(params)}`);
}

export async function createComment(data: API.Comment) {
  return request<API.Comment>('/api/comments', {
    method: 'POST',
    body: JSON.stringify(data),
    revalidate: { paths: [`/articles/${data['articleId']}`] },
  });
}

export async function updateComment(id: number | string, data: Record<string, any>) {
  return request<API.Comment>(`/api/comments/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
    revalidate: { paths: [`/articles/${data['articleId']}`] },
  });
}

export async function upvoteComment(id: string | number) {
  return request<API.Comment>(`/api/comments/${id}/upvote`, {
    method: 'POST',
    revalidate: { paths: [`/articles/${id}`] },
  });
}