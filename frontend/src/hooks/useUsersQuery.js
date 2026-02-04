import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { queryKeys, invalidateUsers } from '../services/queryClient';
import { userService } from '../services/userService';

/**
 * TanStack Query hooks for user operations
 * These provide automatic caching, background refetching, and cache invalidation
 */

/**
 * Hook to fetch paginated users with caching
 */
export function useUsersQuery(filters = {}) {
    return useQuery({
        queryKey: queryKeys.users.list(filters),
        queryFn: () => userService.getUsers(filters),
        placeholderData: (previousData) => previousData,
        select: (data) => ({
            users: data.items || [],
            pagination: data.pagination || { page: 1, limit: 10, total_items: 0, total_pages: 0 }
        }),
    });
}

/**
 * Hook to fetch a single user by ID
 */
export function useUserQuery(id) {
    return useQuery({
        queryKey: queryKeys.users.detail(id),
        queryFn: () => userService.getUserById(id),
        enabled: !!id,
    });
}

/**
 * Hook to create a new user with automatic cache invalidation
 */
export function useCreateUserMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (userData) => userService.createUser(userData),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: queryKeys.users.all });
        },
    });
}

/**
 * Hook to update an existing user
 */
export function useUpdateUserMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ id, data }) => userService.updateUser(id, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: queryKeys.users.all });
            queryClient.invalidateQueries({ queryKey: queryKeys.users.detail(variables.id) });
        },
    });
}

/**
 * Hook to delete users
 */
export function useDeleteUsersMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (ids) => userService.deleteUsers(ids),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: queryKeys.users.all });
        },
    });
}

/**
 * Hook to reset user password
 */
export function useResetPasswordMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (email) => userService.resetPassword(email),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: queryKeys.users.all });
        },
    });
}

/**
 * Hook to set user password
 */
export function useSetPasswordMutation() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ email, data }) => userService.setPassword(email, data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: queryKeys.users.all });
        },
    });
}
