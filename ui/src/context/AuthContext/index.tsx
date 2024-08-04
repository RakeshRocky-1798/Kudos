import { createContext } from 'react';
import { DarkKnightInstance } from '@navi/dark-knight';

const AuthContext = createContext<DarkKnightInstance | null>(null);
export const AuthProvider = AuthContext.Provider;
export default AuthContext;
