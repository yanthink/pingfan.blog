import request from '@/request';

export async function searchArticles(params: Record<string, any>) {
  return request<API.Article[]>(`/api/articles/search?${new URLSearchParams({
    ...params,
    include: 'User,Tags',
  })}`);
}

export async function createArticle(data: Record<string, any>) {
  return request<API.Article>(`/api/articles?`, {
    method: 'POST',
    body: JSON.stringify(data),
    revalidate: { tags: ['articles'] },
  });
}

export async function updateArticle(id: number | string, data: Record<string, any>) {
  return await request<API.Article>(`/api/articles/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
    revalidate: { tags: ['articles'], paths: [`/articles/${id}`, `/articles/${id}/edit`] },
  });
}

export async function likeArticle(id: string | number) {
  return request<API.Article>(`/api/articles/${id}/like`, {
    method: 'POST',
    revalidate: { tags: ['articles'], paths: [`/articles/${id}`] },
  });
}

export async function favoriteArticle(id: string | number) {
  return request<API.Article>(`/api/articles/${id}/favorite`, {
    method: 'POST',
    revalidate: { tags: ['articles'], paths: [`/articles/${id}`] },
  });
}