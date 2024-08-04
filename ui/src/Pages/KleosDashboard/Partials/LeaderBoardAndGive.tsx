import { FC } from 'react';
import { useSelector } from 'react-redux';
import Lottie from 'lottie-react';

import cardLoader from '@src/assets/lotties/cardLoader.json';
import GiveKleosBlock from '@src/Pages/KleosDashboard/Partials/GiveKleosBlock';
import { selectIsLoading } from '@src/store/Dashboard/DashboardSlice';
import LeaderBoardTable from './LeaderBoardTable';
import styles from './LeaderBoardAndGive.module.scss';

const LeaderBoardAndGive: FC = () => {
  const isLoading = useSelector(selectIsLoading);

  const returnContent = (): JSX.Element => {
    if (isLoading) {
      return <Lottie animationData={cardLoader} loop />;
    }
    return (
      <>
        <div className={styles['give-kleos-wrapper']}>
          <GiveKleosBlock />
        </div>
        <LeaderBoardTable />
      </>
    );
  };

  /*
        TODO: Also, consume the selectLeaderBoardError selector and handle the leaderBoard API error scenario here
        Note: Please create a common retry component under the /components folder for the retry logic and re-use it here
    */

  return <div className={styles['leaderboard-wrapper']}>{returnContent()}</div>;
};

export default LeaderBoardAndGive;
