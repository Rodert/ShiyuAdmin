import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { Button, Descriptions, Modal, Tag } from 'antd';
import React, { useRef, useState } from 'react';
import type { OperationLog } from '@/services/shiyu-api/operation_log';
import { getOperationLogList } from '@/services/shiyu-api/operation_log';

const OperationLogPage: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [detailVisible, setDetailVisible] = useState(false);
  const [currentLog, setCurrentLog] = useState<OperationLog | null>(null);

  const columns: ProColumns<OperationLog>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '用户编码',
      dataIndex: 'user_code',
      key: 'user_code',
      width: 120,
      search: true,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      width: 120,
      search: true,
    },
    {
      title: '模块',
      dataIndex: 'module',
      key: 'module',
      width: 140,
      search: true,
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      width: 100,
      search: true,
    },
    {
      title: '方法',
      dataIndex: 'method',
      key: 'method',
      width: 90,
      search: false,
    },
    {
      title: '路径',
      dataIndex: 'path',
      key: 'path',
      width: 220,
      ellipsis: true,
      search: false,
    },
    {
      title: 'IP',
      dataIndex: 'ip',
      key: 'ip',
      width: 140,
      search: false,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 90,
      valueType: 'select',
      valueEnum: {
        1: { text: '成功', status: 'Success' },
        0: { text: '失败', status: 'Error' },
      },
      render: (_: any, record: OperationLog) =>
        record.status === 1 ? (
          <Tag color="green">成功</Tag>
        ) : (
          <Tag color="red">失败</Tag>
        ),
    },
    {
      title: '耗时(ms)',
      dataIndex: 'latency_ms',
      key: 'latency_ms',
      width: 100,
      search: false,
    },
    {
      title: '错误信息',
      dataIndex: 'error_msg',
      key: 'error_msg',
      width: 260,
      ellipsis: true,
      search: false,
    },
    {
      title: '时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      valueType: 'dateTime',
      search: false,
      render: (_, record) => {
        if (!record.created_at) return '-';
        const d = new Date(record.created_at * 1000);
        return d.toLocaleString();
      },
    },
    {
      title: '时间范围',
      dataIndex: 'created_at_range',
      valueType: 'dateTimeRange',
      hideInTable: true,
      search: {
        transform: (value: [string, string]) => {
          const [start, end] = value || [];
          return {
            start_time: start ? Math.floor(new Date(start).getTime() / 1000) : undefined,
            end_time: end ? Math.floor(new Date(end).getTime() / 1000) : undefined,
          };
        },
      },
    },
    {
      title: '操作',
      key: 'option',
      width: 100,
      valueType: 'option',
      render: (_: any, record: OperationLog) => [
        <Button
          key="detail"
          type="link"
          size="small"
          onClick={() => {
            setCurrentLog(record);
            setDetailVisible(true);
          }}
        >
          详情
        </Button>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<OperationLog>
        headerTitle="操作日志"
        actionRef={actionRef}
        rowKey="id"
        search={{ labelWidth: 'auto' }}
        pagination={{
          defaultPageSize: 10,
          showSizeChanger: true,
        }}
        request={async (params) => {
          const res = await getOperationLogList({
            page: params.current || 1,
            page_size: params.pageSize || 10,
            user_code: params.user_code as string,
            username: params.username as string,
            module: params.module as string,
            action: params.action as string,
            status: params.status as number,
            start_time: params.start_time as number,
            end_time: params.end_time as number,
          });
          if (res.code === 200 && res.data) {
            return {
              data: res.data.items,
              success: true,
              total: res.data.total,
            };
          }
          return {
            data: [],
            success: false,
            total: 0,
          };
        }}
        columns={columns}
      />

      <Modal
        title="日志详情"
        open={detailVisible}
        onCancel={() => {
          setDetailVisible(false);
          setCurrentLog(null);
        }}
        footer={null}
        width={700}
      >
        {currentLog && (
          <Descriptions column={1} bordered size="small">
            <Descriptions.Item label="ID">{currentLog.id}</Descriptions.Item>
            <Descriptions.Item label="用户编码">
              {currentLog.user_code || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="用户名">
              {currentLog.username || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="模块">
              {currentLog.module || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="操作">
              {currentLog.action || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="方法">
              {currentLog.method || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="路径">
              {currentLog.path || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="IP">{currentLog.ip || '-'}</Descriptions.Item>
            <Descriptions.Item label="状态">
              {currentLog.status === 1 ? (
                <Tag color="green">成功</Tag>
              ) : (
                <Tag color="red">失败</Tag>
              )}
            </Descriptions.Item>
            <Descriptions.Item label="耗时(ms)">
              {currentLog.latency_ms ?? '-'}
            </Descriptions.Item>
            <Descriptions.Item label="时间">
              {currentLog.created_at
                ? new Date(currentLog.created_at * 1000).toLocaleString()
                : '-'}
            </Descriptions.Item>
            <Descriptions.Item label="错误信息">
              {currentLog.error_msg || '-'}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </PageContainer>
  );
};

export default OperationLogPage;
