import type { OnlineUser } from "@/services/shiyu-api/monitor";
import {
  forceLogoutOnlineUser,
  getOnlineUsers,
} from "@/services/shiyu-api/monitor";
import type { ActionType, ProColumns } from "@ant-design/pro-components";
import { PageContainer, ProTable } from "@ant-design/pro-components";
import { Button, Descriptions, Drawer, message, Modal, Space } from "antd";
import React, { useRef, useState } from "react";

const formatTime = (timestamp?: number) => {
  if (!timestamp) return "-";
  return new Date(timestamp * 1000).toLocaleString();
};

const OnlineUsersPage: React.FC = () => {
  const actionRef = useRef<ActionType>(null);
  const [currentUser, setCurrentUser] = useState<OnlineUser>();
  const [drawerOpen, setDrawerOpen] = useState(false);

  const columns: ProColumns<OnlineUser>[] = [
    {
      title: "用户编码",
      dataIndex: "user_code",
      key: "user_code",
      width: 220,
      ellipsis: true,
      copyable: true,
    },
    {
      title: "用户名",
      dataIndex: "username",
      key: "username",
      width: 120,
    },
    {
      title: "部门名称",
      dataIndex: "dept_name",
      key: "dept_name",
      width: 140,
      ellipsis: true,
      renderText: (text) => text || "-",
    },
    {
      title: "会话编号",
      dataIndex: "session_id",
      key: "session_id",
      width: 170,
      ellipsis: true,
      copyable: true,
    },
    {
      title: "主机IP",
      dataIndex: "host_ip",
      key: "host_ip",
      width: 150,
      renderText: (text) => text || "-",
    },
    {
      title: "IP",
      dataIndex: "ip",
      key: "ip",
      width: 150,
      renderText: (text) => text || "-",
    },
    {
      title: "登录地点",
      dataIndex: "login_location",
      key: "login_location",
      width: 190,
      ellipsis: true,
      renderText: (text) => text || "-",
    },
    {
      title: "操作系统",
      dataIndex: "os",
      key: "os",
      width: 120,
      renderText: (text) => text || "-",
    },
    {
      title: "浏览器详细版本",
      dataIndex: "browser_detail",
      key: "browser_detail",
      width: 190,
      ellipsis: true,
      renderText: (text) => text || "-",
    },
    {
      title: "登录时间",
      dataIndex: "login_time",
      key: "login_time",
      width: 180,
      render: (_, record) => formatTime(record.login_time),
    },
    {
      title: "最后访问时间",
      dataIndex: "last_active",
      key: "last_active",
      width: 180,
      render: (_, record) => formatTime(record.last_active),
    },
    {
      title: "操作",
      key: "option",
      valueType: "option",
      fixed: "right",
      width: 140,
      render: (_, record) => [
        <Button
          key="view"
          type="link"
          size="small"
          onClick={() => {
            setCurrentUser(record);
            setDrawerOpen(true);
          }}
        >
          查看
        </Button>,
        <Button
          key="force"
          type="link"
          size="small"
          danger
          onClick={() => {
            Modal.confirm({
              title: "确认强退",
              content: `确定强制登出用户 ${
                record.username || record.user_code
              } 吗？`,
              onOk: async () => {
                const res = await forceLogoutOnlineUser(record.session_id);
                if (res.code === 200) {
                  message.success("已强制登出");
                  actionRef.current?.reload();
                }
              },
            });
          }}
        >
          强退
        </Button>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<OnlineUser>
        actionRef={actionRef}
        headerTitle="在线用户"
        rowKey="session_id"
        search={false}
        scroll={{ x: 1600 }}
        toolBarRender={false}
        request={async () => {
          const res = await getOnlineUsers();
          return {
            data: res.code === 200 && res.data ? res.data : [],
            success: res.code === 200,
          };
        }}
        columns={columns}
      />
      <Drawer
        title="在线用户详情"
        width={560}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
      >
        {currentUser && (
          <Descriptions column={1} bordered size="small">
            <Descriptions.Item label="用户名">
              {currentUser.username || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="用户编码">
              {currentUser.user_code || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="部门名称">
              {currentUser.dept_name || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="会话编号">
              {currentUser.session_id || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="主机IP">
              {currentUser.host_ip || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="IP">
              {currentUser.ip || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="登录地点">
              {currentUser.login_location || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="操作系统">
              {currentUser.os || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="浏览器名称">
              {currentUser.browser || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="浏览器版本">
              {currentUser.browser_version || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="浏览器详细版本">
              {currentUser.browser_detail || "-"}
            </Descriptions.Item>
            <Descriptions.Item label="登录时间">
              {formatTime(currentUser.login_time)}
            </Descriptions.Item>
            <Descriptions.Item label="最后访问时间">
              {formatTime(currentUser.last_active)}
            </Descriptions.Item>
            <Descriptions.Item label="User-Agent">
              <Space style={{ wordBreak: "break-all" }}>
                {currentUser.user_agent || "-"}
              </Space>
            </Descriptions.Item>
          </Descriptions>
        )}
      </Drawer>
    </PageContainer>
  );
};

export default OnlineUsersPage;
