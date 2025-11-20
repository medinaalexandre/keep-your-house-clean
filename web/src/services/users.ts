export interface User {
  id: number;
  name: string;
  email: string;
  tenant_id: number;
  points: number;
  role: string;
  status: string;
  last_login_at: string | null;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}

const API_BASE_URL = 'http://localhost:8080';

const getAuthToken = (): string | null => {
  return localStorage.getItem('token');
};

const getHeaders = (): HeadersInit => {
  const token = getAuthToken();
  return {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
  };
};

export const getUsers = async (): Promise<User[]> => {
  const response = await fetch(`${API_BASE_URL}/users`, {
    method: 'GET',
    headers: getHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to fetch users');
  }

  return response.json();
};

export const getTopUsers = async (): Promise<User[]> => {
  const response = await fetch(`${API_BASE_URL}/users/ranking`, {
    method: 'GET',
    headers: getHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to fetch top users');
  }

  return response.json();
};

export const getCurrentUser = (): { id: number; name: string; tenant_id: number } | null => {
  const userId = localStorage.getItem('user_id');
  const userName = localStorage.getItem('user_name');
  const tenantId = localStorage.getItem('tenant_id');
  
  if (!userId || !tenantId) return null;

  return {
    id: parseInt(userId, 10),
    name: userName || 'Usuário',
    tenant_id: parseInt(tenantId, 10),
  };
};

export interface CreateUserRequest {
  name: string;
  email: string;
  password: string;
  tenant_id: number;
  points?: number;
  role?: string;
  status?: string;
}

export const createUser = async (req: CreateUserRequest): Promise<User> => {
  const response = await fetch(`${API_BASE_URL}/users`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({
      name: req.name,
      email: req.email,
      password: req.password,
      tenant_id: req.tenant_id,
      points: req.points || 0,
      role: req.role || 'user',
      status: req.status || 'active',
    }),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to create user');
  }

  return response.json();
};

