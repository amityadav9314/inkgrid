import React, { createContext, useState, useEffect, useContext } from 'react';
import { authService } from '../services/authService';

interface User {
  id: number;
  email: string;
  name: string;
}

interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, name: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Check if user is already logged in
    const storedToken = localStorage.getItem('token');
    const storedUser = localStorage.getItem('user');
    
    console.log('AuthContext initialization:', { storedToken, storedUser });
    
    if (storedToken && storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser);
        console.log('Parsed user:', parsedUser);
        
        // Important: Set both token and user in state
        setToken(storedToken);
        setUser(parsedUser);
        
        console.log('Auth state set from localStorage');
      } catch (error) {
        console.error('Failed to parse stored user data:', error);
        // Clear invalid data
        localStorage.removeItem('token');
        localStorage.removeItem('user');
      }
    } else {
      console.log('No stored authentication data found');
    }
    
    setIsLoading(false);
  }, []);
  
  // Debug effect to monitor authentication state
  useEffect(() => {
    const isAuth = !!user && !!token;
    console.log('Auth state updated:', { 
      user, 
      token, 
      isAuthenticated: isAuth 
    });
    
    // If we have a token but no auth state, something is wrong - try to fix it
    const storedToken = localStorage.getItem('token');
    const storedUser = localStorage.getItem('user');
    if (storedToken && storedUser && !isAuth) {
      console.log('Auth state inconsistency detected, attempting to fix...');
      try {
        setToken(storedToken);
        setUser(JSON.parse(storedUser));
      } catch (error) {
        console.error('Failed to restore auth state:', error);
      }
    }
  }, [user, token]);

  const login = async (email: string, password: string) => {
    setIsLoading(true);
    try {
      console.log('Attempting login for:', email);
      const response = await authService.login(email, password);
      console.log('Login response:', response);
      
      // Store in localStorage first
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
      
      // Then update state
      setToken(response.token);
      setUser(response.user);
      
      console.log('Login successful, auth state updated');
    } catch (error) {
      console.error('Login failed:', error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (email: string, password: string, name: string) => {
    setIsLoading(true);
    try {
      await authService.register(email, password, name);
      // After registration, log the user in
      await login(email, password);
    } catch (error) {
      console.error('Registration failed:', error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    setUser(null);
    setToken(null);
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  };

  // Calculate isAuthenticated based on both user and token
  const isAuthenticated = !!user && !!token;
  
  const value = {
    user,
    token,
    isAuthenticated,
    isLoading,
    login,
    register,
    logout
  };
  
  // Debug the final auth state before rendering
  console.log('AuthContext providing state:', { user, token, isAuthenticated });

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};
