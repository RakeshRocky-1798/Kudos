import { FC, useEffect } from 'react';

import { Grid } from '@navi/web-ui/lib/layouts';

import MyDetailsDiv from '@src/Pages/KleosDashboard/Partials/MyDetailsDiv';
import LeaderBoardAndGive from '@src/Pages/KleosDashboard/Partials/LeaderBoardAndGive';

import useKleosDashboardApi from './useKleosDashboardApi';
import styles from './FinalKleosDashboard.module.scss';

const { GridContainer, GridRow, GridColumn } = Grid;

const KleosDashboard: FC = () => {
  const { startFetchingDashboardData } = useKleosDashboardApi();

  useEffect(() => {
    startFetchingDashboardData();
  }, []);

  return (
    <GridContainer>
      <GridRow className={styles['user-details-container']}>
        <GridColumn xs className={styles['user-details-content']}>
          <MyDetailsDiv />
        </GridColumn>
        <GridColumn xs className={styles['user-details-content']}>
          <LeaderBoardAndGive />
        </GridColumn>
      </GridRow>
    </GridContainer>
  );
};

export default KleosDashboard;
