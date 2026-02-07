import React from 'react';
import styles from './Button.module.css';

const Button = ({ children, variant = 'primary', size = 'md', loading = false, disabled, className = '', ...props }) => {
    // Construct class names
    const buttonClasses = [
        styles.button,
        styles[variant],
        styles[size],
        loading ? styles.loading : '',
        className
    ].filter(Boolean).join(' ');

    return (
        <button
            className={buttonClasses}
            disabled={disabled || loading}
            {...props}
        >
            {loading && <span className={styles.spinner} />}
            {children}
        </button>
    );
};

export default Button;
