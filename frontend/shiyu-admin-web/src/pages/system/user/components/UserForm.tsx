import { ProForm, ProFormText, ProFormSelect, ProFormTreeSelect } from '@ant-design/pro-components';
import type { ProFormInstance } from '@ant-design/pro-components';
import { Button, Modal, message } from 'antd';
import React, { useEffect, useMemo, useRef, useState } from 'react';
import type { CreateUserRequest, UpdateUserRequest, User } from '@/services/shiyu-api/user';
import { getUserRoles, setUserRoles } from '@/services/shiyu-api/user';
import { getRoleList, type Role } from '@/services/shiyu-api/role';
import { getDeptTree, type Dept } from '@/services/shiyu-api/dept';

type TreeOption = {
  title: string;
  value: string;
  children?: TreeOption[];
};

interface UserFormProps {
  visible: boolean;
  onCancel: () => void;
  onSubmit: (values: CreateUserRequest | UpdateUserRequest) => Promise<User | undefined>;
  title: string;
  initialValues?: User;
}

const UserForm: React.FC<UserFormProps> = ({
  visible,
  onCancel,
  onSubmit,
  title,
  initialValues,
}) => {
  const isEdit = !!initialValues;
  const [roleOptions, setRoleOptions] = useState<{ label: string; value: string }[]>([]);
  const [deptTreeOptions, setDeptTreeOptions] = useState<TreeOption[]>([]);
  const [selectedRoles, setSelectedRoles] = useState<string[]>([]);
  const formRef = useRef<ProFormInstance>(null);
  const memoizedInitialValues = useMemo(
    () => (initialValues ? { ...initialValues, dept_code: initialValues.dept_code || undefined } : undefined),
    [initialValues],
  );

  // 加载角色列表
  useEffect(() => {
    if (!visible) {
      formRef.current?.resetFields();
      setSelectedRoles([]);
      return;
    }

    formRef.current?.resetFields();
    if (memoizedInitialValues) {
      formRef.current?.setFieldsValue(memoizedInitialValues);
    }

    loadRoles();
    loadDepts();
    if (isEdit && initialValues?.user_code) {
      loadUserRoles(initialValues.user_code);
    } else {
      setSelectedRoles([]);
      formRef.current?.setFieldsValue({ role_codes: [] });
    }
  }, [visible, isEdit, initialValues?.user_code]);

  const loadRoles = async () => {
    try {
      const res = await getRoleList({ page: 1, page_size: 1000 });
      if (res.code === 200 && res.data) {
        const options = res.data.items.map((role: Role) => ({
          label: role.role_name,
          value: role.role_code,
        }));
        setRoleOptions(options);
      }
    } catch (error) {
      console.error('加载角色列表失败:', error);
    }
  };

  const loadUserRoles = async (userCode: string) => {
    try {
      const res = await getUserRoles(userCode);
      if (res.code === 200 && res.data) {
        const roleCodes = res.data.map((role: Role) => role.role_code);
        setSelectedRoles(roleCodes);
        formRef.current?.setFieldsValue({ role_codes: roleCodes });
      }
    } catch (error) {
      console.error('加载用户角色失败:', error);
    }
  };

  const loadDepts = async () => {
    try {
      const res = await getDeptTree();
      if (res.code === 200 && res.data) {
        setDeptTreeOptions(convertDeptToOptions(res.data));
      }
    } catch (error) {
      console.error('加载部门树失败:', error);
    }
  };

  const convertDeptToOptions = (depts: Dept[]): TreeOption[] => {
    return depts.map((dept) => ({
      title: `${dept.dept_name} (${dept.dept_code})`,
      value: dept.dept_code,
      children: dept.children ? convertDeptToOptions(dept.children) : undefined,
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
        key={isEdit ? initialValues?.user_code : 'create'}
        initialValues={memoizedInitialValues}
        onFinish={async (values) => {
          try {
            // 先提交用户信息
            const savedUser = await onSubmit(values as CreateUserRequest | UpdateUserRequest);
            
            // 获取用户编码（新建时从表单值获取，编辑时从 initialValues 获取）
            const userCode = isEdit 
              ? initialValues?.user_code 
              : savedUser?.user_code;
            
            // 设置用户角色
            if (userCode && selectedRoles.length >= 0) {
              try {
                const res = await setUserRoles(userCode, selectedRoles);
                if (res.code === 200) {
                  if (selectedRoles.length > 0) {
                    message.success('角色分配成功');
                  }
                }
              } catch (error) {
                message.warning('用户信息已保存，但角色分配失败');
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
              name="username"
              label="用户名"
              extra="需唯一，仅允许字母、数字、点、下划线或中划线"
              fieldProps={{ maxLength: 64, showCount: true }}
              rules={[
                { required: true, whitespace: true, message: '请输入用户名' },
                { pattern: /^[A-Za-z0-9._-]+$/, message: '用户名格式不正确' },
              ]}
            />
            <ProFormText.Password
              name="password"
              label="密码"
              rules={[{ required: true, message: '请输入密码' }]}
            />
          </>
        )}
        <ProFormText name="nickname" label="昵称" />
        <ProFormText name="email" label="邮箱" />
        <ProFormText name="phone" label="手机号" />
        <ProFormTreeSelect
          name="dept_code"
          label="所属部门"
          fieldProps={{
            treeData: deptTreeOptions,
            allowClear: true,
            treeDefaultExpandAll: true,
            placeholder: '请选择所属部门',
          }}
        />
        <ProFormSelect
          name="status"
          label="状态"
          options={[
            { label: '启用', value: 1 },
            { label: '禁用', value: 0 },
          ]}
        />
        {isEdit && (
          <ProFormText.Password
            name="password"
            label="新密码（留空不修改）"
          />
        )}
        <ProFormSelect
          name="role_codes"
          label="角色"
          mode="multiple"
          options={roleOptions}
          fieldProps={{
            value: selectedRoles,
            onChange: (value) => {
              setSelectedRoles(value as string[]);
              formRef.current?.setFieldsValue({ role_codes: value });
            },
          }}
          placeholder="请选择角色"
        />
      </ProForm>
    </Modal>
  );
};

export default UserForm;

