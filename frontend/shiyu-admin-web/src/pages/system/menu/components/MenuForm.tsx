import { ProForm, ProFormText, ProFormSelect } from '@ant-design/pro-components';
import { Button, Modal } from 'antd';
import React, { useMemo } from 'react';
import type { CreateMenuRequest, Menu, UpdateMenuRequest } from '@/services/shiyu-api/menu';

interface MenuFormProps {
  visible: boolean;
  onCancel: () => void;
  onSubmit: (values: CreateMenuRequest | UpdateMenuRequest) => void;
  title: string;
  initialValues?: Menu;
}

const MenuForm: React.FC<MenuFormProps> = ({
  visible,
  onCancel,
  onSubmit,
  title,
  initialValues,
}) => {
  const isEdit = !!initialValues;
  
  // 使用 useMemo 确保 initialValues 引用稳定
  const memoizedInitialValues = useMemo(() => {
    return initialValues ? { ...initialValues } : undefined;
  }, [initialValues]);

  return (
    <Modal
      title={title}
      open={visible}
      onCancel={onCancel}
      footer={null}
      width={600}
    >
      <ProForm
        key={isEdit ? initialValues?.menu_code : 'create'}
        initialValues={memoizedInitialValues}
        onFinish={async (values) => {
          onSubmit(values as CreateMenuRequest | UpdateMenuRequest);
        }}
        submitter={{
          render: (props, doms) => {
            return [
              <Button key="cancel" onClick={onCancel}>
                取消
              </Button>,
              <Button key="submit" type="primary" onClick={() => props.form?.submit?.()}>
                确定
              </Button>,
            ];
          },
        }}
      >
        {!isEdit && (
          <>
            <ProFormText
              name="menu_code"
              label="菜单编码"
              rules={[{ required: true, message: '请输入菜单编码' }]}
            />
            <ProFormSelect
              name="menu_type"
              label="菜单类型"
              rules={[{ required: true, message: '请选择菜单类型' }]}
              options={[
                { label: '目录', value: 'M' },
                { label: '菜单', value: 'C' },
                { label: '按钮', value: 'F' },
              ]}
            />
          </>
        )}
        {isEdit && (
          <ProFormSelect
            name="menu_type"
            label="菜单类型"
            options={[
              { label: '目录', value: 'M' },
              { label: '菜单', value: 'C' },
              { label: '按钮', value: 'F' },
            ]}
          />
        )}
        <ProFormText name="parent_code" label="父菜单编码" />
        <ProFormText name="menu_name" label="菜单名称" />
        <ProFormText name="perms" label="权限标识" />
        <ProFormText name="path" label="路径" />
        <ProFormText name="component" label="组件" />
        <ProFormSelect
          name="status"
          label="状态"
          options={[
            { label: '启用', value: 1 },
            { label: '禁用', value: 0 },
          ]}
        />
      </ProForm>
    </Modal>
  );
};

export default MenuForm;

