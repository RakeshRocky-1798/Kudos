import { useDispatch } from 'react-redux';
import { AxiosResponse } from 'axios';
import { toast } from '@navi/web-ui/lib/primitives/Toast';
import { getFromLocalStorage } from '@src/service/storage';
import { ApiService } from '@src/service/api';
import {
  FETCH_ALL_ACHIEVEMENTS,
  FETCH_ALL_USERS,
  FETCH_FIRST_PAGE_DATA,
} from '@src/Pages/KleosDashboard/constants';
import { HomeData } from '@src/Pages/FirstPage/types';
import {
  setAllAchievementData,
  setAllUserData,
  setCurrentUser,
  setIsGiven,
  setIsLoading,
} from '@store/Dashboard/DashboardSlice';
import { AchievementData, UserData } from '@src/Pages/Common/types';

interface useFirstPageApis {
  fetchAndSetFirstPageDataAsync: () => Promise<void>;
  fetchAndSetAllUsersAsync: () => Promise<void>;
  fetchAndSetAllAchievementsAsync: () => Promise<void>;
  callAllApis: () => Promise<void>;
}

const useFirstPageApis = (): useFirstPageApis => {
  const dispatch = useDispatch();

  const fetchAndSetFirstPageDataAsync = async (): Promise<void> => {
    try {
      const emailId: string = getFromLocalStorage('email-id') || '';
      const res: AxiosResponse<HomeData> = await ApiService.get(
        FETCH_FIRST_PAGE_DATA(emailId),
      );
      dispatch(setCurrentUser(res.data.data['userId']));
      dispatch(setIsGiven(res.data.data['givenCount'] !== '0'));
    } catch (err) {
      dispatch(setIsGiven(false));
    }
  };

  const fetchAndSetAllUsersAsync = async (): Promise<void> => {
    try {
      const emailId: string = getFromLocalStorage('email-id') || '';
      const res: AxiosResponse<UserData> = await ApiService.get(
        FETCH_ALL_USERS(emailId),
      );
      dispatch(setAllUserData(res.data));
    } catch (err) {
      toast.error('Error While Fetching Data');
    }
  };

  const fetchAndSetAllAchievementsAsync = async (): Promise<void> => {
    try {
      const res: AxiosResponse<AchievementData> = await ApiService.get(
        FETCH_ALL_ACHIEVEMENTS(),
      );
      dispatch(setAllAchievementData(res.data));
    } catch (err) {
      toast.error('Error While Fetching Data');
    }
  };

  const callAllApis = async (): Promise<void> => {
    dispatch(setIsLoading(true));
    try {
      await Promise.all([
        fetchAndSetFirstPageDataAsync(),
        fetchAndSetAllUsersAsync(),
        fetchAndSetAllAchievementsAsync(),
      ]);
    } catch (error) {
      toast.error('Error While Fetching Data');
    }
  };

  return {
    fetchAndSetFirstPageDataAsync,
    fetchAndSetAllUsersAsync,
    fetchAndSetAllAchievementsAsync,
    callAllApis,
  };
};

export default useFirstPageApis;
