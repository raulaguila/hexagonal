import { QueryClient } from '@tanstack/react-query';

/**
 * React Query client configuration
 * Provides caching, background refetching, and error handling
 */
export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            // Data is considered fresh for 30 seconds
            staleTime: 30 * 1000,
            // Cache data for 5 minutes
            gcTime: 5 * 60 * 1000,
            // Retry failed requests up to 2 times
            retry: 2,
            // Retry with exponential backoff
            retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
            // Refetch on window focus
            refetchOnWindowFocus: true,
            // Don't refetch on mount if data exists
            refetchOnMount: false,
        },
        mutations: {
            // Retry mutations once
            retry: 1,
            // Show error notifications handled by components
            onError: (error) => {
                console.error('Mutation error:', error);
            },
        },
    },
});

/**
 * Query keys factory for consistent key management
 */
export const queryKeys = {
    // User queries
    users: {
        all: ['users'],
        list: (filters) => ['users', 'list', filters],
        detail: (id) => ['users', 'detail', id],
    },
    // Role queries
    roles: {
        all: ['roles'],
        list: (filters) => ['roles', 'list', filters],
        detail: (id) => ['roles', 'detail', id],
    },
    // Auth queries
    auth: {
        user: ['auth', 'user'],
    },
};

/**
 * Invalidate all user-related queries
 */
export const invalidateUsers = () => {
    queryClient.invalidateQueries({ queryKey: queryKeys.users.all });
};

/**
 * Invalidate all role-related queries
 */
export const invalidateRoles = () => {
    queryClient.invalidateQueries({ queryKey: queryKeys.roles.all });
};

export default queryClient;
