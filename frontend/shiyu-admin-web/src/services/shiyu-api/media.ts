import { request } from '@umijs/max';

export interface MediaFile {
  id: number;
  file_code: string;
  original_name: string;
  mime_type: string;
  size: number;
  sha256: string;
  access_level: string;
  uploader_code: string;
  created_at: string;
}

export async function getFiles(params: { page?: number; page_size?: number; keyword?: string; recycled?: boolean }) {
  return request<{ code: number; data: { items: MediaFile[]; total: number } }>('/api/v1/system/files' + (params.recycled ? '/recycle-bin' : ''), { method: 'GET', params });
}

export async function uploadFile(file: File) {
  const form = new FormData();
  form.append('file', file);
  return request<{ code: number; data: MediaFile; message?: string }>('/api/v1/system/files/upload', { method: 'POST', data: form, requestType: 'form' });
}

export async function deleteFile(code: string) { return request('/api/v1/system/files/' + code, { method: 'DELETE' }); }
export async function restoreFile(code: string) { return request('/api/v1/system/files/' + code + '/restore', { method: 'POST' }); }
export const fileDownloadURL = (code: string) => '/api/v1/system/files/' + code + '/download';
