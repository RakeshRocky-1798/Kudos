import { configureStore } from '@reduxjs/toolkit';
import logger from 'redux-logger';

import DashBoardReducer, {
  DashboardState,
} from '@src/store/Dashboard/DashboardSlice';
import UserKleosReducer, {
  UserKleosState,
} from '@src/store/UserKudos/UserSlice';
import AdminReducer, { AdminState } from '@src/store/Admin/AdminSlice';

export interface AppState {
  dashboard: DashboardState;
  userKleos: UserKleosState;
  admin: AdminState;
}

const store = configureStore({
  reducer: {
    dashboard: DashBoardReducer,
    userKleos: UserKleosReducer,
    admin: AdminReducer,
  },
  middleware: getDefaultMiddleware =>
    getDefaultMiddleware({
      serializableCheck: false,
      logger: window.config.ENV !== 'PRODUCTION' ? logger : false,
    }),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export { store };
