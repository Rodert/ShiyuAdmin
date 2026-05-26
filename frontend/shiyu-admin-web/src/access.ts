/**
 * @see https://umijs.org/docs/max/access#access
 * */
export default function access(
  initialState: { currentUser?: API.CurrentUser } | undefined,
) {
  const { currentUser } = initialState ?? {};
  const hasPermission = (permission: string) => {
    if (!currentUser) {
      return false;
    }
    if (currentUser.isSuperAdmin) {
      return true;
    }
    return currentUser.permissions?.includes(permission) ?? false;
  };

  return {
    authenticated: !!currentUser,
    canAdmin: !!currentUser && currentUser.access === 'admin',
    canViewWelcome: hasPermission('welcome:view'),
    canViewDashboard: hasPermission('dashboard:view'),
    canViewSystem:
      hasPermission('system:user:list') ||
      hasPermission('system:role:list') ||
      hasPermission('system:menu:list') ||
      hasPermission('system:dept:list') ||
      hasPermission('system:operation-log:list'),
    canViewSystemUser: hasPermission('system:user:list'),
    canViewSystemRole: hasPermission('system:role:list'),
    canViewSystemMenu: hasPermission('system:menu:list'),
    canViewSystemDept: hasPermission('system:dept:list'),
    canViewOperationLog: hasPermission('system:operation-log:list'),
    canViewMonitorSection:
      hasPermission('system:monitor:view') ||
      hasPermission('system:cache:list') ||
      currentUser?.isSuperAdmin ||
      hasPermission('system:data:view'),
    canViewMonitor: hasPermission('system:monitor:view'),
    canViewCache: hasPermission('system:cache:list'),
    canViewData: currentUser?.isSuperAdmin || hasPermission('system:data:view'),
  };
}
