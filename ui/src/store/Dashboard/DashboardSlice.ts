import { createSlice, PayloadAction } from '@reduxjs/toolkit';

import {
  UserKleosData,
  DashboardLeaderBoardData,
  DashboardData,
  leaderBoardData,
} from '@src/Pages/KleosDashboard/types';
import { AchievementData, UserData } from '@src/Pages/Common/types';
import { AppState } from '../index';

// NOTE: If the type / interface is increasing move this to separate file and import it
export interface DashboardState {
  userKleosData: UserKleosData;
  leaderBoardData: DashboardLeaderBoardData;
  isLoading: boolean;
  isGiven: boolean;
  currentUser: string;
  allUserData: UserData;
  allAchievementData: AchievementData;
  successModalState: boolean;
  giveKleosModalState: boolean;
}

const initialState: DashboardState = {
  userKleosData: {
    data: {
      myData: {
        email: '',
        userNumber: '',
        profileUrl: '',
      },
      kleosMetrics: {
        givenCount: 0,
        receivedCount: 0,
      },
      achievementDropdown: [],
      totalAchievement: [],
      recentRecognition: [],
    },
    error: {},
    status: '',
  },
  leaderBoardData: {
    data: {
      leaderBoardData: [],
    },
    error: {},
    status: '',
  },
  isLoading: false,
  isGiven: false,
  currentUser: '',
  allUserData: {
    data: [],
    error: {},
    status: '',
  },
  allAchievementData: {
    data: [],
    error: {},
    status: '',
  },
  successModalState: false,
  giveKleosModalState: false,
};

const DashboardSlice = createSlice({
  name: 'Dashboard',
  initialState,
  reducers: {
    setDashboardData: (state, action: PayloadAction<UserKleosData>): void => {
      state.userKleosData = action.payload;
    },
    setLeaderBoardData: (
      state,
      action: PayloadAction<DashboardLeaderBoardData>,
    ): void => {
      state.leaderBoardData = action.payload;
    },
    setIsLoading: (state, action: PayloadAction<boolean>): void => {
      state.isLoading = action.payload;
    },
    setIsGiven: (state, action: PayloadAction<boolean>): void => {
      state.isGiven = action.payload;
    },
    setCurrentUser: (state, action: PayloadAction<string>): void => {
      state.currentUser = action.payload;
    },
    setAllUserData: (state, action: PayloadAction<UserData>): void => {
      state.allUserData = action.payload;
    },
    setAllAchievementData: (
      state,
      action: PayloadAction<AchievementData>,
    ): void => {
      state.allAchievementData = action.payload;
    },
    setSuccessModalState: (state, action: PayloadAction<boolean>): void => {
      state.successModalState = action.payload;
    },
    setGiveKleosModalState: (state, action: PayloadAction<boolean>): void => {
      state.giveKleosModalState = action.payload;
    },
  },
});

export const {
  setDashboardData,
  setLeaderBoardData,
  setIsLoading,
  setIsGiven,
  setCurrentUser,
  setAllUserData,
  setAllAchievementData,
  setSuccessModalState,
  setGiveKleosModalState,
} = DashboardSlice.actions;

export const selectDashboardData = (state: AppState): DashboardData =>
  state.dashboard.userKleosData.data;
export const selectDashboardError = (state: AppState): unknown =>
  state.dashboard.userKleosData.error;
export const selectLeaderBoardData = (state: AppState): leaderBoardData[] =>
  state.dashboard.leaderBoardData.data?.leaderBoardData;

export const selectLeaderBoardError = (state: AppState): unknown =>
  state.dashboard.leaderBoardData.error;
export const selectIsLoading = (state: AppState): boolean =>
  state.dashboard.isLoading;
export const selectIsGiven = (state: AppState): boolean =>
  state.dashboard.isGiven;
export const selectCurrentUser = (state: AppState): string =>
  state.dashboard.currentUser;
export const selectAllUserData = (state: AppState): UserData =>
  state.dashboard.allUserData;
export const selectAllAchievementData = (state: AppState): AchievementData =>
  state.dashboard.allAchievementData;
export const selectSuccessModalState = (state: AppState): boolean =>
  state.dashboard.successModalState;
export const selectGiveKleosModalState = (state: AppState): boolean =>
  state.dashboard.giveKleosModalState;

export default DashboardSlice.reducer;
