import React from 'react';
import styles from './Card.module.css';

/**
 * Card component for content sections
 */
function Card({
    children,
    padding = '1.5rem',
    className = '',
    style = {},
    ...props
}) {
    return (
        <div
            className={`${styles.card} ${className}`}
            style={{ padding, ...style }}
            {...props}
        >
            {children}
        </div>
    );
}

/**
 * Stat Card for dashboard statistics
 */
function StatCard({ title, value, icon: Icon, color = 'var(--color-primary)', trend }) {
    return (
        <Card className={styles.statCard}>
            <div
                className={styles.iconWrapper}
                style={{
                    backgroundColor: `${color}15`,
                    color: color
                }}
            >
                {Icon && <Icon size={24} />}
            </div>
            <div className={styles.content}>
                <h3 className={styles.title}>
                    {title}
                </h3>
                <p className={styles.value}>
                    {value}
                </p>
            </div>
            {trend && (
                <div className={`${styles.trend} ${trend > 0 ? styles.trendUp : styles.trendDown}`}>
                    {trend > 0 ? '+' : ''}{trend}%
                </div>
            )}
        </Card>
    );
}

export { Card, StatCard };
export default Card;
