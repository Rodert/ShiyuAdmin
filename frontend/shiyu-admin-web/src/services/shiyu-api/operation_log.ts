import { request } from '@umijs/max';

export interface OperationLog {
  id: number;
  user_code?: string;
  username?: string;
  module?: string;
  action?: string;
  method?: string;
  path?: string;
  ip?: string;
  status: number;
  error_msg?: string;
  latency_ms?: number;
  created_at: number;
}

export interface OperationLogListResponse {
  items: OperationLog[];
  page: number;
  size: number;
  total: number;
}

export interface OperationLogQueryParams {
  page?: number;
  page_size?: number;
  user_code?: string;
  username?: string;
  module?: string;
  action?: string;
  status?: number;
  start_time?: number;
  end_time?: number;
}

/** 获取操作日志列表（分页） */
export async function getOperationLogList(params: OperationLogQueryParams) {
  return request<{
    code: number;
    data: OperationLogListResponse;
    message?: string;
  }>('/api/v1/system/operation-logs', {
    method: 'GET',
    params,
  });
}
