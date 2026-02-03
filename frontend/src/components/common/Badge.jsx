import React from 'react';

/**
 * Badge component for status indicators and labels
 */
function Badge({
    children,
    variant = 'default',
    size = 'md',
    dot = false,
    style = {},
    ...props
}) {
    const variants = {
        default: {
            backgroundColor: 'var(--color-background)',
            color: 'var(--color-text-secondary)',
            border: '1px solid var(--color-border)'
        },
        primary: {
            backgroundColor: 'var(--color-primary-light)',
            color: 'var(--color-primary)',
            border: '1px solid transparent'
        },
        success: {
            backgroundColor: 'rgba(16, 185, 129, 0.1)',
            color: '#10b981',
            border: '1px solid transparent'
        },
        error: {
            backgroundColor: 'rgba(239, 68, 68, 0.1)',
            color: '#ef4444',
            border: '1px solid transparent'
        },
        warning: {
            backgroundColor: 'rgba(245, 158, 11, 0.1)',
            color: '#f59e0b',
            border: '1px solid transparent'
        }
    };

    const sizes = {
        sm: { padding: '0.125rem 0.375rem', fontSize: '0.625rem' },
        md: { padding: '0.25rem 0.5rem', fontSize: '0.75rem' },
        lg: { padding: '0.375rem 0.75rem', fontSize: '0.875rem' }
    };

    return (
        <span
            style={{
                display: 'inline-flex',
                alignItems: 'center',
                gap: '0.25rem',
                borderRadius: 'var(--radius-full)',
                fontWeight: 600,
                whiteSpace: 'nowrap',
                ...variants[variant],
                ...sizes[size],
                ...style
            }}
            {...props}
        >
            {dot && (
                <span
                    style={{
                        width: '6px',
                        height: '6px',
                        borderRadius: '50%',
                        backgroundColor: 'currentColor'
                    }}
                />
            )}
            {children}
        </span>
    );
}

export default Badge;
