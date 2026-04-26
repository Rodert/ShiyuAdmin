import React from 'react';

const Footer: React.FC = () => (
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

export default Footer;
