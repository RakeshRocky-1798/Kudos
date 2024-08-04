import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '@src/hooks/useAuth';

interface SecureProps {
  permissions: string[];
  element: React.ReactElement;
}

const Secure: React.FC<SecureProps> = ({ permissions, element }) => {
  const { hasAnyPermission } = useAuth();

  const hasRequiredPermissions = hasAnyPermission(permissions);

  if (!hasRequiredPermissions) {
    return <Navigate to="/unauthorized" replace />;
  }

  return element;
};

export default Secure;
