import { FC, Suspense, useEffect } from 'react';
import { Routes, Route, Outlet } from 'react-router-dom';
import LeftNav from './components/LeftNavbar';

import routes from './router';
import useFirstPageApis from '@src/Pages/FirstPage/FirstPageApi';
import Secure from '@components/Secure';
import { useDispatch } from 'react-redux';
import { setIsLoading } from '@store/Dashboard/DashboardSlice';

const AppRouter: FC = () => {
  const { callAllApis } = useFirstPageApis();
  const dispatch = useDispatch();

  useEffect(() => {
    callAllApis().then(r => {
      dispatch(setIsLoading(false));
    });
  }, []);

  return (
    <Suspense>
      <Routes>
        <Route
          element={
            <LeftNav>
              <Outlet />
            </LeftNav>
          }
        >
          {routes.map(component => (
            <Route
              key={component?.id}
              path={component?.path}
              element={
                component?.permissions?.length ? (
                  <Secure
                    permissions={component?.permissions || []}
                    element={component?.element}
                  />
                ) : (
                  component?.element
                )
              }
            />
          ))}
        </Route>
      </Routes>
    </Suspense>
  );
};

export default AppRouter;
