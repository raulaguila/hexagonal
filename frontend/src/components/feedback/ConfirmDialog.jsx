import React from 'react';
import { createPortal } from 'react-dom';
import { AlertTriangle } from 'lucide-react';
import Button from '../common/Button';

/**
 * Confirmation dialog component - replaces window.confirm()
 */
function ConfirmDialog({
    isOpen,
    onClose,
    onConfirm,
    title = 'Confirm Action',
    message = 'Are you sure you want to proceed?',
    confirmText = 'Confirm',
    cancelText = 'Cancel',
    variant = 'danger',
    loading = false
}) {
    // Simple approach: don't render if not open
    if (!isOpen) return null;

    const variantColors = {
        danger: 'var(--color-error)',
        warning: 'var(--color-warning)',
        primary: 'var(--color-primary)'
    };

    const iconBgColors = {
        danger: 'rgba(239, 68, 68, 0.1)',
        warning: 'rgba(245, 158, 11, 0.1)',
        primary: 'rgba(99, 102, 241, 0.1)'
    };

    return createPortal(
        <div
            style={{
                position: 'fixed',
                inset: 0,
                zIndex: 60,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                backgroundColor: 'rgba(0, 0, 0, 0.5)',
                backdropFilter: 'blur(4px)'
            }}
            onClick={onClose}
        >
            {/* Dialog */}
            <div
                style={{
                    position: 'relative',
                    backgroundColor: 'var(--color-surface)',
                    borderRadius: 'var(--radius-xl)',
                    boxShadow: 'var(--shadow-xl)',
                    width: '100%',
                    maxWidth: '400px',
                    margin: '1rem',
                    border: '1px solid var(--color-border)',
                    overflow: 'hidden',
                    animation: 'dialogFadeIn 0.2s ease-out'
                }}
                onClick={e => e.stopPropagation()}
            >
                {/* Content */}
                <div style={{ padding: '1.5rem', textAlign: 'center' }}>
                    {/* Icon */}
                    <div
                        style={{
                            width: '48px',
                            height: '48px',
                            borderRadius: '50%',
                            backgroundColor: iconBgColors[variant],
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            margin: '0 auto 1rem'
                        }}
                    >
                        <AlertTriangle size={24} style={{ color: variantColors[variant] }} />
                    </div>

                    {/* Title */}
                    <h3 style={{
                        margin: '0 0 0.5rem',
                        fontSize: '1.125rem',
                        fontWeight: 600,
                        color: 'var(--color-text-main)'
                    }}>
                        {title}
                    </h3>

                    {/* Message */}
                    <p style={{
                        margin: 0,
                        fontSize: '0.875rem',
                        color: 'var(--color-text-secondary)',
                        lineHeight: 1.5
                    }}>
                        {message}
                    </p>
                </div>

                {/* Actions */}
                <div
                    style={{
                        display: 'flex',
                        gap: '0.75rem',
                        padding: '1rem 1.5rem',
                        borderTop: '1px solid var(--color-border)',
                        backgroundColor: 'var(--color-background)'
                    }}
                >
                    <Button
                        variant="secondary"
                        onClick={onClose}
                        disabled={loading}
                        style={{ flex: 1 }}
                    >
                        {cancelText}
                    </Button>
                    <Button
                        variant={variant}
                        onClick={onConfirm}
                        loading={loading}
                        style={{ flex: 1 }}
                    >
                        {confirmText}
                    </Button>
                </div>
            </div>
            <style>{`
        @keyframes dialogFadeIn {
          from { opacity: 0; transform: scale(0.95) translateY(10px); }
          to { opacity: 1; transform: scale(1) translateY(0); }
        }
      `}</style>
        </div>,
        document.body
    );
}

export default ConfirmDialog;
