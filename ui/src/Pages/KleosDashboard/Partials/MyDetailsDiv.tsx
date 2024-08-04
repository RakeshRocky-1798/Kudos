import React, { FC, useCallback, useMemo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import Lottie from 'lottie-react';

import { Button, Typography } from '@navi/web-ui/lib/primitives';
import { ArrowRightIcon } from '@navi/web-ui/lib/icons';
import { Grid } from '@navi/web-ui/lib/layouts';

import ClapIcon from '@src/assets/icons/clap.svg';
import ThinkIcon from '@src/assets/icons/thinking.svg';
import RaiseHand from '@src/assets/icons/raiseHand.svg';
import cardLoader from '@src/assets/lotties/cardLoader.json';
import { getFromLocalStorage } from '@src/service/storage';
import {
  selectDashboardData,
  selectIsLoading,
} from '@store/Dashboard/DashboardSlice';
import {
  AchievementData,
  recentRecognitionData,
} from '@src/Pages/KleosDashboard/types';
import { setCurrentTab } from '@store/UserKudos/UserSlice';

import NoKudosData from './NoKudosData';
import RecentRecognition from './RecentRecognition';
import TotalAchievements from './TotalAchievements';
import styles from './MyDetailsDiv.module.scss';

const { GridContainer, GridRow, GridColumn } = Grid;

const MyDetailsDiv: FC = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const isLoading: boolean = useSelector(selectIsLoading);
  const myRecognition: Array<recentRecognitionData> =
    useSelector(selectDashboardData).recentRecognition || [];
  const myKleosData: AchievementData[] =
    useSelector(selectDashboardData).totalAchievement || [];
  const userKudosReceivedCount =
    useSelector(selectDashboardData).kleosMetrics.receivedCount || 0;
  const userKudosGivenCount =
    useSelector(selectDashboardData).kleosMetrics.givenCount || 0;
  const userName: string = getFromLocalStorage('user-name') || '';

  const onLeftButtonClick = (): void => {
    dispatch(setCurrentTab('received'));
    navigate('/user-kudos?fromDashboard=true');
  };

  const onRightButtonClick = (): void => {
    dispatch(setCurrentTab('given'));
    navigate('/user-kudos?fromDashboard=true');
  };

  const CardLoader: JSX.Element = useMemo(() => {
    return (
      <div className={styles['lottie-container']}>
        <Lottie animationData={cardLoader} loop />
      </div>
    );
  }, []);

  const returnCardData = (isRecognition = true): JSX.Element => {
    if (isLoading) {
      return CardLoader;
    }
    return isRecognition ? <RecentRecognition /> : <TotalAchievements />;
  };

  const returnReviewedBlock = useCallback(
    (userKudosReceived: string, isReceived = true) => {
      const msg: string = isReceived ? 'Received' : 'Given';
      const fromto: string = isReceived ? 'from' : 'to';
      if (userKudosReceived.toString() === '0') {
        return (
          <>
            <img
              src={ThinkIcon}
              alt="giveKleosAlt"
              className={styles['thinkWrapper']}
            />
            <Typography
              variant="p1"
              color="var(--navi-color-gray-c2)"
              style={{ fontSize: '14px' }}
            >
              No Kudos {msg} Yet
            </Typography>
          </>
        );
      }

      return (
        <>
          <div className={styles['button-wrapper']}>
            <Button
              variant="text"
              className={styles['button-container']}
              onClick={isReceived ? onLeftButtonClick : onRightButtonClick}
            >
              View
              <ArrowRightIcon color="var(--navi-color-blue-base)" />
            </Button>
          </div>
          <img
            src={isReceived ? ClapIcon : RaiseHand}
            alt="giveKleosAlt"
            className={styles['clapWrapper']}
          />
          <div>
            <Typography
              variant="h3"
              color="var(--navi-color-gray-c1)"
              style={{ display: 'inline' }}
            >
              {userKudosReceived}
            </Typography>{' '}
            <Typography variant="p2" style={{ display: 'inline' }}>
              Kudos
            </Typography>
          </div>

          <Typography variant="p3" color="var(--navi-color-gray-c2)">
            {msg} {fromto} others
          </Typography>
        </>
      );
    },
    [],
  );

  const returnHeaderBlock = useCallback(
    (userKudosReceived: string) => {
      if (userKudosReceived.toString() === '0') {
        return (
          <div className={styles['headerWrapper']}>
            <Typography
              variant="h1"
              className={styles['yourNameTypography']}
              color="var(--navi-color-gray-c1)"
              style={{ fontSize: '18px' }}
            >
              Hey, {userName} 👋
            </Typography>
            <Typography
              variant="p3"
              color="var(--navi-color-navigation-blue-border)"
            >
              You&apos;ve just started a culture of appreciation and recognition
              at Navi! Yay!
            </Typography>
          </div>
        );
      }
      return (
        <div className={styles['headerWrapper']}>
          <Typography
            variant="h1"
            className={styles['yourNameTypography']}
            color="var(--navi-color-gray-c1)"
            style={{ fontSize: '18px' }}
          >
            Hey, {userName} 👋
          </Typography>
          <Typography
            variant="p3"
            className={styles['yourDescTypography']}
            color="var(--navi-color-navigation-blue-border)"
          >
            Let&apos;s keep the appreciation flowing freely!
          </Typography>
        </div>
      );
    },
    [userName],
  );

  return (
    <GridContainer>
      <GridRow>
        <GridColumn xs={12}>
          {isLoading ? (
            <></>
          ) : (
            returnHeaderBlock(userKudosReceivedCount.toString())
          )}
        </GridColumn>
      </GridRow>

      <GridRow className={styles['kudos-content-block-wrapper']}>
        <GridColumn xs className={styles['kudos-count-block']}>
          {!isLoading
            ? returnReviewedBlock(userKudosReceivedCount.toString())
            : CardLoader}
        </GridColumn>
        <GridColumn xs className={styles['kudos-count-block']}>
          {!isLoading
            ? returnReviewedBlock(userKudosGivenCount.toString(), false)
            : CardLoader}
        </GridColumn>
      </GridRow>
      {!myRecognition.length || !myKleosData.length ? (
        <GridRow
          className={styles['kudos-content-block-wrapper']}
          style={{ margin: '0 0 16px' }}
        >
          <GridColumn xs>
            <NoKudosData />
          </GridColumn>
        </GridRow>
      ) : (
        <GridRow className={styles['kudos-content-block-wrapper']}>
          <GridColumn xs className={styles['kudos-data-block']}>
            {returnCardData()}
          </GridColumn>
          <GridColumn xs className={styles['kudos-data-block']}>
            {returnCardData(false)}
          </GridColumn>
        </GridRow>
      )}
    </GridContainer>
  );
};

export default MyDetailsDiv;
