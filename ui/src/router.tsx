import React, { lazy } from 'react';
import { useSelector } from 'react-redux';

import {
  selectCurrentUser,
  selectIsGiven,
  selectIsLoading,
} from '@store/Dashboard/DashboardSlice';
import FinalDashboard from '@src/Pages/KleosDashboard/index';
import { USER_PERMISSION } from '@src/Pages/KleosDashboard/constants';

import Unauthorized from '@components/Unauthorized';
import styles from '@src/Pages/FirstPage/WelcomeFirstPage.module.scss';
import cardLoader from '@src/assets/lotties/cardLoader.json';
import Lottie from 'lottie-react';

// const FinalDashboard = lazy(() => import('./Pages/KleosDashboard/index'));
const FinalUserKleos = lazy(() => import('./Pages/UserKleos/index'));
const FinalKleosAdmin = lazy(() => import('./Pages/KleosAdmin/index'));
const FinalFirstPage = lazy(() => import('./Pages/FirstPage/index'));
const WelcomeCarousel = lazy(
  () => import('./Pages/FirstPage/WelcomeFirstPage'),
);

export interface CustomRouteObject {
  id: string;
  path: string;
  element: JSX.Element;
  permissions?: string[]; // Add in-case of protected routes
}

const CheckIsGiven = (): JSX.Element => {
  const isGiven = useSelector(selectIsGiven);
  const currentUser = useSelector(selectCurrentUser);
  const isLoading = useSelector(selectIsLoading);

  return (
    <>
      {!isLoading ? (
        isGiven ? (
          <FinalDashboard />
        ) : (
          <FinalFirstPage userId={currentUser} />
        )
      ) : (
        <Lottie animationData={cardLoader} loop />
      )}
    </>
  );
};

const routes: CustomRouteObject[] = [
  {
    id: 'DASHBOARD',
    path: '/',
    element: <CheckIsGiven />,
  },
  {
    id: 'USER-KUDOS',
    path: '/user-kudos',
    element: <FinalUserKleos />,
  },
  {
    id: 'ADMIN',
    path: '/admin',
    element: <FinalKleosAdmin />,
    permissions: [USER_PERMISSION.REQUEST_ADMIN_BOARD_READ_PERMISSION],
  },
  {
    id: 'UNAUTHORIZED',
    path: '/unauthorized',
    element: <Unauthorized variant="p3" />,
  },
  {
    id: 'ABOUT_KUDOS',
    path: '/about-kudos',
    element: (
      <div className={styles.wrapper}>
        <div className={styles.container}>
          <WelcomeCarousel />
        </div>
      </div>
    ),
  },
];

export default routes;
