import { createSlice, PayloadAction } from '@reduxjs/toolkit';

import {
  userKudosApiResponse,
  PaginatedDataType,
  receivedTabData,
} from '@src/Pages/UserKleos/types';
import { transformUserKudosData } from './utils';
import { AppState } from '../index';

export interface UserKleosState {
  receivedTabData: userKudosApiResponse;
  givenTabData: userKudosApiResponse;
  isLoading: boolean;
  currentTab: PaginatedDataType;
  isListLoading: boolean;
}

interface SetCurrentTabProps {
  data: userKudosApiResponse;
  tabData: receivedTabData[];
}

export const initialState: UserKleosState = {
  receivedTabData: {
    data: [],
    pages: {
      pageNumber: 0,
      pageSize: 10,
      totalPages: 0,
      totalElements: 0,
      hasData: false,
    },
  },
  givenTabData: {
    data: [],
    pages: {
      pageNumber: 0,
      pageSize: 10,
      totalPages: 0,
      totalElements: 0,
      hasData: false,
    },
  },
  isLoading: false,
  currentTab: 'given',
  isListLoading: false,
};

const UserKleosSlice = createSlice({
  name: 'UserKleos',
  initialState,
  reducers: {
    setCurrentTabData: (
      state,
      action: PayloadAction<SetCurrentTabProps>,
    ): void => {
      if (state.currentTab === 'given') {
        state.givenTabData = {
          // This ensure we are appending the data to the existing data in the state
          data: transformUserKudosData(
            action.payload.tabData,
            action.payload.data,
          )?.data,
          pages: action.payload.data?.pages,
        };
        return;
      }
      state.receivedTabData = {
        // This ensure we are appending the data to the existing data in the state
        data: transformUserKudosData(
          action.payload.tabData,
          action.payload.data,
        )?.data,
        pages: action.payload.data?.pages,
      };
    },
    setIsLoading: (state, action: PayloadAction<boolean>): void => {
      state.isLoading = action.payload;
    },
    setCurrentTab: (state, action: PayloadAction<PaginatedDataType>): void => {
      state.currentTab = action.payload;
    },
    resetTabsData: (state): void => {
      state.givenTabData = initialState.givenTabData;
      state.receivedTabData = initialState.receivedTabData;
    },
    setIsListLoading: (
      state: UserKleosState,
      action: PayloadAction<boolean>,
    ): void => {
      state.isListLoading = action.payload;
    },
  },
});

export const selectReceivedTabData = (state: AppState): userKudosApiResponse =>
  state.userKleos.receivedTabData;
export const selectGivenTabData = (state: AppState): userKudosApiResponse =>
  state.userKleos.givenTabData;
export const selectIsLoading = (state: AppState): boolean =>
  state.userKleos.isLoading;
export const selectCurrentTab = (state: AppState): PaginatedDataType =>
  state.userKleos.currentTab;
export const selectIsListLoading = (state: AppState): boolean =>
  state.userKleos.isListLoading;

export const {
  setCurrentTabData,
  setIsLoading,
  setCurrentTab,
  resetTabsData,
  setIsListLoading,
} = UserKleosSlice.actions;

export default UserKleosSlice.reducer;
