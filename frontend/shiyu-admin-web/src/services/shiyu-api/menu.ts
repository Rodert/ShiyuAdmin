import { request } from '@umijs/max';

export interface Menu {
  menu_code: string;
  parent_code: string;
  menu_type: string; // M=目录, C=菜单, F=按钮
  menu_name: string;
  perms: string;
  path: string;
  component: string;
  status: number;
  children?: Menu[];
}

export interface CreateMenuRequest {
  menu_code: string;
  parent_code?: string;
  menu_type: string;
  menu_name?: string;
  perms?: string;
  path?: string;
  component?: string;
  status?: number;
}

export interface UpdateMenuRequest {
  parent_code?: string;
  menu_type?: string;
  menu_name?: string;
  perms?: string;
  path?: string;
  component?: string;
  status?: number;
}

/** 获取菜单列表 */
export async function getMenuList() {
  return request<{
    code: number;
    data: Menu[];
    message?: string;
  }>('/api/v1/system/menus', {
    method: 'GET',
  });
}

/** 获取菜单树 */
export async function getMenuTree() {
  return request<{
    code: number;
    data: Menu[];
    message?: string;
  }>('/api/v1/system/menus/tree', {
    method: 'GET',
  });
}

/** 获取菜单详情 */
export async function getMenu(menuCode: string) {
  return request<{
    code: number;
    data: Menu;
    message?: string;
  }>(`/api/v1/system/menus/${menuCode}`, {
    method: 'GET',
  });
}

/** 创建菜单 */
export async function createMenu(data: CreateMenuRequest) {
  return request<{
    code: number;
    data: Menu;
    message?: string;
  }>('/api/v1/system/menus', {
    method: 'POST',
    data,
  });
}

/** 更新菜单 */
export async function updateMenu(menuCode: string, data: UpdateMenuRequest) {
  return request<{
    code: number;
    data: Menu;
    message?: string;
  }>(`/api/v1/system/menus/${menuCode}`, {
    method: 'PUT',
    data,
  });
}

/** 删除菜单 */
export async function deleteMenu(menuCode: string) {
  return request<{
    code: number;
    data: { deleted: boolean };
    message?: string;
  }>(`/api/v1/system/menus/${menuCode}`, {
    method: 'DELETE',
  });
}

