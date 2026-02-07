import React from 'react';
import { ArrowUp, ArrowDown, ArrowUpDown } from 'lucide-react';
import styles from './Table.module.css';

export const Table = ({ children }) => (
    <div className={styles.container}>
        <table className={styles.table}>
            {children}
        </table>
    </div>
);

export const Thead = ({ children }) => (
    <thead className={styles.thead}>
        {children}
    </thead>
);

export const Tbody = ({ children }) => (
    <tbody>
        {children}
    </tbody>
);

export const Tr = ({ children, className = '', ...props }) => (
    <tr className={`${styles.tr} ${className}`} {...props}>
        {children}
    </tr>
);

/**
 * Sortable Table Header
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

    // Determine justification based on text-align
    const justifyContent = style.textAlign === 'center' ? 'center' :
        style.textAlign === 'right' ? 'flex-end' : 'flex-start';

    return (
        <th
            onClick={handleClick}
            className={`${styles.th} ${canSort ? styles.sortable : ''}`}
            style={style}
            {...props}
        >
            <div className={styles.headerContent} style={{ justifyContent }}>
                <span>{children}</span>
                {canSort && (
                    <span className={`${styles.sortIcon} ${isSorted ? styles.active : ''}`}>
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
        className={styles.td}
        style={style}
        {...props}
    >
        {children}
    </td>
);
