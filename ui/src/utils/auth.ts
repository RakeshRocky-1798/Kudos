import { toast } from '@navi/web-ui/lib/primitives/Toast';
import DarkKnight, { DarkKnightInstance } from '@navi/dark-knight';
import { removeFromLocalStorage, setLocalStorage } from '@src/utils/storage';
import { setToken, getToken } from '@src/utils/authUtils';

interface InitAuth {
  renderApp: (auth: DarkKnightInstance) => void;
}

const config = window.config ?? {};

const initOptions = {
  url: config.AUTH_BASE_URL,
  clientId: config.AUTH_CLIENT_ID,
  onLoad: 'login-required',
};

const initAuth = ({ renderApp }: InitAuth): void => {
  const auth = DarkKnight(initOptions);

  if (!getToken()) {
    removeFromLocalStorage('search-with-filter');
  }

  const onAuthSuccess = (auth: DarkKnightInstance): void => {
    setToken(auth.sessionToken || '');
    setLocalStorage('email-id', `${auth?.idTokenParsed?.emailId}`);
    setLocalStorage('user-name', `${auth?.idTokenParsed?.name}`);
    setLocalStorage('user-data', `${JSON.stringify(auth?.idTokenParsed)}`);
    setLocalStorage(
      'user-profile-url',
      `${auth?.idTokenParsed?.profilePictureUrl}` || 'null',
    );
  };

  const onAuthError = (error: {
    response: { status: number; [x: string]: unknown };
  }): void => {
    renderApp(auth);

    if (error?.response?.status === 403)
      toast.error('You are not authorized. Permission denied!');

    if (error?.response?.status === 404)
      toast.error('You are not authorized. Please get access!');

    toast.error('You are not authorized. Permission denied!');

    setTimeout(() => {
      window?.location?.reload();
    }, 2000);

    return;
  };

  auth
    .init({
      onLoad: 'check-sso',
      sessionToken: getToken() || undefined,
      enableLogging: true,
    })
    .then((authenticated: boolean) => {
      if (!authenticated)
        auth.login({
          redirectUri: `${window.location.protocol}//${window.location.host}/agent/login/?redirectUri=${window.location.href}`,
        });

      onAuthSuccess(auth);
      renderApp(auth);
    })
    .catch(error => {
      onAuthError(error);
    });
};

export { initAuth };
