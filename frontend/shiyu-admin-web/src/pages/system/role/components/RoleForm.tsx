import { getMenuTree, type Menu } from "@/services/shiyu-api/menu";
import type {
  CreateRoleRequest,
  Role,
  UpdateRoleRequest,
} from "@/services/shiyu-api/role";
import { getRoleMenus, setRoleMenus } from "@/services/shiyu-api/role";
import type { ProFormInstance } from "@ant-design/pro-components";
import {
  ProForm,
  ProFormSelect,
  ProFormText,
} from "@ant-design/pro-components";
import { Button, message, Modal, Tree } from "antd";
import type { DataNode } from "antd/es/tree";
import React, { useEffect, useMemo, useRef, useState } from "react";

interface RoleFormProps {
  visible: boolean;
  onCancel: () => void;
  onSubmit: (values: CreateRoleRequest | UpdateRoleRequest) => Promise<void>;
  title: string;
  initialValues?: Role;
}

const RoleForm: React.FC<RoleFormProps> = ({
  visible,
  onCancel,
  onSubmit,
  title,
  initialValues,
}) => {
  const isEdit = !!initialValues;
  const [menuTree, setMenuTree] = useState<DataNode[]>([]);
  const [selectedMenuKeys, setSelectedMenuKeys] = useState<React.Key[]>([]);
  const formRef = useRef<ProFormInstance>(null);
  const memoizedInitialValues = useMemo(
    () => (initialValues ? { ...initialValues } : undefined),
    [initialValues]
  );

  // 加载菜单树
  useEffect(() => {
    if (!visible) {
      formRef.current?.resetFields();
      setSelectedMenuKeys([]);
      return;
    }

    loadMenuTree();
    if (isEdit && initialValues?.role_code) {
      formRef.current?.resetFields();
      formRef.current?.setFieldsValue(memoizedInitialValues);
      loadRoleMenus(initialValues.role_code);
    } else {
      formRef.current?.resetFields();
      setSelectedMenuKeys([]);
    }
  }, [visible, isEdit, initialValues?.role_code, memoizedInitialValues]);

  const loadMenuTree = async () => {
    try {
      const res = await getMenuTree();
      if (res.code === 200 && res.data) {
        const tree = convertMenuToTree(res.data);
        setMenuTree(tree);
      }
    } catch (error) {
      console.error("加载菜单树失败:", error);
    }
  };

  const loadRoleMenus = async (roleCode: string) => {
    try {
      const res = await getRoleMenus(roleCode);
      if (res.code === 200 && res.data) {
        const menuCodes = res.data.map((menu: Menu) => menu.menu_code);
        setSelectedMenuKeys(menuCodes);
      }
    } catch (error) {
      console.error("加载角色菜单失败:", error);
    }
  };

  // 将菜单数据转换为树形结构
  const convertMenuToTree = (menus: Menu[]): DataNode[] => {
    return menus.map((menu) => ({
      title: `${menu.menu_name} (${menu.menu_code})`,
      key: menu.menu_code,
      children: menu.children ? convertMenuToTree(menu.children) : undefined,
    }));
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
        key={isEdit ? initialValues?.role_code : "create"}
        initialValues={memoizedInitialValues}
        onFinish={async (values) => {
          try {
            const submitValues = { ...values };
            if (isEdit) {
              delete submitValues.role_key;
            }
            // 先提交角色信息
            await onSubmit(
              submitValues as CreateRoleRequest | UpdateRoleRequest
            );

            // 获取角色编码（新建时从表单值获取，编辑时从 initialValues 获取）
            const roleCode = isEdit
              ? initialValues?.role_code
              : (values as CreateRoleRequest).role_code;

            // 设置角色菜单
            if (roleCode && selectedMenuKeys.length >= 0) {
              try {
                const menuCodes = selectedMenuKeys.map((key) => String(key));
                const res = await setRoleMenus(roleCode, menuCodes);
                if (res.code === 200) {
                  if (selectedMenuKeys.length > 0) {
                    message.success("菜单权限分配成功");
                  }
                }
              } catch (error) {
                message.warning("角色信息已保存，但菜单权限分配失败");
              }
            }
          } catch (error) {
            // 错误已在 onSubmit 中处理
          }
        }}
        submitter={{
          render: (props, doms) => {
            return [
              <Button key="cancel" onClick={onCancel}>
                取消
              </Button>,
              <Button
                key="submit"
                type="primary"
                onClick={() => props.form?.submit?.()}
              >
                确定
              </Button>,
            ];
          },
        }}
      >
        {!isEdit && (
          <>
            <ProFormText
              name="role_code"
              label="角色编码"
              extra="需唯一，仅允许字母、数字、下划线或中划线"
              fieldProps={{ maxLength: 32, showCount: true }}
              rules={[
                { required: true, whitespace: true, message: "请输入角色编码" },
                { pattern: /^[A-Za-z0-9_-]+$/, message: "角色编码格式不正确" },
              ]}
            />
            <ProFormText
              name="role_name"
              label="角色名称"
              rules={[
                { required: true, whitespace: true, message: "请输入角色名称" },
              ]}
            />
            <ProFormText
              name="role_key"
              label="角色标识"
              extra="需唯一，仅允许字母、数字、冒号、下划线或中划线"
              fieldProps={{ maxLength: 64, showCount: true }}
              rules={[
                { required: true, whitespace: true, message: "请输入角色标识" },
                { pattern: /^[A-Za-z0-9:_-]+$/, message: "角色标识格式不正确" },
              ]}
            />
          </>
        )}
        {isEdit && (
          <>
            <ProFormText name="role_name" label="角色名称" />
            <ProFormText
              name="role_key"
              label="角色标识"
              disabled
              extra="角色标识用于权限判断，创建后不可修改"
            />
          </>
        )}
        <ProFormSelect
          name="data_scope"
          label="数据权限"
          options={[
            { label: "全部数据", value: "all" },
            { label: "部门数据", value: "dept" },
            { label: "部门及以下", value: "dept_and_child" },
            { label: "仅本人", value: "self" },
          ]}
        />
        <ProFormSelect
          name="status"
          label="状态"
          options={[
            { label: "启用", value: 1 },
            { label: "禁用", value: 0 },
          ]}
        />
        <ProForm.Item label="菜单权限">
          <Tree
            checkable
            treeData={menuTree}
            checkedKeys={selectedMenuKeys}
            onCheck={(checkedKeys) => {
              setSelectedMenuKeys(checkedKeys as React.Key[]);
            }}
            style={{ maxHeight: "300px", overflow: "auto" }}
          />
        </ProForm.Item>
      </ProForm>
    </Modal>
  );
};

export default RoleForm;
