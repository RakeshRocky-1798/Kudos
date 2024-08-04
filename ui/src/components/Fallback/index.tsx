import React from 'react';
import cx from 'classnames';
import LoadingIcon from '@navi/web-ui/lib/icons/LoadingIcon/LoadingIcon';
import styles from './Fallback.module.scss';

interface FallbackComponentProps {
  customClassName?: string;
  listLoader?: boolean;
  [x: string | number | symbol]: unknown;
}

const FallbackComponent: React.FC<FallbackComponentProps> = ({
  customClassName = '',
  listLoader = false,
  ...restProps
}) => {
  return (
    <div
      className={cx(styles['fallback-container'], customClassName, {
        [styles['fallback-container--list-loader']]: listLoader,
      })}
      {...restProps}
    >
      <LoadingIcon size="lg" />
    </div>
  );
};

export default FallbackComponent;
