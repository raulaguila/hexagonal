import axios from 'axios';

// Default to local backend if not specified
const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:9999';
const API_URL = `${BASE_URL}/v1`;

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor: Attach Token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    const language = localStorage.getItem('language') || 'en-US';

    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    config.headers['Accept-Language'] = language;
    return config;
  },
  (error) => Promise.reject(error)
);

// Response Interceptor: Handle 401/Refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // If 401 and not already retrying
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        if (!refreshToken) throw new Error('No refresh token');

        // Call refresh endpoint
        // Call refresh endpoint
        // NOTE: Backend expects PUT /v1/auth with Refresh Token in Authorization header
        const response = await axios.put(`${API_URL}/auth`, {}, {
          headers: {
            Authorization: `Bearer ${refreshToken}`
          }
        });

        const { access_token } = response.data;

        // Save new token
        localStorage.setItem('access_token', access_token);

        // Retry original request
        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return api(originalRequest);
      } catch (err) {
        // Refresh failed - Logout
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/login';
        return Promise.reject(err);
      }
    }
    return Promise.reject(error);
  }
);

export default api;
