import React from 'react';
import styles from './Badge.module.css';

/**
 * Badge component for status indicators and labels
 */
function Badge({
    children,
    variant = 'default',
    size = 'md',
    dot = false,
    className = '',
    style = {},
    ...props
}) {
    const badgeClasses = [
        styles.badge,
        styles[variant],
        styles[size],
        className
    ].filter(Boolean).join(' ');

    return (
        <span
            className={badgeClasses}
            style={style}
            {...props}
        >
            {dot && <span className={styles.dot} />}
            {children}
        </span>
    );
}

export default Badge;
