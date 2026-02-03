import React from 'react';

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
            className={className}
            style={{
                backgroundColor: 'var(--color-surface)',
                borderRadius: 'var(--radius-lg)',
                border: '1px solid var(--color-border)',
                boxShadow: 'var(--shadow-sm)',
                padding,
                ...style
            }}
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
        <Card style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
            <div
                style={{
                    padding: '0.75rem',
                    borderRadius: 'var(--radius-md)',
                    backgroundColor: `${color}15`,
                    color: color
                }}
            >
                {Icon && <Icon size={24} />}
            </div>
            <div style={{ flex: 1 }}>
                <h3
                    style={{
                        margin: 0,
                        fontSize: '0.875rem',
                        color: 'var(--color-text-secondary)',
                        fontWeight: 500
                    }}
                >
                    {title}
                </h3>
                <p
                    style={{
                        margin: '0.25rem 0 0 0',
                        fontSize: '1.5rem',
                        fontWeight: 700,
                        color: 'var(--color-text-main)'
                    }}
                >
                    {value}
                </p>
            </div>
            {trend && (
                <div
                    style={{
                        fontSize: '0.75rem',
                        fontWeight: 600,
                        color: trend > 0 ? '#10b981' : '#ef4444'
                    }}
                >
                    {trend > 0 ? '+' : ''}{trend}%
                </div>
            )}
        </Card>
    );
}

export { Card, StatCard };
export default Card;
