import React, { FC } from 'react';
import cx from 'classnames';

import { Typography } from '@navi/web-ui/lib/primitives';
import { TypographyProps } from '@navi/web-ui/lib/primitives/Typography/types';

import styles from './Unauthorized.module.scss';
import SomethingWentWrongIcon from '@src/assets/icons/SomethingWentWrongIcon';

type CustomTypographyProps = Omit<TypographyProps, 'children'>;

interface UnauthorizedProps extends CustomTypographyProps {
  Content?: React.ReactNode;
  message?: string;
  containerClass?: string;
  messageClass?: string;
}

const Unauthorized: FC<UnauthorizedProps> = ({
  Content = null,
  message = 'Please contact administrator to get access to this page.',
  containerClass = '',
  messageClass = '',
  ...typographyProps
}) => {
  if (Content) {
    return <>{Content}</>;
  }

  return (
    <div
      className={cx(styles['unauthorized-container'], {
        [containerClass]: containerClass,
      })}
    >
      <SomethingWentWrongIcon />
      <Typography
        className={cx(styles['unauthorized-message'], {
          [messageClass]: messageClass,
        })}
        {...typographyProps}
      >
        {message || 'Unauthorized Access'}
      </Typography>
    </div>
  );
};

export default Unauthorized;
