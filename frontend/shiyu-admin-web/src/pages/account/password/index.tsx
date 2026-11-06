import { changeCurrentPassword } from "@/services/shiyu-api/user";
import {
  PageContainer,
  ProCard,
  ProForm,
  ProFormText,
} from "@ant-design/pro-components";
import { Alert, message } from "antd";
import React from "react";

const ChangePassword: React.FC = () => {
  return (
    <PageContainer>
      <ProCard>
        <Alert type="info" showIcon message="使用默认密码的账号须先设置独立密码。修改后所有旧登录会话将失效。" style={{ marginBottom: 24 }} />
        <ProForm
          style={{ maxWidth: 520 }}
          onFinish={async (values) => {
            const res = await changeCurrentPassword({
              old_password: values.old_password,
              new_password: values.new_password,
            });
            if (res.code === 200) {
              message.success("密码修改成功，请重新登录");
              localStorage.removeItem("shiyu_token");
              window.location.assign("/user/login");
              return true;
            }
            return false;
          }}
        >
          <ProFormText.Password
            name="old_password"
            label="原密码"
            rules={[{ required: true, message: "请输入原密码" }]}
          />
          <ProFormText.Password
            name="new_password"
            label="新密码"
            rules={[
              { required: true, message: "请输入新密码" },
              { min: 12, max: 72, message: "密码至少 12 位，最多 72 字节" },
            ]}
          />
          <ProFormText.Password
            name="confirm_password"
            label="确认新密码"
            dependencies={["new_password"]}
            rules={[
              { required: true, message: "请确认新密码" },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue("new_password") === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error("两次输入的密码不一致"));
                },
              }),
            ]}
          />
        </ProForm>
      </ProCard>
    </PageContainer>
  );
};

export default ChangePassword;
