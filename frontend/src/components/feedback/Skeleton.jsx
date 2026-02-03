import React from 'react';

/**
 * Skeleton loading component for better UX during data fetching
 */
function Skeleton({
    width = '100%',
    height = '1rem',
    borderRadius = 'var(--radius-md)',
    className = '',
    style = {}
}) {
    return (
        <div
            className={className}
            style={{
                width,
                height,
                borderRadius,
                backgroundColor: 'var(--color-border)',
                animation: 'skeleton-pulse 1.5s ease-in-out infinite',
                ...style
            }}
        />
    );
}

/**
 * Skeleton for text lines
 */
function SkeletonText({ lines = 3, gap = '0.5rem' }) {
    return (
        <div style={{ display: 'flex', flexDirection: 'column', gap }}>
            {Array.from({ length: lines }).map((_, i) => (
                <Skeleton
                    key={i}
                    width={i === lines - 1 ? '60%' : '100%'}
                    height="0.875rem"
                />
            ))}
        </div>
    );
}

/**
 * Skeleton for avatar/circle
 */
function SkeletonAvatar({ size = 40 }) {
    return (
        <Skeleton
            width={size}
            height={size}
            borderRadius="50%"
        />
    );
}

/**
 * Skeleton for table row
 */
function SkeletonTableRow({ columns = 4 }) {
    return (
        <tr style={{ borderBottom: '1px solid var(--color-border)' }}>
            {Array.from({ length: columns }).map((_, i) => (
                <td key={i} style={{ padding: '1rem' }}>
                    <Skeleton height="1rem" width={i === 0 ? '60%' : '80%'} />
                </td>
            ))}
        </tr>
    );
}

/**
 * Skeleton for card
 */
function SkeletonCard() {
    return (
        <div style={{
            backgroundColor: 'var(--color-surface)',
            padding: '1.5rem',
            borderRadius: 'var(--radius-lg)',
            border: '1px solid var(--color-border)'
        }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                <Skeleton width={48} height={48} borderRadius="var(--radius-md)" />
                <div style={{ flex: 1 }}>
                    <Skeleton width="40%" height="0.75rem" style={{ marginBottom: '0.5rem' }} />
                    <Skeleton width="60%" height="1.5rem" />
                </div>
            </div>
        </div>
    );
}

// Add keyframes to document if not already present
if (typeof document !== 'undefined') {
    const styleId = 'skeleton-styles';
    if (!document.getElementById(styleId)) {
        const style = document.createElement('style');
        style.id = styleId;
        style.textContent = `
      @keyframes skeleton-pulse {
        0%, 100% { opacity: 1; }
        50% { opacity: 0.4; }
      }
    `;
        document.head.appendChild(style);
    }
}

export { Skeleton, SkeletonText, SkeletonAvatar, SkeletonTableRow, SkeletonCard };
export default Skeleton;
