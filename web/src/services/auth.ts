import { translateError } from '../utils/errorTranslator';
import { API_BASE_URL } from '../config/api';

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  tenant_name: string;
  tenant_domain: string;
  user_name: string;
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user_id: number;
  tenant_id: number;
  email: string;
  name: string;
}

export const login = async (req: LoginRequest): Promise<LoginResponse> => {
  const response = await fetch(`${API_BASE_URL}/api/v1/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      email: req.email,
      password: req.password,
    }),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    const errorMessage = errorData.error || 'Login failed';
    throw new Error(translateError(errorMessage));
  }

  const loginResponse = await response.json();
  localStorage.setItem('token', loginResponse.token);
  localStorage.setItem('user_id', loginResponse.user_id.toString());
  localStorage.setItem('user_name', loginResponse.name);
  localStorage.setItem('tenant_id', loginResponse.tenant_id.toString());
  return loginResponse;
};

export const register = async (req: RegisterRequest): Promise<LoginResponse> => {
  const response = await fetch(`${API_BASE_URL}/api/v1/auth/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      tenant_name: req.tenant_name,
      tenant_domain: req.tenant_domain,
      user_name: req.user_name,
      email: req.email,
      password: req.password,
    }),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    const errorMessage = errorData.error || 'Registration failed';
    throw new Error(translateError(errorMessage));
  }

  const registerResponse = await response.json();
  localStorage.setItem('token', registerResponse.token);
  localStorage.setItem('user_id', registerResponse.user_id.toString());
  localStorage.setItem('user_name', registerResponse.name);
  localStorage.setItem('tenant_id', registerResponse.tenant_id.toString());
  return registerResponse;
};

export const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user_id');
  localStorage.removeItem('user_name');
  localStorage.removeItem('tenant_id');
};

