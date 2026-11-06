import {
  getCleanupCandidates,
  getFileUsage,
  getUnusedAttachments,
  purgeFiles,
  recycleAttachment,
  restoreAttachment,
  type FileUsage,
  type MediaFile,
} from "@/services/shiyu-api/media";
import { formatDateTime } from "@/utils/dateTime";
import { useModel } from "@umijs/max";
import {
  Alert,
  Button,
  Col,
  Drawer,
  message,
  Modal,
  Row,
  Space,
  Statistic,
  Table,
  Tabs,
  Typography,
} from "antd";
import { useEffect, useState } from "react";

export function formatFileSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KiB`;
  if (bytes < 1024 ** 3) return `${(bytes / 1024 ** 2).toFixed(1)} MiB`;
  return `${(bytes / 1024 ** 3).toFixed(2)} GiB`;
}
export function FileUsageCards({
  work = false,
  revision = 0,
}: {
  work?: boolean;
  revision?: number;
}) {
  const [data, setData] = useState<FileUsage>();
  const [failed, setFailed] = useState(false);
  useEffect(() => {
    let active = true;
    setFailed(false);
    getFileUsage(work)
      .then(({ data }) => {
        if (active) setData(data);
      })
      .catch(() => {
        if (active) setFailed(true);
      });
    return () => {
      active = false;
    };
  }, [work, revision]);
  if (failed)
    return <Alert type="error" message="容量统计暂不可用，请刷新重试" />;
  return (
    <>
      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        {(
          [
            ["有效文件", data?.active],
            ["回收站", data?.recycled],
            ["未使用的事项附件", data?.unused],
            ["到期可清理", data?.cleanable],
          ] as const
        ).map(([label, amount]) => (
          <Col xs={12} lg={6} key={label}>
            <Statistic
              title={label}
              value={amount ? formatFileSize(amount.bytes) : "—"}
            />
            <Typography.Text type="secondary">
              {amount?.count ?? "—"} 个文件
            </Typography.Text>
          </Col>
        ))}
      </Row>
      <Typography.Paragraph type="secondary">
        统计已登记文件，不含日志和备份；未使用附件包含在有效文件内。回收站至少保留{" "}
        {data?.retention_days ?? "—"} 天，到期后由管理员选择清理。
      </Typography.Paragraph>
    </>
  );
}
export function UnusedAttachments({
  open,
  onClose,
  onPick,
}: {
  open: boolean;
  onClose: () => void;
  onPick?: (file: MediaFile) => void;
}) {
  const { initialState } = useModel("@@initialState");
  const [recycled, setRecycled] = useState(false);
  const [page, setPage] = useState(1);
  const [revision, setRevision] = useState(0);
  const [items, setItems] = useState<MediaFile[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [failed, setFailed] = useState(false);
  const reload = () => setRevision((n) => n + 1);
  useEffect(() => {
    if (!open) return;
    let active = true;
    setLoading(true);
    setFailed(false);
    getUnusedAttachments({ page, page_size: 10, recycled })
      .then(({ data }) => {
        if (active) {
          setItems(data.items);
          setTotal(data.total);
        }
      })
      .catch(() => {
        if (active) setFailed(true);
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
  }, [open, recycled, page, revision]);
  return (
    <Drawer
      title="我的未使用附件"
      zIndex={1100}
      open={open}
      onClose={onClose}
      width={800}
      destroyOnHidden
    >
      {open && <FileUsageCards work revision={revision} />}
      <Alert
        type="info"
        showIcon
        message="这里只显示自己上传且尚未关联事项的附件。移入回收站后仍可恢复；审批凭证不会出现在清理列表中。"
        style={{ marginBottom: 16 }}
      />
      <Space>
        <Button onClick={reload}>刷新附件</Button>
      </Space>
      <Tabs
        activeKey={recycled ? "trash" : "unused"}
        onChange={(key) => {
          setRecycled(key === "trash");
          setPage(1);
        }}
        items={[
          { key: "unused", label: "未使用" },
          { key: "trash", label: "回收站" },
        ]}
      />
      {failed ? (
        <Alert type="error" message="附件加载失败，请刷新重试" />
      ) : (
        <Table<MediaFile>
          rowKey="file_code"
          dataSource={items}
          loading={loading}
          scroll={{ x: 650 }}
          pagination={{
            current: page,
            total,
            pageSize: 10,
            showSizeChanger: false,
            onChange: setPage,
          }}
          columns={[
            { title: "文件名", dataIndex: "original_name" },
            { title: "大小", dataIndex: "size", render: formatFileSize },
            {
              title: "上传时间",
              dataIndex: "created_at",
              render: (value) =>
                formatDateTime(value, initialState?.installation?.timezone),
            },
            {
              title: "操作",
              render: (_, file) => (
                <Space>
                  {file.purge_pending ? (
                    <Typography.Text type="secondary">
                      清理处理中
                    </Typography.Text>
                  ) : recycled ? (
                    <Button
                      type="link"
                      onClick={async () => {
                        try {
                          await restoreAttachment(file.file_code);
                          reload();
                          message.success("附件已恢复");
                        } catch {}
                      }}
                    >
                      恢复附件
                    </Button>
                  ) : (
                    <>
                      {onPick && (
                        <Button type="link" onClick={() => onPick(file)}>
                          添加至申请
                        </Button>
                      )}
                      <Button
                        type="link"
                        danger
                        onClick={() =>
                          Modal.confirm({
                            title: "附件移入回收站",
                            content: `「${file.original_name}」将进入回收站，在永久清理前可以恢复。`,
                            okText: "移入回收站",
                            cancelText: "取消",
                            onOk: async () => {
                              await recycleAttachment(file.file_code);
                              reload();
                            },
                          })
                        }
                      >
                        移入回收站
                      </Button>
                    </>
                  )}
                </Space>
              ),
            },
          ]}
        />
      )}
    </Drawer>
  );
}
export function CleanupDrawer({
  open,
  onClose,
  onChanged,
}: {
  open: boolean;
  onClose: () => void;
  onChanged: () => void;
}) {
  const { initialState } = useModel("@@initialState");
  const [items, setItems] = useState<MediaFile[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [revision, setRevision] = useState(0);
  const [selected, setSelected] = useState<React.Key[]>([]);
  const [loading, setLoading] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [result, setResult] = useState<string>();
  useEffect(() => {
    if (!open) return;
    let active = true;
    setLoading(true);
    setError("");
    setSelected([]);
    getCleanupCandidates({ page, page_size: 20 })
      .then(({ data }) => {
        if (active) {
          setItems(data.items);
          setTotal(data.total);
        }
      })
      .catch(() => {
        if (active) setError("清理清单加载失败，请刷新重试");
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
  }, [open, page, revision]);
  const chosen = items.filter((item) => selected.includes(item.file_code));
  return (
    <Drawer
      title="到期文件清理"
      open={open}
      onClose={() => {
        if (!busy) onClose();
      }}
      width={820}
      destroyOnHidden
    >
      <Alert
        type="warning"
        showIcon
        message="永久清理会删除物理文件，无法从回收站恢复。"
        description="仅列出超过保留期、未被审批引用的本地文件。处理失败的文件保留清理状态，可以重试。"
        style={{ marginBottom: 16 }}
      />
      {result && (
        <Alert type="info" message={result} style={{ marginBottom: 16 }} />
      )}
      {error && (
        <Alert type="error" message={error} style={{ marginBottom: 16 }} />
      )}
      <Space style={{ marginBottom: 16 }}>
        <Button disabled={busy} onClick={() => setRevision((n) => n + 1)}>
          刷新清单
        </Button>
        <Button
          danger
          disabled={!chosen.length || busy || loading}
          loading={busy}
          onClick={() =>
            Modal.confirm({
              title: `永久清理所选 ${chosen.length} 个文件？`,
              content: `将删除 ${formatFileSize(
                chosen.reduce((sum, f) => sum + f.size, 0)
              )} 文件内容，此操作不可恢复。`,
              okText: "永久清理",
              cancelText: "取消",
              okButtonProps: { danger: true },
              onOk: async () => {
                setBusy(true);
                try {
                  const { data } = await purgeFiles(
                    chosen.map((f) => f.file_code)
                  );
                  const failed = data.filter((r) => !r.purged);
                  setResult(
                    `已清理 ${data.length - failed.length} 个，未完成 ${
                      failed.length
                    } 个${
                      failed.length
                        ? `：${failed.map((r) => r.message).join("；")}`
                        : ""
                    }`
                  );
                  setRevision((n) => n + 1);
                  onChanged();
                } finally {
                  setBusy(false);
                }
              },
            })
          }
        >
          清理所选文件
        </Button>
      </Space>
      <Table<MediaFile>
        rowKey="file_code"
        loading={loading}
        dataSource={error ? [] : items}
        rowSelection={{
          selectedRowKeys: selected,
          onChange: setSelected,
          getCheckboxProps: () => ({ disabled: busy }),
        }}
        pagination={{
          current: page,
          total,
          pageSize: 20,
          showSizeChanger: false,
          onChange: setPage,
        }}
        scroll={{ x: 650 }}
        columns={[
          { title: "文件名", dataIndex: "original_name" },
          { title: "大小", dataIndex: "size", render: formatFileSize },
          {
            title: "进入回收站",
            dataIndex: "deleted_at",
            render: (v) =>
              v ? formatDateTime(v, initialState?.installation?.timezone) : "—",
          },
          {
            title: "状态",
            dataIndex: "purge_pending",
            render: (v) => (v ? "待重试清理" : "已超过保留期"),
          },
        ]}
      />
    </Drawer>
  );
}
