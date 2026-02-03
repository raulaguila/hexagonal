import { useState, useCallback } from 'react';
import { roleService } from '../services/roleService';

/**
 * Custom hook for role data management
 */
export function useRoles() {
    const [roles, setRoles] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const fetchRoles = useCallback(async ({ page = 1, limit = 100, search = '' } = {}) => {
        setLoading(true);
        setError(null);
        try {
            const response = await roleService.getRoles({ page, limit, search });
            const items = Array.isArray(response) ? response : response.items || [];
            setRoles(items);
        } catch (err) {
            setError(err.response?.data?.message || err.message);
        } finally {
            setLoading(false);
        }
    }, []);

    const createRole = useCallback(async (roleData) => {
        const newRole = await roleService.createRole(roleData);
        setRoles(prev => [...prev, newRole]);
        return newRole;
    }, []);

    const updateRole = useCallback(async (id, roleData) => {
        const updatedRole = await roleService.updateRole(id, roleData);
        setRoles(prev => prev.map(r => r.id === id ? updatedRole : r));
        return updatedRole;
    }, []);

    const deleteRole = useCallback(async (id) => {
        await roleService.deleteRoles([id]);
        setRoles(prev => prev.filter(r => r.id !== id));
    }, []);

    return {
        roles,
        loading,
        error,
        fetchRoles,
        createRole,
        updateRole,
        deleteRole,
        refetch: fetchRoles
    };
}

export default useRoles;
