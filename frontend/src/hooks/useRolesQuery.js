import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { queryKeys } from '../services/queryClient';
import { roleService } from '../services/roleService';

/**
 * TanStack Query hooks for role operations
 */

/**
 * Hook to fetch paginated roles with caching
 */
export function useRolesQuery(filters = {}) {
    return useQuery({
        queryKey: queryKeys.roles.list(filters),
        queryFn: () => roleService.getRoles(filters),
        placeholderData: (previousData) => previousData,
        select: (data) => ({
            roles: data.items || [],
            pagination: data.pagination || { page: 1, limit: 10, total_items: 0, total_pages: 0 }
        }),
    });
}

/**
 * Hook to fetch a single role by ID
 */
export function useRoleQuery(id) {
    return useQuery({
        queryKey: queryKeys.roles.detail(id),
        queryFn: () => roleService.getRoleById(id),
        enabled: !!id,
    });
}

/**
 * Hook to create a new role
 */
export function useCreateRoleMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (roleData) => roleService.createRole(roleData),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: queryKeys.roles.all });
        },
    });
}

/**
 * Hook to update an existing role
 */
export function useUpdateRoleMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ id, data }) => roleService.updateRole(id, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: queryKeys.roles.all });
            queryClient.invalidateQueries({ queryKey: queryKeys.roles.detail(variables.id) });
        },
    });
}

/**
 * Hook to delete roles
 */
export function useDeleteRolesMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (ids) => roleService.deleteRoles(ids),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: queryKeys.roles.all });
        },
    });
}
