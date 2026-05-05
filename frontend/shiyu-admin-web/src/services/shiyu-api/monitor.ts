import { request } from '@umijs/max';

export interface CacheStats {
  redis_version?: string;
  mode?: string;
  used_memory?: number;
  used_memory_human?: string;
  db_size?: number;
  connected_clients?: number;
  keyspace_hits?: number;
  keyspace_misses?: number;
  hit_rate?: number;
}

export interface OnlineUser {
  user_code: string;
  username: string;
  ip?: string;
  user_agent?: string;
  last_active: number;
}

export interface DatabaseStats {
  status?: string;
  driver?: string;
  database?: string;
  version?: string;
  table_count?: number;
  open_connections?: number;
  in_use?: number;
  idle?: number;
}

/** 获取缓存监控数据 */
export async function getCacheStats(opts?: { skipErrorHandler?: boolean }) {
  return request<{
    code: number;
    data: CacheStats;
    message?: string;
  }>('/api/v1/system/monitor/cache', {
    method: 'GET',
    ...opts,
  });
}

/** 获取数据库监控数据 */
export async function getDatabaseStats(opts?: { skipErrorHandler?: boolean }) {
  return request<{
    code: number;
    data: DatabaseStats;
    message?: string;
  }>('/api/v1/system/monitor/database', {
    method: 'GET',
    ...opts,
  });
}

/** 获取在线用户列表 */
export async function getOnlineUsers(opts?: { skipErrorHandler?: boolean }) {
  return request<{
    code: number;
    data: OnlineUser[];
    message?: string;
  }>('/api/v1/system/monitor/online-users', {
    method: 'GET',
    ...opts,
  });
}
