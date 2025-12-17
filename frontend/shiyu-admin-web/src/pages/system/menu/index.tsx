import { PlusOutlined } from '@ant-design/icons';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { Button, message, Modal, Space } from 'antd';
import React, { useRef, useState } from 'react';
import {
  createMenu,
  deleteMenu,
  getMenuTree,
  updateMenu,
  type CreateMenuRequest,
  type Menu,
  type UpdateMenuRequest,
} from '@/services/shiyu-api/menu';
import MenuForm from './components/MenuForm';

const MenuManagement: React.FC = () => {
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [updateModalVisible, setUpdateModalVisible] = useState(false);
  const [editingRecord, setEditingRecord] = useState<Menu | null>(null);
  const actionRef = useRef<ActionType>();

  const handleCreate = async (values: CreateMenuRequest) => {
    try {
      const res = await createMenu(values);
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

  const handleUpdate = async (values: UpdateMenuRequest) => {
    if (!editingRecord) return;
    try {
      const res = await updateMenu(editingRecord.menu_code, values);
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

  const handleDelete = (record: Menu) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除菜单 "${record.menu_name}" 吗？`,
      onOk: async () => {
        try {
          const res = await deleteMenu(record.menu_code);
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

  const columns: ProColumns<Menu>[] = [
    {
      title: '菜单编码',
      dataIndex: 'menu_code',
      key: 'menu_code',
      width: 120,
    },
    {
      title: '菜单名称',
      dataIndex: 'menu_name',
      key: 'menu_name',
      width: 150,
    },
    {
      title: '菜单类型',
      dataIndex: 'menu_type',
      key: 'menu_type',
      width: 100,
      valueEnum: {
        M: { text: '目录' },
        C: { text: '菜单' },
        F: { text: '按钮' },
      },
    },
    {
      title: '权限标识',
      dataIndex: 'perms',
      key: 'perms',
      width: 150,
    },
    {
      title: '路径',
      dataIndex: 'path',
      key: 'path',
      width: 200,
    },
    {
      title: '组件',
      dataIndex: 'component',
      key: 'component',
      width: 200,
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
      <ProTable<Menu>
        headerTitle="菜单管理"
        actionRef={actionRef}
        rowKey="menu_code"
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
          const res = await getMenuTree();
          if (res.code === 200 && res.data) {
            // Flatten tree for table display
            const flatten = (menus: Menu[]): Menu[] => {
              const result: Menu[] = [];
              menus.forEach((menu) => {
                result.push(menu);
                if (menu.children && menu.children.length > 0) {
                  result.push(...flatten(menu.children));
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

      <MenuForm
        visible={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        onSubmit={handleCreate}
        title="新建菜单"
      />

      {editingRecord && (
        <MenuForm
          visible={updateModalVisible}
          onCancel={() => {
            setUpdateModalVisible(false);
            setEditingRecord(null);
          }}
          onSubmit={handleUpdate}
          title="编辑菜单"
          initialValues={editingRecord}
        />
      )}
    </PageContainer>
  );
};

export default MenuManagement;

