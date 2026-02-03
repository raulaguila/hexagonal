import React from 'react';
import { Inbox } from 'lucide-react';
import Button from './Button';

/**
 * Empty state component for "No results" scenarios
 */
function EmptyState({
    icon: Icon = Inbox,
    title = 'No items found',
    description = 'There are no items to display.',
    action,
    actionLabel = 'Create New',
    onAction,
    style = {}
}) {
    return (
        <div
            style={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                padding: '3rem 2rem',
                textAlign: 'center',
                ...style
            }}
        >
            <div
                style={{
                    width: '64px',
                    height: '64px',
                    borderRadius: '50%',
                    backgroundColor: 'var(--color-background)',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    marginBottom: '1rem'
                }}
            >
                {Icon && <Icon size={28} style={{ color: 'var(--color-text-muted)' }} />}
            </div>

            <h3
                style={{
                    margin: '0 0 0.5rem',
                    fontSize: '1rem',
                    fontWeight: 600,
                    color: 'var(--color-text-main)'
                }}
            >
                {title}
            </h3>

            <p
                style={{
                    margin: 0,
                    fontSize: '0.875rem',
                    color: 'var(--color-text-secondary)',
                    maxWidth: '300px'
                }}
            >
                {description}
            </p>

            {(action || onAction) && (
                <Button
                    variant="primary"
                    onClick={onAction}
                    style={{ marginTop: '1.5rem' }}
                >
                    {actionLabel}
                </Button>
            )}
        </div>
    );
}

export default EmptyState;
