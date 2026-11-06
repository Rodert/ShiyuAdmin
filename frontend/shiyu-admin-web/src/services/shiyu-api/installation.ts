import { request } from '@umijs/max';

export type InstallationStatus = {
  installed: boolean;
  logo_url: string;
  timezone: string;
  max_upload_mb: number;
  file_retention_days: number;
  name: string;
  checks: { database: boolean; storage: boolean };
};

export const getInstallation = () =>
  request<{ code: number; data: InstallationStatus }>('/api/v1/system/installation', { skipErrorHandler: true });

export const installSystem = (data: { token: string; name: string; username: string; password: string }) =>
  request('/api/v1/system/installation', { method: 'POST', data });

export type VersionInfo = {
  version: string;
  commit: string;
  build_time: string;
  schema_version: number;
  release_notes: string;
  upgrade_guide: string;
};

export const getVersion = () =>
  request<{ code: number; data: VersionInfo }>('/api/v1/system/version');
