import React from 'react';
import styles from './Input.module.css';

const Input = ({ label, error, icon: Icon, className = '', containerStyle = {}, ...props }) => {
    // Construct input class names
    const inputClasses = [
        styles.input,
        Icon ? styles.hasIcon : '',
        error ? styles.errorBorder : '',
        className
    ].filter(Boolean).join(' ');

    return (
        <div className={styles.container} style={containerStyle}>
            {label && (
                <label className={styles.label}>
                    {label}
                </label>
            )}
            <div className={styles.inputWrapper}>
                {Icon && (
                    <div className={styles.iconWrapper}>
                        <Icon size={18} />
                    </div>
                )}
                <input
                    className={inputClasses}
                    {...props}
                />
            </div>
            {error && (
                <span className={styles.errorMessage}>
                    {error}
                </span>
            )}
        </div>
    );
};

export default Input;
