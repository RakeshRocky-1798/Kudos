import { createSlice, PayloadAction } from '@reduxjs/toolkit';

import {
  adminBoardData,
  AdminPortalBoardData,
} from '@src/Pages/KleosDashboard/types';
import { AppState } from '../index';

export interface AdminState {
  adminBoardData: AdminPortalBoardData;
  isLoading: boolean;
}

const initialState: AdminState = {
  adminBoardData: {
    data: {
      adminAllUser: [],
    },
    error: {},
    status: '',
  },
  isLoading: false,
};

const AdminSlice = createSlice({
  name: 'Admin',
  initialState,
  reducers: {
    setAdminBoardData: (
      state,
      action: PayloadAction<AdminPortalBoardData>,
    ): void => {
      state.adminBoardData = action.payload;
    },
    setIsLoading: (state, action: PayloadAction<boolean>): void => {
      state.isLoading = action.payload;
    },
  },
});

export const { setAdminBoardData, setIsLoading } = AdminSlice.actions;

export const selectAdminBoardData = (state: AppState): adminBoardData[] =>
  state.admin.adminBoardData.data?.adminAllUser;

export const selectAdminBoardError = (state: AppState): unknown =>
  state.dashboard.leaderBoardData.error;
export const selectIsLoading = (state: AppState): boolean =>
  state.admin.isLoading;

export default AdminSlice.reducer;
