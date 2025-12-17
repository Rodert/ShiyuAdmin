import { request } from '@umijs/max';

export interface Role {
  role_code: string;
  role_name: string;
  role_key: string;
  data_scope: string;
  status: number;
}

export interface RoleListResponse {
  items: Role[];
  page: number;
  size: number;
  total: number;
}

export interface CreateRoleRequest {
  role_code: string;
  role_name: string;
  role_key: string;
  data_scope?: string;
  status?: number;
}

export interface UpdateRoleRequest {
  role_name?: string;
  role_key?: string;
  data_scope?: string;
  status?: number;
}

/** 获取角色列表 */
export async function getRoleList(params: {
  page?: number;
  page_size?: number;
}) {
  return request<{
    code: number;
    data: RoleListResponse;
    message?: string;
  }>('/api/v1/system/roles', {
    method: 'GET',
    params,
  });
}

/** 获取角色详情 */
export async function getRole(roleCode: string) {
  return request<{
    code: number;
    data: Role;
    message?: string;
  }>(`/api/v1/system/roles/${roleCode}`, {
    method: 'GET',
  });
}

/** 创建角色 */
export async function createRole(data: CreateRoleRequest) {
  return request<{
    code: number;
    data: Role;
    message?: string;
  }>('/api/v1/system/roles', {
    method: 'POST',
    data,
  });
}

/** 更新角色 */
export async function updateRole(roleCode: string, data: UpdateRoleRequest) {
  return request<{
    code: number;
    data: Role;
    message?: string;
  }>(`/api/v1/system/roles/${roleCode}`, {
    method: 'PUT',
    data,
  });
}

/** 删除角色 */
export async function deleteRole(roleCode: string) {
  return request<{
    code: number;
    data: { deleted: boolean };
    message?: string;
  }>(`/api/v1/system/roles/${roleCode}`, {
    method: 'DELETE',
  });
}

/** 获取角色菜单 */
export async function getRoleMenus(roleCode: string) {
  return request<{
    code: number;
    data: any[];
    message?: string;
  }>(`/api/v1/system/roles/${roleCode}/menus`, {
    method: 'GET',
  });
}

/** 设置角色菜单 */
export async function setRoleMenus(roleCode: string, menuCodes: string[]) {
  return request<{
    code: number;
    data: { updated: boolean };
    message?: string;
  }>(`/api/v1/system/roles/${roleCode}/menus`, {
    method: 'PUT',
    data: { menu_codes: menuCodes },
  });
}

