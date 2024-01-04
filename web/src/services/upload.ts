import request from '@/request';

export async function upload(files: File[], type: string) {
  return Promise.all(files.map<Promise<{ url: string }>>(file => {
    const formData = new FormData();
    formData.append('type', type);
    formData.append('file', file);

    return request<{ url: string }>('/api/resources/upload', {
      method: 'POST',
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      body: formData,
    } as any).then(({ data }) => data);
  }))
}