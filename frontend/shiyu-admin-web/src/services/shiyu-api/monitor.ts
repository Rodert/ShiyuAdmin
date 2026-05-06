import { request } from "@umijs/max";

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
  session_id: string;
  user_code: string;
  username: string;
  dept_name?: string;
  ip?: string;
  login_location?: string;
  host_ip?: string;
  user_agent?: string;
  os?: string;
  browser?: string;
  browser_version?: string;
  browser_detail?: string;
  login_time?: number;
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
  }>("/api/v1/system/monitor/cache", {
    method: "GET",
    ...opts,
  });
}

/** 获取数据库监控数据 */
export async function getDatabaseStats(opts?: { skipErrorHandler?: boolean }) {
  return request<{
    code: number;
    data: DatabaseStats;
    message?: string;
  }>("/api/v1/system/monitor/database", {
    method: "GET",
    ...opts,
  });
}

/** 获取在线用户列表 */
export async function getOnlineUsers(opts?: { skipErrorHandler?: boolean }) {
  return request<{
    code: number;
    data: OnlineUser[];
    message?: string;
  }>("/api/v1/system/monitor/online-users", {
    method: "GET",
    ...opts,
  });
}

/** 强退在线用户 */
export async function forceLogoutOnlineUser(sessionId: string) {
  return request<{
    code: number;
    data: { forced: boolean };
    message?: string;
  }>(`/api/v1/system/monitor/online-users/${sessionId}`, {
    method: "DELETE",
  });
}
