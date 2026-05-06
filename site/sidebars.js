// @ts-check

const sidebars = {
  tutorialSidebar: [
    'intro',
    {
      type: 'category',
      label: '开箱使用',
      items: [
        'getting-started/quick-start',
        'getting-started/accounts',
        'getting-started/screenshots',
        'getting-started/project-structure',
      ],
    },
    {
      type: 'category',
      label: '功能模块',
      items: [
        'features/overview',
        'features/rbac',
        'features/dashboard',
        'features/system-management',
        'features/monitor-cache-log',
      ],
    },
    {
      type: 'category',
      label: '开发指南',
      items: [
        'development/architecture',
        'development/backend',
        'development/frontend',
        'development/database',
        'development/testing',
      ],
    },
    {
      type: 'category',
      label: '部署上线',
      items: [
        'deployment/docker-compose',
        'deployment/github-pages',
        'deployment/render',
        'deployment/fly',
        'deployment/free-deployment',
      ],
    },
    {
      type: 'category',
      label: '资源与创收',
      items: [
        'resources/cloud-offers',
        'resources/official-links',
      ],
    },
  ],
};

module.exports = sidebars;
