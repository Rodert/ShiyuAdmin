// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 获取当前的用户 GET /api/v1/system/profile */
export async function currentUser(options?: { [key: string]: any }) {
  const res = await request<{
    code: number;
    data?: {
      user_code?: string;
      username?: string;
    };
    message?: string;
  }>('/api/v1/system/profile', {
    method: 'GET',
    ...(options || {}),
  });

  if (res.code !== 200 || !res.data) {
    throw new Error(res.message || '获取用户信息失败');
  }

  const user: API.CurrentUser = {
    name: res.data.username,
    userid: res.data.user_code,
    access: 'admin',
  };

  return { data: user };
}

/** 退出登录接口 POST /api/login/outLogin */
export async function outLogin(options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/login/outLogin', {
    method: 'POST',
    ...(options || {}),
  });
}

/** 登录接口 POST /api/v1/system/auth/login */
export async function login(body: API.LoginParams, options?: { [key: string]: any }) {
  const res = await request<{
    code: number;
    data?: {
      token?: string;
      token_type?: string;
      expire_in?: number;
    };
    message?: string;
  }>('/api/v1/system/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: {
      username: body.username,
      password: body.password,
    },
    ...(options || {}),
  });

  if (res.code === 200 && res.data?.token) {
    // 持久化 Token，后续请求拦截器会自动携带 Authorization 头
    localStorage.setItem('shiyu_token', res.data.token);
    const result: API.LoginResult = {
      status: 'ok',
      type: body.type || 'account',
      currentAuthority: 'admin',
    };
    return result;
  }

  const result: API.LoginResult = {
    status: 'error',
    type: body.type || 'account',
    currentAuthority: 'guest',
  };
  return result;
}

/** 此处后端没有提供注释 GET /api/notices */
export async function getNotices(options?: { [key: string]: any }) {
  return request<API.NoticeIconList>('/api/notices', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取规则列表 GET /api/rule */
export async function rule(
  params: {
    // query
    /** 当前的页码 */
    current?: number;
    /** 页面的容量 */
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<API.RuleList>('/api/rule', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 更新规则 PUT /api/rule */
export async function updateRule(options?: { [key: string]: any }) {
  return request<API.RuleListItem>('/api/rule', {
    method: 'POST',
    data: {
      method: 'update',
      ...(options || {}),
    },
  });
}

/** 新建规则 POST /api/rule */
export async function addRule(options?: { [key: string]: any }) {
  return request<API.RuleListItem>('/api/rule', {
    method: 'POST',
    data: {
      method: 'post',
      ...(options || {}),
    },
  });
}

/** 删除规则 DELETE /api/rule */
export async function removeRule(options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/rule', {
    method: 'POST',
    data: {
      method: 'delete',
      ...(options || {}),
    },
  });
}
