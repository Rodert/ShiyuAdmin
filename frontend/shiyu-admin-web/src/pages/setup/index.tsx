import { getInstallation, installSystem, type InstallationStatus } from '@/services/shiyu-api/installation';
import { Button, Card, Form, Input, Alert, Space, Tag, Typography, message } from 'antd';
import React, { useEffect, useState } from 'react';

const Setup: React.FC = () => {
  const [status, setStatus] = useState<InstallationStatus>();
  const [error, setError] = useState(false);
  const [saving, setSaving] = useState(false);
  const check = async () => {
    setError(false);
    try {
      const res = await getInstallation();
      setStatus(res.data);
      if (res.data.installed) window.location.replace('/user/login');
    } catch { setError(true); }
  };
  useEffect(() => { void check(); }, []);
  return (
    <main style={{ minHeight: '100vh', padding: '48px 20px', background: '#f0f5ff' }}>
      <Card style={{ maxWidth: 620, margin: '0 auto' }}>
        <Typography.Title level={2}>欢迎使用 Shiyu Admin</Typography.Title>
        <Typography.Paragraph type="secondary">完成首次设置，创建属于你的管理系统。</Typography.Paragraph>
        <Space style={{ marginBottom: 24 }}>
          <Tag color={status?.checks.database ? 'success' : 'default'}>数据库：{status?.checks.database ? '已连接' : '检查中'}</Tag>
          <Tag color={status?.checks.storage ? 'success' : 'default'}>文件存储：{status?.checks.storage ? '可用' : '待检查'}</Tag>
          <Button onClick={check}>重新检查</Button>
        </Space>
        {error && <Alert type="error" showIcon message="无法连接服务，请检查后端是否启动后重试。" style={{ marginBottom: 16 }} />}
        {status && !status.checks.storage && <Alert type="error" showIcon message="文件目录不可写，请检查持久化目录权限。" style={{ marginBottom: 16 }} />}
        <Form layout="vertical" initialValues={{ name: 'Shiyu Admin', username: 'admin' }} onFinish={async (values) => {
          setSaving(true);
          try {
            await installSystem({ token: values.token, name: values.name, username: values.username, password: values.password });
            localStorage.removeItem('shiyu_token');
            message.success('安装成功，请使用刚创建的管理员账号登录');
            window.location.assign('/user/login');
          } catch { /* Global request handler displays the server error. */ }
          finally { setSaving(false); }
        }}>
          <Form.Item label="安装授权码" name="token" rules={[{ required: true, message: '请输入服务器上的安装授权码' }]} extra={<>在服务器执行 <Typography.Text code>docker compose exec shiyu-app cat /app/data/setup-token</Typography.Text> 获取。使用其他部署配置时添加对应的 -f 参数；本地开发请读取 data 目录下的 setup-token 文件。</>}>
            <Input.Password autoComplete="off" />
          </Form.Item>
          <Form.Item label="系统名称" name="name" rules={[{ required: true, whitespace: true }, { max: 64 }]}><Input maxLength={64} /></Form.Item>
          <Form.Item label="管理员账号" name="username" rules={[{ required: true }, { pattern: /^[A-Za-z0-9_][A-Za-z0-9_.-]{2,63}$/, message: '请输入 3–64 位字母、数字或 ._-，以字母、数字或下划线开头' }]}><Input autoComplete="username" /></Form.Item>
          <Form.Item label="管理员密码" name="password" rules={[{ required: true }, { validator: async (_, value) => { if (value && (new TextEncoder().encode(value).length < 12 || new TextEncoder().encode(value).length > 72)) throw new Error('密码须为 12–72 字节'); } }]}><Input.Password autoComplete="new-password" /></Form.Item>
          <Form.Item label="确认密码" name="confirm" dependencies={['password']} rules={[{ required: true }, ({ getFieldValue }) => ({ validator: async (_, value) => { if (value !== getFieldValue('password')) throw new Error('两次输入的密码不一致'); } })]}><Input.Password autoComplete="new-password" /></Form.Item>
          <Button type="primary" htmlType="submit" block size="large" loading={saving} disabled={!status?.checks.database || !status?.checks.storage}>完成安装</Button>
        </Form>
      </Card>
    </main>
  );
};

export default Setup;
