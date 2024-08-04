import { FC } from 'react';
import { useSelector } from 'react-redux';
import { Typography } from '@navi/web-ui/lib/primitives';
import GiveKleosForm from '@src/Pages/Common/GiveKleosForm';
import { getFromLocalStorage } from '@src/service/storage';
import starStruckIcon from '@src/assets/icons/starStruck.svg';
import {
  selectAllAchievementData,
  selectAllUserData,
} from '@store/Dashboard/DashboardSlice';
import styles from './FirstPage.module.scss';

interface FirstPageProps {
  userId: string;
}
const userName: string = getFromLocalStorage('user-name') || '';

const FirstGiveKudos: FC<FirstPageProps> = ({ userId }) => {
  const allUsers = useSelector(selectAllUserData).data;
  const allAchievements = useSelector(selectAllAchievementData).data;

  return (
    <div className={styles['mainDiv']}>
      <div className={styles['childDiv']}>
        <img
          src={starStruckIcon}
          alt="giveKleosAlt"
          className={styles['starWrapper']}
        />
        <div style={{ margin: '16px 0' }}></div>

        <div className={styles['myNameDiv']}>
          <Typography variant="h2" color={'#1C1C1C'}>
            {userName},
          </Typography>
        </div>
        <div className={styles['paraDiv']}>
          <Typography variant={'p2'} color={'#657384'}>
            Start your Kudos adventure by giving Kudos to your rockstar
            shipmates! 🚀
          </Typography>
        </div>

        <div style={{ margin: '12px 0' }}></div>

        <div className={styles['giveKleosForm']}>
          <GiveKleosForm
            toOptions={allUsers}
            achievementOptions={allAchievements}
            userId={userId}
          />
        </div>
      </div>
    </div>
  );
};

export default FirstGiveKudos;
