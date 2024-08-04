import React from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import { DarkKnightInstance } from '@navi/dark-knight';
import '@navi/web-ui/lib/styles/index.css';

import ToastWrapper from '@components/ToastWrapper';
import { initAuth } from '@src/utils/auth';

import App from './App';
import './index.css';

const renderAppInsideRoot = (auth: DarkKnightInstance): void => {
  const container = document.getElementById('root');

  if (container) {
    const root = createRoot(container);
    const isAuthenticated =
      (auth.sessionToken && auth.sessionToken?.length > 0) || false;
    root.render(
      <React.StrictMode>
        <ToastWrapper />
        <BrowserRouter>
          <App auth={auth} isAuthenticated={isAuthenticated} />
        </BrowserRouter>
      </React.StrictMode>,
    );
  }
};

initAuth({
  renderApp: renderAppInsideRoot,
});
