import api from './api';

/**
 * User Service - API layer for user operations
 */
export const userService = {
    /**
     * Get paginated list of users
     */
    async getUsers({ page = 1, limit = 10, search = '', order = '', sort = '' } = {}) {
        const { data } = await api.get('/user', {
            params: { page, limit, search, order, sort }
        });
        return data;
    },

    /**
     * Get user by ID
     */
    async getUserById(id) {
        const { data } = await api.get(`/user/${id}`);
        return data;
    },

    /**
     * Create new user
     */
    async createUser(userData) {
        const { data } = await api.post('/user', userData);
        return data;
    },

    /**
     * Update existing user
     */
    async updateUser(id, userData) {
        const { data } = await api.put(`/user/${id}`, userData);
        return data;
    },

    /**
     * Delete users by IDs
     */
    async deleteUsers(ids) {
        const { data } = await api.delete('/user', { data: { ids } });
        return data;
    },

    /**
     * Reset user password
     */
    async resetPassword(email) {
        const { data } = await api.post('/user/reset-password', { email });
        return data;
    }
};

export default userService;
