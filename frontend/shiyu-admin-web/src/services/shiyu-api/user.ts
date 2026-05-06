import { request } from "@umijs/max";

export interface User {
  user_code: string;
  username: string;
  nickname: string;
  email?: string;
  phone?: string;
  avatar?: string;
  dept_code?: string;
  status: number;
  is_super_admin?: boolean;
}

export interface UserListResponse {
  items: User[];
  page: number;
  size: number;
  total: number;
}

export interface CreateUserRequest {
  username: string;
  password: string;
  nickname?: string;
  email?: string;
  phone?: string;
  avatar?: string;
  dept_code?: string;
  status?: number;
}

export interface UpdateUserRequest {
  nickname?: string;
  email?: string;
  phone?: string;
  avatar?: string;
  dept_code?: string;
  status?: number;
  password?: string;
}

export interface UpdateProfileRequest {
  nickname?: string;
  email?: string;
  phone?: string;
  avatar?: string;
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

/** 获取用户列表 */
export async function getUserList(
  params: { page?: number; page_size?: number },
  opts?: { skipErrorHandler?: boolean }
) {
  return request<{
    code: number;
    data: UserListResponse;
    message?: string;
  }>("/api/v1/system/users", {
    method: "GET",
    params,
    ...opts,
  });
}

/** 创建用户 */
export async function createUser(data: CreateUserRequest) {
  return request<{
    code: number;
    data: User;
    message?: string;
  }>("/api/v1/system/users", {
    method: "POST",
    data,
  });
}

/** 更新用户 */
export async function updateUser(userCode: string, data: UpdateUserRequest) {
  return request<{
    code: number;
    data: User;
    message?: string;
  }>(`/api/v1/system/users/${userCode}`, {
    method: "PUT",
    data,
  });
}

/** 更新当前登录用户资料 */
export async function updateCurrentProfile(data: UpdateProfileRequest) {
  return request<{
    code: number;
    data: User;
    message?: string;
  }>("/api/v1/system/profile", {
    method: "PUT",
    data,
  });
}

/** 修改当前登录用户密码 */
export async function changeCurrentPassword(data: ChangePasswordRequest) {
  return request<{
    code: number;
    data: { updated: boolean };
    message?: string;
  }>("/api/v1/system/profile/password", {
    method: "PUT",
    data,
  });
}

/** 删除用户 */
export async function deleteUser(userCode: string) {
  return request<{
    code: number;
    data: { deleted: boolean };
    message?: string;
  }>(`/api/v1/system/users/${userCode}`, {
    method: "DELETE",
  });
}

/** 获取用户角色 */
export async function getUserRoles(userCode: string) {
  return request<{
    code: number;
    data: any[];
    message?: string;
  }>(`/api/v1/system/users/${userCode}/roles`, {
    method: "GET",
  });
}

/** 设置用户角色 */
export async function setUserRoles(userCode: string, roleCodes: string[]) {
  return request<{
    code: number;
    data: { updated: boolean };
    message?: string;
  }>(`/api/v1/system/users/${userCode}/roles`, {
    method: "PUT",
    data: { role_codes: roleCodes },
  });
}
