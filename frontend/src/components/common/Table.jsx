import React from 'react';
import { ArrowUp, ArrowDown, ArrowUpDown } from 'lucide-react';

export const Table = ({ children }) => (
    <div style={{
        width: '100%',
        overflowX: 'auto',
    }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left', minWidth: '600px' }}>
            {children}
        </table>
    </div>
);

export const Thead = ({ children }) => (
    <thead style={{
        backgroundColor: 'var(--color-background)',
        borderBottom: '1px solid var(--color-border)'
    }}>
        {children}
    </thead>
);

export const Tbody = ({ children }) => (
    <tbody>
        {children}
    </tbody>
);

export const Tr = ({ children, className = '', ...props }) => (
    <tr className={`table-row ${className}`} {...props}>
        {children}
    </tr>
);

/**
 * Sortable Table Header
 * @param {string} sortKey - Key to sort by when clicked
 * @param {string} currentSort - Current sort key
 * @param {'asc'|'desc'} sortOrder - Current sort order
 * @param {function} onSort - Callback with (sortKey, order)
 */
export const Th = ({
    children,
    sortKey,
    currentSort,
    sortOrder,
    onSort,
    style = {},
    ...props
}) => {
    const isSorted = sortKey && currentSort === sortKey;
    const canSort = !!sortKey && !!onSort;

    const handleClick = () => {
        if (!canSort) return;
        const newOrder = isSorted && sortOrder === 'asc' ? 'desc' : 'asc';
        onSort(sortKey, newOrder);
    };

    return (
        <th
            onClick={handleClick}
            style={{
                padding: '0.75rem 1rem',
                fontSize: '0.75rem',
                textTransform: 'uppercase',
                letterSpacing: '0.05em',
                color: 'var(--color-text-secondary)',
                fontWeight: 600,
                cursor: canSort ? 'pointer' : 'default',
                userSelect: canSort ? 'none' : 'auto',
                transition: 'background-color 0.15s',
                whiteSpace: 'nowrap',
                ...style
            }}
            onMouseEnter={(e) => canSort && (e.currentTarget.style.backgroundColor = 'rgba(99, 102, 241, 0.05)')}
            onMouseLeave={(e) => canSort && (e.currentTarget.style.backgroundColor = 'transparent')}
            {...props}
        >
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <span>{children}</span>
                {canSort && (
                    <span style={{
                        opacity: isSorted ? 1 : 0.3,
                        display: 'flex',
                        alignItems: 'center'
                    }}>
                        {isSorted && sortOrder === 'asc' && <ArrowUp size={14} />}
                        {isSorted && sortOrder === 'desc' && <ArrowDown size={14} />}
                        {!isSorted && <ArrowUpDown size={14} />}
                    </span>
                )}
            </div>
        </th>
    );
};

export const Td = ({ children, colSpan, style = {}, ...props }) => (
    <td
        colSpan={colSpan}
        style={{
            padding: '1rem',
            fontSize: '0.875rem',
            color: 'var(--color-text-main)',
            ...style
        }}
        {...props}
    >
        {children}
    </td>
);

// Responsive styles - inject once
if (typeof document !== 'undefined') {
    const styleId = 'table-responsive-styles';
    if (!document.getElementById(styleId)) {
        const style = document.createElement('style');
        style.id = styleId;
        style.textContent = `
            @media (max-width: 768px) {
                .table-row td:first-child {
                    position: sticky;
                    left: 0;
                    background: var(--color-surface);
                    z-index: 1;
                }
                .table-row:hover td:first-child {
                    background: var(--color-background);
                }
            }
        `;
        document.head.appendChild(style);
    }
}
