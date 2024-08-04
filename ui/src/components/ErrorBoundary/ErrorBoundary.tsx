import React, { ReactElement, useState, useEffect } from 'react';

interface ErrorBoundaryProps {
  children: ReactElement;
  fallbackComponent?: ReactElement;
}

const ErrorBoundary: React.FC<ErrorBoundaryProps> = ({
  children,
  fallbackComponent,
}) => {
  const [hasError, setHasError] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [info, setInfo] = useState<React.ErrorInfo | null>(null);

  useEffect(() => {
    window.addEventListener('error', errorBoundary);
    return () => {
      window.removeEventListener('error', errorBoundary);
    };
  }, []);

  const handleError = (error: Error, info: React.ErrorInfo): void => {
    setHasError(true);
    setError(error);
    setInfo(info);
  };

  const errorBoundary = (event: ErrorEvent): void => {
    handleError(event.error, { componentStack: '' });
  };

  return hasError ? (fallbackComponent ? fallbackComponent : null) : children;
};

export default ErrorBoundary;
