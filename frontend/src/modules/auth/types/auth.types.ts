/**
 * Authentication related types
 */

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  username: string;
}

export interface AuthState {
  token: string | null;
  username: string | null;
  isAuthenticated: boolean;
}
