import React from 'react';

const Button = ({ children, variant = 'primary', className = '', ...props }) => {
    const baseStyles = {
        padding: '0.625rem 1rem',
        borderRadius: 'var(--radius-md)',
        fontWeight: 500,
        width: '100%',
        transition: 'all 0.2s ease',
        fontSize: '0.875rem',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: '0.5rem',
    };

    const variants = {
        primary: {
            backgroundColor: 'var(--color-primary)',
            color: '#fff',
        },
        secondary: {
            backgroundColor: 'var(--color-surface)',
            color: 'var(--color-text-main)',
            border: '1px solid var(--color-border)',
        },
        ghost: {
            backgroundColor: 'transparent',
            color: 'var(--color-text-secondary)',
        }
    };

    const combinedStyles = {
        ...baseStyles,
        ...variants[variant],
        ...(props.style || {})
    };

    // Remove style from props to avoid overwriting
    const { style, ...restProps } = props;

    return (
        <button
            style={combinedStyles}
            {...restProps}
            onMouseOver={(e) => {
                if (variant === 'primary') e.currentTarget.style.backgroundColor = 'var(--color-primary-hover)';
            }}
            onMouseOut={(e) => {
                if (variant === 'primary') e.currentTarget.style.backgroundColor = 'var(--color-primary)';
            }}
        >
            {children}
        </button>
    );
};

export default Button;
