import request from '@/request';

export async function login(code: string, data?: Record<string, any>) {
  return request.post<API.User & { token?: string }>('/api/login/wx', { code, ...data });
}

export async function scanLogin(uuid: string) {
  return request.post('/api/login/wx_scan', { uuid });
}

export async function getUser(params?: Record<string, any>) {
  return request.get<API.User>('/api/user', params);
}

export async function updateUserProfile(data: Record<string, any>) {
  return request.put<API.User>('/api/user', data);
}

export async function getUserUnreadNotificationCount() {
  return request.get<number>('/api/user/unread_notification_count');
}

export async function getUserNotifications(params?: Record<string, any>) {
  return request.get<API.Notification[]>('/api/user/notifications', params);
}

export async function userNotificationsMarkAsRead() {
  return request.post<{ rows: number }>('/api/user/notifications_mark_as_read');
}

export async function getUserFavorites(params?: Record<string, any>) {
  return request.get<API.Favorite[]>('/api/user/favorites', params);
}

export async function getUserComments(params?: Record<string, any>) {
  return request.get<API.Comment[]>('/api/user/comments', params);
}

export async function getUserLikes(params?: Record<string, any>) {
  return request.get<API.Like[]>('/api/user/likes', params);
}

export async function getUserUpvotes(params?: Record<string, any>) {
  return request.get<API.Upvote[]>('/api/user/upvotes', params);
}