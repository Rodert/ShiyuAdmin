import type { ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProCard, ProTable } from '@ant-design/pro-components';
import React from 'react';
import type { CacheStats, OnlineUser } from '@/services/shiyu-api/monitor';
import { getCacheStats, getOnlineUsers } from '@/services/shiyu-api/monitor';

const MonitorPage: React.FC = () => {
  const onlineColumns: ProColumns<OnlineUser>[] = [
    {
      title: '用户编码',
      dataIndex: 'user_code',
      key: 'user_code',
      width: 140,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      width: 140,
    },
    {
      title: 'IP',
      dataIndex: 'ip',
      key: 'ip',
      width: 160,
    },
    {
      title: 'User-Agent',
      dataIndex: 'user_agent',
      key: 'user_agent',
      ellipsis: true,
    },
    {
      title: '最后活跃时间',
      dataIndex: 'last_active',
      key: 'last_active',
      width: 200,
      render: (_: any, record: OnlineUser) => {
        if (!record.last_active) return '-';
        const d = new Date(record.last_active * 1000);
        return d.toLocaleString();
      },
    },
  ];

  return (
    <PageContainer>
      <ProCard title="缓存监控" bordered style={{ marginBottom: 16 }}>
        <ProTable<CacheStats>
          rowKey="redis_version"
          search={false}
          options={false}
          pagination={false}
          toolBarRender={false}
          request={async () => {
            const res = await getCacheStats();
            if (res.code === 200 && res.data) {
              return {
                data: [res.data],
                success: true,
              };
            }
            return {
              data: [],
              success: false,
            };
          }}
          columns={[
            {
              title: 'Redis 版本',
              dataIndex: 'redis_version',
              key: 'redis_version',
              render: (text: string | undefined) => text || '-',
            },
            {
              title: '模式',
              dataIndex: 'mode',
              key: 'mode',
              render: (text: string | undefined) =>
                text === 'cluster' ? '集群' : '单机',
            },
            {
              title: '已用内存',
              dataIndex: 'used_memory_human',
              key: 'used_memory_human',
              render: (text: string | undefined) => text || '-',
            },
            {
              title: '键数量',
              dataIndex: 'db_size',
              key: 'db_size',
            },
            {
              title: '连接数',
              dataIndex: 'connected_clients',
              key: 'connected_clients',
            },
            {
              title: '命中率',
              dataIndex: 'hit_rate',
              key: 'hit_rate',
              render: (value: number | undefined) => {
                if (value === undefined || value === null) return '-';
                const percent = Math.round(Number(value) * 10000) / 100;
                return `${percent}%`;
              },
            },
          ]}
        />
      </ProCard>

      <ProCard title="在线用户" bordered>
        <ProTable<OnlineUser>
          headerTitle="在线用户列表"
          rowKey="user_code"
          search={false}
          toolBarRender={false}
          request={async () => {
            const res = await getOnlineUsers();
            if (res.code === 200 && res.data) {
              return {
                data: res.data,
                success: true,
              };
            }
            return {
              data: [],
              success: false,
            };
          }}
          columns={onlineColumns}
        />
      </ProCard>
    </PageContainer>
  );
};

export default MonitorPage;
