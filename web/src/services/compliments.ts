import { API_BASE_URL } from '../config/api';

export interface Compliment {
  id: number;
  title: string;
  description: string;
  points: number;
  from_user_id: number;
  to_user_id: number;
  tenant_id: number;
  created_at: string;
  created_by_id: number;
  updated_at: string;
  updated_by_id: number | null;
  deleted_at: string | null;
  viewed_at: string | null;
}

export interface ComplimentWithUser extends Compliment {
  from_user_name: string | null;
}

export interface MarkAsViewedRequest {
  ids: number[];
}

export interface CreateComplimentRequest {
  title: string;
  description: string;
  points: number;
  to_user_id: number;
}

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

export const createCompliment = async (req: CreateComplimentRequest): Promise<Compliment> => {
  const response = await fetch(`${API_BASE_URL}/api/v1/compliments`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(req),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to create compliment');
  }

  return response.json();
};

export const getCompliments = async (): Promise<Compliment[]> => {
  const response = await fetch(`${API_BASE_URL}/api/v1/compliments`, {
    method: 'GET',
    headers: getHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to fetch compliments');
  }

  return response.json();
};

export const getCompliment = async (id: number): Promise<Compliment> => {
  const response = await fetch(`${API_BASE_URL}/api/v1/compliments/${id}`, {
    method: 'GET',
    headers: getHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to fetch compliment');
  }

  return response.json();
};

export const getLastReceivedCompliment = async (): Promise<ComplimentWithUser | null> => {
  const response = await fetch(`${API_BASE_URL}/api/v1/compliments/last-received`, {
    method: 'GET',
    headers: getHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to fetch last received compliment');
  }

  const data = await response.json();
  return data || null;
};

export const getUserComplimentsHistory = async (): Promise<ComplimentWithUser[]> => {
  const response = await fetch(`${API_BASE_URL}/api/v1/compliments/history`, {
    method: 'GET',
    headers: getHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to fetch compliments history');
  }

  return response.json();
};

export const getUnviewedReceivedCompliments = async (): Promise<ComplimentWithUser[]> => {
  const response = await fetch(`${API_BASE_URL}/api/v1/compliments/unviewed`, {
    method: 'GET',
    headers: getHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to fetch unviewed compliments');
  }

  return response.json();
};

export const markComplimentsAsViewed = async (ids: number[]): Promise<void> => {
  const response = await fetch(`${API_BASE_URL}/api/v1/compliments/mark-viewed`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ ids }),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to mark compliments as viewed');
  }
};

export const deleteCompliment = async (id: number): Promise<void> => {
  const response = await fetch(`${API_BASE_URL}/api/v1/compliments/${id}`, {
    method: 'DELETE',
    headers: getHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to delete compliment');
  }
};

