import {
  actOnRequest,
  createRequest,
  downloadAttachment,
  getPeople,
  getRequests,
  getWorkDetail,
  notifyWorkChanged,
  uploadAttachment,
  type Person,
  type WorkDetail,
  type WorkItem,
} from "@/services/shiyu-api/work";
import { formatDateTime } from "@/utils/dateTime";
import { PageContainer } from "@ant-design/pro-components";
import { history, useLocation, useModel } from "@umijs/max";
import {
  Alert,
  Button,
  Card,
  Descriptions,
  Drawer,
  Empty,
  Form,
  Input,
  List,
  message,
  Modal,
  Select,
  Space,
  Spin,
  Table,
  Tabs,
  Tag,
  Timeline,
  Typography,
  Upload,
} from "antd";
import { useEffect, useRef, useState } from "react";

const statusLabels: Record<string, string> = {
  pending: "待审批",
  approved: "已通过",
  rejected: "已驳回",
  withdrawn: "已撤回",
  submitted: "提交申请",
};
const colors: Record<string, string> = {
  pending: "processing",
  approved: "success",
  rejected: "error",
  withdrawn: "default",
};
export default function Requests() {
  const { initialState } = useModel("@@initialState");
  const user = initialState?.currentUser?.userid;
  const timezone = initialState?.installation?.timezone;
  const location = useLocation();
  const params = new URLSearchParams(location.search);
  const requestedView = params.get("view") || "mine";
  const view = ["mine", "pending", "handled"].includes(requestedView)
    ? requestedView
    : "mine";
  const selectedId = params.get("id");
  const [rows, setRows] = useState<WorkItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);
  const [detail, setDetail] = useState<WorkDetail>();
  const [detailError, setDetailError] = useState(false);
  const [creating, setCreating] = useState(false);
  const [saving, setSaving] = useState(false);
  const [people, setPeople] = useState<Person[]>([]);
  const [files, setFiles] = useState<{ file_code: string; name: string }[]>([]);
  const [uploading, setUploading] = useState(false);
  const [action, setAction] = useState<string>();
  const [comment, setComment] = useState("");
  const [revision, setRevision] = useState(0);
  const [form] = Form.useForm();
  const searchRevision = useRef(0);
  const searchTimer = useRef<ReturnType<typeof setTimeout> | undefined>(
    undefined
  );
  const changeView = (next: string) => {
    setPage(1);
    history.push(`/work/requests?view=${next}`);
  };
  const reload = () => {
    setRevision((n) => n + 1);
    notifyWorkChanged();
  };
  useEffect(() => {
    let active = true;
    setLoading(true);
    setError(false);
    getRequests({ view, page, page_size: 20 })
      .then(({ data }) => {
        if (active) {
          setRows(data.items);
          setTotal(data.total);
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
  }, [view, page, revision]);
  useEffect(() => {
    let active = true;
    setDetail(undefined);
    setDetailError(false);
    if (selectedId)
      getWorkDetail(selectedId)
        .then(({ data }) => {
          if (active) setDetail(data);
        })
        .catch(() => {
          if (active) setDetailError(true);
        });
    return () => {
      active = false;
    };
  }, [selectedId, revision]);
  useEffect(() => () => clearTimeout(searchTimer.current), []);
  const searchPeople = (keyword: string) => {
    const sequence = ++searchRevision.current;
    clearTimeout(searchTimer.current);
    searchTimer.current = setTimeout(() => {
      getPeople(keyword)
        .then(({ data }) => {
          if (sequence === searchRevision.current) setPeople(data);
        })
        .catch(() => {});
    }, 250);
  };
  const closeDetail = () => history.push(`/work/requests?view=${view}`);
  return (
    <PageContainer
      title="事项与审批"
      subTitle="提交申请、处理待办，并保留完整审批记录"
    >
      <Card>
        <Space style={{ marginBottom: 16 }}>
          <Button
            type="primary"
            onClick={() => {
              form.resetFields();
              setFiles([]);
              setCreating(true);
              searchPeople("");
            }}
          >
            新建申请
          </Button>
          <Button onClick={reload}>刷新</Button>
        </Space>
        <Tabs
          activeKey={view}
          onChange={changeView}
          items={[
            { key: "mine", label: "我的申请" },
            { key: "pending", label: "待我审批" },
            { key: "handled", label: "已结束的审批" },
          ]}
        />
        {error ? (
          <Alert type="error" message="事项加载失败，请点击刷新重试" />
        ) : (
          <Table<WorkItem>
            rowKey="id"
            loading={loading}
            dataSource={rows}
            scroll={{ x: 800 }}
            pagination={{
              current: page,
              total,
              pageSize: 20,
              showSizeChanger: false,
              onChange: setPage,
            }}
            locale={{ emptyText: <Empty description="暂无事项" /> }}
            columns={[
              {
                title: "事项",
                dataIndex: "title",
                render: (text, row) => (
                  <Button
                    type="link"
                    onClick={() =>
                      history.push(`/work/requests?view=${view}&id=${row.id}`)
                    }
                  >
                    {text}
                  </Button>
                ),
              },
              { title: "申请人", dataIndex: "applicant_name" },
              { title: "审批人", dataIndex: "approver_name" },
              {
                title: "状态",
                dataIndex: "status",
                render: (status) => (
                  <Tag color={colors[status]}>{statusLabels[status]}</Tag>
                ),
              },
              {
                title: "提交时间",
                dataIndex: "created_at",
                render: (value) => formatDateTime(value, timezone),
              },
            ]}
          />
        )}
      </Card>
      <Modal
        title="新建申请"
        open={creating}
        onCancel={() => {
          if (!saving && !uploading) setCreating(false);
        }}
        onOk={() => form.submit()}
        confirmLoading={saving}
        okButtonProps={{ disabled: uploading }}
        okText="提交申请"
        cancelText="取消"
        width={640}
        destroyOnHidden
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={async (values) => {
            if (saving || uploading) return;
            setSaving(true);
            try {
              const { data } = await createRequest({
                ...values,
                file_codes: files.map((f) => f.file_code),
              });
              setCreating(false);
              setPage(1);
              history.push(`/work/requests?view=mine&id=${data.id}`);
              reload();
              message.success("申请已提交，审批人将收到站内通知");
            } catch {
              /* request handler displays the error */
            } finally {
              setSaving(false);
            }
          }}
        >
          <Form.Item
            name="title"
            label="事项标题"
            rules={[{ required: true, whitespace: true }]}
          >
            <Input maxLength={120} showCount />
          </Form.Item>
          <Form.Item
            name="approver_code"
            label="审批人"
            rules={[{ required: true }]}
            extra="输入姓名或账号搜索；只能选择其他启用用户。"
          >
            <Select
              showSearch
              filterOption={false}
              onSearch={searchPeople}
              options={people.map((p) => ({
                value: p.user_code,
                label: `${p.nickname || p.username}（${p.username}）`,
              }))}
              placeholder="选择审批人"
              notFoundContent="暂无匹配用户"
            />
          </Form.Item>
          <Form.Item
            name="content"
            label="申请说明"
            rules={[{ required: true, whitespace: true }]}
          >
            <Input.TextArea rows={5} maxLength={10000} showCount />
          </Form.Item>
          <Form.Item
            label="附件"
            extra={`最多 10 个文件，每个不超过 ${
              initialState?.installation?.max_upload_mb || 2048
            } MiB。提交后附件作为审批凭证保留。`}
          >
            <Upload
              showUploadList={false}
              disabled={uploading || saving || files.length >= 10}
              beforeUpload={async (file) => {
                if (
                  file.size >
                  (initialState?.installation?.max_upload_mb || 2048) *
                    1024 ** 2
                ) {
                  message.error("文件超过当前上传上限");
                  return Upload.LIST_IGNORE;
                }
                setUploading(true);
                try {
                  const { data } = await uploadAttachment(file);
                  setFiles((old) => [
                    ...old,
                    { file_code: data.file_code, name: data.original_name },
                  ]);
                } catch {
                  /* handled globally */
                } finally {
                  setUploading(false);
                }
                return false;
              }}
            >
              <Button loading={uploading}>上传附件</Button>
            </Upload>
            <List
              size="small"
              dataSource={files}
              renderItem={(f) => (
                <List.Item
                  actions={[
                    <Button
                      key="remove"
                      type="link"
                      disabled={saving}
                      onClick={() =>
                        setFiles((old) =>
                          old.filter((x) => x.file_code !== f.file_code)
                        )
                      }
                    >
                      移除
                    </Button>,
                  ]}
                >
                  {f.name}
                </List.Item>
              )}
            />
          </Form.Item>
        </Form>
      </Modal>
      <Drawer
        title="事项详情"
        open={!!selectedId}
        onClose={closeDetail}
        width={680}
      >
        {detailError ? (
          <Alert
            type="error"
            message="事项不存在、无权访问或暂时不可用"
            action={<Button onClick={reload}>重试</Button>}
          />
        ) : !detail ? (
          <Spin />
        ) : (
          <>
            <Typography.Title level={4}>{detail.title}</Typography.Title>
            <Descriptions
              column={1}
              items={[
                {
                  key: "status",
                  label: "状态",
                  children: (
                    <Tag color={colors[detail.status]}>
                      {statusLabels[detail.status]}
                    </Tag>
                  ),
                },
                {
                  key: "applicant",
                  label: "申请人",
                  children: detail.applicant_name,
                },
                {
                  key: "approver",
                  label: "审批人",
                  children: detail.approver_name,
                },
                {
                  key: "created",
                  label: "提交时间",
                  children: formatDateTime(detail.created_at, timezone),
                },
              ]}
            />
            <Typography.Paragraph
              style={{
                whiteSpace: "pre-wrap",
                overflowWrap: "anywhere",
                marginTop: 16,
              }}
            >
              {detail.content}
            </Typography.Paragraph>
            <List
              header="附件"
              dataSource={detail.attachments}
              locale={{ emptyText: "无附件" }}
              renderItem={(f) => (
                <List.Item>
                  <Button
                    type="link"
                    onClick={() =>
                      downloadAttachment(detail.id, f.file_code, f.name).catch(
                        (err) => message.error(err.message)
                      )
                    }
                  >
                    {f.name}
                  </Button>
                </List.Item>
              )}
            />
            {detail.status === "pending" && (
              <Space style={{ margin: "16px 0" }}>
                {detail.approver_code === user && (
                  <>
                    <Button
                      type="primary"
                      onClick={() => {
                        setComment("");
                        setAction("approved");
                      }}
                    >
                      通过申请
                    </Button>
                    <Button
                      danger
                      onClick={() => {
                        setComment("");
                        setAction("rejected");
                      }}
                    >
                      驳回申请
                    </Button>
                  </>
                )}
                {detail.applicant_code === user && (
                  <Button
                    onClick={() => {
                      setComment("");
                      setAction("withdrawn");
                    }}
                  >
                    撤回申请
                  </Button>
                )}
              </Space>
            )}
            <Typography.Title level={5}>审批记录</Typography.Title>
            <Timeline
              items={detail.events.map((e) => ({
                children: (
                  <>
                    <div>
                      {e.actor_name} · {statusLabels[e.action]} ·{" "}
                      {formatDateTime(e.created_at, timezone)}
                    </div>
                    {e.comment && (
                      <Typography.Paragraph style={{ whiteSpace: "pre-wrap" }}>
                        {e.comment}
                      </Typography.Paragraph>
                    )}
                  </>
                ),
              }))}
            />
          </>
        )}
      </Drawer>
      <Modal
        title={
          action === "withdrawn"
            ? "确认撤回申请"
            : action === "rejected"
            ? "驳回申请"
            : "通过申请"
        }
        open={!!action}
        onCancel={() => {
          if (!saving) setAction(undefined);
        }}
        confirmLoading={saving}
        okText="确认处理"
        cancelText="取消"
        onOk={async () => {
          if (!detail || !action || saving) return;
          if (action === "rejected" && !comment.trim()) {
            message.warning("请填写驳回原因");
            return;
          }
          setSaving(true);
          try {
            await actOnRequest(detail.id, action, comment);
            setAction(undefined);
            reload();
            message.success("处理成功");
          } catch {
            reload();
          } finally {
            setSaving(false);
          }
        }}
      >
        <Input.TextArea
          aria-label="处理意见"
          placeholder={
            action === "rejected"
              ? "请填写驳回原因（必填）"
              : "处理意见（选填）"
          }
          value={comment}
          onChange={(e) => setComment(e.target.value)}
          maxLength={1000}
          rows={4}
        />
      </Modal>
    </PageContainer>
  );
}
