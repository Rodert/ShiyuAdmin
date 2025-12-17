import { request } from '@umijs/max';

export interface Dept {
  dept_code: string;
  parent_code: string;
  dept_name: string;
  status: number;
  children?: Dept[];
}

export interface CreateDeptRequest {
  dept_code: string;
  parent_code?: string;
  dept_name: string;
  status?: number;
}

export interface UpdateDeptRequest {
  parent_code?: string;
  dept_name?: string;
  status?: number;
}

/** 获取部门列表 */
export async function getDeptList() {
  return request<{
    code: number;
    data: Dept[];
    message?: string;
  }>('/api/v1/system/depts', {
    method: 'GET',
  });
}

/** 获取部门树 */
export async function getDeptTree() {
  return request<{
    code: number;
    data: Dept[];
    message?: string;
  }>('/api/v1/system/depts/tree', {
    method: 'GET',
  });
}

/** 获取部门详情 */
export async function getDept(deptCode: string) {
  return request<{
    code: number;
    data: Dept;
    message?: string;
  }>(`/api/v1/system/depts/${deptCode}`, {
    method: 'GET',
  });
}

/** 创建部门 */
export async function createDept(data: CreateDeptRequest) {
  return request<{
    code: number;
    data: Dept;
    message?: string;
  }>('/api/v1/system/depts', {
    method: 'POST',
    data,
  });
}

/** 更新部门 */
export async function updateDept(deptCode: string, data: UpdateDeptRequest) {
  return request<{
    code: number;
    data: Dept;
    message?: string;
  }>(`/api/v1/system/depts/${deptCode}`, {
    method: 'PUT',
    data,
  });
}

/** 删除部门 */
export async function deleteDept(deptCode: string) {
  return request<{
    code: number;
    data: { deleted: boolean };
    message?: string;
  }>(`/api/v1/system/depts/${deptCode}`, {
    method: 'DELETE',
  });
}

