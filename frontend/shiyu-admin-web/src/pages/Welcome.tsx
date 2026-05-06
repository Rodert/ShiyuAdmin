import type { MenuDataItem } from "@ant-design/pro-components";
import { PageContainer } from "@ant-design/pro-components";
import { request, useModel } from "@umijs/max";
import { Card, Col, List, Row, Space, Tag, Typography, theme } from "antd";
import React, { useEffect, useState } from "react";

type SystemStatus = "normal" | "abnormal" | "checking";

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

interface ChangelogVersion {
  version: string;
  items: string[];
}

const formatTime = (timestamp?: number) => {
  if (!timestamp) {
    return "暂无记录";
  }
  return new Date(timestamp * 1000).toLocaleString();
};

const countMenuItems = (menus?: MenuDataItem[]): number => {
  if (!menus?.length) {
    return 0;
  }
  return menus.reduce(
    (total, item) => total + 1 + countMenuItems(item.children),
    0
  );
};

const getSystemStatusText = (status: SystemStatus) => {
  if (status === "checking") {
    return "检测中";
  }
  return status === "normal" ? "正常" : "异常";
};

const getStatusColor = (status?: SystemStatus) => {
  if (status === "normal") {
    return "success";
  }
  if (status === "abnormal") {
    return "error";
  }
  return "processing";
};

const projectFeatures = [
  "内置用户、角色、菜单、部门等后台基础管理能力",
  "提供操作日志、在线用户、服务监控、缓存监控、数据监控等运维视图",
  "支持 JWT 登录认证、权限菜单、按钮权限和超级管理员模式",
  "适合快速搭建中小型管理后台，也可作为 Go + React 全栈项目脚手架",
];

const techStacks = [
  {
    title: "前端",
    items: ["React 19", "Umi Max", "Ant Design Pro", "Ant Design 5", "ECharts"],
  },
  {
    title: "后端",
    items: ["Go", "Gin", "GORM", "JWT", "分层架构"],
  },
  {
    title: "数据与缓存",
    items: ["PostgreSQL", "MySQL", "SQLite", "Redis"],
  },
  {
    title: "工程能力",
    items: ["Docker", "自动迁移", "种子数据", "多环境配置"],
  },
];

const contactItems = [
  {
    label: "QQ 交流群",
    value: "ShiyuAdmin 通用后台管理系统",
    href: "https://qm.qq.com/q/jviU7TCgWQ",
  },
  {
    label: "QQ 交流群 2",
    value: "待补充",
  },
  {
    label: "QQ 交流群 3",
    value: "待补充",
  },
  {
    label: "反馈渠道",
    value: "待补充",
  },
];

const fallbackChangelogItems: ChangelogVersion[] = [
  {
    version: "暂无更新日志",
    items: ["请维护 public/changelog.md。"],
  },
];

const fetchSystemStatus = async (): Promise<SystemStatus> => {
  try {
    const res = await request<SystemHealthResponse>("/api/v1/system/health", {
      method: "GET",
      skipErrorHandler: true,
    });
    return res.code === 200 && res.data?.status === "ok"
      ? "normal"
      : "abnormal";
  } catch (_error) {
    return "abnormal";
  }
};

const parseChangelog = (markdown: string): ChangelogVersion[] => {
  const result: ChangelogVersion[] = [];
  let current: ChangelogVersion | undefined;

  markdown.split(/\r?\n/).forEach((line) => {
    const versionMatch = line.match(/^##\s+(.+?)\s*$/);
    if (versionMatch) {
      current = {
        version: versionMatch[1],
        items: [],
      };
      result.push(current);
      return;
    }

    const itemMatch = line.match(/^-\s+(.+?)\s*$/);
    if (itemMatch && current) {
      current.items.push(itemMatch[1]);
    }
  });

  return result.length > 0 ? result : fallbackChangelogItems;
};

const fetchChangelog = async (): Promise<ChangelogVersion[]> => {
  try {
    const res = await fetch("/changelog.md", { cache: "no-cache" });
    if (!res.ok) {
      return fallbackChangelogItems;
    }
    return parseChangelog(await res.text());
  } catch (_error) {
    return fallbackChangelogItems;
  }
};

const buildWelcomeItems = (
  currentUser: API.CurrentUser | undefined,
  menuData: MenuDataItem[] | undefined,
  systemStatus: SystemStatus
): WelcomeItem[] => {
  const isSuperAdmin = Boolean(currentUser?.isSuperAdmin);
  const permissionCount = countMenuItems(menuData);
  const roleName = isSuperAdmin
    ? "超级管理员"
    : currentUser?.access || "普通用户";
  const permissionText = isSuperAdmin
    ? "全部权限"
    : `菜单权限 ${permissionCount} 项`;

  return [
    {
      key: "role",
      label: "当前角色 / 权限",
      value: roleName,
      description: permissionText,
    },
    {
      key: "systemStatus",
      label: "系统状态",
      value: getSystemStatusText(systemStatus),
      description:
        systemStatus === "normal" ? "核心服务响应正常" : "核心服务检测未通过",
      status: systemStatus,
    },
    {
      key: "lastLogin",
      label: "最近登录时间",
      value: formatTime(currentUser?.loginAt),
      description: "基于当前登录令牌签发时间",
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
        height: "100%",
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
        {item.status ? (
          <Tag color={getStatusColor(item.status)}>{item.value}</Tag>
        ) : (
          item.value
        )}
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
  const { initialState } = useModel("@@initialState");
  const [systemStatus, setSystemStatus] = useState<SystemStatus>("checking");
  const [changelogItems, setChangelogItems] = useState<ChangelogVersion[]>(
    fallbackChangelogItems
  );

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

  useEffect(() => {
    let mounted = true;
    fetchChangelog().then((items) => {
      if (mounted) {
        setChangelogItems(items);
      }
    });
    return () => {
      mounted = false;
    };
  }, []);

  const currentUser = initialState?.currentUser;
  const userName = currentUser?.name || "管理员";
  const overviewItems = buildWelcomeItems(
    currentUser,
    initialState?.menuData,
    systemStatus
  );

  return (
    <PageContainer title="首页">
      <Space direction="vertical" size={20} style={{ width: "100%" }}>
        <Card
          bordered={false}
          style={{
            borderRadius: 18,
            overflow: "hidden",
            background:
              initialState?.settings?.navTheme === "realDark"
                ? "linear-gradient(135deg, #151922 0%, #202633 100%)"
                : "linear-gradient(135deg, #f7fbff 0%, #eef4ff 48%, #fffaf1 100%)",
          }}
          styles={{
            body: {
              padding: 28,
            },
          }}
        >
          <Typography.Title level={2} style={{ marginBottom: 12 }}>
            {userName}，欢迎使用 ShiyuAdmin
          </Typography.Title>
          <Typography.Paragraph
            style={{
              maxWidth: 720,
              color: token.colorTextSecondary,
              fontSize: 15,
              marginBottom: 0,
            }}
          >
            ShiyuAdmin
            是一套通用后台管理系统，面向业务管理、权限配置、系统监控和数据运维等常见后台场景。
          </Typography.Paragraph>
        </Card>

        <Row gutter={[16, 16]}>
          {overviewItems.map((item) => (
            <Col key={item.key} xs={24} md={8}>
              <OverviewCard item={item} />
            </Col>
          ))}
        </Row>

        <Row gutter={[16, 16]}>
          <Col xs={24} lg={12}>
            <Card
              bordered={false}
              title="项目介绍"
              style={{ height: "100%", borderRadius: 12 }}
            >
              <List
                split={false}
                dataSource={projectFeatures}
                renderItem={(item) => (
                  <List.Item style={{ paddingLeft: 0, paddingRight: 0 }}>
                    <Typography.Text>{item}</Typography.Text>
                  </List.Item>
                )}
              />
            </Card>
          </Col>
          <Col xs={24} lg={12}>
            <Card
              bordered={false}
              title="技术选型"
              style={{ height: "100%", borderRadius: 12 }}
            >
              <Space direction="vertical" size={14} style={{ width: "100%" }}>
                {techStacks.map((stack) => (
                  <div key={stack.title}>
                    <Typography.Text strong>{stack.title}</Typography.Text>
                    <div style={{ marginTop: 8 }}>
                      {stack.items.map((item) => (
                        <Tag key={item} style={{ marginBottom: 8 }}>
                          {item}
                        </Tag>
                      ))}
                    </div>
                  </div>
                ))}
              </Space>
            </Card>
          </Col>
        </Row>

        <Row gutter={[16, 16]}>
          <Col xs={24} lg={12}>
            <Card
              bordered={false}
              title="联系信息"
              style={{ height: "100%", borderRadius: 12 }}
            >
              <List
                dataSource={contactItems}
                renderItem={(item) => (
                  <List.Item>
                    <List.Item.Meta
                      title={item.label}
                      description={
                        item.href ? (
                          <Typography.Link
                            href={item.href}
                            target="_blank"
                            rel="noreferrer"
                          >
                            点击链接加入群聊【{item.value}】
                          </Typography.Link>
                        ) : (
                          item.value
                        )
                      }
                    />
                  </List.Item>
                )}
              />
            </Card>
          </Col>
          <Col xs={24} lg={12}>
            <Card
              bordered={false}
              title="更新日志"
              style={{ height: "100%", borderRadius: 12 }}
            >
              <List
                dataSource={changelogItems}
                renderItem={(item) => (
                  <List.Item>
                    <List.Item.Meta
                      title={
                        <Space size={8} wrap>
                          <Tag color="blue">{item.version}</Tag>
                        </Space>
                      }
                      description={
                        <List
                          split={false}
                          dataSource={
                            item.items.length > 0 ? item.items : ["暂无内容"]
                          }
                          renderItem={(entry) => (
                            <List.Item style={{ padding: "2px 0" }}>
                              <Typography.Text>{entry}</Typography.Text>
                            </List.Item>
                          )}
                        />
                      }
                    />
                  </List.Item>
                )}
              />
            </Card>
          </Col>
        </Row>
      </Space>
    </PageContainer>
  );
};

export default Welcome;
