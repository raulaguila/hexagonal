import React from 'react';
import { ChevronLeft, ChevronRight, MoreHorizontal } from 'lucide-react';

const Pagination = ({ currentPage, totalPages, onPageChange, totalItems, itemsPerPage, onLimitChange, t }) => {

    // Generate page numbers to display
    const getPageNumbers = () => {
        const delta = 1; // Number of pages to show around current page
        const range = [];
        const rangeWithDots = [];
        let l;

        for (let i = 1; i <= totalPages; i++) {
            if (i === 1 || i === totalPages || (i >= currentPage - delta && i <= currentPage + delta)) {
                range.push(i);
            }
        }

        range.forEach(i => {
            if (l) {
                if (i - l === 2) {
                    rangeWithDots.push(l + 1);
                } else if (i - l !== 1) {
                    rangeWithDots.push('...');
                }
            }
            rangeWithDots.push(i);
            l = i;
        });

        return rangeWithDots;
    };

    const handlePageClick = (page) => {
        if (page !== '...' && page !== currentPage) {
            onPageChange(page);
        }
    };

    if (totalItems === 0) return null;

    return (
        <div style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            padding: '0.75rem 1.5rem',
            borderTop: '1px solid var(--color-border)',
            backgroundColor: 'var(--color-surface)',
            color: 'var(--color-text-secondary)',
            fontSize: '0.875rem'
        }}>
            {/* Left side: Results info */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                <span>{t ? t('pagination.showing') : 'Showing'}</span>
                <span style={{ fontWeight: 600, color: 'var(--color-text-main)' }}>
                    {((currentPage - 1) * itemsPerPage) + 1}
                </span>
                <span>{t ? t('pagination.to') : 'to'}</span>
                <span style={{ fontWeight: 600, color: 'var(--color-text-main)' }}>
                    {Math.min(currentPage * itemsPerPage, totalItems)}
                </span>
                <span>{t ? t('pagination.of') : 'of'}</span>
                <span style={{ fontWeight: 600, color: 'var(--color-text-main)' }}>
                    {totalItems}
                </span>
                <span>{t ? t('pagination.results') : 'results'}</span>
            </div>

            {/* Right side: Page controls */}
            <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
                {/* Limit Selector */}
                {onLimitChange && (
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <span style={{ fontSize: '0.875rem' }}>{t ? t('pagination.rows_per_page') : 'Rows per page'}:</span>
                        <select
                            value={itemsPerPage}
                            onChange={(e) => onLimitChange(Number(e.target.value))}
                            style={{
                                padding: '0.25rem 0.5rem',
                                borderRadius: 'var(--radius-md)',
                                border: '1px solid var(--color-border)',
                                backgroundColor: 'var(--color-surface)',
                                color: 'var(--color-text-main)',
                                fontSize: '0.875rem',
                                outline: 'none',
                                cursor: 'pointer'
                            }}
                        >
                            <option value={5}>5</option>
                            <option value={10}>10</option>
                            <option value={20}>20</option>
                            <option value={50}>50</option>
                            <option value={100}>100</option>
                        </select>
                    </div>
                )}

                <div style={{ display: 'flex', gap: '0.25rem', alignItems: 'center' }}>
                    <button
                        className="pagination-button"
                        onClick={() => onPageChange(Math.max(1, currentPage - 1))}
                        disabled={currentPage === 1}
                        aria-label="Previous Page"
                    >
                        <ChevronLeft size={16} />
                    </button>

                    {getPageNumbers().map((page, index) => (
                        <button
                            key={index}
                            className={`pagination-button ${page === currentPage ? 'active' : ''}`}
                            onClick={() => handlePageClick(page)}
                            disabled={page === '...'}
                            style={page === '...' ? { cursor: 'default', border: 'none' } : {}}
                        >
                            {page === '...' ? <MoreHorizontal size={14} /> : page}
                        </button>
                    ))}

                    <button
                        className="pagination-button"
                        onClick={() => onPageChange(Math.min(totalPages, currentPage + 1))}
                        disabled={currentPage === totalPages}
                        aria-label="Next Page"
                    >
                        <ChevronRight size={16} />
                    </button>
                </div>
            </div>
        </div>
    );
};

export default Pagination;
