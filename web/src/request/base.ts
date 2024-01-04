export default async function baseRequest<T = unknown, M = Record<string, any>>(url: string, init: RequestInit = {}) {
  const fullURL = url.startsWith('http') ? url : `${process.env.NEXT_PUBLIC_BASE_URL}${url}`;

  const requestInit: Record<string, any> = { ...init };

  requestInit.headers = {
    'Accept': 'application/json',
    'Content-Type': 'application/json; charset=utf-8',
    ...init.headers,
  };

  // multipart/form-data fetch 会自动设置 boundary，
  // 没有 boundary 接口会报 no multipart boundary param in Content-Type 错误
  if (requestInit.headers['Content-Type'] === 'multipart/form-data') {
    delete requestInit.headers['Content-Type'];
  }

  const response = await fetch(fullURL, requestInit);

  let body: RequestResponse<T, M> | undefined;

  try {
    body = await response.json();
  } catch {

  }

  if (!response.ok) {
    return Promise.reject(body);
  }

  return body!;
}
