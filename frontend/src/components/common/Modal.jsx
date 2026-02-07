import React from 'react';
import { X } from 'lucide-react';
import { createPortal } from 'react-dom';
import styles from './Modal.module.css';

const Modal = ({ isOpen, onClose, title, children, size = 'md' }) => {
    // Simple approach: always render if isOpen is true
    if (!isOpen) return null;

    return createPortal(
        <div
            className={styles.overlay}
            onClick={onClose}
        >
            {/* Content */}
            <div
                className={`${styles.modal} ${styles[size]}`}
                onClick={e => e.stopPropagation()}
            >
                {/* Header */}
                <div className={styles.header}>
                    <h3 className={styles.title}>{title}</h3>
                    <button
                        onClick={onClose}
                        className={styles.closeButton}
                    >
                        <X size={20} />
                    </button>
                </div>

                {/* Body */}
                <div className={styles.body}>
                    {children}
                </div>
            </div>
        </div>,
        document.body
    );
};

export default Modal;
