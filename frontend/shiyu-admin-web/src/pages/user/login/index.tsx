import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { ProForm, ProFormCheckbox, ProFormText } from '@ant-design/pro-components';
import { Helmet, history, useModel } from '@umijs/max';
import type { MenuDataItem } from '@ant-design/pro-components';
import { Alert, App, Button } from 'antd';
import { createStyles } from 'antd-style';
import React, { useState } from 'react';
import { flushSync } from 'react-dom';
import { login } from '@/services/ant-design-pro/api';
import { getMenuTree, type Menu as APIMenu } from '@/services/shiyu-api/menu';
import Settings from '../../../../config/defaultSettings';
const useStyles = createStyles(() => {
  return {
    container: {
      display: 'flex',
      minHeight: '100vh',
      overflow: 'auto',
      color: '#e5edf7',
      background:
        'radial-gradient(circle at 20% 20%, rgba(59, 130, 246, 0.22), transparent 28%), linear-gradient(135deg, #111827 0%, #020617 100%)',
      '@media (max-width: 768px)': {
        display: 'block',
      },
    },
    infoPanel: {
      position: 'relative',
      flex: 1,
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      padding: '80px clamp(48px, 7vw, 96px)',
      borderRight: '1px solid rgba(148, 163, 184, 0.12)',
      overflow: 'hidden',
      '&::after': {
        position: 'absolute',
        right: '-120px',
        bottom: '-120px',
        width: 320,
        height: 320,
        content: '""',
        background: 'linear-gradient(135deg, rgba(37, 99, 235, 0.28), rgba(14, 165, 233, 0.08))',
        borderRadius: '50%',
        filter: 'blur(2px)',
      },
      '@media (max-width: 768px)': {
        display: 'none',
      },
    },
    heroEyebrow: {
      marginBottom: 20,
      color: '#60a5fa',
      fontSize: 13,
      fontWeight: 700,
      letterSpacing: 2,
      textTransform: 'uppercase',
    },
    heroTitle: {
      maxWidth: 560,
      margin: 0,
      color: '#f8fafc',
      fontSize: 'clamp(36px, 5vw, 58px)',
      fontWeight: 800,
      lineHeight: 1.08,
    },
    heroDesc: {
      maxWidth: 460,
      marginTop: 24,
      color: '#94a3b8',
      fontSize: 17,
      lineHeight: 1.8,
    },
    features: {
      display: 'grid',
      gap: 14,
      marginTop: 42,
      padding: 0,
      color: '#cbd5e1',
      listStyle: 'none',
      '& li': {
        display: 'flex',
        alignItems: 'center',
        gap: 12,
      },
      '& li::before': {
        display: 'inline-flex',
        alignItems: 'center',
        justifyContent: 'center',
        width: 22,
        height: 22,
        color: '#bfdbfe',
        fontWeight: 800,
        content: '"✓"',
        background: 'rgba(37, 99, 235, 0.22)',
        border: '1px solid rgba(96, 165, 250, 0.28)',
        borderRadius: 999,
      },
    },
    loginPanel: {
      flex: 1,
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: 40,
      '@media (max-width: 768px)': {
        minHeight: '100vh',
        padding: 20,
      },
    },
    loginCard: {
      width: '100%',
      maxWidth: 420,
      padding: 32,
      background: 'rgba(15, 23, 42, 0.72)',
      border: '1px solid rgba(148, 163, 184, 0.18)',
      borderRadius: 22,
      boxShadow: '0 24px 80px rgba(2, 6, 23, 0.45)',
      backdropFilter: 'blur(18px)',
      '@media (max-width: 768px)': {
        padding: 24,
      },
    },
    brand: {
      display: 'flex',
      alignItems: 'center',
      gap: 14,
      marginBottom: 34,
    },
    logo: {
      width: 48,
      height: 48,
      padding: 6,
      objectFit: 'contain',
      background: '#f8fafc',
      borderRadius: 14,
      boxShadow: '0 12px 28px rgba(37, 99, 235, 0.24)',
    },
    brandTitle: {
      margin: 0,
      color: '#f8fafc',
      fontSize: 24,
      fontWeight: 800,
      lineHeight: 1.1,
    },
    brandSubtitle: {
      marginTop: 6,
      color: '#94a3b8',
      fontSize: 13,
    },
    loginForm: {
      '& .ant-form-item': {
        marginBottom: 18,
      },
      '& .ant-input-affix-wrapper': {
        padding: '12px 14px',
        color: '#e5edf7',
        background: '#111827',
        borderColor: '#334155',
        borderRadius: 12,
      },
      '& .ant-input-affix-wrapper-focused': {
        borderColor: '#3b82f6',
        boxShadow: '0 0 0 3px rgba(59, 130, 246, 0.16)',
      },
      '& .ant-input': {
        color: '#e5edf7',
        background: 'transparent',
      },
      '& .ant-input::placeholder': {
        color: '#64748b',
      },
      '& .anticon': {
        color: '#64748b',
      },
      '& .ant-checkbox-wrapper': {
        color: '#94a3b8',
      },
    },
    submitButton: {
      height: 46,
      marginTop: 4,
      fontWeight: 700,
      background: '#2563eb',
      borderRadius: 12,
      boxShadow: '0 14px 30px rgba(37, 99, 235, 0.28)',
    },
    copyright: {
      marginTop: 22,
      color: '#64748b',
      fontSize: 12,
      textAlign: 'center',
      '& a': {
        color: '#93c5fd',
      },
      '& a:hover': {
        color: '#bfdbfe',
      },
    },
  };
});
const LoginMessage: React.FC<{
  content: string;
}> = ({ content }) => {
  return (
    <Alert
      style={{
        marginBottom: 24,
      }}
      message={content}
      type="error"
      showIcon
    />
  );
};
const Login: React.FC = () => {
  const [userLoginState, setUserLoginState] = useState<API.LoginResult>({});
  const { initialState, setInitialState } = useModel('@@initialState');
  const { styles } = useStyles();
  const { message } = App.useApp();
  const buildMenuData = (menus: APIMenu[]): MenuDataItem[] => {
    return (menus || [])
      .filter((m) => m.status === 1 && m.menu_type !== 'F')
      .map((m) => ({
        name: m.menu_name,
        path: m.path || `/${m.menu_code}`,
        children: m.children ? buildMenuData(m.children) : undefined,
      }));
  };
  const fetchUserInfoAndMenu = async () => {
    const userInfo = await initialState?.fetchUserInfo?.();
    if (!userInfo) {
      return;
    }
    let menuData: MenuDataItem[] = [];
    try {
      const menuRes = await getMenuTree();
      if (menuRes.code === 200 && Array.isArray(menuRes.data)) {
        menuData = buildMenuData(menuRes.data);
      }
    } catch (_error) {
      menuData = [];
    }
    flushSync(() => {
      setInitialState((s) => ({
        ...s,
        currentUser: userInfo,
        menuData,
      }));
    });
  };
  const handleSubmit = async (values: API.LoginParams) => {
    try {
      // 登录
      const msg = await login({
        ...values,
        type: 'account',
      });
      if (msg.status === 'ok') {
        const defaultLoginSuccessMessage = '登录成功！';
        message.success(defaultLoginSuccessMessage);
        await fetchUserInfoAndMenu();
        const urlParams = new URL(window.location.href).searchParams;
        const redirect = urlParams.get('redirect');
        // 使用 history.push 而不是 window.location.href，这样会考虑 base 配置
        if (redirect) {
          // redirect 可能是绝对路径，需要处理
          const redirectPath = redirect.startsWith('http') ? new URL(redirect).pathname : redirect;
          history.push(redirectPath);
        } else {
          history.push('/welcome');
        }
        return;
      }
      console.log(msg);
      // 如果失败去设置用户错误信息
      setUserLoginState(msg);
    } catch (error) {
      const defaultLoginFailureMessage = '登录失败，请重试！';
      console.log(error);
      message.error(defaultLoginFailureMessage);
    }
  };
  const { status } = userLoginState;
  return (
    <div className={styles.container}>
      <Helmet>
        <title>
          {'登录'}
          {Settings.title && ` - ${Settings.title}`}
        </title>
      </Helmet>
      <section className={styles.infoPanel}>
        <div className={styles.heroEyebrow}>Open Source Admin</div>
        <h1 className={styles.heroTitle}>
          Shiyu Admin
          <br />
          构建下一代后台
        </h1>
        <p className={styles.heroDesc}>
          一款极简、高效的开源通用后台管理系统，内置用户、角色、菜单、部门和操作日志能力，适合作为业务中后台的起点。
        </p>
        <ul className={styles.features}>
          <li>极简美学，专注核心业务</li>
          <li>高性能架构，响应迅速</li>
          <li>完善的权限管理，安全可靠</li>
        </ul>
      </section>

      <section className={styles.loginPanel}>
        <div className={styles.loginCard}>
          <div className={styles.brand}>
            <img className={styles.logo} alt="Shiyu Admin logo" src="/logo.png" />
            <div>
              <h2 className={styles.brandTitle}>Shiyu Admin</h2>
              <div className={styles.brandSubtitle}>仕宇通用管理后台</div>
            </div>
          </div>
          {status === 'error' && (
            <LoginMessage content={'错误的用户名或密码'} />
          )}
          <ProForm
            className={styles.loginForm}
            initialValues={{
              username: 'admin',
              password: 'Admin@123',
              autoLogin: true,
            }}
            submitter={{
              render: (_, doms) => (
                <Button
                  className={styles.submitButton}
                  type="primary"
                  htmlType="submit"
                  block
                  loading={doms?.[1]?.props?.loading}
                >
                  进入控制台
                </Button>
              ),
            }}
            onFinish={async (values) => {
              await handleSubmit(values as API.LoginParams);
            }}
          >
            <ProFormText
              name="username"
              fieldProps={{
                size: 'large',
                prefix: <UserOutlined />,
              }}
              placeholder={'用户名'}
              rules={[
                {
                  required: true,
                  message: '用户名是必填项！',
                },
              ]}
            />
            <ProFormText.Password
              name="password"
              fieldProps={{
                size: 'large',
                prefix: <LockOutlined />,
              }}
              placeholder={'密码'}
              rules={[
                {
                  required: true,
                  message: '密码是必填项！',
                },
              ]}
            />
            <ProFormCheckbox noStyle name="autoLogin">
              自动登录
            </ProFormCheckbox>
          </ProForm>
          <div className={styles.copyright}>
            © 2026 Shiyu Platform |{' '}
            <a href="https://github.com/Rodert/ShiyuAdmin" target="_blank" rel="noreferrer">
              GitHub 开源
            </a>
          </div>
        </div>
      </section>
    </div>
  );
};
export default Login;
