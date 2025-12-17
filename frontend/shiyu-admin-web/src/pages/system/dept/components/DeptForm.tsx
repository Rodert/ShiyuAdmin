import { ProForm, ProFormText, ProFormSelect } from '@ant-design/pro-components';
import { Button, Modal } from 'antd';
import React, { useMemo } from 'react';
import type { CreateDeptRequest, Dept, UpdateDeptRequest } from '@/services/shiyu-api/dept';

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
            rules={[{ required: true, message: '请输入部门编码' }]}
          />
        )}
        <ProFormText name="parent_code" label="父部门编码" />
        {!isEdit && (
          <ProFormText
            name="dept_name"
            label="部门名称"
            rules={[{ required: true, message: '请输入部门名称' }]}
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

