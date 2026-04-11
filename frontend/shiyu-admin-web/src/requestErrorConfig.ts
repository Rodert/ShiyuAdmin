import type { RequestOptions } from '@@/plugin-request/request';
import type { RequestConfig } from '@umijs/max';
import { message, notification } from 'antd';

// 与后端约定的响应数据格式（统一 code / data / message）
interface ResponseStructure<T = any> {
  code: number;
  data: T;
  message?: string;
  error_code?: string;
}

interface ErrorPresentation {
  title: string;
  description: string;
}

type BackendErrorCode =
  | 'bad_request'
  | 'unauthorized'
  | 'forbidden'
  | 'not_found'
  | 'conflict'
  | 'service_unavailable'
  | 'internal_server_error'
  | 'duplicate_user_code'
  | 'duplicate_username'
  | 'user_creation_conflict'
  | 'duplicate_user_identity'
  | 'duplicate_role_code'
  | 'duplicate_role_key'
  | 'role_creation_conflict'
  | 'role_update_conflict'
  | 'duplicate_menu_code'
  | 'duplicate_permission_key'
  | 'menu_creation_conflict'
  | 'menu_update_conflict'
  | 'duplicate_department_code'
  | 'department_creation_conflict'
  | 'department_update_conflict';

const errorPresentations: Record<BackendErrorCode, ErrorPresentation> = {
  bad_request: {
    title: 'Bad Request',
    description: '请求参数不正确，请检查输入内容。',
  },
  unauthorized: {
    title: 'Unauthorized',
    description: '当前登录状态无效或已过期，请重新登录。',
  },
  forbidden: {
    title: 'Forbidden',
    description: '当前账号没有执行该操作的权限。',
  },
  not_found: {
    title: 'Not Found',
    description: '目标数据不存在，或已被删除。',
  },
  conflict: {
    title: 'Conflict',
    description: '存在数据冲突，请检查后重试。',
  },
  service_unavailable: {
    title: 'Service Unavailable',
    description: '服务暂时不可用，请稍后重试。',
  },
  internal_server_error: {
    title: 'Server Error',
    description: '服务端处理失败，请稍后重试。',
  },
  duplicate_user_code: {
    title: 'Duplicate User Code',
    description: '用户编码不能重复，请更换后重试。',
  },
  duplicate_username: {
    title: 'Duplicate Username',
    description: '用户名不能重复，请更换后重试。',
  },
  user_creation_conflict: {
    title: 'User Creation Conflict',
    description: '用户名重复，或系统生成用户编码时发生冲突，请重试。',
  },
  duplicate_user_identity: {
    title: 'Duplicate User Identity',
    description: '用户编码或用户名已存在，请检查后重试。',
  },
  duplicate_role_code: {
    title: 'Duplicate Role Code',
    description: '角色编码不能重复，请更换后重试。',
  },
  duplicate_role_key: {
    title: 'Duplicate Role Key',
    description: '角色标识不能重复，请更换后重试。',
  },
  role_creation_conflict: {
    title: 'Role Creation Conflict',
    description: '角色编码或角色标识已存在，请检查后重试。',
  },
  role_update_conflict: {
    title: 'Role Update Conflict',
    description: '角色编码或角色标识已存在，请检查后重试。',
  },
  duplicate_menu_code: {
    title: 'Duplicate Menu Code',
    description: '菜单编码不能重复，请更换后重试。',
  },
  duplicate_permission_key: {
    title: 'Duplicate Permission Key',
    description: '权限标识不能重复，请更换后重试。',
  },
  menu_creation_conflict: {
    title: 'Menu Creation Conflict',
    description: '菜单编码已存在，请更换后重试。',
  },
  menu_update_conflict: {
    title: 'Menu Update Conflict',
    description: '菜单编码已存在，请更换后重试。',
  },
  duplicate_department_code: {
    title: 'Duplicate Department Code',
    description: '部门编码不能重复，请更换后重试。',
  },
  department_creation_conflict: {
    title: 'Department Creation Conflict',
    description: '部门编码已存在，请更换后重试。',
  },
  department_update_conflict: {
    title: 'Department Update Conflict',
    description: '部门编码已存在，请更换后重试。',
  },
};

const statusPresentationMap: Record<number, ErrorPresentation> = {
  400: {
    title: 'Bad Request',
    description: '请求参数不正确，请检查输入内容。',
  },
  401: {
    title: 'Unauthorized',
    description: '当前登录状态无效或已过期，请重新登录。',
  },
  403: {
    title: 'Forbidden',
    description: '当前账号没有执行该操作的权限。',
  },
  404: {
    title: 'Not Found',
    description: '目标数据不存在，或已被删除。',
  },
  409: {
    title: 'Conflict',
    description: '存在重复数据，请检查编码、标识等字段是否唯一。',
  },
  500: {
    title: 'Server Error',
    description: '服务端处理失败，请稍后重试。',
  },
};

const formatErrorPresentation = (
  statusCode?: number,
  errorMessage?: string,
  backendErrorCode?: string,
): ErrorPresentation => {
  if (backendErrorCode && backendErrorCode in errorPresentations) {
    return errorPresentations[backendErrorCode as BackendErrorCode];
  }

  if (statusCode && statusPresentationMap[statusCode]) {
    const fallback = statusPresentationMap[statusCode];
    return {
      title: fallback.title,
      description: errorMessage ? `${errorMessage}。${fallback.description}` : fallback.description,
    };
  }

  return {
    title: 'Request Failed',
    description: errorMessage ? `${errorMessage}。请检查后重试。` : '请求失败，请稍后重试。',
  };
};

/**
 * @name 错误处理
 * pro 自带的错误处理， 可以在这里做自己的改动
 * @doc https://umijs.org/docs/max/request#配置
 */
export const errorConfig: RequestConfig = {
  // 错误处理： umi@3 的错误处理方案。
  errorConfig: {
    // 错误抛出：code !== 200 视为业务错误
    errorThrower: (res) => {
      const { code, data, message, error_code } = res as ResponseStructure;
      if (code !== 200) {
        const error: any = new Error(message || '请求错误');
        error.name = 'BizError';
        error.info = { statusCode: code, errorCode: error_code, errorMessage: message, data };
        throw error;
      }
    },
    // 错误接收及处理
    errorHandler: (error: any, opts: any) => {
      if (opts?.skipErrorHandler) throw error;
      // 我们的 errorThrower 抛出的错误。
      if (error.name === 'BizError') {
        const errorInfo = error.info as {
          statusCode?: number;
          errorCode?: string;
          errorMessage?: string;
        };
        if (errorInfo) {
          const { errorMessage, errorCode, statusCode } = errorInfo;
          const presentation = formatErrorPresentation(statusCode, errorMessage, errorCode);
          notification.error({
            message: statusCode ? `${presentation.title} (${statusCode})` : presentation.title,
            description: presentation.description,
          });
        }
      } else if (error.response) {
        // Axios 的错误
        // 请求成功发出且服务器也响应了状态码，但状态代码超出了 2xx 的范围
        const presentation = formatErrorPresentation(error.response.status);
        notification.error({
          message: `${presentation.title} (${error.response.status})`,
          description: presentation.description,
        });
      } else if (error.request) {
        // 请求已经成功发起，但没有收到响应
        // \`error.request\` 在浏览器中是 XMLHttpRequest 的实例，
        // 而在node.js中是 http.ClientRequest 的实例
        notification.error({
          message: 'No Response',
          description: '服务暂时没有返回结果，请检查网络或稍后重试。',
        });
      } else {
        // 发送请求时出了点问题
        message.error('Request Failed，请稍后重试');
      }
    },
  },

  // 请求拦截器
  requestInterceptors: [
    (config: RequestOptions) => {
      // 为所有请求自动附加 Authorization 头（如果本地有 token）
      const token = localStorage.getItem('shiyu_token');
      const headers = {
        ...(config.headers || {}),
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      };
      return { ...config, headers };
    },
  ],

  // 响应拦截器
  responseInterceptors: [
    (response) => {
      return response;
    },
  ],
};
