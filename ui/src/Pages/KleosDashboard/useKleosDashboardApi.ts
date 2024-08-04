import { useDispatch } from 'react-redux';
import { AxiosResponse } from 'axios';
import { toast } from '@navi/web-ui/lib/primitives/Toast/index';

import { ApiService } from '@src/service/api';
import { getFromLocalStorage } from '@src/service/storage';
import {
  setIsLoading,
  setDashboardData,
  setLeaderBoardData,
} from '@src/store/Dashboard/DashboardSlice';
import {
  FETCH_USER_KLEOS_DATA,
  FETCH_DASHBOARD_LEADERBOARD_DATA,
} from './constants';
import { UserKleosData, DashboardLeaderBoardData } from './types';

interface IUseKleosDashboardApis {
  fetchAndSetDashboardData: () => void;
  fetchAndSetLeaderBoardData: () => void;
  startFetchingDashboardData: () => void;
}

const useKleosDashboardApis = (): IUseKleosDashboardApis => {
  const dispatch = useDispatch();

  const fetchAndSetDashboardData = (): void => {
    const emailId: string = getFromLocalStorage('email-id') || '';

    ApiService.get(FETCH_USER_KLEOS_DATA(emailId))
      .then((res: AxiosResponse<UserKleosData>) => {
        dispatch(setDashboardData(res?.data));
      })
      .catch(err => {
        toast.error(err.message);
      });
  };

  const fetchAndSetLeaderBoardData = (): void => {
    const emailId: string = getFromLocalStorage('email-id') || '';
    const isReceived: boolean = window?.config?.LEADERBOARD_TYPE == 'RECEIVED';

    ApiService.get(FETCH_DASHBOARD_LEADERBOARD_DATA(isReceived, emailId))
      .then((res: AxiosResponse<DashboardLeaderBoardData>) => {
        dispatch(setLeaderBoardData(res?.data));
      })
      .catch(err => {
        toast.error(err.message);
      });
  };

  const startFetchingDashboardData = (): void => {
    fetchAndSetDashboardData();
    fetchAndSetLeaderBoardData();
  };

  return {
    fetchAndSetDashboardData,
    fetchAndSetLeaderBoardData,
    startFetchingDashboardData,
  };
};

export default useKleosDashboardApis;
