import { ProForm, ProFormText, ProFormSelect, ProFormTreeSelect } from '@ant-design/pro-components';
import type { ProFormInstance } from '@ant-design/pro-components';
import { Button, Modal } from 'antd';
import React, { useEffect, useMemo, useRef, useState } from 'react';
import { getMenuTree, type CreateMenuRequest, type Menu, type UpdateMenuRequest } from '@/services/shiyu-api/menu';

type TreeOption = {
  title: string;
  value: string;
  disabled?: boolean;
  children?: TreeOption[];
};

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
  const [menuTreeOptions, setMenuTreeOptions] = useState<TreeOption[]>([]);
  const formRef = useRef<ProFormInstance>(null);
  
  // 使用 useMemo 确保 initialValues 引用稳定
  const memoizedInitialValues = useMemo(() => {
    return initialValues
      ? { ...initialValues, parent_code: initialValues.parent_code || undefined }
      : undefined;
  }, [initialValues]);

  useEffect(() => {
    if (!visible) {
      formRef.current?.resetFields();
      return;
    }

    formRef.current?.resetFields();
    if (memoizedInitialValues) {
      formRef.current?.setFieldsValue(memoizedInitialValues);
    }

    const loadMenuOptions = async () => {
      const res = await getMenuTree();
      if (res.code === 200 && res.data) {
        setMenuTreeOptions(buildMenuOptions(res.data, initialValues?.menu_code));
      }
    };

    loadMenuOptions().catch((error) => {
      console.error('加载菜单树失败:', error);
      setMenuTreeOptions([]);
    });
  }, [visible, initialValues?.menu_code]);

  const buildMenuOptions = (menus: Menu[], currentMenuCode?: string): TreeOption[] => {
    const excludedCodes = currentMenuCode ? collectMenuDescendantCodes(menus, currentMenuCode) : new Set<string>();
    return menus
      .filter((menu) => menu.menu_code !== currentMenuCode)
      .map((menu) => convertMenuToOption(menu, excludedCodes));
  };

  const convertMenuToOption = (menu: Menu, excludedCodes: Set<string>): TreeOption => ({
    title: `${menu.menu_name || menu.menu_code} (${menu.menu_code})`,
    value: menu.menu_code,
    disabled: excludedCodes.has(menu.menu_code),
    children: menu.children?.map((child) => convertMenuToOption(child, excludedCodes)),
  });

  const collectMenuDescendantCodes = (menus: Menu[], currentMenuCode: string): Set<string> => {
    const target = findMenuNode(menus, currentMenuCode);
    const codes = new Set<string>([currentMenuCode]);

    const visit = (menu?: Menu) => {
      if (!menu) {
        return;
      }
      codes.add(menu.menu_code);
      menu.children?.forEach(visit);
    };

    visit(target);
    return codes;
  };

  const findMenuNode = (menus: Menu[], targetCode: string): Menu | undefined => {
    for (const menu of menus) {
      if (menu.menu_code === targetCode) {
        return menu;
      }
      const child = menu.children ? findMenuNode(menu.children, targetCode) : undefined;
      if (child) {
        return child;
      }
    }
    return undefined;
  };

  return (
    <Modal
      title={title}
      open={visible}
      onCancel={onCancel}
      footer={null}
      width={600}
      destroyOnHidden
    >
      <ProForm
        formRef={formRef}
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
              extra="需唯一，仅允许字母、数字、下划线或中划线"
              fieldProps={{ maxLength: 32, showCount: true }}
              rules={[
                { required: true, whitespace: true, message: '请输入菜单编码' },
                { pattern: /^[A-Za-z0-9_-]+$/, message: '菜单编码格式不正确' },
              ]}
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
        <ProFormTreeSelect
          name="parent_code"
          label="父菜单"
          fieldProps={{
            treeData: menuTreeOptions,
            allowClear: true,
            treeDefaultExpandAll: true,
            placeholder: '不选则为顶级菜单',
          }}
        />
        <ProFormText name="menu_name" label="菜单名称" />
        <ProFormText
          name="perms"
          label="权限标识"
          extra="非空时需唯一，仅允许字母、数字、冒号、下划线或中划线"
          fieldProps={{ maxLength: 128, showCount: true }}
          rules={[{ pattern: /^[A-Za-z0-9:_-]*$/, message: '权限标识格式不正确' }]}
        />
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

