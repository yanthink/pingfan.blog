import request from '@/request';

export async function getArticles(params?: Record<string, any>) {
  return request.get<API.Article[]>('/api/articles/cursor_paginate', params);
}

export async function searchArticles(params?: Record<string, any>) {
  return request.get<API.Article[]>('/api/articles/search', params);
}

export async function getArticle(id: number | string, params?: Record<string, any>) {
  return request.get<API.Article>(`/api/articles/${id}`, params);
}


export async function likeArticle(id: string | number) {
  return request.post<API.Article>(`/api/articles/${id}/like`);
}

export async function favoriteArticle(id: string | number) {
  return request.post<API.Article>(`/api/articles/${id}/favorite`);
}