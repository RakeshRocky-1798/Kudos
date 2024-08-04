import { FC, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Button, Card } from '@navi/web-ui/lib/primitives';
import { ArrowLeftIcon } from '@navi/web-ui/lib/icons';

import UserKudosTabs from '@src/Pages/UserKleos/partials/UserKudosTabs';
import FallbackComponent from '@src/components/Fallback';
import {
  selectIsLoading,
  selectCurrentTab,
  resetTabsData,
  setCurrentTab,
} from '@src/store/UserKudos/UserSlice';

import useUserKudosApi from './useUserKudosApi';
import GiveKleosTypoButton from './partials/GiveKleosTypoButton';
import styles from './FinalGiveKleos.module.scss';

const UserKleos: FC = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const isLoading = useSelector(selectIsLoading);
  const currentTab = useSelector(selectCurrentTab);
  const { fetchUserKudosData } = useUserKudosApi();
  const [searchParams] = useSearchParams();

  useEffect(() => {
    fetchUserKudosData(0, currentTab, false);

    return () => {
      if (window.location.href !== '/user-data') {
        dispatch(resetTabsData());
      }
    };
  }, []);

  const onUserLink = (): void => {
    navigate('/');
  };

  const returnDashboardBreadCrumb = (): JSX.Element | null => {
    if (searchParams?.get('fromDashboard')) {
      return (
        <div className={styles['button-wrapper']}>
          <Button
            variant="text"
            className={styles['button-container']}
            onClick={onUserLink}
          >
            <ArrowLeftIcon width={24} height={24} />
            Dashboard
          </Button>
        </div>
      );
    }
    return null;
  };

  return (
    <>
      {returnDashboardBreadCrumb()}
      <div className={styles['give-kleos-board-wrapper']}>
        <div className={styles['green-div']}></div>
        <div className={styles['blue-div']}></div>
        <div className={styles['div-wrapper']}>
          <div className={styles['upper-div']}>
            <GiveKleosTypoButton />
          </div>

          <Card className={styles['card-wrapper']}>
            {isLoading ? <FallbackComponent /> : <UserKudosTabs />}
          </Card>
        </div>
      </div>
    </>
  );
};

export default UserKleos;
