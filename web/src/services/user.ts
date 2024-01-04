import request from '@/request';
import { ReadonlyURLSearchParams } from 'next/navigation';

export async function getAuthUser() {
  return request<API.User>('/api/user');
}

export async function updateUserProfile(data: Record<string, any>) {
  return request<API.User>('/api/user', { method: 'PUT', body: JSON.stringify(data) });
}

export async function updateUserMeta(data: Record<string, any>) {
  return request<API.User>('/api/user/meta', { method: 'PUT', body: JSON.stringify(data) });
}

export async function updateUserPassword(data: Record<string, any>) {
  return request<API.User>('/api/user/password', { method: 'PUT', body: JSON.stringify(data) });
}

export async function getUserUnreadNotificationCount() {
  return request<number>('/api/user/unread_notification_count');
}

export async function getUserNotifications(params: ConstructorParameters<typeof URLSearchParams>[0] | ReadonlyURLSearchParams) {
  return request<API.Notification[]>(`/api/user/notifications?${new URLSearchParams(params)}`);
}

export async function userNotificationsMarkAsRead() {
  return request<{ rows: number }>('/api/user/notifications_mark_as_read', { method: 'POST' });
}

export async function getUserFavorites(params: ConstructorParameters<typeof URLSearchParams>[0] | ReadonlyURLSearchParams) {
  return request<API.Favorite[]>(`/api/user/favorites?${new URLSearchParams(params)}`);
}

export async function getUserComments(params: ConstructorParameters<typeof URLSearchParams>[0] | ReadonlyURLSearchParams) {
  return request<API.Comment[]>(`/api/user/comments?${new URLSearchParams(params)}`);
}

export async function getUserLikes(params: ConstructorParameters<typeof URLSearchParams>[0] | ReadonlyURLSearchParams) {
  return request<API.Like[]>(`/api/user/likes?${new URLSearchParams(params)}`);
}

export async function getUserUpvotes(params: ConstructorParameters<typeof URLSearchParams>[0] | ReadonlyURLSearchParams) {
  return request<API.Upvote[]>(`/api/user/upvotes?${new URLSearchParams(params)}`);
}

export async function getUsers(params: Record<string, any>) {
  return request<API.User[]>(`/api/users?${new URLSearchParams(params)}`);
}

export async function updateUser(id: number|string, data: Record<string, any>) {
  return request<API.User>(`/api/users/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}