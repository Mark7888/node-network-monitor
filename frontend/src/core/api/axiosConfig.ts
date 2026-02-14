import axios from 'axios';
import env from '../config/env';

/**
 * Configure axios instance with base URL and common settings
 */
const api = axios.create({
  baseURL: env.apiUrl,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

export default api;
