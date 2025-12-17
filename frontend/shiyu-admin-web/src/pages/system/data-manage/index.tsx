import type { ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProCard, ProTable } from '@ant-design/pro-components';
import React, { useState } from 'react';
import type { TableMeta, ColumnMeta } from '@/services/shiyu-api/data_manage';
import { getTables, getTableColumns, getTableRows } from '@/services/shiyu-api/data_manage';

const DataManagePage: React.FC = () => {
  const [currentTable, setCurrentTable] = useState<string | undefined>();
  const [columnsMeta, setColumnsMeta] = useState<ColumnMeta[]>([]);

  const tableColumns: ProColumns<TableMeta>[] = [
    {
      title: '表名',
      dataIndex: 'table_name',
      key: 'table_name',
    },
    {
      title: '类型',
      dataIndex: 'table_type',
      key: 'table_type',
    },
  ];

  const columnColumns: ProColumns<ColumnMeta>[] = [
    {
      title: '字段名',
      dataIndex: 'column_name',
      key: 'column_name',
    },
    {
      title: '数据类型',
      dataIndex: 'data_type',
      key: 'data_type',
    },
    {
      title: '可为空',
      dataIndex: 'is_nullable',
      key: 'is_nullable',
      render: (value: boolean | undefined) => (value ? '是' : '否'),
    },
    {
      title: '最大长度',
      dataIndex: 'max_length',
      key: 'max_length',
      render: (value: number | undefined) => (value != null ? value : '-'),
    },
    {
      title: '默认值',
      dataIndex: 'column_default',
      key: 'column_default',
      ellipsis: true,
      render: (value: string | undefined) => value || '-',
    },
  ];

  return (
    <PageContainer>
      <ProCard split="vertical" bordered>
        <ProCard title="数据表" colSpan="40%">
          <ProTable<TableMeta>
            rowKey="table_name"
            search={false}
            pagination={false}
            toolBarRender={false}
            request={async () => {
              const res = await getTables();
              if (res.code === 200 && res.data) {
                if (!currentTable && res.data.length > 0) {
                  setCurrentTable(res.data[0].table_name);
                }
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
            columns={tableColumns}
            onRow={(record: TableMeta) => {
              return {
                onClick: () => {
                  setCurrentTable(record.table_name);
                },
              };
            }}
          />
        </ProCard>
        <ProCard title={currentTable ? `表信息（${currentTable}）` : '表信息'} split="horizontal">
          <ProCard title="字段">
            <ProTable<ColumnMeta>
              rowKey="column_name"
              search={false}
              pagination={false}
              toolBarRender={false}
              request={async () => {
                if (!currentTable) {
                  setColumnsMeta([]);
                  return {
                    data: [],
                    success: true,
                  };
                }
                const res = await getTableColumns(currentTable);
                if (res.code === 200 && res.data) {
                  setColumnsMeta(res.data);
                  return {
                    data: res.data,
                    success: true,
                  };
                }
                setColumnsMeta([]);
                return {
                  data: [],
                  success: false,
                };
              }}
              columns={columnColumns}
            />
          </ProCard>
          <ProCard title="数据">
            <ProTable<Record<string, any>>
              rowKey="_row_key"
              search={false}
              pagination={{
                defaultPageSize: 10,
              }}
              toolBarRender={false}
              request={async (params) => {
                if (!currentTable) {
                  return {
                    data: [],
                    success: true,
                    total: 0,
                  };
                }
                const res = await getTableRows(currentTable, {
                  page: params.current || 1,
                  page_size: params.pageSize || 10,
                });
                if (res.code === 200 && res.data) {
                  const itemsWithKey = (res.data.items || []).map((row, index) => ({
                    _row_key:
                      row.id !== undefined && row.id !== null
                        ? String(row.id)
                        : `${res.data.page}-${index}`,
                    ...row,
                  }));
                  return {
                    data: itemsWithKey,
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
              columns={
                columnsMeta.length > 0
                  ? (columnsMeta.map((col) => ({
                      title: col.column_name,
                      dataIndex: col.column_name,
                      key: col.column_name,
                      ellipsis: true,
                    })) as ProColumns<Record<string, any>>[])
                  : [
                      {
                        title: '无字段信息',
                        dataIndex: '_',
                        render: () => '-',
                      } as ProColumns<Record<string, any>>,
                    ]
              }
            />
          </ProCard>
        </ProCard>
      </ProCard>
    </PageContainer>
  );
};

export default DataManagePage;
