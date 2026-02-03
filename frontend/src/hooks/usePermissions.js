import { useCallback } from 'react';
import { useAuth } from '../context/AuthContext';

/**
 * Custom hook for checking user permissions
 * 
 * Usage:
 *   const { hasPermission, hasAnyPermission, hasAllPermissions } = usePermissions();
 *   
 *   if (hasPermission('users:create')) { ... }
 *   if (hasAnyPermission(['users:edit', 'users:delete'])) { ... }
 */
export function usePermissions() {
    const { user } = useAuth();

    /**
     * Get all permissions from user's roles
     */
    const getAllPermissions = useCallback(() => {
        if (!user?.roles) return [];

        const permissions = new Set();
        user.roles.forEach(role => {
            if (role.permissions) {
                role.permissions.forEach(p => permissions.add(p));
            }
        });
        return Array.from(permissions);
    }, [user]);

    /**
     * Check if user has a specific permission
     */
    const hasPermission = useCallback((permission) => {
        if (!user?.roles) return false;

        return user.roles.some(role =>
            role.permissions?.includes(permission)
        );
    }, [user]);

    /**
     * Check if user has any of the specified permissions
     */
    const hasAnyPermission = useCallback((permissions) => {
        return permissions.some(p => hasPermission(p));
    }, [hasPermission]);

    /**
     * Check if user has all of the specified permissions
     */
    const hasAllPermissions = useCallback((permissions) => {
        return permissions.every(p => hasPermission(p));
    }, [hasPermission]);

    /**
     * Check if user has a specific role
     */
    const hasRole = useCallback((roleName) => {
        if (!user?.roles) return false;
        return user.roles.some(role => role.name === roleName);
    }, [user]);

    /**
     * Check if user is root/admin
     */
    const isRoot = useCallback(() => {
        return hasRole('ROOT') || hasRole('ADMIN');
    }, [hasRole]);

    return {
        permissions: getAllPermissions(),
        hasPermission,
        hasAnyPermission,
        hasAllPermissions,
        hasRole,
        isRoot
    };
}

export default usePermissions;
