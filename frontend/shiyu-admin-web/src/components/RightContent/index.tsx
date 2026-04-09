import { QuestionCircleOutlined } from '@ant-design/icons';
import { DownOutlined } from '@ant-design/icons';
import { getLocale, setLocale } from '@umijs/max';
import { Dropdown, Space } from 'antd';
import type { MenuProps } from 'antd';

export type SiderTheme = 'light' | 'dark';

export const SelectLang: React.FC = () => {
  const locale = getLocale();
  const labels: Record<string, string> = {
    'zh-CN': '简体中文',
    'en-US': 'English',
  };

  const onMenuClick: MenuProps['onClick'] = ({ key }) => {
    // 切换语言时不做整页刷新，避免丢失当前端口（如 :18000）
    setLocale(String(key), false);
  };

  return (
    <Dropdown
      menu={{
        selectedKeys: [locale],
        onClick: onMenuClick,
        items: [
          { key: 'zh-CN', label: '简体中文' },
          { key: 'en-US', label: 'English' },
        ],
      }}
      style={{
        padding: 4,
      }}
      trigger={['click']}
    >
      <a onClick={(e) => e.preventDefault()}>
        <Space>
          {labels[locale] || locale}
          <DownOutlined />
        </Space>
      </a>
    </Dropdown>
  );
};

export const Question: React.FC = () => {
  return (
    <a
      href="https://pro.ant.design/docs/getting-started"
      target="_blank"
      rel="noreferrer"
      style={{
        display: 'inline-flex',
        padding: '4px',
        fontSize: '18px',
        color: 'inherit',
      }}
    >
      <QuestionCircleOutlined />
    </a>
  );
};
