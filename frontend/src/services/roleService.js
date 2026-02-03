import api from './api';

/**
 * Role Service - API layer for role operations
 */
export const roleService = {
    /**
     * Get paginated list of roles
     */
    async getRoles({ page = 1, limit = 10, search = '' } = {}) {
        const { data } = await api.get('/role', {
            params: { page, limit, search }
        });
        return data;
    },

    /**
     * Get role by ID
     */
    async getRoleById(id) {
        const { data } = await api.get(`/role/${id}`);
        return data;
    },

    /**
     * Create new role
     */
    async createRole(roleData) {
        const { data } = await api.post('/role', roleData);
        return data;
    },

    /**
     * Update existing role
     */
    async updateRole(id, roleData) {
        const { data } = await api.put(`/role/${id}`, roleData);
        return data;
    },

    /**
     * Delete roles by IDs
     */
    async deleteRoles(ids) {
        const { data } = await api.delete('/role', { data: { ids } });
        return data;
    }
};

export default roleService;
