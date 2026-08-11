import { deleteFile, fileDownloadURL, getFiles, restoreFile, uploadFile, type MediaFile } from '@/services/shiyu-api/media';
import { hasPermission } from '@/utils/permission';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, message, Modal, Space, Tag, Upload } from 'antd';
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
  const { initialState } = useModel('@@initialState');
  const canUpload = hasPermission(initialState?.currentUser, 'system:file:upload');
  const canDelete = hasPermission(initialState?.currentUser, 'system:file:delete');
  const columns: ProColumns<MediaFile>[] = [
    { title: '文件名', dataIndex: 'original_name', ellipsis: true, copyable: true },
    { title: '类型', dataIndex: 'mime_type', width: 220, ellipsis: true, render: (_, record) => <Tag>{record.mime_type}</Tag> },
    { title: '大小', dataIndex: 'size', width: 110, search: false, render: (_, record) => formatSize(record.size) },
    { title: '上传时间', dataIndex: 'created_at', width: 180, valueType: 'dateTime', search: false },
    { title: '操作', valueType: 'option', width: 180, render: (_, record) => [
      !recycled && <Button key="download" type="link" href={fileDownloadURL(record.file_code)}>下载</Button>,
      !recycled && canDelete && <Button key="delete" type="link" danger onClick={() => Modal.confirm({ title: '移入回收站', content: `确定删除「${record.original_name}」吗？`, onOk: async () => { await deleteFile(record.file_code); message.success('已移入回收站'); actionRef.current?.reload(); } })}>删除</Button>,
      recycled && <Button key="restore" type="link" onClick={async () => { await restoreFile(record.file_code); message.success('已恢复'); actionRef.current?.reload(); }}>恢复</Button>,
    ].filter(Boolean) },
  ];
  return <PageContainer extra={<Space><Button onClick={() => { setRecycled(!recycled); actionRef.current?.reload(); }}>{recycled ? '返回文件库' : '回收站'}</Button>{canUpload && <Upload showUploadList={false} beforeUpload={async (file) => { if (file.size > 2 * 1024 ** 3) { message.error('单文件不能超过 2 GiB'); return Upload.LIST_IGNORE; } await uploadFile(file); message.success('上传成功'); actionRef.current?.reload(); return false; }}><Button type="primary">上传文件</Button></Upload>}</Space>}>
    <ProTable<MediaFile> actionRef={actionRef} rowKey="file_code" columns={columns} search={{ labelWidth: 80 }} request={async (params) => { const res = await getFiles({ page: params.current, page_size: params.pageSize, keyword: params.original_name as string, recycled }); return { data: res.data?.items || [], total: res.data?.total || 0, success: res.code === 200 }; }} />
  </PageContainer>;
};

export default FilesPage;
