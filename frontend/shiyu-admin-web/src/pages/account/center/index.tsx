import { updateCurrentProfile } from "@/services/shiyu-api/user";
import {
  PageContainer,
  ProCard,
  ProForm,
  type ProFormInstance,
  ProFormText,
} from "@ant-design/pro-components";
import { useModel } from "@umijs/max";
import { Avatar, Descriptions, message, Space, Tag, Typography } from "antd";
import React, { useEffect, useRef, useState } from "react";

const AccountCenter: React.FC = () => {
  const { initialState, setInitialState } = useModel("@@initialState");
  const currentUser = initialState?.currentUser;
  const formRef = useRef<ProFormInstance>(null);
  const [avatarPreview, setAvatarPreview] = useState(
    currentUser?.avatar || "/logo-v2.png"
  );

  useEffect(() => {
    const avatar = currentUser?.avatar || "/logo-v2.png";
    setAvatarPreview(avatar);
    formRef.current?.setFieldsValue({
      nickname: currentUser?.nickname,
      email: currentUser?.email,
      phone: currentUser?.phone,
      avatar,
    });
  }, [currentUser]);

  const updateAvatar = (avatar: string) => {
    const nextAvatar = avatar || "/logo-v2.png";
    setAvatarPreview(nextAvatar);
    formRef.current?.setFieldValue("avatar", nextAvatar);
  };

  return (
    <PageContainer>
      <ProCard>
        <div
          style={{
            display: "flex",
            gap: 32,
            alignItems: "flex-start",
            flexWrap: "wrap",
          }}
        >
          <Space direction="vertical" align="center" size="middle">
            <Avatar size={112} src={avatarPreview} />
            <Typography.Text type="secondary">
              填写头像 URL 后保存
            </Typography.Text>
          </Space>
          <ProForm
            formRef={formRef}
            style={{ maxWidth: 560, flex: 1 }}
            initialValues={{
              nickname: currentUser?.nickname,
              email: currentUser?.email,
              phone: currentUser?.phone,
              avatar: currentUser?.avatar || "/logo-v2.png",
            }}
            onFinish={async (values) => {
              const res = await updateCurrentProfile(values);
              if (res.code === 200 && res.data) {
                setInitialState((state) => ({
                  ...state,
                  currentUser: {
                    ...state?.currentUser,
                    nickname: res.data.nickname,
                    email: res.data.email,
                    phone: res.data.phone,
                    avatar: res.data.avatar || "/logo-v2.png",
                  },
                }));
                setAvatarPreview(res.data.avatar || "/logo-v2.png");
                message.success("保存成功");
                return true;
              }
              return false;
            }}
          >
            <ProFormText name="nickname" label="昵称" />
            <ProFormText
              name="email"
              label="邮箱"
              rules={[{ type: "email", message: "邮箱格式不正确" }]}
            />
            <ProFormText name="phone" label="手机号" />
            <ProFormText
              name="avatar"
              label="头像"
              placeholder="/logo-v2.png"
              extra="填写图片 URL 或站内静态资源路径，例如 /logo-v2.png；不要粘贴 base64 图片内容。"
              rules={[{ max: 255, message: "头像地址不能超过 255 个字符" }]}
              fieldProps={{
                onChange: (event) => updateAvatar(event.target.value),
              }}
            />
          </ProForm>
        </div>
        <Descriptions
          bordered
          column={{ xs: 1, sm: 1, md: 2 }}
          size="middle"
          style={{ marginTop: 24 }}
        >
          <Descriptions.Item label="用户名">
            {currentUser?.name || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="用户编码">
            {currentUser?.userid || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="部门">
            {currentUser?.deptName || currentUser?.deptCode || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            {currentUser?.status === undefined ? (
              "-"
            ) : (
              <Tag color={currentUser.status === 1 ? "green" : "red"}>
                {currentUser.status === 1 ? "启用" : "停用"}
              </Tag>
            )}
          </Descriptions.Item>
          <Descriptions.Item label="账号类型">
            <Tag color={currentUser?.isSuperAdmin ? "blue" : "default"}>
              {currentUser?.isSuperAdmin ? "超级管理员" : "普通用户"}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="角色" span={2}>
            {currentUser?.roles && currentUser.roles.length > 0
              ? currentUser.roles.map((role) => (
                  <Tag key={role.role_code || role.role_key}>
                    {role.role_name || role.role_key || role.role_code}
                  </Tag>
                ))
              : "-"}
          </Descriptions.Item>
        </Descriptions>
      </ProCard>
    </PageContainer>
  );
};

export default AccountCenter;
