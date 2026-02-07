import React from 'react';
import { createPortal } from 'react-dom';
import { AlertTriangle } from 'lucide-react';
import Button from '../common/Button';
import styles from './ConfirmDialog.module.css';

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

    // Map variant to style classes
    const wrapperClass = styles[`${variant}Wrapper`];
    const iconClass = styles[`${variant}Icon`];

    return createPortal(
        <div
            className={styles.overlay}
            onClick={onClose}
        >
            {/* Dialog */}
            <div
                className={styles.dialog}
                onClick={e => e.stopPropagation()}
            >
                {/* Content */}
                <div className={styles.content}>
                    {/* Icon */}
                    <div className={`${styles.iconWrapper} ${wrapperClass}`}>
                        <AlertTriangle size={24} className={iconClass} />
                    </div>

                    {/* Title */}
                    <h3 className={styles.title}>
                        {title}
                    </h3>

                    {/* Message */}
                    <p className={styles.message}>
                        {message}
                    </p>
                </div>

                {/* Actions */}
                <div className={styles.actions}>
                    <Button
                        variant="secondary"
                        onClick={onClose}
                        disabled={loading}
                        className={styles.button}
                    >
                        {cancelText}
                    </Button>
                    <Button
                        variant={variant}
                        onClick={onConfirm}
                        loading={loading}
                        className={styles.button}
                    >
                        {confirmText}
                    </Button>
                </div>
            </div>
        </div>,
        document.body
    );
}

export default ConfirmDialog;
