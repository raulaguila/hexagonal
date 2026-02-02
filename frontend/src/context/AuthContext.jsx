import React, { createContext, useContext, useState, useEffect } from 'react';
import api from '../utils/api';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    // Check if user is logged in on mount
    useEffect(() => {
        const checkAuth = async () => {
            const token = localStorage.getItem('access_token');
            if (token) {
                try {
                    // Verify token and get user info
                    const { data } = await api.get('/auth');
                    setUser(data);
                } catch (error) {
                    console.error("Session check failed", error);
                    // If 401, the interceptor might have tried to refresh.
                    // If it failed, tokens are cleared by interceptor.
                    if (!localStorage.getItem('access_token')) {
                        setUser(null);
                    }
                }
            }
            setLoading(false);
        };

        checkAuth();
    }, []);

    const login = async (username, password) => {
        // Backend expects 'login' field not 'email'
        const { data } = await api.post('/auth', { login: username, password });

        // Save tokens
        localStorage.setItem('access_token', data.accesstoken);
        localStorage.setItem('refresh_token', data.refreshtoken);

        // Fetch full user details immediately
        // Use user from response if available, otherwise fetch
        if (data.user) {
            setUser(data.user);
            return data.user;
        }

        const userRes = await api.get('/auth');
        setUser(userRes.data);
        return userRes.data;
        setUser(userRes.data);

        return userRes.data;
    };

    const logout = async () => {
        try {
            await api.post('/auth/logout');
        } catch (err) {
            console.warn("Logout failed on server", err);
        } finally {
            localStorage.removeItem('access_token');
            localStorage.removeItem('refresh_token');
            setUser(null);
        }
    };

    const value = {
        user,
        loading,
        isAuthenticated: !!user,
        login,
        logout,
    };

    return (
        <AuthContext.Provider value={value}>
            {!loading && children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};
