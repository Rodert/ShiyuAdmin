import type { DatabaseStats } from "@/services/shiyu-api/monitor";
import { getDatabaseStats } from "@/services/shiyu-api/monitor";
import { renderCellText } from "@/utils/tableRender";
import type { ProColumns } from "@ant-design/pro-components";
import { PageContainer, ProTable } from "@ant-design/pro-components";
import React from "react";

const ServiceMonitorPage: React.FC = () => {
  const columns: ProColumns<DatabaseStats>[] = [
    {
      title: "状态",
      dataIndex: "status",
      key: "status",
      valueEnum: {
        normal: { text: "正常", status: "Success" },
        abnormal: { text: "异常", status: "Error" },
        disabled: { text: "未启用", status: "Default" },
      },
    },
    {
      title: "驱动",
      dataIndex: "driver",
      key: "driver",
    },
    {
      title: "数据库",
      dataIndex: "database",
      key: "database",
      render: renderCellText<DatabaseStats, "database">(
        "database",
        (value) => value || "-"
      ),
    },
    {
      title: "版本",
      dataIndex: "version",
      key: "version",
      ellipsis: true,
      render: renderCellText<DatabaseStats, "version">(
        "version",
        (value) => value || "-"
      ),
    },
    {
      title: "表数量",
      dataIndex: "table_count",
      key: "table_count",
    },
    {
      title: "连接数",
      dataIndex: "open_connections",
      key: "open_connections",
    },
    {
      title: "使用中",
      dataIndex: "in_use",
      key: "in_use",
    },
    {
      title: "空闲",
      dataIndex: "idle",
      key: "idle",
    },
  ];

  return (
    <PageContainer>
      <ProTable<DatabaseStats>
        headerTitle="服务监控"
        rowKey="database"
        search={false}
        options={false}
        pagination={false}
        toolBarRender={false}
        request={async () => {
          const res = await getDatabaseStats();
          return {
            data: res.code === 200 && res.data ? [res.data] : [],
            success: res.code === 200,
          };
        }}
        columns={columns}
      />
    </PageContainer>
  );
};

export default ServiceMonitorPage;
