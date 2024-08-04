import { FC, useCallback } from 'react';
import { useSelector } from 'react-redux';

import { Typography, Tag, Avatar, Tooltip } from '@navi/web-ui/lib/primitives';
import { TagColors } from '@navi/web-ui/lib/primitives/Tag/types';

import AchievementIcon from '@src/assets/icons/achievement.svg';
import { recentRecognitionData } from '@src/Pages/KleosDashboard/types';
import { selectDashboardData } from '@src/store/Dashboard/DashboardSlice';
import {
  returnAchievementIcon,
  returnAchieveTagVariantColor,
} from '../constants';
import styles from './RecentRecognition.module.scss';

const NoKudosReceived: FC = () => {
  return (
    <div className={styles['no-kudos-received-wrapper']}>
      <div>
        <Typography variant="h3">No Kudos received yet, eh?</Typography>
      </div>
      <div>
        <Typography variant="p3">
          Here&apos;s something to keep your spirits high:
        </Typography>
      </div>
      <div>
        <Typography variant="p3">
          {`Sun peeks through, shadows fly, Chin up high, reach for the sky! 🌞🐦`}
        </Typography>
        <Typography variant="p3">
          {`Blooms unfurl, bees abuzz, Happy heart, life's a sweet fuzz. 🌸🐝`}
        </Typography>
        <Typography variant="p3">
          {`Raindrops dance, laughter rings, Joyful soul, on rainbow wings. 🌈🎶`}
        </Typography>
      </div>
    </div>
  );
};

const RecentRecognition: FC = () => {
  const myRecognition: Array<recentRecognitionData> =
    useSelector(selectDashboardData).recentRecognition || [];

  const returnRecentRecognition = useCallback(() => {
    if (!myRecognition.length) {
      return (
        <div className={styles['recent-recognition-container']}>
          <NoKudosReceived />
        </div>
      );
    }

    return myRecognition.map(
      (recognition: recentRecognitionData, index: number) => (
        <div
          key={`${recognition.achievement.aLabel}-${index}`}
          className={styles['recent-recognition-container']}
        >
          <Typography
            variant="p4"
            color={'#585757'}
            className={styles['kudos-received-message']}
          >
            {recognition.message}
          </Typography>
          <div className={styles['recent-recognition-tag-wrapper']}>
            <Tag
              label={recognition.achievement.aLabel}
              variant="transparent"
              color={
                returnAchieveTagVariantColor(
                  recognition.achievement.aEmoji,
                ) as TagColors
              }
              endAdornment={returnAchievementIcon(
                recognition.achievement.aEmoji,
              )}
              size="sm"
            />
            <div
              className={styles['recent-recognition-tag-wrapper--user-data']}
            >
              {recognition.achievement.aFrom.profileUrl.length > 0 ? (
                <Avatar
                  src={recognition.achievement.aFrom.profileUrl}
                  alt="User Avatar"
                  sizeVariant="sm"
                  isImage
                />
              ) : (
                <Avatar size={24}>
                  {recognition.achievement.aFrom.userName
                    ?.charAt(0)
                    ?.toLocaleUpperCase()}
                </Avatar>
              )}
              <Tooltip
                text={recognition.achievement.aFrom.userName}
                position="bottom"
              >
                <Typography
                  variant="p5"
                  style={{ fontSize: '10px' }}
                  color={'#969696'}
                  className={styles['recent-recognition-user']}
                >
                  {recognition.achievement.aFrom.userName}
                </Typography>
              </Tooltip>
            </div>
          </div>
        </div>
      ),
    );
  }, [myRecognition]);

  // TODO: Handle the empty scenario
  return (
    <>
      <div className={styles['recent-recognition-wrapper']}>
        <div className={styles['imageWrap']}>
          <Typography
            variant={'p4'}
            className={styles['receivedTypo']}
            color={'#969696'}
            style={{ fontSize: '13px' }}
          >
            Here&apos;s what&apos;s getting you the appreciation!
          </Typography>
        </div>
        {returnRecentRecognition()}
      </div>
    </>
  );
};

export default RecentRecognition;
