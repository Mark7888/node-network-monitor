import api from './axiosConfig';
import toast from 'react-hot-toast';

/**
 * Setup axios interceptors for authentication and error handling
 */
export const setupInterceptors = (
  onUnauthorized: () => void
) => {
  // Request interceptor: Add auth token to requests
  api.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // Response interceptor: Handle auth errors
  api.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response?.status === 401) {
        // Clear token and redirect to login
        localStorage.removeItem('token');
        localStorage.removeItem('username');
        onUnauthorized();
        toast.error('Session expired. Please login again.');
      } else if (error.response?.status >= 500) {
        toast.error('Server error. Please try again later.');
      } else if (!error.response) {
        toast.error('Network error. Please check your connection.');
      }
      
      return Promise.reject(error);
    }
  );
};

export default api;
