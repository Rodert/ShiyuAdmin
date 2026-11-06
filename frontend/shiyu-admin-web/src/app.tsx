import NotificationBell from '@/components/NotificationBell';
import { getInstallation, type InstallationStatus } from '@/services/shiyu-api/installation';
import {
  AvatarDropdown,
  AvatarName,
  Footer,
  UserProfileAction,
} from '@/components';
import { currentUser as queryCurrentUser } from '@/services/ant-design-pro/api';
import { getMenuTree, type Menu as APIMenu } from '@/services/shiyu-api/menu';
import { LinkOutlined } from '@ant-design/icons';
import type {
  Settings as LayoutSettings,
  MenuDataItem,
} from '@ant-design/pro-components';
import { SettingDrawer } from '@ant-design/pro-components';
import '@ant-design/v5-patch-for-react-19';
import type { RequestConfig, RunTimeLayoutConfig } from '@umijs/max';
import { history, Link } from '@umijs/max';
import defaultSettings from '../config/defaultSettings';
import { errorConfig } from './requestErrorConfig';
const isDev = process.env.NODE_ENV === 'development';
const loginPath = '/user/login';
const themeSettingsStorageKey = 'shiyu_theme_settings';

const themeColorList = [
  { key: 'shiyu-blue', color: '#1677FF', title: '科技蓝' },
  { key: 'shiyu-cyan', color: '#08979C', title: '青绿' },
  { key: 'shiyu-green', color: '#389E0D', title: '活力绿' },
  { key: 'shiyu-orange', color: '#D46B08', title: '暖橙' },
  { key: 'shiyu-red', color: '#CF1322', title: '赤红' },
];

function getThemeSettings(): Partial<LayoutSettings> {
  if (typeof window === 'undefined') {
    return defaultSettings as Partial<LayoutSettings>;
  }

  try {
    const stored = localStorage.getItem(themeSettingsStorageKey);
    return stored
      ? { ...defaultSettings, ...JSON.parse(stored) }
      : (defaultSettings as Partial<LayoutSettings>);
  } catch (_error) {
    return defaultSettings as Partial<LayoutSettings>;
  }
}

function saveThemeSettings(settings: Partial<LayoutSettings>) {
  localStorage.setItem(themeSettingsStorageKey, JSON.stringify(settings));
}

/**
 * @see https://umijs.org/docs/api/runtime-config#getinitialstate
 * */
export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  installation?: InstallationStatus;
  currentUser?: API.CurrentUser;
  menuData?: MenuDataItem[];
  loading?: boolean;
  fetchUserInfo?: () => Promise<API.CurrentUser | undefined>;
}> {
  const settings = getThemeSettings();
  let installation: InstallationStatus;
  try { installation = (await getInstallation()).data; }
  catch {
    if (history.location.pathname !== '/setup') history.replace('/setup');
    return { settings };
  }
  if (!installation.installed) {
    if (history.location.pathname !== '/setup') history.replace('/setup');
    return { settings, installation };
  }
  if (history.location.pathname === '/setup') history.replace('/user/login');
  settings.title = installation.name;
  const buildMenuData = (menus: APIMenu[]): MenuDataItem[] => {
    return (menus || [])
      .filter((m) => m.status === 1 && m.menu_type !== 'F')
      .map((m) => ({
        name: m.menu_name,
        path: m.path || `/${m.menu_code}`,
        children: m.children ? buildMenuData(m.children) : undefined,
      }));
  };
  const fetchMenuData = async () => {
    try {
      const res = await getMenuTree();
      if (res.code === 200 && Array.isArray(res.data)) {
        return buildMenuData(res.data);
      }
    } catch (_e) {
      // ignore and fallback to static menu
    }
    return [];
  };
  const fetchUserInfo = async () => {
    try {
      const msg = await queryCurrentUser({
        skipErrorHandler: true,
      });
      return msg.data;
    } catch (_error) {
      history.push(loginPath);
    }
    return undefined;
  };
  // Public entry pages must not request authenticated APIs before login.
  if (!localStorage.getItem('shiyu_token')) {
    if (history.location.pathname !== loginPath) history.replace(loginPath);
    return { fetchUserInfo, installation, settings };
  }
  // 如果不是登录页面，执行
  const { location } = history;
  if (
    ![loginPath, '/user/register', '/user/register-result'].includes(
      location.pathname,
    )
  ) {
    const currentUser = await fetchUserInfo();
    const menuData = currentUser && !currentUser.requirePasswordChange ? await fetchMenuData() : [];
    return {
      fetchUserInfo,
      installation,
      currentUser,
      menuData,
      settings,
    };
  }
  return {
    fetchUserInfo,
    installation,
    settings,
  };
}

// ProLayout 支持的api https://procomponents.ant.design/components/layout
export const layout: RunTimeLayoutConfig = ({
  initialState,
  setInitialState,
}) => {
  return {
    actionsRender: () => [<NotificationBell key="notifications" />, <UserProfileAction key="profile" />],
    avatarProps: {
      src: initialState?.currentUser?.avatar,
      title: <AvatarName />,
      render: (_, avatarChildren) => {
        return <AvatarDropdown menu>{avatarChildren}</AvatarDropdown>;
      },
    },
    waterMarkProps: {
      content: initialState?.currentUser?.name,
    },
    footerRender: () => <Footer />,
    onPageChange: () => {
      const { location } = history;
      // 如果没有登录，重定向到 login
      if (!initialState?.installation?.installed) {
        if (location.pathname !== '/setup') history.replace('/setup');
        return;
      }
      if (initialState?.currentUser?.requirePasswordChange && location.pathname !== '/account/password') {
        history.replace('/account/password'); return;
      }
      if (!initialState?.currentUser && location.pathname !== loginPath) {
        history.push(loginPath);
      }
    },
    bgLayoutImgList: [
      {
        src: 'https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/D2LWSqNny4sAAAAAAAAAAAAAFl94AQBr',
        left: 85,
        bottom: 100,
        height: '303px',
      },
      {
        src: 'https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/C2TWRpJpiC0AAAAAAAAAAAAAFl94AQBr',
        bottom: -68,
        right: -45,
        height: '303px',
      },
      {
        src: 'https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/F6vSTbj8KpYAAAAAAAAAAAAAFl94AQBr',
        bottom: 0,
        left: 0,
        width: '331px',
      },
    ],
    links: isDev
      ? [
          <Link key="openapi" to="/umi/plugin/openapi" target="_blank">
            <LinkOutlined />
            <span>OpenAPI 文档</span>
          </Link>,
        ]
      : [],
    menuHeaderRender: undefined,
    menuDataRender: () => {
      if (!initialState?.currentUser || initialState.currentUser.requirePasswordChange) return [];
      const menus = (initialState.menuData || []).filter((m) => !['/welcome', '/work'].includes(m.path || '')).map((m) => ({ ...m }));
      if (initialState.currentUser.isSuperAdmin) {
        const system = menus.find((m) => m.path === '/system');
        const settings = { name: '系统设置', path: '/system/settings' };
        if (system) system.children = [...(system.children || []).filter((m) => m.path !== settings.path), settings];
        else menus.push({ name: '系统管理', path: '/system', children: [settings] });
      }
      return [{ name: '工作台', path: '/welcome' }, { name: '事项与审批', path: '/work/requests' }, { name: '站内通知', path: '/work/notifications' }, ...menus];
    },
    // 自定义 403 页面
    // unAccessible: <div>unAccessible</div>,
    // 增加一个 loading 的状态
    childrenRender: (children) => {
      // if (initialState?.loading) return <PageLoading />;
      return (
        <>
          {children}
          <SettingDrawer
            disableUrlParams
            enableDarkTheme
            themeOnly
            hideHintAlert
            hideCopyButton
            colorList={themeColorList}
            settings={initialState?.settings}
            onSettingChange={(settings) => {
              saveThemeSettings(settings);
              setInitialState((preInitialState) => ({
                ...preInitialState,
                settings,
              }));
            }}
          />
        </>
      );
    },
    ...initialState?.settings,
    title: initialState?.installation?.name || defaultSettings.title,
    logo: initialState?.installation?.logo_url || '/logo-v2.png',
  };
};

/**
 * @name request 配置，可以配置错误处理
 * 它基于 axios 和 ahooks 的 useRequest 提供了一套统一的网络请求和错误处理方案。
 * @doc https://umijs.org/docs/max/request#配置
 */
export const request: RequestConfig = {
  // 开发环境使用代理，生产环境需要配置实际的后端地址
  baseURL:
    process.env.NODE_ENV === 'production'
      ? process.env.REACT_APP_API_BASE_URL || ''
      : '',
  ...errorConfig,
};
