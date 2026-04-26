import { PageContainer } from '@ant-design/pro-components';
import type { MenuDataItem } from '@ant-design/pro-components';
import { request, useModel } from '@umijs/max';
import { Card, Col, Row, Space, Tag, Typography, theme } from 'antd';
import React, { useEffect, useState } from 'react';

type SystemStatus = 'normal' | 'abnormal' | 'checking';

interface WelcomeItem {
  key: string;
  label: string;
  value: React.ReactNode;
  description?: React.ReactNode;
  status?: SystemStatus;
}

interface SystemHealthResponse {
  code: number;
  data?: {
    status?: string;
    time?: number;
  };
}

const formatTime = (timestamp?: number) => {
  if (!timestamp) {
    return '暂无记录';
  }
  return new Date(timestamp * 1000).toLocaleString();
};

const countMenuItems = (menus?: MenuDataItem[]): number => {
  if (!menus?.length) {
    return 0;
  }
  return menus.reduce((total, item) => total + 1 + countMenuItems(item.children), 0);
};

const getSystemStatusText = (status: SystemStatus) => {
  if (status === 'checking') {
    return '检测中';
  }
  return status === 'normal' ? '正常' : '异常';
};

const getStatusColor = (status?: SystemStatus) => {
  if (status === 'normal') {
    return 'success';
  }
  if (status === 'abnormal') {
    return 'error';
  }
  return 'processing';
};

const fetchSystemStatus = async (): Promise<SystemStatus> => {
  try {
    const res = await request<SystemHealthResponse>('/api/v1/system/health', {
      method: 'GET',
      skipErrorHandler: true,
    });
    return res.code === 200 && res.data?.status === 'ok' ? 'normal' : 'abnormal';
  } catch (_error) {
    return 'abnormal';
  }
};

const buildWelcomeItems = (
  currentUser: API.CurrentUser | undefined,
  menuData: MenuDataItem[] | undefined,
  systemStatus: SystemStatus,
): WelcomeItem[] => {
  const isSuperAdmin = Boolean(currentUser?.isSuperAdmin);
  const permissionCount = countMenuItems(menuData);
  const roleName = isSuperAdmin ? '超级管理员' : currentUser?.access || '普通用户';
  const permissionText = isSuperAdmin ? '全部权限' : `菜单权限 ${permissionCount} 项`;

  return [
    {
      key: 'role',
      label: '当前角色 / 权限',
      value: roleName,
      description: permissionText,
    },
    {
      key: 'systemStatus',
      label: '系统状态',
      value: getSystemStatusText(systemStatus),
      description: systemStatus === 'normal' ? '核心服务响应正常' : '核心服务检测未通过',
      status: systemStatus,
    },
    {
      key: 'lastLogin',
      label: '最近登录时间',
      value: formatTime(currentUser?.loginAt),
      description: '基于当前登录令牌签发时间',
    },
  ];
};

const OverviewCard: React.FC<{ item: WelcomeItem }> = ({ item }) => {
  const { token } = theme.useToken();

  return (
    <Card
      bordered={false}
      styles={{
        body: {
          minHeight: 132,
          padding: 20,
        },
      }}
      style={{
        height: '100%',
        borderRadius: 14,
        boxShadow: token.boxShadowTertiary,
      }}
    >
      <Typography.Text type="secondary">{item.label}</Typography.Text>
      <div
        style={{
          marginTop: 14,
          marginBottom: 8,
          color: token.colorTextHeading,
          fontSize: 24,
          fontWeight: 600,
          lineHeight: 1.25,
        }}
      >
        {item.status ? <Tag color={getStatusColor(item.status)}>{item.value}</Tag> : item.value}
      </div>
      {item.description && (
        <Typography.Text type="secondary" style={{ fontSize: 13 }}>
          {item.description}
        </Typography.Text>
      )}
    </Card>
  );
};

const Welcome: React.FC = () => {
  const { token } = theme.useToken();
  const { initialState } = useModel('@@initialState');
  const [systemStatus, setSystemStatus] = useState<SystemStatus>('checking');

  useEffect(() => {
    let mounted = true;
    fetchSystemStatus().then((status) => {
      if (mounted) {
        setSystemStatus(status);
      }
    });
    return () => {
      mounted = false;
    };
  }, []);

  const currentUser = initialState?.currentUser;
  const userName = currentUser?.name || '管理员';
  const overviewItems = buildWelcomeItems(currentUser, initialState?.menuData, systemStatus);

  return (
    <PageContainer>
      <Space direction="vertical" size={20} style={{ width: '100%' }}>
        <Card
          bordered={false}
          style={{
            borderRadius: 18,
            overflow: 'hidden',
            background:
              initialState?.settings?.navTheme === 'realDark'
                ? 'linear-gradient(135deg, #151922 0%, #202633 100%)'
                : 'linear-gradient(135deg, #f7fbff 0%, #eef4ff 48%, #fffaf1 100%)',
          }}
          styles={{
            body: {
              padding: 28,
            },
          }}
        >
          <Typography.Title level={2} style={{ marginBottom: 12 }}>
            {userName}，欢迎使用 Shiyu Admin
          </Typography.Title>
          <Typography.Paragraph
            style={{
              maxWidth: 720,
              color: token.colorTextSecondary,
              fontSize: 15,
              marginBottom: 0,
            }}
          >
            仕宇通用管理后台提供用户、角色、菜单、部门、日志、监控等常用后台管理能力。
          </Typography.Paragraph>
        </Card>

        <Row gutter={[16, 16]}>
          {overviewItems.map((item) => (
            <Col key={item.key} xs={24} md={8}>
              <OverviewCard item={item} />
            </Col>
          ))}
        </Row>
      </Space>
    </PageContainer>
  );
};

export default Welcome;
