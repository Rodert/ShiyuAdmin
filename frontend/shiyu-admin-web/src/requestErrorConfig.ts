import type { RequestOptions } from '@@/plugin-request/request';
import type { RequestConfig } from '@umijs/max';
import { message, notification } from 'antd';

// 与后端约定的响应数据格式（统一 code / data / message）
interface ResponseStructure<T = any> {
  code: number;
  data: T;
  message?: string;
}

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
      const { code, data, message } = res as ResponseStructure;
      if (code !== 200) {
        const error: any = new Error(message || '请求错误');
        error.name = 'BizError';
        error.info = { errorCode: code, errorMessage: message, data };
        throw error;
      }
    },
    // 错误接收及处理
    errorHandler: (error: any, opts: any) => {
      if (opts?.skipErrorHandler) throw error;
      // 我们的 errorThrower 抛出的错误。
      if (error.name === 'BizError') {
        const errorInfo = error.info as {
          errorCode?: number;
          errorMessage?: string;
        };
        if (errorInfo) {
          const { errorMessage, errorCode } = errorInfo;
          notification.error({
            message: `错误码：${errorCode}`,
            description: errorMessage || '请求失败',
          });
        }
      } else if (error.response) {
        // Axios 的错误
        // 请求成功发出且服务器也响应了状态码，但状态代码超出了 2xx 的范围
        message.error(`Response status:${error.response.status}`);
      } else if (error.request) {
        // 请求已经成功发起，但没有收到响应
        // \`error.request\` 在浏览器中是 XMLHttpRequest 的实例，
        // 而在node.js中是 http.ClientRequest 的实例
        message.error('None response! Please retry.');
      } else {
        // 发送请求时出了点问题
        message.error('Request error, please retry.');
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
