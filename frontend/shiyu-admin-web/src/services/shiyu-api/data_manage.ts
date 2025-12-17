import { request } from '@umijs/max';

export interface TableMeta {
  table_name: string;
  table_type: string;
}

export interface ColumnMeta {
  column_name: string;
  data_type: string;
  is_nullable: boolean;
  max_length?: number;
  column_default?: string;
}

export interface TableDataPage {
  items: Record<string, any>[];
  page: number;
  size: number;
  total: number;
}

/** 获取所有表列表 */
export async function getTables() {
  return request<{
    code: number;
    data: TableMeta[];
    message?: string;
  }>('/api/v1/system/data/tables', {
    method: 'GET',
  });
}

/** 获取指定表的字段列表 */
export async function getTableColumns(table: string) {
  return request<{
    code: number;
    data: ColumnMeta[];
    message?: string;
  }>(`/api/v1/system/data/tables/${table}/columns`, {
    method: 'GET',
  });
}

/** 分页获取指定表的数据 */
export async function getTableRows(table: string, params?: { page?: number; page_size?: number }) {
  return request<{
    code: number;
    data: TableDataPage;
    message?: string;
  }>(`/api/v1/system/data/tables/${table}/rows`, {
    method: 'GET',
    params,
  });
}
