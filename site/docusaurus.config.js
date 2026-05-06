// @ts-check

const config = {
  title: 'Shiyu Admin',
  tagline: 'Go + Gin + Gorm + React + Ant Design Pro + RBAC 的开源通用后台管理系统',
  favicon: 'img/logo.png',
  url: 'https://rodert.github.io',
  baseUrl: '/ShiyuAdmin/',
  organizationName: 'Rodert',
  projectName: 'ShiyuAdmin',
  deploymentBranch: 'gh-pages',
  trailingSlash: false,
  onBrokenLinks: 'throw',
  markdown: {
    hooks: {
      onBrokenMarkdownLinks: 'warn',
    },
  },
  i18n: {
    defaultLocale: 'zh-CN',
    locales: ['zh-CN'],
  },
  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.js',
          routeBasePath: 'docs',
          editUrl: 'https://github.com/Rodert/ShiyuAdmin/tree/main/site/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      },
    ],
  ],
  themeConfig: {
    image: 'img/home-img.png',
    navbar: {
      title: 'Shiyu Admin',
      logo: {
        alt: 'Shiyu Admin Logo',
        src: 'img/logo.png',
      },
      items: [
        {to: '/docs/intro', label: '文档', position: 'left'},
        {to: '/docs/getting-started/quick-start', label: '快速开始', position: 'left'},
        {to: '/docs/deployment/docker-compose', label: '部署', position: 'left'},
        {to: '/docs/resources/cloud-offers', label: '云服务器优惠', position: 'left'},
        {to: '/docs/resources/official-links', label: '组件官网', position: 'left'},
        {
          href: 'https://github.com/Rodert/ShiyuAdmin',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'light',
      links: [
        {
          title: '项目',
          items: [
            {label: '快速开始', to: '/docs/getting-started/quick-start'},
            {label: '演示图', to: '/docs/getting-started/screenshots'},
            {label: 'GitHub', href: 'https://github.com/Rodert/ShiyuAdmin'},
          ],
        },
        {
          title: '资源',
          items: [
            {label: '云服务器优惠', to: '/docs/resources/cloud-offers'},
            {label: '组件官网引用', to: '/docs/resources/official-links'},
            {label: '免费部署方案', to: '/docs/deployment/free-deployment'},
          ],
        },
        {
          title: '作者',
          items: [
            {label: 'JavaPub', href: 'https://javapub.net.cn/'},
            {label: 'Gitee 镜像', href: 'https://gitee.com/rodert/ShiyuAdmin'},
            {label: 'GitCode 镜像', href: 'https://gitcode.com/JavaPub/ShiyuAdmin'},
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Shiyu Admin. Apache-2.0 License.`,
    },
    prism: {
      theme: require('prism-react-renderer').themes.github,
      darkTheme: require('prism-react-renderer').themes.dracula,
    },
    colorMode: {
      defaultMode: 'light',
      disableSwitch: false,
      respectPrefersColorScheme: true,
    },
  },
};

module.exports = config;
