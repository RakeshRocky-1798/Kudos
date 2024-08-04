import React from 'react';
import { Provider } from 'react-redux';

import { DarkKnightInstance } from '@navi/dark-knight';
import { AuthProvider } from '@src/context/AuthContext';

import FallbackComponent from '@components/Fallback';
import { store } from '@src/store';
import AppRouter from './AppRouter';

interface AppProps {
  auth: DarkKnightInstance;
  isAuthenticated: boolean;
}

//home page
const App: React.FC<AppProps> = props => {
  const { isAuthenticated } = props;
  if (!isAuthenticated) {
    return <FallbackComponent />;
  }

  return (
    <Provider store={store}>
      <AuthProvider value={props.auth}>
        <AppRouter />
      </AuthProvider>
    </Provider>
  );
};

export default App;
