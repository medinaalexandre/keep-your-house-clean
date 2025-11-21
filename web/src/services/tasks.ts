import { API_BASE_URL } from '../config/api';

export interface Task {
  id: number;
  title: string;
  description: string;
  points: number;
  status: string;
  scheduled_to: string | null;
  scheduled_by_id: number | null;
  frequency_value: number;
  frequency_unit: string;
  completed: boolean;
  completed_by_id: number | null;
  tenant_id: number;
  created_at: string;
  created_by_id: number;
  updated_at: string;
  updated_by_id: number | null;
  deleted_at: string | null;
}

export interface TaskWithUser extends Task {
  completed_by_name: string | null;
}

export interface CreateTaskRequest {
  title: string;
  description: string;
  points: number;
  status: string;
  scheduled_to: string | null;
  scheduled_by_id: number | null;
  frequency_value: number;
  frequency_unit: 'days' | 'weeks' | 'months';
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

export const getUpcomingTasks = async (limit: number = 5, offset: number = 0): Promise<Task[]> => {
  const params = new URLSearchParams({
    limit: limit.toString(),
    offset: offset.toString(),
  });

  const response = await fetch(`${API_BASE_URL}/tasks/upcoming?${params}`, {
    method: 'GET',
    headers: getHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to fetch upcoming tasks');
  }

  return response.json();
};

export const getCompletedTasksHistory = async (limit: number = 5): Promise<TaskWithUser[]> => {
  const params = new URLSearchParams({
    limit: limit.toString(),
  });

  const response = await fetch(`${API_BASE_URL}/tasks/history?${params}`, {
    method: 'GET',
    headers: getHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to fetch completed tasks history');
  }

  return response.json();
};

export const createTask = async (req: CreateTaskRequest): Promise<Task> => {
  const response = await fetch(`${API_BASE_URL}/tasks`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(req),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to create task');
  }

  return response.json();
};

export const completeTask = async (taskId: number, completedById?: number): Promise<Task> => {
  const body: { completed_by_id?: number } = {};
  if (completedById) {
    body.completed_by_id = completedById;
  }
  
  const response = await fetch(`${API_BASE_URL}/tasks/${taskId}/complete`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(body),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to complete task');
  }

  return response.json();
};

export const undoCompleteTask = async (taskId: number): Promise<Task> => {
  const response = await fetch(`${API_BASE_URL}/tasks/${taskId}/undo`, {
    method: 'POST',
    headers: getHeaders(),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
    throw new Error(errorData.error || 'Failed to undo task completion');
  }

  return response.json();
};
