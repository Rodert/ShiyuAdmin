import {
  InfoCircleOutlined,
  LockOutlined,
  LogoutOutlined,
  UserOutlined,
} from '@ant-design/icons';
import { outLogin } from '@/services/ant-design-pro/api';
import { history, useModel } from '@umijs/max';
import type { MenuProps } from 'antd';
import { Button, Descriptions, Modal, Spin, Tag, Tooltip } from 'antd';
import { createStyles } from 'antd-style';
import React from 'react';
import { flushSync } from 'react-dom';
import HeaderDropdown from '../HeaderDropdown';

export type GlobalHeaderRightProps = {
  menu?: boolean;
  children?: React.ReactNode;
};

export const AvatarName = () => {
  const { initialState } = useModel('@@initialState');
  const { currentUser } = initialState || {};
  return <span className="anticon">{currentUser?.name}</span>;
};

const formatLoginTime = (loginAt?: number) => {
  if (!loginAt) {
    return '-';
  }
  const timestamp = loginAt > 10_000_000_000 ? loginAt : loginAt * 1000;
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).format(new Date(timestamp));
};

const useStyles = createStyles(({ token }) => {
  return {
    action: {
      display: 'flex',
      height: '48px',
      marginLeft: 'auto',
      overflow: 'hidden',
      alignItems: 'center',
      padding: '0 8px',
      cursor: 'pointer',
      borderRadius: token.borderRadius,
      '&:hover': {
        backgroundColor: token.colorBgTextHover,
      },
    },
    iconAction: {
      display: 'inline-flex',
      alignItems: 'center',
      justifyContent: 'center',
      width: 40,
      height: 40,
      color: token.colorText,
    },
  };
});

export const UserProfileAction: React.FC = () => {
  const { styles } = useStyles();
  const { initialState } = useModel('@@initialState');
  const [open, setOpen] = React.useState(false);
  const currentUser = initialState?.currentUser;

  if (!currentUser) {
    return null;
  }

  const roleText = currentUser.isSuperAdmin
    ? '超级管理员'
    : currentUser.access || '普通用户';

  return (
    <>
      <Tooltip title="个人信息">
        <Button
          aria-label="查看个人信息"
          className={styles.iconAction}
          icon={<InfoCircleOutlined />}
          type="text"
          onClick={() => setOpen(true)}
        />
      </Tooltip>
      <Modal
        title="个人信息"
        open={open}
        footer={null}
        width={520}
        onCancel={() => setOpen(false)}
      >
        <Descriptions column={1} size="middle" bordered>
          <Descriptions.Item label="用户名">
            {currentUser.name || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="昵称">
            {currentUser.nickname || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="邮箱">
            {currentUser.email || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="手机号">
            {currentUser.phone || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="部门">
            {currentUser.deptName || currentUser.deptCode || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="用户编码">
            {currentUser.userid || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="账号类型">
            <Tag color={currentUser.isSuperAdmin ? 'blue' : 'default'}>
              {roleText}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="权限标识">
            {currentUser.access || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="角色">
            {currentUser.roles && currentUser.roles.length > 0
              ? currentUser.roles.map((role) => (
                  <Tag key={role.role_code || role.role_key}>
                    {role.role_name || role.role_key || role.role_code}
                  </Tag>
                ))
              : '-'}
          </Descriptions.Item>
          <Descriptions.Item label="最近登录">
            {formatLoginTime(currentUser.loginAt)}
          </Descriptions.Item>
        </Descriptions>
      </Modal>
    </>
  );
};

export const AvatarDropdown: React.FC<GlobalHeaderRightProps> = ({
  menu,
  children,
}) => {
  /**
   * 退出登录，并且将当前的 url 保存
   */
  const loginOut = async () => {
    try {
      await outLogin({ skipErrorHandler: true });
    } catch (_error) {
      // 即使后端暂时不可用，也要清理本地登录态。
    }
    localStorage.removeItem('shiyu_token');

    const { search, pathname } = window.location;
    const urlParams = new URL(window.location.href).searchParams;
    const searchParams = new URLSearchParams({
      redirect: pathname + search,
    });
    /** 此方法会跳转到 redirect 参数所在的位置 */
    const redirect = urlParams.get('redirect');
    // Note: There may be security issues, please note
    if (window.location.pathname !== '/user/login' && !redirect) {
      history.replace({
        pathname: '/user/login',
        search: searchParams.toString(),
      });
    }
  };
  const { styles } = useStyles();

  const { initialState, setInitialState } = useModel('@@initialState');

  const onMenuClick: MenuProps['onClick'] = (event) => {
    const { key } = event;
    if (key === 'logout') {
      flushSync(() => {
        setInitialState((s) => ({ ...s, currentUser: undefined }));
      });
      loginOut();
      return;
    }
    history.push(`/account/${key}`);
  };

  const loading = (
    <span className={styles.action}>
      <Spin
        size="small"
        style={{
          marginLeft: 8,
          marginRight: 8,
        }}
      />
    </span>
  );

  if (!initialState) {
    return loading;
  }

  const { currentUser } = initialState;

  if (!currentUser?.name) {
    return loading;
  }

  const menuItems = [
    ...(menu
      ? [
          {
            key: 'center',
            icon: <UserOutlined />,
            label: '个人中心',
          },
          {
            key: 'password',
            icon: <LockOutlined />,
            label: '修改密码',
          },
          {
            type: 'divider' as const,
          },
        ]
      : []),
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
    },
  ];

  return (
    <HeaderDropdown
      menu={{
        selectedKeys: [],
        onClick: onMenuClick,
        items: menuItems,
      }}
    >
      {children}
    </HeaderDropdown>
  );
};
