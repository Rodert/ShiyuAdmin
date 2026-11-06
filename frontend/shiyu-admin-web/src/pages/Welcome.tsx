import {
  getNotifications,
  getRequests,
  type Notice,
  type WorkItem,
} from "@/services/shiyu-api/work";
import { formatDateTime } from "@/utils/dateTime";
import { PageContainer } from "@ant-design/pro-components";
import { history, useModel } from "@umijs/max";
import {
  Alert,
  Button,
  Card,
  Col,
  Empty,
  List,
  Row,
  Space,
  Statistic,
  Typography,
} from "antd";
import { useEffect, useState } from "react";
export default function Welcome() {
  const { initialState } = useModel("@@initialState");
  const [data, setData] = useState<{
    mine: number;
    pending: number;
    unread: number;
    requests: WorkItem[];
    notices: Notice[];
  }>();
  const [error, setError] = useState(false);
  const [revision, setRevision] = useState(0);
  useEffect(() => {
    let active = true;
    setError(false);
    Promise.all([
      getRequests({ view: "mine", page_size: 1 }),
      getRequests({ view: "pending", page_size: 5 }),
      getNotifications({ page_size: 5 }),
    ])
      .then(([mine, pending, notices]) => {
        if (active)
          setData({
            mine: mine.data.total,
            pending: pending.data.total,
            unread: notices.data.unread,
            requests: pending.data.items,
            notices: notices.data.items,
          });
      })
      .catch(() => {
        if (active) setError(true);
      });
    return () => {
      active = false;
    };
  }, [revision]);
  return (
    <PageContainer
      title="工作台"
      subTitle={`你好，${
        initialState?.currentUser?.name || "欢迎回来"
      }。从这里开始今天的工作。`}
      extra={<Button onClick={() => setRevision((n) => n + 1)}>刷新</Button>}
    >
      {error && (
        <Alert
          type="error"
          message="工作台加载失败，请刷新重试"
          style={{ marginBottom: 16 }}
        />
      )}
      <Row gutter={[16, 16]}>
        {[
          {
            title: "待我审批",
            value: data?.pending,
            path: "/work/requests?view=pending",
          },
          {
            title: "我的申请",
            value: data?.mine,
            path: "/work/requests?view=mine",
          },
          {
            title: "未读通知",
            value: data?.unread,
            path: "/work/notifications",
          },
        ].map((item) => (
          <Col xs={24} md={8} key={item.title}>
            <Card loading={!data && !error}>
              <Statistic title={item.title} value={item.value ?? 0} />
              <Button type="link" onClick={() => history.push(item.path)}>
                查看{item.title}
              </Button>
            </Card>
          </Col>
        ))}
        <Col xs={24} lg={12}>
          <Card title="待处理事项" loading={!data && !error}>
            <List
              dataSource={data?.requests || []}
              locale={{ emptyText: <Empty description="暂时没有待审批事项" /> }}
              renderItem={(item) => (
                <List.Item
                  actions={[
                    <Button
                      key="open"
                      type="link"
                      onClick={() =>
                        history.push(
                          `/work/requests?view=pending&id=${item.id}`
                        )
                      }
                    >
                      处理
                    </Button>,
                  ]}
                >
                  <List.Item.Meta
                    title={item.title}
                    description={`${item.applicant_name} · ${formatDateTime(
                      item.created_at,
                      initialState?.installation?.timezone
                    )}`}
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card
            title="最近通知"
            loading={!data && !error}
            extra={
              <Button
                type="link"
                onClick={() => history.push("/work/notifications")}
              >
                全部通知
              </Button>
            }
          >
            <List
              dataSource={data?.notices || []}
              locale={{ emptyText: <Empty description="暂无通知" /> }}
              renderItem={(item) => (
                <List.Item>
                  <List.Item.Meta
                    title={item.title}
                    description={formatDateTime(
                      item.created_at,
                      initialState?.installation?.timezone
                    )}
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
        {initialState?.currentUser?.isSuperAdmin && (
          <Col span={24}>
            <Card title="团队初始化">
              <Typography.Paragraph type="secondary">
                配置系统信息，建立部门与用户，再按岗位分配角色。事项申请与通知对所有启用账号开放，审批数据仅对参与人可见。
              </Typography.Paragraph>
              <Space wrap>
                {[
                  ["系统设置", "/system/settings"],
                  ["部门管理", "/system/dept"],
                  ["用户管理", "/system/user"],
                  ["角色权限", "/system/role"],
                ].map(([title, path]) => (
                  <Button key={path} onClick={() => history.push(path)}>
                    {title}
                  </Button>
                ))}
              </Space>
            </Card>
          </Col>
        )}
      </Row>
    </PageContainer>
  );
}
