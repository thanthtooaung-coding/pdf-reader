import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react';
import { useNavigate } from 'react-router-dom';
import { setUnauthorizedHandler } from '@/lib/apiError';
import {
  clearAuthSession,
  getAccessToken,
  getStoredUser,
  setAuthSession,
} from '@/lib/authStorage';
import type { UserResponse } from '@/schemas/apiSchema';

interface AuthContextValue {
  user: UserResponse | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  setSession: (token: string, user: UserResponse) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const navigate = useNavigate();
  const [user, setUser] = useState<UserResponse | null>(() => getStoredUser());
  const [isLoading, setIsLoading] = useState(true);

  const logout = useCallback(() => {
    clearAuthSession();
    setUser(null);
    navigate('/login', { replace: true });
  }, [navigate]);

  useEffect(() => {
    setUnauthorizedHandler(logout);
  }, [logout]);

  useEffect(() => {
    const token = getAccessToken();
    if (!token) {
      setUser(null);
      setIsLoading(false);
      return;
    }
    setUser(getStoredUser());
    setIsLoading(false);
  }, []);

  const setSession = useCallback((token: string, nextUser: UserResponse) => {
    setAuthSession(token, nextUser);
    setUser(nextUser);
  }, []);

  const value = useMemo(
    () => ({
      user,
      isAuthenticated: Boolean(getAccessToken() && user),
      isLoading,
      setSession,
      logout,
    }),
    [user, isLoading, setSession, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}
