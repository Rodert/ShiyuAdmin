import type { ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProCard, ProTable } from '@ant-design/pro-components';
import React from 'react';
import type { CacheStats, DatabaseStats, OnlineUser } from '@/services/shiyu-api/monitor';
import { getCacheStats, getDatabaseStats, getOnlineUsers } from '@/services/shiyu-api/monitor';
import { renderCellText } from '@/utils/tableRender';

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
      <ProCard title="数据库监控" bordered style={{ marginBottom: 16 }}>
        <ProTable<DatabaseStats>
          rowKey="database"
          search={false}
          options={false}
          pagination={false}
          toolBarRender={false}
          request={async () => {
            const res = await getDatabaseStats();
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
              title: '状态',
              dataIndex: 'status',
              key: 'status',
              valueEnum: {
                normal: { text: '正常', status: 'Success' },
                abnormal: { text: '异常', status: 'Error' },
                disabled: { text: '未启用', status: 'Default' },
              },
            },
            {
              title: '驱动',
              dataIndex: 'driver',
              key: 'driver',
            },
            {
              title: '数据库',
              dataIndex: 'database',
              key: 'database',
              render: renderCellText<DatabaseStats, 'database'>(
                'database',
                (value) => value || '-',
              ),
            },
            {
              title: '版本',
              dataIndex: 'version',
              key: 'version',
              ellipsis: true,
              render: renderCellText<DatabaseStats, 'version'>(
                'version',
                (value) => value || '-',
              ),
            },
            {
              title: '表数量',
              dataIndex: 'table_count',
              key: 'table_count',
            },
            {
              title: '连接数',
              dataIndex: 'open_connections',
              key: 'open_connections',
            },
            {
              title: '使用中',
              dataIndex: 'in_use',
              key: 'in_use',
            },
            {
              title: '空闲',
              dataIndex: 'idle',
              key: 'idle',
            },
          ]}
        />
      </ProCard>

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
              render: renderCellText<CacheStats, 'redis_version'>(
                'redis_version',
                (text) => text || '-',
              ),
            },
            {
              title: '模式',
              dataIndex: 'mode',
              key: 'mode',
              render: renderCellText<CacheStats, 'mode'>('mode', (text) =>
                text === 'cluster' ? '集群' : '单机',
              ),
            },
            {
              title: '已用内存',
              dataIndex: 'used_memory_human',
              key: 'used_memory_human',
              render: renderCellText<CacheStats, 'used_memory_human'>(
                'used_memory_human',
                (text) => text || '-',
              ),
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
              render: renderCellText<CacheStats, 'hit_rate'>('hit_rate', (value) => {
                if (value === undefined || value === null) return '-';
                const percent = Math.round(Number(value) * 10000) / 100;
                return `${percent}%`;
              }),
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
