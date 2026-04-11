export default [
  { path: '/user', layout: false, routes: [{ path: '/user/login', component: './user/login' }] },
  { path: '/welcome', icon: 'smile', component: './Welcome' },
  {
    path: '/system',
    icon: 'setting',
    routes: [
      { path: '/system/user', component: './system/user' },
      { path: '/system/role', component: './system/role' },
      { path: '/system/menu', component: './system/menu' },
      { path: '/system/dept', component: './system/dept' },
      { path: '/system/operation-log', component: './system/operation-log' },
      { path: '/system/monitor', component: './system/monitor' },
      { path: '/system/data-manage', component: './system/data-manage' },
    ],
  },
  { path: '/', redirect: '/welcome' },
  { path: '*', layout: false, component: './404' },
];
