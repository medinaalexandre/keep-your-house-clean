const getApiBaseUrl = (): string => {
  if (import.meta.env.VITE_API_BASE_URL) {
    return import.meta.env.VITE_API_BASE_URL;
  }
  
  if (import.meta.env.DEV) {
    return 'http://localhost:8080';
  }
  
  return '';
};

export const API_BASE_URL = getApiBaseUrl();

