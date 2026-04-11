import { ProForm, ProFormText, ProFormSelect, ProFormTreeSelect } from '@ant-design/pro-components';
import type { ProFormInstance } from '@ant-design/pro-components';
import { Button, Modal } from 'antd';
import React, { useEffect, useMemo, useRef, useState } from 'react';
import { getDeptTree, type CreateDeptRequest, type Dept, type UpdateDeptRequest } from '@/services/shiyu-api/dept';

type TreeOption = {
  title: string;
  value: string;
  disabled?: boolean;
  children?: TreeOption[];
};

interface DeptFormProps {
  visible: boolean;
  onCancel: () => void;
  onSubmit: (values: CreateDeptRequest | UpdateDeptRequest) => void;
  title: string;
  initialValues?: Dept;
}

const DeptForm: React.FC<DeptFormProps> = ({
  visible,
  onCancel,
  onSubmit,
  title,
  initialValues,
}) => {
  const isEdit = !!initialValues;
  const [deptTreeOptions, setDeptTreeOptions] = useState<TreeOption[]>([]);
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

    const loadDeptOptions = async () => {
      const res = await getDeptTree();
      if (res.code === 200 && res.data) {
        setDeptTreeOptions(buildDeptOptions(res.data, initialValues?.dept_code));
      }
    };

    loadDeptOptions().catch((error) => {
      console.error('加载部门树失败:', error);
      setDeptTreeOptions([]);
    });
  }, [visible, initialValues?.dept_code]);

  const buildDeptOptions = (depts: Dept[], currentDeptCode?: string): TreeOption[] => {
    const excludedCodes = currentDeptCode ? collectDeptDescendantCodes(depts, currentDeptCode) : new Set<string>();
    return depts
      .filter((dept) => dept.dept_code !== currentDeptCode)
      .map((dept) => convertDeptToOption(dept, excludedCodes));
  };

  const convertDeptToOption = (dept: Dept, excludedCodes: Set<string>): TreeOption => ({
    title: `${dept.dept_name} (${dept.dept_code})`,
    value: dept.dept_code,
    disabled: excludedCodes.has(dept.dept_code),
    children: dept.children?.map((child) => convertDeptToOption(child, excludedCodes)),
  });

  const collectDeptDescendantCodes = (depts: Dept[], currentDeptCode: string): Set<string> => {
    const target = findDeptNode(depts, currentDeptCode);
    const codes = new Set<string>([currentDeptCode]);

    const visit = (dept?: Dept) => {
      if (!dept) {
        return;
      }
      codes.add(dept.dept_code);
      dept.children?.forEach(visit);
    };

    visit(target);
    return codes;
  };

  const findDeptNode = (depts: Dept[], targetCode: string): Dept | undefined => {
    for (const dept of depts) {
      if (dept.dept_code === targetCode) {
        return dept;
      }
      const child = dept.children ? findDeptNode(dept.children, targetCode) : undefined;
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
        key={isEdit ? initialValues?.dept_code : 'create'}
        initialValues={memoizedInitialValues}
        onFinish={async (values) => {
          onSubmit(values as CreateDeptRequest | UpdateDeptRequest);
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
          <ProFormText
            name="dept_code"
            label="部门编码"
            extra="需唯一，仅允许字母、数字、下划线或中划线"
            fieldProps={{ maxLength: 32, showCount: true }}
            rules={[
              { required: true, whitespace: true, message: '请输入部门编码' },
              { pattern: /^[A-Za-z0-9_-]+$/, message: '部门编码格式不正确' },
            ]}
          />
        )}
        <ProFormTreeSelect
          name="parent_code"
          label="父部门"
          fieldProps={{
            treeData: deptTreeOptions,
            allowClear: true,
            treeDefaultExpandAll: true,
            placeholder: '不选则为顶级部门',
          }}
        />
        {!isEdit && (
          <ProFormText
            name="dept_name"
            label="部门名称"
            rules={[{ required: true, whitespace: true, message: '请输入部门名称' }]}
          />
        )}
        {isEdit && <ProFormText name="dept_name" label="部门名称" />}
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

export default DeptForm;

