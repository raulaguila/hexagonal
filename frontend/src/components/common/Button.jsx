import React from 'react';

const Button = ({ children, variant = 'primary', size = 'md', loading = false, disabled, ...props }) => {
    const baseStyles = {
        display: 'inline-flex',
        alignItems: 'center',
        justifyContent: 'center',
        fontWeight: 500,
        borderRadius: 'var(--radius-md)',
        transition: 'all var(--transition-fast)',
        cursor: 'pointer',
        border: '1px solid transparent',
        gap: '0.5rem',
        position: 'relative',
        overflow: 'hidden'
    };

    const sizes = {
        sm: { padding: '0.375rem 0.75rem', fontSize: '0.75rem' },
        md: { padding: '0.625rem 1rem', fontSize: '0.875rem' },
        lg: { padding: '0.75rem 1.5rem', fontSize: '1rem' }
    };

    const variants = {
        primary: {
            backgroundColor: 'var(--color-primary)',
            color: '#fff',
            boxShadow: 'var(--shadow-sm)',
        },
        secondary: {
            backgroundColor: 'var(--color-surface)',
            color: 'var(--color-text-main)',
            borderColor: 'var(--color-border)',
            boxShadow: 'var(--shadow-sm)',
        },
        danger: {
            backgroundColor: 'var(--color-error)',
            color: '#fff',
            boxShadow: 'var(--shadow-sm)',
        },
        ghost: {
            backgroundColor: 'transparent',
            color: 'var(--color-text-secondary)',
        }
    };

    const combinedStyles = {
        ...baseStyles,
        ...sizes[size],
        ...variants[variant],
        opacity: disabled || loading ? 0.7 : 1,
        pointerEvents: disabled || loading ? 'none' : 'auto',
        ...(props.style || {})
    };

    const { style: _, onMouseOver, onMouseOut, ...restProps } = props;

    // Hover effects logic
    const handleMouseOver = (e) => {
        if (disabled || loading) return;
        if (variant === 'primary') e.currentTarget.style.backgroundColor = 'var(--color-primary-hover)';
        if (variant === 'secondary') e.currentTarget.style.backgroundColor = 'var(--color-surface-hover)';
        if (variant === 'ghost') e.currentTarget.style.backgroundColor = 'var(--color-background)';
        if (onMouseOver) onMouseOver(e);
    };

    const handleMouseOut = (e) => {
        if (disabled || loading) return;
        const v = variants[variant];
        e.currentTarget.style.backgroundColor = v.backgroundColor;
        if (onMouseOut) onMouseOut(e);
    };

    return (
        <button
            style={combinedStyles}
            onMouseOver={handleMouseOver}
            onMouseOut={handleMouseOut}
            disabled={disabled || loading}
            {...restProps}
        >
            {loading ? (
                <span className="animate-spin" style={{
                    border: '2px solid currentColor',
                    borderTopColor: 'transparent',
                    borderRadius: '50%',
                    width: '1em',
                    height: '1em',
                    display: 'inline-block',
                    animation: 'spin 1s linear infinite'
                }} />
            ) : null}
            {children}
            <style>{`@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }`}</style>
        </button>
    );
};

export default Button;
