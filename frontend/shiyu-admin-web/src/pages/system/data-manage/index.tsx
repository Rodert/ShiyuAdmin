import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProCard, ProTable } from '@ant-design/pro-components';
import { Button, Modal } from 'antd';
import React, { useEffect, useRef, useState } from 'react';
import type { TableMeta, ColumnMeta } from '@/services/shiyu-api/data_manage';
import { getTables, getTableColumns, getTableRows } from '@/services/shiyu-api/data_manage';
import { renderCellText } from '@/utils/tableRender';

const DataManagePage: React.FC = () => {
  const [currentTable, setCurrentTable] = useState<string | undefined>();
  const [columnsMeta, setColumnsMeta] = useState<ColumnMeta[]>([]);
  const [dataModalOpen, setDataModalOpen] = useState(false);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const columnsActionRef = useRef<ActionType>(null);
  const rowsActionRef = useRef<ActionType>(null);

  useEffect(() => {
    columnsActionRef.current?.reload();
  }, [currentTable]);

  useEffect(() => {
    if (dataModalOpen) {
      rowsActionRef.current?.reload();
    }
  }, [currentTable, dataModalOpen]);

  const openDataModal = (tableName: string) => {
    setCurrentTable(tableName);
    setSelectedRowKeys([tableName]);
    setDataModalOpen(true);
  };

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
    {
      title: '表注释',
      dataIndex: 'table_comment',
      key: 'table_comment',
      ellipsis: true,
      render: renderCellText<TableMeta, 'table_comment'>(
        'table_comment',
        (value) => value || '-',
      ),
    },
    {
      title: '操作',
      key: 'option',
      valueType: 'option',
      render: (_, record) => [
        <Button
          key="preview"
          type="link"
          onClick={(event) => {
            event.stopPropagation();
            openDataModal(record.table_name);
          }}
        >
          查看数据
        </Button>,
      ],
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
      render: renderCellText<ColumnMeta, 'is_nullable'>('is_nullable', (value) =>
        value ? '是' : '否',
      ),
    },
    {
      title: '最大长度',
      dataIndex: 'max_length',
      key: 'max_length',
      render: renderCellText<ColumnMeta, 'max_length'>('max_length', (value) =>
        value != null ? value : '-',
      ),
    },
    {
      title: '默认值',
      dataIndex: 'column_default',
      key: 'column_default',
      ellipsis: true,
      render: renderCellText<ColumnMeta, 'column_default'>(
        'column_default',
        (value) => value || '-',
      ),
    },
    {
      title: '字段注释',
      dataIndex: 'column_comment',
      key: 'column_comment',
      ellipsis: true,
      render: renderCellText<ColumnMeta, 'column_comment'>(
        'column_comment',
        (value) => value || '-',
      ),
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
                  const firstTable = res.data[0].table_name;
                  setCurrentTable(firstTable);
                  setSelectedRowKeys([firstTable]);
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
            rowSelection={{
              type: 'radio',
              selectedRowKeys,
              onChange: (keys) => {
                setSelectedRowKeys(keys);
              },
            }}
            onRow={(record: TableMeta) => {
              return {
                onClick: () => {
                  setCurrentTable(record.table_name);
                  setSelectedRowKeys([record.table_name]);
                },
              };
            }}
          />
        </ProCard>
        <ProCard
          title={currentTable ? `表信息（${currentTable}）` : '表信息'}
          extra={
            currentTable ? (
              <Button type="primary" onClick={() => openDataModal(currentTable)}>
                查看数据
              </Button>
            ) : undefined
          }
        >
          <ProCard title="字段">
            <ProTable<ColumnMeta>
              actionRef={columnsActionRef}
              rowKey="column_name"
              search={false}
              pagination={false}
              toolBarRender={false}
              params={{ tableName: currentTable }}
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
        </ProCard>
      </ProCard>
      <Modal
        title={currentTable ? `表数据预览（${currentTable}）` : '表数据预览'}
        open={dataModalOpen}
        width={1100}
        footer={null}
        destroyOnClose
        onCancel={() => setDataModalOpen(false)}
      >
        <ProTable<Record<string, any>>
          actionRef={rowsActionRef}
          rowKey="_row_key"
          search={false}
          pagination={{
            defaultPageSize: 10,
          }}
          toolBarRender={false}
          params={{ tableName: currentTable }}
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
      </Modal>
    </PageContainer>
  );
};

export default DataManagePage;
