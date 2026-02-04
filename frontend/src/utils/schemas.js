import { z } from 'zod';

/**
 * Zod validation schemas for forms
 */

// Login schema
export const loginSchema = z.object({
    login: z
        .string()
        .min(1, 'Username or email is required')
        .min(3, 'Must be at least 3 characters'),
    password: z
        .string()
        .min(1, 'Password is required')
        .min(6, 'Password must be at least 6 characters'),
});

// User creation/update schema
export const userSchema = z.object({
    name: z
        .string()
        .min(1, 'Name is required')
        .min(5, 'Name must be at least 5 characters')
        .max(100, 'Name must be less than 100 characters'),
    username: z
        .string()
        .min(1, 'Username is required')
        .min(5, 'Username must be at least 5 characters')
        .max(50, 'Username must be less than 50 characters')
        .regex(/^[a-zA-Z0-9_]+$/, 'Username can only contain letters, numbers, and underscores'),
    email: z
        .string()
        .min(1, 'Email is required')
        .email('Invalid email format'),
    status: z.boolean().optional(),
    role_ids: z.array(z.string().uuid('Invalid role ID')).optional(),
});

// User update schema (all fields optional)
export const userUpdateSchema = userSchema.partial();

// Password schema
export const passwordSchema = z.object({
    password: z
        .string()
        .min(1, 'Password is required')
        .min(6, 'Password must be at least 6 characters')
        .max(128, 'Password must be less than 128 characters'),
    password_confirm: z
        .string()
        .min(1, 'Password confirmation is required'),
}).refine((data) => data.password === data.password_confirm, {
    message: "Passwords don't match",
    path: ['password_confirm'],
});

// Role schema
export const roleSchema = z.object({
    name: z
        .string()
        .min(1, 'Name is required')
        .min(4, 'Name must be at least 4 characters')
        .max(100, 'Name must be less than 100 characters'),
    permissions: z.array(z.string()).optional(),
    enabled: z.boolean().optional(),
});

// Role update schema (all fields optional)
export const roleUpdateSchema = roleSchema.partial();

// Type exports for TypeScript-like usage
export const schemaTypes = {
    login: loginSchema,
    user: userSchema,
    userUpdate: userUpdateSchema,
    password: passwordSchema,
    role: roleSchema,
    roleUpdate: roleUpdateSchema,
};

export default schemaTypes;
