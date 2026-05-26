export default [
  {
    path: '/user',
    layout: false,
    routes: [{ path: '/user/login', component: './user/login' }],
  },
  {
    path: '/account',
    hideInMenu: true,
    access: 'authenticated',
    routes: [
      { path: '/account/center', component: './account/center' },
      { path: '/account/password', component: './account/password' },
    ],
  },
  {
    path: '/welcome',
    icon: 'smile',
    component: './Welcome',
    access: 'canViewWelcome',
  },
  {
    path: '/dashboard',
    icon: 'dashboard',
    component: './Dashboard',
    access: 'canViewDashboard',
  },
  {
    path: '/monitor',
    icon: 'monitor',
    access: 'canViewMonitorSection',
    routes: [
      {
        path: '/monitor/online-users',
        component: './monitor/online-users',
        access: 'canViewMonitor',
      },
      {
        path: '/monitor/service',
        component: './monitor/service',
        access: 'canViewMonitor',
      },
      {
        path: '/monitor/cache',
        component: './monitor/cache',
        access: 'canViewCache',
      },
      {
        path: '/monitor/data',
        component: './monitor/data',
        access: 'canViewData',
      },
    ],
  },
  { path: '/cockpit', redirect: '/dashboard' },
  {
    path: '/system',
    icon: 'setting',
    access: 'canViewSystem',
    routes: [
      {
        path: '/system/user',
        component: './system/user',
        access: 'canViewSystemUser',
      },
      {
        path: '/system/role',
        component: './system/role',
        access: 'canViewSystemRole',
      },
      {
        path: '/system/menu',
        component: './system/menu',
        access: 'canViewSystemMenu',
      },
      {
        path: '/system/dept',
        component: './system/dept',
        access: 'canViewSystemDept',
      },
      {
        path: '/system/operation-log',
        component: './system/operation-log',
        access: 'canViewOperationLog',
      },
    ],
  },
  { path: '/', redirect: '/welcome' },
  { path: '*', layout: false, component: './404' },
];
