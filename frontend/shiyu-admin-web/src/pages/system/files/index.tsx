import { deleteFile, getFileBlob, getFiles, restoreFile, uploadFile, type MediaFile } from '@/services/shiyu-api/media';
import { hasPermission } from '@/utils/permission';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Image, message, Modal, Space, Tag, Typography, Upload } from 'antd';
import React, { useRef, useState } from 'react';
import { useModel } from '@umijs/max';

const formatSize = (size: number) => {
  if (size < 1024) return `${size} B`;
  if (size < 1024 ** 2) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / 1024 ** 2).toFixed(1)} MB`;
};

const FilesPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [recycled, setRecycled] = useState(false);
  const [preview, setPreview] = useState<{ file: MediaFile; url: string }>();
  const { initialState } = useModel('@@initialState');
  const canUpload = hasPermission(initialState?.currentUser, 'system:file:upload');
  const canDelete = hasPermission(initialState?.currentUser, 'system:file:delete');
  const openPreview = async (file: MediaFile) => {
    try {
      const blob = await getFileBlob(file.file_code, true);
      if (preview) URL.revokeObjectURL(preview.url);
      setPreview({ file, url: URL.createObjectURL(blob) });
    } catch {
      message.error('该文件暂时无法预览');
    }
  };
  const columns: ProColumns<MediaFile>[] = [
    {
      title: '文件名',
      dataIndex: 'original_name',
      ellipsis: true,
      copyable: true,
      render: (_, record) => (
        <Typography.Link onClick={() => openPreview(record)} title={`预览 ${record.original_name}`}>
          {record.original_name}
        </Typography.Link>
      ),
    },
    { title: '类型', dataIndex: 'mime_type', width: 220, ellipsis: true, render: (_, record) => <Tag>{record.mime_type}</Tag> },
    { title: '大小', dataIndex: 'size', width: 110, search: false, render: (_, record) => formatSize(record.size) },
    { title: '上传时间', dataIndex: 'created_at', width: 180, valueType: 'dateTime', search: false },
    { title: '操作', valueType: 'option', width: 140, render: (_, record) => [
      !recycled && <Button key="download" type="link" onClick={async () => { try { const blob = await getFileBlob(record.file_code); const url = URL.createObjectURL(blob); const link = document.createElement('a'); link.href = url; link.download = record.original_name; link.click(); URL.revokeObjectURL(url); } catch { message.error('下载失败'); } }}>下载</Button>,
      !recycled && canDelete && <Button key="delete" type="link" danger onClick={() => Modal.confirm({ title: '移入回收站', content: `确定删除「${record.original_name}」吗？`, onOk: async () => { await deleteFile(record.file_code); message.success('已移入回收站'); actionRef.current?.reload(); } })}>删除</Button>,
      recycled && <Button key="restore" type="link" onClick={async () => { await restoreFile(record.file_code); message.success('已恢复'); actionRef.current?.reload(); }}>恢复</Button>,
    ].filter(Boolean) },
  ];
  return <PageContainer extra={<Space><Button onClick={() => { setRecycled(!recycled); actionRef.current?.reload(); }}>{recycled ? '返回文件库' : '回收站'}</Button>{canUpload && <Upload showUploadList={false} beforeUpload={async (file) => { if (file.size > 2 * 1024 ** 3) { message.error('单文件不能超过 2 GiB'); return Upload.LIST_IGNORE; } await uploadFile(file); message.success('上传成功'); actionRef.current?.reload(); return false; }}><Button type="primary">上传文件</Button></Upload>}</Space>}>
    <ProTable<MediaFile> actionRef={actionRef} rowKey="file_code" columns={columns} search={{ labelWidth: 80 }} request={async (params) => { const res = await getFiles({ page: params.current, page_size: params.pageSize, keyword: params.original_name as string, recycled }); return { data: res.data?.items || [], total: res.data?.total || 0, success: res.code === 200 }; }} />
    <Modal title={preview?.file.original_name} open={!!preview} footer={null} width={900} onCancel={() => { if (preview) URL.revokeObjectURL(preview.url); setPreview(undefined); }}>
      {preview?.file.mime_type.startsWith('image/') ? <Image src={preview.url} preview={false} style={{ maxHeight: '70vh', display: 'block', margin: 'auto' }} /> : preview?.file.mime_type === 'application/pdf' || preview?.file.mime_type.startsWith('text/') ? <iframe title={preview.file.original_name} src={preview.url} style={{ width: '100%', height: '70vh', border: 0 }} /> : <Typography.Text type="secondary">该格式暂不支持在线预览，请下载后查看。</Typography.Text>}
    </Modal>
  </PageContainer>;
};

export default FilesPage;
