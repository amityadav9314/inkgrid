import { api } from './api';
import axios from 'axios';

// Use the AUTH_URL for authentication endpoints
const AUTH_URL = process.env.REACT_APP_AUTH_URL || 'http://localhost:8034/goinkgrid/auth';

// Create a separate axios instance for auth requests
const authClient = axios.create({
  baseURL: AUTH_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

interface LoginResponse {
  token: string;
  refresh_token: string;
  expires_at: string;
}

interface User {
  id: number;
  email: string;
  name: string;
}

interface AuthResponse {
  token: string;
  user: User;
}

interface RegisterData {
  email: string;
  password: string;
  name: string;
}

interface LoginData {
  email: string;
  password: string;
}

class AuthService {
  async login(email: string, password: string): Promise<AuthResponse> {
    try {
      const response = await authClient.post<LoginResponse>('/login', { email, password });
      console.log('Login response:', response.data);
      
      // Extract user info from JWT token
      const token = response.data.token;
      const tokenParts = token.split('.');
      const payload = JSON.parse(atob(tokenParts[1]));
      
      console.log('Token payload:', payload);
      
      // Create user object from token payload
      const user: User = {
        id: payload.id || 1,
        email: payload.email || email,
        name: email.split('@')[0] // Use part of email as name if not in token
      };
      
      // Return both token and user object
      return {
        token: response.data.token,
        user
      };
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  }

  async register(email: string, password: string, name: string): Promise<void> {
    try {
      console.log('Registering user:', { email, name });
      await authClient.post('/register', { email, password, name });
      // Registration is successful, but we don't get a token yet
      // User will need to login after registration
    } catch (error) {
      console.error('Registration error:', error);
      throw error;
    }
  }

  async refreshToken(refreshToken: string): Promise<LoginResponse> {
    try {
      const response = await authClient.post<LoginResponse>('/refresh', { refresh_token: refreshToken });
      return response.data;
    } catch (error) {
      console.error('Token refresh error:', error);
      throw error;
    }
  }

  isAuthenticated(): boolean {
    return !!localStorage.getItem('token');
  }
}

export const authService = new AuthService();
