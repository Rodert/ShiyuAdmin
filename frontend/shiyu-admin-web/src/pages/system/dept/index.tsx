import { PlusOutlined } from '@ant-design/icons';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { Button, message, Modal, Space } from 'antd';
import React, { useRef, useState } from 'react';
import {
  createDept,
  deleteDept,
  getDeptTree,
  updateDept,
  type CreateDeptRequest,
  type Dept,
  type UpdateDeptRequest,
} from '@/services/shiyu-api/dept';
import DeptForm from './components/DeptForm';

const DeptManagement: React.FC = () => {
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [updateModalVisible, setUpdateModalVisible] = useState(false);
  const [editingRecord, setEditingRecord] = useState<Dept | null>(null);
  const actionRef = useRef<ActionType>();

  const handleCreate = async (values: CreateDeptRequest) => {
    try {
      const res = await createDept(values);
      if (res.code === 200) {
        message.success('创建成功');
        setCreateModalVisible(false);
        actionRef.current?.reload();
      } else {
        message.error(res.message || '创建失败');
      }
    } catch (error) {
      message.error('创建失败');
    }
  };

  const handleUpdate = async (values: UpdateDeptRequest) => {
    if (!editingRecord) return;
    try {
      const res = await updateDept(editingRecord.dept_code, values);
      if (res.code === 200) {
        message.success('更新成功');
        setUpdateModalVisible(false);
        setEditingRecord(null);
        actionRef.current?.reload();
      } else {
        message.error(res.message || '更新失败');
      }
    } catch (error) {
      message.error('更新失败');
    }
  };

  const handleDelete = (record: Dept) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除部门 "${record.dept_name}" 吗？`,
      onOk: async () => {
        try {
          const res = await deleteDept(record.dept_code);
          if (res.code === 200) {
            message.success('删除成功');
            actionRef.current?.reload();
          } else {
            message.error(res.message || '删除失败');
          }
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  };

  const columns: ProColumns<Dept>[] = [
    {
      title: '部门编码',
      dataIndex: 'dept_code',
      key: 'dept_code',
      width: 120,
    },
    {
      title: '部门名称',
      dataIndex: 'dept_name',
      key: 'dept_name',
      width: 200,
    },
    {
      title: '父部门编码',
      dataIndex: 'parent_code',
      key: 'parent_code',
      width: 120,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      valueEnum: {
        1: { text: '启用', status: 'Success' },
        0: { text: '禁用', status: 'Error' },
      },
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      fixed: 'right',
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            size="small"
            onClick={() => {
              setEditingRecord(record);
              setUpdateModalVisible(true);
            }}
          >
            编辑
          </Button>
          <Button
            type="link"
            size="small"
            danger
            onClick={() => handleDelete(record)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <PageContainer>
      <ProTable<Dept>
        headerTitle="部门管理"
        actionRef={actionRef}
        rowKey="dept_code"
        search={false}
        toolBarRender={() => [
          <Button
            type="primary"
            key="primary"
            onClick={() => {
              setCreateModalVisible(true);
            }}
          >
            <PlusOutlined /> 新建
          </Button>,
        ]}
        request={async () => {
          const res = await getDeptTree();
          if (res.code === 200 && res.data) {
            // Flatten tree for table display
            const flatten = (depts: Dept[]): Dept[] => {
              const result: Dept[] = [];
              depts.forEach((dept) => {
                result.push(dept);
                if (dept.children && dept.children.length > 0) {
                  result.push(...flatten(dept.children));
                }
              });
              return result;
            };
            return {
              data: flatten(res.data),
              success: true,
              total: flatten(res.data).length,
            };
          }
          return {
            data: [],
            success: false,
            total: 0,
          };
        }}
        columns={columns}
        pagination={false}
      />

      <DeptForm
        visible={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        onSubmit={handleCreate}
        title="新建部门"
      />

      {editingRecord && (
        <DeptForm
          visible={updateModalVisible}
          onCancel={() => {
            setUpdateModalVisible(false);
            setEditingRecord(null);
          }}
          onSubmit={handleUpdate}
          title="编辑部门"
          initialValues={editingRecord}
        />
      )}
    </PageContainer>
  );
};

export default DeptManagement;

