import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { Button, message, Modal } from 'antd';
import React, { useRef, useState } from 'react';
import {
  createRole,
  deleteRole,
  getRoleList,
  updateRole,
  type CreateRoleRequest,
  type Role,
  type UpdateRoleRequest,
} from '@/services/shiyu-api/role';
import RoleForm from './components/RoleForm';

const RoleManagement: React.FC = () => {
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [updateModalVisible, setUpdateModalVisible] = useState(false);
  const [editingRecord, setEditingRecord] = useState<Role | null>(null);
  const actionRef = useRef<ActionType>();

  const handleCreate = async (values: CreateRoleRequest): Promise<void> => {
    try {
      const res = await createRole(values);
      if (res.code === 200) {
        message.success('创建成功');
        setCreateModalVisible(false);
        actionRef.current?.reload();
      } else {
        message.error(res.message || '创建失败');
        throw new Error(res.message || '创建失败');
      }
    } catch (error) {
      message.error('创建失败');
      throw error;
    }
  };

  const handleUpdate = async (values: UpdateRoleRequest): Promise<void> => {
    if (!editingRecord) return;
    try {
      const res = await updateRole(editingRecord.role_code, values);
      if (res.code === 200) {
        message.success('更新成功');
        setUpdateModalVisible(false);
        setEditingRecord(null);
        actionRef.current?.reload();
      } else {
        message.error(res.message || '更新失败');
        throw new Error(res.message || '更新失败');
      }
    } catch (error) {
      message.error('更新失败');
      throw error;
    }
  };

  const handleDelete = (record: Role) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除角色 "${record.role_name}" 吗？`,
      onOk: async () => {
        try {
          const res = await deleteRole(record.role_code);
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

  const columns: ProColumns<Role>[] = [
    {
      title: '角色编码',
      dataIndex: 'role_code',
      key: 'role_code',
      width: 120,
    },
    {
      title: '角色名称',
      dataIndex: 'role_name',
      key: 'role_name',
      width: 120,
    },
    {
      title: '角色标识',
      dataIndex: 'role_key',
      key: 'role_key',
      width: 120,
    },
    {
      title: '数据权限',
      dataIndex: 'data_scope',
      key: 'data_scope',
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
      render: (_, record) => [
        <Button
          key="edit"
          type="link"
          size="small"
          onClick={() => {
            setEditingRecord(record);
            setUpdateModalVisible(true);
          }}
        >
          编辑
        </Button>,
        <Button
          key="delete"
          type="link"
          size="small"
          danger
          onClick={() => handleDelete(record)}
        >
          删除
        </Button>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<Role>
        headerTitle="角色管理"
        actionRef={actionRef}
        rowKey="role_code"
        search={{
          labelWidth: 'auto',
        }}
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
        request={async (params) => {
          const res = await getRoleList({
            page: params.current || 1,
            page_size: params.pageSize || 10,
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

      <RoleForm
        visible={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        onSubmit={handleCreate}
        title="新建角色"
      />

      {editingRecord && (
        <RoleForm
          visible={updateModalVisible}
          onCancel={() => {
            setUpdateModalVisible(false);
            setEditingRecord(null);
          }}
          onSubmit={handleUpdate}
          title="编辑角色"
          initialValues={editingRecord}
        />
      )}
    </PageContainer>
  );
};

export default RoleManagement;

