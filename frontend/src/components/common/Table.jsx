import React from 'react';

export const Table = ({ children }) => (
    <div style={{
        width: '100%',
        overflowX: 'auto',
        backgroundColor: 'var(--color-surface)',
        borderRadius: 'var(--radius-lg)',
        boxShadow: 'var(--shadow-sm)',
        border: '1px solid var(--color-border)'
    }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left' }}>
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

export const Tr = ({ children }) => (
    <tr style={{ borderBottom: '1px solid var(--color-border)' }}>
        {children}
    </tr>
);

export const Th = ({ children }) => (
    <th style={{
        padding: '0.75rem 1rem',
        fontSize: '0.75rem',
        textTransform: 'uppercase',
        letterSpacing: '0.05em',
        color: 'var(--color-text-secondary)',
        fontWeight: 600
    }}>
        {children}
    </th>
);

export const Td = ({ children }) => (
    <td style={{
        padding: '1rem',
        fontSize: '0.875rem',
        color: 'var(--color-text-main)'
    }}>
        {children}
    </td>
);
