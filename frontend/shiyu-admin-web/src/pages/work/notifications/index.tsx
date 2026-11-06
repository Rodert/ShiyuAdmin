import {
  getNotifications,
  notifyWorkChanged,
  readAllNotifications,
  readNotification,
  type Notice,
} from "@/services/shiyu-api/work";
import { formatDateTime } from "@/utils/dateTime";
import { PageContainer } from "@ant-design/pro-components";
import { history, useModel } from "@umijs/max";
import {
  Alert,
  Button,
  Card,
  List,
  Space,
  Switch,
  Tag,
  Typography,
} from "antd";
import { useEffect, useState } from "react";
export default function Notifications() {
  const { initialState } = useModel("@@initialState");
  const [items, setItems] = useState<Notice[]>([]);
  const [total, setTotal] = useState(0);
  const [unread, setUnread] = useState(false);
  const [count, setCount] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);
  const [revision, setRevision] = useState(0);
  const [saving, setSaving] = useState(false);
  useEffect(() => {
    let active = true;
    setLoading(true);
    setError(false);
    getNotifications({ page, page_size: 20, unread })
      .then(({ data }) => {
        if (active) {
          setItems(data.items);
          setTotal(data.total);
          setCount(data.unread);
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
  }, [page, unread, revision]);
  return (
    <PageContainer title="站内通知" subTitle="审批进展与待办提醒">
      <Card>
        <Space wrap style={{ marginBottom: 24 }}>
          <Typography.Text>未读 {count} 条</Typography.Text>
          <Switch
            checked={unread}
            checkedChildren="仅未读"
            unCheckedChildren="全部消息"
            onChange={(value) => {
              setUnread(value);
              setPage(1);
            }}
          />
          <Button
            loading={saving}
            disabled={count === 0}
            onClick={async () => {
              setSaving(true);
              try {
                await readAllNotifications();
                setPage(1);
                setRevision((n) => n + 1);
                notifyWorkChanged();
              } catch {
              } finally {
                setSaving(false);
              }
            }}
          >
            全部标为已读
          </Button>
          <Button onClick={() => setRevision((n) => n + 1)}>刷新</Button>
        </Space>
        {error ? (
          <Alert type="error" message="通知加载失败，请刷新重试" />
        ) : (
          <List
            loading={loading}
            dataSource={items}
            locale={{ emptyText: "暂无通知" }}
            pagination={{
              current: page,
              total,
              pageSize: 20,
              showSizeChanger: false,
              onChange: setPage,
            }}
            renderItem={(item) => (
              <List.Item
                actions={[
                  <Button
                    key="open"
                    type="link"
                    onClick={async () => {
                      try {
                        await readNotification(item.id);
                        notifyWorkChanged();
                        history.push(`/work/requests?id=${item.request_id}`);
                      } catch {}
                    }}
                  >
                    查看事项
                  </Button>,
                ]}
              >
                <List.Item.Meta
                  title={
                    <Space>
                      {!item.read_at && <Tag color="blue">未读</Tag>}
                      {item.title}
                    </Space>
                  }
                  description={formatDateTime(
                    item.created_at,
                    initialState?.installation?.timezone
                  )}
                />
              </List.Item>
            )}
          />
        )}
      </Card>
    </PageContainer>
  );
}
