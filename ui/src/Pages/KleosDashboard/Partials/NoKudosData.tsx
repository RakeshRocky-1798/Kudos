import { FC } from 'react';

import { Typography } from '@navi/web-ui/lib/primitives';

import styles from './MyDetailsDiv.module.scss';

const NoKudosData: FC = () => {
  return (
    <div className={styles['no-kudos-received-wrapper']}>
      <div className={styles['noTypoWrapper']}>
        <Typography variant="h3" className={styles['noTypo-emoji']}>
          😅
        </Typography>
        <Typography variant="h3" className={styles['noTypo']}>
          Uh oh!
        </Typography>
      </div>
      <div className={styles['noTypo2Wrapper']}>
        <Typography
          variant="p3"
          className={styles['noTypo2']}
          color={'var(--navi-color-navigation-blue-border)'}
        >
          It seems your kudos mailbox is empty for now.
        </Typography>
        <Typography
          variant="p3"
          color={'var(--navi-color-navigation-blue-border)'}
        >
          Keep being amazing, and those virtual high-fives will start rolling
          in!
        </Typography>
      </div>
      {/* <div className={styles['noTypo3Wrapper']}>
        <div className={styles['noTypo4Wrapper']}>
          <Typography
            variant="p3"
            className={styles['noTypo2']}
            color={'#657384'}
          >
            {`Sun peeks through, shadows fly, Chin up high, reach for the sky! 🌞🐦`}
          </Typography>
        </div>

        <div className={styles['noTypo5Wrapper']}>
          <Typography
            variant="p3"
            className={styles['noTypo2']}
            color={'#657384'}
          >
            {`Blooms unfurl, bees abuzz, Happy heart, life's a sweet fuzz. 🌸🐝`}
          </Typography>
        </div>

        <div className={styles['noTypo6Wrapper']}>
          <Typography
            variant="p3"
            className={styles['noTypo2']}
            color={'#657384'}
          >
            {`Raindrops dance, laughter rings, Joyful soul, on rainbow wings. 🌈🎶`}
          </Typography>
        </div>
      </div> */}
    </div>
  );
};

export default NoKudosData;
