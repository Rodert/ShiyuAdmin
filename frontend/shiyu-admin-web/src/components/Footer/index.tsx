import { useLocation } from '@umijs/max';
import React from 'react';

/** 仪表盘全屏沉浸式展示时隐藏全局页脚，避免底边出一条浅色横带 */
const Footer: React.FC = () => {
  const { pathname } = useLocation();
  const isDashboard = pathname === '/dashboard' || pathname.endsWith('/dashboard');
  if (isDashboard) {
    return null;
  }

  return (
    <footer
      style={{
        padding: '12px 0 20px',
        textAlign: 'center',
        fontSize: 12,
        color: 'rgba(0, 0, 0, 0.35)',
      }}
    >
      项目来源：
      <a
        href="https://github.com/Rodert/ShiyuAdmin"
        target="_blank"
        rel="noreferrer"
        style={{
          color: 'rgba(0, 0, 0, 0.45)',
        }}
      >
        GitHub
      </a>
    </footer>
  );
};

export default Footer;
