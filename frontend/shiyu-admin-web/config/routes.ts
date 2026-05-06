export default [
  {
    path: "/user",
    layout: false,
    routes: [{ path: "/user/login", component: "./user/login" }],
  },
  {
    path: "/account",
    hideInMenu: true,
    routes: [
      { path: "/account/center", component: "./account/center" },
      { path: "/account/password", component: "./account/password" },
    ],
  },
  { path: "/welcome", icon: "smile", component: "./Welcome" },
  { path: "/dashboard", icon: "dashboard", component: "./Dashboard" },
  {
    path: "/monitor",
    icon: "monitor",
    routes: [
      { path: "/monitor/online-users", component: "./monitor/online-users" },
      { path: "/monitor/service", component: "./monitor/service" },
      { path: "/monitor/cache", component: "./monitor/cache" },
      { path: "/monitor/data", component: "./monitor/data" },
    ],
  },
  { path: "/cockpit", redirect: "/dashboard" },
  {
    path: "/system",
    icon: "setting",
    routes: [
      { path: "/system/user", component: "./system/user" },
      { path: "/system/role", component: "./system/role" },
      { path: "/system/menu", component: "./system/menu" },
      { path: "/system/dept", component: "./system/dept" },
      { path: "/system/operation-log", component: "./system/operation-log" },
    ],
  },
  { path: "/", redirect: "/welcome" },
  { path: "*", layout: false, component: "./404" },
];
