import { useState, useCallback } from 'react';
import { userService } from '../services/userService';

/**
 * Custom hook for user data management
 */
export function useUsers() {
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 10,
        total: 0,
        totalPages: 0
    });

    const fetchUsers = useCallback(async ({ page = 1, limit = 10, search = '', order = '', sort = '' } = {}) => {
        setLoading(true);
        setError(null);
        try {
            const response = await userService.getUsers({ page, limit, search, order, sort });
            const items = Array.isArray(response) ? response : response.items || [];
            const paginationData = response.pagination || {};

            setUsers(items);
            setPagination({
                page: paginationData.page || page,
                limit: paginationData.limit || limit,
                total: paginationData.total_items || items.length,
                totalPages: paginationData.total_pages || 1
            });
        } catch (err) {
            setError(err.response?.data?.message || err.message);
        } finally {
            setLoading(false);
        }
    }, []);

    const createUser = useCallback(async (userData) => {
        const newUser = await userService.createUser(userData);
        setUsers(prev => [...prev, newUser]);
        return newUser;
    }, []);

    const updateUser = useCallback(async (id, userData) => {
        const updatedUser = await userService.updateUser(id, userData);
        setUsers(prev => prev.map(u => u.id === id ? updatedUser : u));
        return updatedUser;
    }, []);

    const deleteUser = useCallback(async (id) => {
        await userService.deleteUsers([id]);
        setUsers(prev => prev.filter(u => u.id !== id));
    }, []);

    const deleteUsers = useCallback(async (ids) => {
        await userService.deleteUsers(ids);
        setUsers(prev => prev.filter(u => !ids.includes(u.id)));
    }, []);

    return {
        users,
        loading,
        error,
        pagination,
        fetchUsers,
        createUser,
        updateUser,
        deleteUser,
        deleteUsers,
        refetch: fetchUsers
    };
}

export default useUsers;
