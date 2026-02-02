import React from 'react';

const Input = ({ label, error, ...props }) => {
    return (
        <div style={{ marginBottom: '1rem' }}>
            {label && (
                <label style={{
                    display: 'block',
                    fontSize: '0.875rem',
                    fontWeight: 500,
                    marginBottom: '0.375rem',
                    color: 'var(--color-text-main)'
                }}>
                    {label}
                </label>
            )}
            <input
                style={{
                    width: '100%',
                    padding: '0.625rem 0.875rem',
                    borderRadius: 'var(--radius-md)',
                    border: `1px solid ${error ? 'var(--color-error)' : 'var(--color-border)'}`,
                    backgroundColor: 'var(--color-surface)',
                    color: 'var(--color-text-main)',
                    fontSize: '0.875rem',
                    outline: 'none',
                    transition: 'border-color 0.2s',
                }}
                onFocus={(e) => e.currentTarget.style.borderColor = 'var(--color-primary)'}
                onBlur={(e) => e.currentTarget.style.borderColor = error ? 'var(--color-error)' : 'var(--color-border)'}
                {...props}
            />
            {error && (
                <span style={{ display: 'block', marginTop: '0.25rem', fontSize: '0.75rem', color: 'var(--color-error)' }}>
                    {error}
                </span>
            )}
        </div>
    );
};

export default Input;
