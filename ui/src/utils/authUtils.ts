const tokenKey = 'auth-token';
const emailKey = 'email';

export const setToken = (token: string): void => {
  localStorage.setItem(tokenKey, token);
};

export const setEmail = (email: string): void => {
  localStorage.setItem(emailKey, email);
};

export const getToken = (): string | null => localStorage.getItem(tokenKey);

export const getEmail = (): string | null => localStorage.getItem(emailKey);

export const removeToken = (): void => {
  localStorage.removeItem(tokenKey);
};

export const removeEmail = (): void => localStorage.removeItem(emailKey);
