import { DarkKnightInstance, DarkKnightTokenParsed } from '@navi/dark-knight';
import { useContext } from 'react';
import AuthContext from '@src/context/AuthContext';
import { ROLES } from '@src/Pages/KleosDashboard/constants';

interface UseAuth {
  auth: DarkKnightInstance | null;
  hasPermission: (permission: string) => boolean;
  hasAnyPermission: (permissions: string[]) => boolean;
  hasPermissions: (permissions: string[]) => boolean;
  hasRole: (role: string) => boolean;
  isAdmin: boolean;
  isAllUsers: boolean;
  isAuthenticated: boolean;
  userInfo: DarkKnightTokenParsed | null;
}

const useAuth = (): UseAuth => {
  const auth = useContext(AuthContext);

  const hasRole = (roleName: string): boolean => {
    const filteredRoles = auth?.idTokenParsed?.roles?.filter(
      role => role === roleName,
    );

    return (filteredRoles || [])?.length > 0;
  };

  const hasPermission = (permissionName: string): boolean => {
    const filteredPermissions = auth?.idTokenParsed?.permissions?.filter(
      permission => permission === permissionName,
    );

    return (filteredPermissions || [])?.length > 0;
  };

  const hasAnyPermission = (permissions: string[]): boolean => {
    const filteredRoles = auth?.idTokenParsed?.permissions?.filter(role =>
      permissions.includes(role),
    );

    return (filteredRoles || [])?.length > 0;
  };

  const hasPermissions = (permissions: string[]): boolean => {
    const filteredPermissions = auth?.idTokenParsed?.permissions?.filter(
      permission => permissions.includes(permission),
    );

    return (filteredPermissions || [])?.length === permissions.length;
  };

  const isAdmin = hasRole(ROLES.ADMIN);
  const isAllUsers = hasRole(ROLES.ALLUSERS);

  return {
    auth,
    isAdmin,
    isAllUsers,
    isAuthenticated: auth?.authenticated || false,
    hasRole,
    hasPermission,
    hasAnyPermission,
    hasPermissions,
    userInfo: auth?.idTokenParsed || null,
  };
};

export { useAuth };
