import { useDispatch } from 'react-redux';
import { AxiosResponse } from 'axios';
import { toast } from '@navi/web-ui/lib/primitives/Toast/index';
import { ApiService } from '@src/service/api';
import { getFromLocalStorage } from '@src/service/storage';

import { setIsLoading, setAdminBoardData } from '@src/store/Admin/AdminSlice';

import {
  ADMIN_DATA,
  FETCH_DASHBOARD_LEADERBOARD_DATA,
} from '@src/Pages/KleosDashboard/constants';

import {
  AdminPortalBoardData,
  DashboardLeaderBoardData,
} from '@src/Pages/KleosDashboard/types';
import { setLeaderBoardData } from '@store/Dashboard/DashboardSlice';

interface IUseKleosAdminApis {
  fetchAndSetAdminBoardData: () => void;
}

const useKleosAdminApis = (): IUseKleosAdminApis => {
  const dispatch = useDispatch();

  const fetchAndSetAdminBoardData = (): void => {
    dispatch(setIsLoading(true));

    ApiService.get(ADMIN_DATA())
      .then((res: AxiosResponse<AdminPortalBoardData>) => {
        dispatch(setAdminBoardData(res?.data));
        dispatch(setIsLoading(false));
      })
      .catch(err => {
        dispatch(setIsLoading(false));
        toast.error(err.message);
      });
  };

  return {
    fetchAndSetAdminBoardData,
  };
};

export default useKleosAdminApis;
