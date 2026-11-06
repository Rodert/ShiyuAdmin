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
import { Alert, Button, message, Modal, Select, Tree } from "antd";
import type { DataNode } from "antd/es/tree";
import React, { useEffect, useMemo, useRef, useState } from "react";

interface RoleFormProps {
  visible: boolean;
  onCancel: () => void;
  onSubmit: (
    values: CreateRoleRequest | UpdateRoleRequest,
    savedCode?: string
  ) => Promise<string>;
  onSaved: () => void;
  title: string;
  initialValues?: Role;
  canAssignMenus?: boolean;
}

const RoleForm: React.FC<RoleFormProps> = ({
  visible,
  onCancel,
  onSubmit,
  onSaved,
  title,
  initialValues,
  canAssignMenus = true,
}) => {
  const isEdit = !!initialValues;
  const [savedCode, setSavedCode] = useState<string>();
  const [saving, setSaving] = useState(false);
  const [menus, setMenus] = useState<Menu[]>([]);
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
      setSavedCode(undefined);
      return;
    }

    if (canAssignMenus) {
      loadMenuTree();
    }
    if (canAssignMenus && isEdit && initialValues?.role_code) {
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
        setMenus(res.data);
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
      onCancel={() => {
        if (!saving) onCancel();
      }}
      footer={null}
      width={600}
      destroyOnHidden
    >
      <ProForm
        formRef={formRef}
        key={isEdit ? initialValues?.role_code : "create"}
        initialValues={
          memoizedInitialValues || { data_scope: "self", status: 1 }
        }
        onFinish={async (values) => {
          if (saving) return;
          setSaving(true);
          let persisted = savedCode;
          try {
            const submitValues = { ...values };
            if (isEdit || savedCode) delete submitValues.role_key;
            persisted = await onSubmit(
              submitValues as CreateRoleRequest | UpdateRoleRequest,
              savedCode
            );
            setSavedCode(persisted);
            if (canAssignMenus)
              await setRoleMenus(persisted, selectedMenuKeys.map(String));
            message.success("角色信息与权限已保存");
            onSaved();
          } catch {
            if (persisted)
              message.warning(
                "角色信息已保存，权限保存未完成，请重试；不会重复创建角色"
              );
          } finally {
            setSaving(false);
          }
        }}
        submitter={{
          render: (props, _doms) => {
            return [
              <Button key="cancel" onClick={onCancel}>
                取消
              </Button>,
              <Button
                key="submit"
                type="primary"
                loading={saving}
                onClick={() => props.form?.submit?.()}
              >
                确定
              </Button>,
            ];
          },
        }}
      >
        <Alert
          showIcon
          type="info"
          style={{ marginBottom: 16 }}
          message="菜单权限决定能做什么，数据权限决定能看到谁的数据"
          description="多个启用角色的权限取并集；全部数据范围会覆盖较小范围。组织与权限配置仍仅超级管理员可修改。事项审批仅按参与人隔离，不随角色扩大范围。"
        />
        {!isEdit && !savedCode && (
          <ProForm.Item
            label="从常用角色开始"
            extra="模板填充后仍可修改；不会赋予超级管理员身份。"
          >
            <Select
              aria-label="角色模板"
              placeholder="选择角色模板（可选）"
              options={[
                { value: "member", label: "团队成员：本人数据与文件" },
                { value: "manager", label: "部门负责人：本部门用户与日志" },
                { value: "viewer", label: "审阅人员：本部门只读" },
              ]}
              onChange={(value) => {
                const config =
                  value === "member"
                    ? {
                        name: "团队成员",
                        scope: "self",
                        perms: [
                          "welcome:view",
                          "system:file:list",
                          "system:file:upload",
                          "system:file:delete",
                        ],
                      }
                    : value === "manager"
                    ? {
                        name: "部门负责人",
                        scope: "dept",
                        perms: [
                          "welcome:view",
                          "system:user:list",
                          "system:user:create",
                          "system:user:update",
                          "system:dept:list",
                          "system:operation-log:list",
                          "system:file:list",
                          "system:file:upload",
                        ],
                      }
                    : {
                        name: "审阅人员",
                        scope: "dept",
                        perms: [
                          "welcome:view",
                          "system:user:list",
                          "system:dept:list",
                          "system:operation-log:list",
                        ],
                      };
                formRef.current?.setFieldsValue({
                  role_name: config.name,
                  data_scope: config.scope,
                  status: 1,
                });
                const keys: string[] = [];
                const collect = (items: Menu[]) => {
                  items.forEach((item) => {
                    if (config.perms.includes(item.perms || ""))
                      keys.push(item.menu_code);
                    if (item.children) collect(item.children);
                  });
                };
                collect(menus);
                setSelectedMenuKeys(keys);
              }}
              disabled={!menus.length}
            />
          </ProForm.Item>
        )}
        {!isEdit && (
          <>
            <ProFormText
              name="role_code"
              disabled={!!savedCode}
              label="角色编码"
              extra="需唯一，仅允许字母、数字、下划线或中划线"
              fieldProps={{ maxLength: 32, showCount: true, "aria-label": "角色编码" }}
              rules={[
                { required: true, whitespace: true, message: "请输入角色编码" },
                { pattern: /^[A-Za-z0-9_-]+$/, message: "角色编码格式不正确" },
              ]}
            />
            <ProFormText
              name="role_name"
              fieldProps={{ "aria-label": "角色名称" }}
              label="角色名称"
              rules={[
                { required: true, whitespace: true, message: "请输入角色名称" },
              ]}
            />
            <ProFormText
              name="role_key"
              disabled={!!savedCode}
              label="角色标识"
              extra="需唯一，仅允许字母、数字、冒号、下划线或中划线"
              fieldProps={{ maxLength: 64, showCount: true, "aria-label": "角色标识" }}
              rules={[
                { required: true, whitespace: true, message: "请输入角色标识" },
                { pattern: /^[A-Za-z0-9:_-]+$/, message: "角色标识格式不正确" },
              ]}
            />
          </>
        )}
        {isEdit && (
          <>
            <ProFormText name="role_name" label="角色名称" fieldProps={{ "aria-label": "角色名称" }} />
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
          rules={[{ required: true }]}
          extra="仅本人：自己的记录；部门：所属部门；部门及以下：包含子部门；全部：所有部门。文件仍仅上传者与超级管理员可见。"
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
        {canAssignMenus && (
          <ProForm.Item label="菜单权限">
            <Tree
              checkable
              checkStrictly
              treeData={menuTree}
              checkedKeys={selectedMenuKeys}
              onCheck={(checkedKeys) => {
                setSelectedMenuKeys(
                  Array.isArray(checkedKeys) ? checkedKeys : checkedKeys.checked
                );
              }}
              style={{ maxHeight: "300px", overflow: "auto" }}
            />
          </ProForm.Item>
        )}
      </ProForm>
    </Modal>
  );
};

export default RoleForm;
