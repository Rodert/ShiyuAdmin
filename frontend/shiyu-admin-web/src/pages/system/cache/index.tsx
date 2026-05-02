import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProCard, ProTable } from '@ant-design/pro-components';
import { Button, Descriptions, Drawer, Space, Tag, Typography } from 'antd';
import React, { useRef, useState } from 'react';
import type {
  CacheDatabase,
  CacheKey,
  CacheValue,
  RedisValueType,
} from '@/services/shiyu-api/cache';
import {
  getCacheDatabases,
  getCacheKeys,
  getCacheValue,
} from '@/services/shiyu-api/cache';
import { renderCellText } from '@/utils/tableRender';

const redisTypes: { label: string; value: RedisValueType }[] = [
  { label: '全部', value: 'all' },
  { label: 'String', value: 'string' },
  { label: 'List', value: 'list' },
  { label: 'Set', value: 'set' },
  { label: 'ZSet', value: 'zset' },
  { label: 'Hash', value: 'hash' },
  { label: 'Stream', value: 'stream' },
];

const typeColorMap: Record<string, string> = {
  string: 'blue',
  list: 'green',
  set: 'cyan',
  zset: 'purple',
  hash: 'orange',
  stream: 'magenta',
  none: 'default',
};

const formatTTL = (ttl?: number) => {
  if (ttl === undefined || ttl === null) return '-';
  if (ttl === -1) return '永久';
  if (ttl === -2) return '不存在';
  if (ttl < 60) return `${ttl} 秒`;
  if (ttl < 3600) return `${Math.floor(ttl / 60)} 分 ${ttl % 60} 秒`;
  const hours = Math.floor(ttl / 3600);
  const minutes = Math.floor((ttl % 3600) / 60);
  return `${hours} 时 ${minutes} 分`;
};

const formatValue = (value: any) => {
  if (value === undefined || value === null) {
    return '';
  }
  if (typeof value === 'string') {
    return value;
  }
  return JSON.stringify(value, null, 2);
};

const CachePage: React.FC = () => {
  const [currentDB, setCurrentDB] = useState(0);
  const [selectedDBKeys, setSelectedDBKeys] = useState<React.Key[]>([0]);
  const [valueDrawerOpen, setValueDrawerOpen] = useState(false);
  const [currentValue, setCurrentValue] = useState<CacheValue | undefined>();
  const keysActionRef = useRef<ActionType>(null);

  const databaseColumns: ProColumns<CacheDatabase>[] = [
    {
      title: '库',
      dataIndex: 'db',
      key: 'db',
      width: 90,
      render: (_, record) => `DB ${record.db}`,
    },
    {
      title: '键数量',
      dataIndex: 'key_count',
      key: 'key_count',
    },
  ];

  const keyColumns: ProColumns<CacheKey>[] = [
    {
      title: 'Key',
      dataIndex: 'key',
      key: 'key',
      ellipsis: true,
      copyable: true,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 120,
      valueType: 'select',
      fieldProps: {
        options: redisTypes,
      },
      render: renderCellText<CacheKey, 'type'>('type', (value) => (
        <Tag color={typeColorMap[value || 'none']}>{value || '-'}</Tag>
      )),
    },
    {
      title: 'TTL',
      dataIndex: 'ttl',
      key: 'ttl',
      width: 140,
      search: false,
      render: renderCellText<CacheKey, 'ttl'>('ttl', formatTTL),
    },
    {
      title: '操作',
      key: 'option',
      valueType: 'option',
      width: 120,
      render: (_, record) => [
        <Button
          key="view"
          type="link"
          onClick={async () => {
            const res = await getCacheValue({ db: currentDB, key: record.key });
            if (res.code === 200 && res.data) {
              setCurrentValue(res.data);
              setValueDrawerOpen(true);
            }
          }}
        >
          查看数据
        </Button>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProCard split="vertical" bordered>
        <ProCard title="Redis 库" colSpan="28%">
          <ProTable<CacheDatabase>
            rowKey="db"
            search={false}
            pagination={false}
            toolBarRender={false}
            columns={databaseColumns}
            request={async () => {
              const res = await getCacheDatabases();
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
            rowSelection={{
              type: 'radio',
              selectedRowKeys: selectedDBKeys,
              onChange: (keys) => {
                const nextDB = Number(keys[0] ?? 0);
                setSelectedDBKeys([nextDB]);
                setCurrentDB(nextDB);
                keysActionRef.current?.reload();
              },
            }}
            onRow={(record) => ({
              onClick: () => {
                setSelectedDBKeys([record.db]);
                setCurrentDB(record.db);
                keysActionRef.current?.reload();
              },
            })}
          />
        </ProCard>
        <ProCard
          title={`键列表（DB ${currentDB}）`}
          extra={
            <Space>
              <Typography.Text type="secondary">
                支持 String/List/Set/ZSet/Hash/Stream
              </Typography.Text>
            </Space>
          }
        >
          <ProTable<CacheKey>
            actionRef={keysActionRef}
            rowKey="key"
            columns={keyColumns}
            params={{ db: currentDB }}
            pagination={{
              pageSize: 20,
            }}
            search={{
              labelWidth: 90,
            }}
            request={async (params) => {
              const res = await getCacheKeys({
                db: currentDB,
                pattern: (params.key as string) || '*',
                type: (params.type as RedisValueType) || 'all',
                count: 1000,
              });
              if (res.code === 200 && res.data) {
                return {
                  data: res.data.items || [],
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
          />
        </ProCard>
      </ProCard>

      <Drawer
        title="缓存数据"
        width={720}
        open={valueDrawerOpen}
        destroyOnClose
        onClose={() => setValueDrawerOpen(false)}
      >
        {currentValue ? (
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <Descriptions column={1} bordered size="small">
              <Descriptions.Item label="库">
                DB {currentValue.db}
              </Descriptions.Item>
              <Descriptions.Item label="Key">
                {currentValue.key}
              </Descriptions.Item>
              <Descriptions.Item label="类型">
                <Tag color={typeColorMap[currentValue.type || 'none']}>
                  {currentValue.type || '-'}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="TTL">
                {formatTTL(currentValue.ttl)}
              </Descriptions.Item>
            </Descriptions>
            <Typography.Paragraph>
              <pre
                style={{
                  maxHeight: 520,
                  overflow: 'auto',
                  padding: 12,
                  margin: 0,
                  background: '#f5f5f5',
                  borderRadius: 6,
                  whiteSpace: 'pre-wrap',
                  wordBreak: 'break-word',
                }}
              >
                {formatValue(currentValue.value)}
              </pre>
            </Typography.Paragraph>
          </Space>
        ) : null}
      </Drawer>
    </PageContainer>
  );
};

export default CachePage;
