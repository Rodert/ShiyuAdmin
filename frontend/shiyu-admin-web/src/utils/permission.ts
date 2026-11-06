export const hasPermission = (
  currentUser: API.CurrentUser | undefined,
  permission: string,
) => {
  if (!currentUser) {
    return false;
  }
  if (currentUser.isSuperAdmin) {
    return true;
  }
  if (permission === 'system:user:assign-role' ||
      (['system:role:', 'system:menu:', 'system:dept:'].some((prefix) => permission.startsWith(prefix)) && !permission.endsWith(':list'))) return false;
  return currentUser.permissions?.includes(permission) ?? false;
};
