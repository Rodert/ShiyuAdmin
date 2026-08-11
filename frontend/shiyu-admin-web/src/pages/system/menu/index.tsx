import { CaretDownOutlined, CaretRightOutlined, PlusOutlined } from '@ant-design/icons';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { Button, message, Modal, Space, Tooltip } from 'antd';
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
import { hasPermission } from '@/utils/permission';
import { useModel } from '@umijs/max';
import MenuForm from './components/MenuForm';

const MenuManagement: React.FC = () => {
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [updateModalVisible, setUpdateModalVisible] = useState(false);
  const [editingRecord, setEditingRecord] = useState<Menu | null>(null);
  const actionRef = useRef<ActionType>(null);
  const { initialState } = useModel('@@initialState');
  const currentUser = initialState?.currentUser;
  const canCreate = hasPermission(currentUser, 'system:menu:create');
  const canUpdate = hasPermission(currentUser, 'system:menu:update');
  const canDelete = hasPermission(currentUser, 'system:menu:delete');

  const countMenus = (menus: Menu[]): number => {
    return menus.reduce(
      (total, menu) => total + 1 + countMenus(menu.children || []),
      0,
    );
  };

  const handleCreate = async (
    values: CreateMenuRequest | UpdateMenuRequest,
  ) => {
    try {
      const res = await createMenu(values as CreateMenuRequest);
      if (res.code === 200) {
        message.success('创建成功');
        setCreateModalVisible(false);
        actionRef.current?.reload();
      }
    } catch (_error) {}
  };

  const handleUpdate = async (
    values: CreateMenuRequest | UpdateMenuRequest,
  ) => {
    if (!editingRecord) return;
    try {
      const res = await updateMenu(
        editingRecord.menu_code,
        values as UpdateMenuRequest,
      );
      if (res.code === 200) {
        message.success('更新成功');
        setUpdateModalVisible(false);
        setEditingRecord(null);
        actionRef.current?.reload();
      }
    } catch (_error) {}
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
          }
        } catch (_error) {}
      },
    });
  };

  const columns: ProColumns<Menu>[] = [
    {
      title: '菜单编码',
      dataIndex: 'menu_code',
      key: 'menu_code',
      width: 220,
      ellipsis: true,
      render: (_, record) => <span className="menu-code">{record.menu_code}</span>,
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
      title: '排序',
      dataIndex: 'sort_order',
      key: 'sort_order',
      width: 72,
      search: false,
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
          {canUpdate && (
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
          )}
          {canDelete && (
            <Button
              type="link"
              size="small"
              danger
              onClick={() => handleDelete(record)}
            >
              删除
            </Button>
          )}
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
        toolBarRender={() =>
          canCreate
            ? [
                <Button
                  type="primary"
                  key="primary"
                  onClick={() => {
                    setCreateModalVisible(true);
                  }}
                >
                  <PlusOutlined /> 新建
                </Button>,
              ]
            : []
        }
        request={async () => {
          const res = await getMenuTree();
          if (res.code === 200 && res.data) {
            return {
              data: res.data,
              success: true,
              total: countMenus(res.data),
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
        expandable={{
          defaultExpandAllRows: false,
          expandIcon: ({ expanded, onExpand, record }) => (
            <Tooltip title={expanded ? '收起子菜单' : '展开子菜单'}>
              <Button
                aria-label={expanded ? '收起子菜单' : '展开子菜单'}
                className="menu-tree-toggle"
                icon={expanded ? <CaretDownOutlined /> : <CaretRightOutlined />}
                onClick={(event) => onExpand(record, event)}
                size="small"
                type="text"
              />
            </Tooltip>
          ),
        }}
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
