import { useDispatch, useSelector } from 'react-redux';
import { AxiosResponse } from 'axios';

import { ApiService } from '@src/service/api';
import {
  setCurrentTabData,
  setIsLoading,
  setIsListLoading,
  selectGivenTabData,
  selectReceivedTabData,
} from '@src/store/UserKudos/UserSlice';
import { getFromLocalStorage } from '@src/service/storage';
import { userKudosApiResponse, PaginatedDataType } from './types';
import { GET_USER_KUDOS_ENDPOINT, DEFAULT_PAGE_SIZE } from './constant';

interface IUserKudosApis {
  fetchUserKudosData: (
    page_number: number,
    dataType: PaginatedDataType,
    isListLoading: boolean,
  ) => void;
}

const useUserKudosApi = (): IUserKudosApis => {
  const dispatch = useDispatch();
  const givenTabData = useSelector(selectGivenTabData);
  const receivedTabData = useSelector(selectReceivedTabData);

  const fetchUserKudosData = (
    page_number = 0,
    dataType: PaginatedDataType = 'given',
    isListLoading = false,
  ): void => {
    const emailId = getFromLocalStorage('email-id') || '';
    dispatch(isListLoading ? setIsListLoading(true) : setIsLoading(true));

    ApiService.get(
      GET_USER_KUDOS_ENDPOINT(
        emailId,
        dataType,
        page_number,
        DEFAULT_PAGE_SIZE,
      ),
    )
      .then((res: AxiosResponse<userKudosApiResponse>) => {
        dispatch(isListLoading ? setIsListLoading(false) : setIsLoading(false));
        const payload = {
          data: res?.data,
          tabData:
            dataType === 'given' ? givenTabData?.data : receivedTabData?.data,
        };
        dispatch(setCurrentTabData(payload));
      })
      .catch(err => {
        dispatch(isListLoading ? setIsListLoading(false) : setIsLoading(false));
      });
  };

  return {
    fetchUserKudosData,
  };
};

export default useUserKudosApi;
