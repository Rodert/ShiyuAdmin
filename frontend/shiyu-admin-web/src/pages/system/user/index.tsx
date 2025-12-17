import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { Button, message, Modal } from 'antd';
import React, { useRef, useState } from 'react';
import {
  createUser,
  deleteUser,
  getUserList,
  updateUser,
  type CreateUserRequest,
  type UpdateUserRequest,
  type User,
} from '@/services/shiyu-api/user';
import UserForm from './components/UserForm';

const UserManagement: React.FC = () => {
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [updateModalVisible, setUpdateModalVisible] = useState(false);
  const [editingRecord, setEditingRecord] = useState<User | null>(null);
  const actionRef = useRef<ActionType>();

  const handleCreate = async (values: CreateUserRequest): Promise<void> => {
    try {
      const res = await createUser(values);
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

  const handleUpdate = async (values: UpdateUserRequest): Promise<void> => {
    if (!editingRecord) return;
    try {
      const res = await updateUser(editingRecord.user_code, values);
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

  const handleDelete = (record: User) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除用户 "${record.username}" 吗？`,
      onOk: async () => {
        try {
          const res = await deleteUser(record.user_code);
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

  const columns: ProColumns<User>[] = [
    {
      title: '用户编码',
      dataIndex: 'user_code',
      key: 'user_code',
      width: 120,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      width: 120,
    },
    {
      title: '昵称',
      dataIndex: 'nickname',
      key: 'nickname',
      width: 120,
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
      width: 180,
    },
    {
      title: '手机号',
      dataIndex: 'phone',
      key: 'phone',
      width: 120,
    },
    {
      title: '部门编码',
      dataIndex: 'dept_code',
      key: 'dept_code',
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
      <ProTable<User>
        headerTitle="用户管理"
        actionRef={actionRef}
        rowKey="user_code"
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
          const res = await getUserList({
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

      <UserForm
        visible={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        onSubmit={handleCreate}
        title="新建用户"
      />

      {editingRecord && (
        <UserForm
          visible={updateModalVisible}
          onCancel={() => {
            setUpdateModalVisible(false);
            setEditingRecord(null);
          }}
          onSubmit={handleUpdate}
          title="编辑用户"
          initialValues={editingRecord}
        />
      )}
    </PageContainer>
  );
};

export default UserManagement;

