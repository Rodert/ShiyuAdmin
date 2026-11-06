import { getVersion, type VersionInfo } from '@/services/shiyu-api/installation';
import { PageContainer, ProCard } from '@ant-design/pro-components';
import { Alert, Button, Descriptions, Tabs, Typography } from 'antd';
import React, { useEffect, useState } from 'react';

const Version: React.FC = () => {
  const [info, setInfo] = useState<VersionInfo>();
  const [failed, setFailed] = useState(false);
  const load = async () => { setFailed(false); try { setInfo((await getVersion()).data); } catch { setFailed(true); } };
  useEffect(() => { void load(); }, []);
  return <PageContainer title="版本与升级"><ProCard loading={!info && !failed}>
    {failed && <Alert type="error" message="读取版本信息失败" action={<Button onClick={load}>重试</Button>} />}
    {info && <>
      <Descriptions column={1} items={[
        { key: 'version', label: '应用版本', children: info.version },
        { key: 'schema', label: '数据库版本', children: info.schema_version },
        { key: 'commit', label: '构建提交', children: info.commit },
        { key: 'time', label: '构建时间', children: info.build_time },
      ]} />
      <Tabs items={[
        { key: 'notes', label: '发布说明', children: <Typography.Paragraph style={{ whiteSpace: 'pre-wrap' }}>{info.release_notes}</Typography.Paragraph> },
        { key: 'upgrade', label: '升级与回退', children: <Typography.Paragraph style={{ whiteSpace: 'pre-wrap' }}>{info.upgrade_guide}</Typography.Paragraph> },
      ]} />
    </>}
  </ProCard></PageContainer>;
};
export default Version;
