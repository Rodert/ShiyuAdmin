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
  return currentUser.permissions?.includes(permission) ?? false;
};
