import { getInstallation } from "@/services/shiyu-api/installation";
import { formatDateTime } from "@/utils/dateTime";
import { PageContainer } from "@ant-design/pro-components";
import { request, useModel } from "@umijs/max";
import {
  Alert,
  Button,
  Card,
  Form,
  Input,
  InputNumber,
  message,
  Select,
  Space,
  Table,
  Typography,
} from "antd";
import { useEffect, useState } from "react";
type Settings = {
  name: string;
  logo_url: string;
  timezone: string;
  max_upload_mb: number;
};
type Change = {
  id: number;
  actor_code: string;
  before: string;
  after: string;
  created_at: string;
};
const names: Record<string, string> = {
  name: "系统名称",
  logo_url: "Logo",
  timezone: "时区",
  max_upload_mb: "上传上限（MiB）",
};
const describeChange = (row: Change) => {
  try {
    const before = JSON.parse(row.before);
    const after = JSON.parse(row.after);
    return (
      Object.keys(names)
        .filter((key) => before[key] !== after[key])
        .map((key) => `${names[key]}：${before[key]} → ${after[key]}`)
        .join("；") || "设置未改变"
    );
  } catch {
    return "历史记录格式不可识别";
  }
};
export default function SystemSettings() {
  const { initialState, setInitialState } = useModel("@@initialState");
  const [form] = Form.useForm<Settings>();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [saving, setSaving] = useState(false);
  const [changes, setChanges] = useState<Change[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [revision, setRevision] = useState(0);
  useEffect(() => {
    let active = true;
    setLoading(true);
    setError(false);
    Promise.all([
      request<{ data: Settings }>("/api/v1/system/settings"),
      request<{ data: { items: Change[]; total: number } }>(
        "/api/v1/system/settings/changes",
        { params: { page, page_size: 10 } }
      ),
    ])
      .then(([settings, history]) => {
        if (active) {
          form.setFieldsValue(settings.data);
          setChanges(history.data.items);
          setTotal(history.data.total);
        }
      })
      .catch(() => {
        if (active) setError(true);
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
  }, [form, page, revision]);
  return (
    <PageContainer title="系统设置" subTitle="配置团队使用的系统外观与基础规则">
      <Alert
        showIcon
        type="info"
        style={{ marginBottom: 16 }}
        message="以下设置保存后生效，无需重启"
        description="名称与 Logo 用于登录页和导航；时区用于工作台、事项、通知和设置记录，其他页面沿用原有时间显示。数据库连接、存储路径和密钥仍由部署配置管理，修改这些配置需要重启服务。"
      />
      {error ? (
        <Alert
          type="error"
          message="设置加载失败"
          action={
            <Button onClick={() => setRevision((n) => n + 1)}>重试</Button>
          }
        />
      ) : (
        <Card loading={loading}>
          <Form<Settings>
            layout="vertical"
            form={form}
            style={{ maxWidth: 640 }}
            onFinish={async (values) => {
              setSaving(true);
              try {
                await request("/api/v1/system/settings", {
                  method: "PUT",
                  data: values,
                });
                const { data } = await getInstallation();
                await setInitialState((state) => ({
                  ...state,
                  installation: data,
                }));
                setPage(1);
                setRevision((n) => n + 1);
                message.success("系统设置已保存");
              } catch {
              } finally {
                setSaving(false);
              }
            }}
          >
            <Form.Item
              name="name"
              label="系统名称"
              rules={[{ required: true, whitespace: true }]}
            >
              <Input maxLength={64} />
            </Form.Item>
            <Form.Item
              name="logo_url"
              label="Logo 地址"
              rules={[{ required: true }]}
              extra="填写本站绝对路径（如 /logo-v2.png）或公开可访问的 HTTPS 图片地址。"
            >
              <Input maxLength={512} />
            </Form.Item>
            <Form.Item
              name="timezone"
              label="显示时区"
              rules={[{ required: true }]}
            >
              <Select
                showSearch
                options={[
                  "Asia/Shanghai",
                  "Asia/Hong_Kong",
                  "Asia/Tokyo",
                  "Asia/Singapore",
                  "UTC",
                  "Europe/London",
                  "America/New_York",
                  "America/Los_Angeles",
                ].map((value) => ({ value, label: value }))}
              />
            </Form.Item>
            <Form.Item
              name="max_upload_mb"
              label="单文件上传上限（MiB）"
              rules={[{ required: true }]}
              extra="适用于文件管理和事项附件。反向代理如有更小限制，仍需同时调整代理配置。"
            >
              <InputNumber min={1} max={2048} precision={0} />
            </Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" loading={saving}>
                保存设置
              </Button>
              <Typography.Text type="secondary">
                仅超级管理员可以修改
              </Typography.Text>
            </Space>
          </Form>
        </Card>
      )}
      <Card title="设置修改记录" style={{ marginTop: 16 }}>
        <Table<Change>
          rowKey="id"
          size="small"
          dataSource={changes}
          scroll={{ x: 700 }}
          pagination={{
            current: page,
            total,
            pageSize: 10,
            showSizeChanger: false,
            onChange: setPage,
          }}
          columns={[
            {
              title: "修改时间",
              dataIndex: "created_at",
              render: (value) =>
                formatDateTime(value, initialState?.installation?.timezone),
            },
            { title: "操作人编码", dataIndex: "actor_code" },
            { title: "修改内容", render: (_, row) => describeChange(row) },
          ]}
        />
      </Card>
    </PageContainer>
  );
}
