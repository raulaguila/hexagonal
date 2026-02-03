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

// Request Interceptor: Attach Token & Language
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

        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;

            try {
                const refreshToken = localStorage.getItem('refresh_token');
                if (!refreshToken) throw new Error('No refresh token');

                const response = await axios.put(`${API_URL}/auth`, {}, {
                    headers: {
                        Authorization: `Bearer ${refreshToken}`
                    }
                });

                const { access_token } = response.data;
                localStorage.setItem('access_token', access_token);

                originalRequest.headers.Authorization = `Bearer ${access_token}`;
                return api(originalRequest);
            } catch (err) {
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
export { API_URL, BASE_URL };
