import React, { FC } from 'react';
import Lottie from 'react-lottie';
import successAnimation from '@src/assets/lotties/lottie.json';
import styles from './SuccessAnimation.module.scss';

const SuccessAnimation: FC = () => {
  const defaultOptions = {
    loop: true,
    autoplay: true,
    animationData: successAnimation,
    rendererSettings: {
      scaleMode: 'fit',
      preserveAspectRatio: 'xMidYMid slice',
    },
  };

  return (
    <div className={styles['lottie']}>
      <Lottie
        options={defaultOptions}
        style={{ width: 'inherit', height: 'inherit' }}
      />
    </div>
  );
};

export default SuccessAnimation;
