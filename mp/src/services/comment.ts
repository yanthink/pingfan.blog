import request from '@/request';

export async function getComments(params: Record<string, any> = {}) {
  return request.get<API.Comment[]>(`/api/comments/cursor_paginate`, params);
}

export async function getComment(id: number|string, params?: Record<string, any>) {
  return request.get<API.Comment>(`/api/comments/${id}`, params);
}

export async function upvoteComment(id: string | number) {
  return request.post<API.Comment>(`/api/comments/${id}/upvote`);
}

export async function createComment(data: API.Comment) {
  return request.post<API.Comment>('/api/comments', data);
}