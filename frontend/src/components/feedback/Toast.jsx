import React, { useEffect, useState } from 'react';
import { CheckCircle, XCircle, AlertTriangle, Info, X } from 'lucide-react';
import styles from './Toast.module.css';

const iconMap = {
    success: CheckCircle,
    error: XCircle,
    warning: AlertTriangle,
    info: Info
};

/**
 * Toast notification component
 */
function Toast({ message, type = 'info', onClose }) {
    const [isVisible, setIsVisible] = useState(false);
    const Icon = iconMap[type];

    useEffect(() => {
        // Trigger enter animation
        requestAnimationFrame(() => setIsVisible(true));
    }, []);

    const handleClose = () => {
        setIsVisible(false);
        setTimeout(onClose, 200);
    };

    return (
        <div
            className={`${styles.toast} ${styles[type]} ${isVisible ? styles.visible : ''}`}
        >
            <Icon size={20} className={styles.icon} />
            <span className={styles.content}>
                {message}
            </span>
            <button
                onClick={handleClose}
                className={styles.closeButton}
            >
                <X size={16} />
            </button>
        </div>
    );
}

export default Toast;
