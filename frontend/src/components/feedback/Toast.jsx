import React, { useEffect, useState } from 'react';
import { CheckCircle, XCircle, AlertTriangle, Info, X } from 'lucide-react';

const iconMap = {
    success: CheckCircle,
    error: XCircle,
    warning: AlertTriangle,
    info: Info
};

const colorMap = {
    success: {
        bg: 'rgba(16, 185, 129, 0.1)',
        border: '#10b981',
        icon: '#10b981'
    },
    error: {
        bg: 'rgba(239, 68, 68, 0.1)',
        border: '#ef4444',
        icon: '#ef4444'
    },
    warning: {
        bg: 'rgba(245, 158, 11, 0.1)',
        border: '#f59e0b',
        icon: '#f59e0b'
    },
    info: {
        bg: 'rgba(99, 102, 241, 0.1)',
        border: '#6366f1',
        icon: '#6366f1'
    }
};

/**
 * Toast notification component
 */
function Toast({ message, type = 'info', onClose }) {
    const [isVisible, setIsVisible] = useState(false);
    const Icon = iconMap[type];
    const colors = colorMap[type];

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
            style={{
                display: 'flex',
                alignItems: 'center',
                gap: '0.75rem',
                padding: '0.875rem 1rem',
                backgroundColor: 'var(--color-surface)',
                borderRadius: 'var(--radius-lg)',
                boxShadow: 'var(--shadow-lg)',
                border: `1px solid ${colors.border}`,
                borderLeft: `4px solid ${colors.border}`,
                minWidth: '280px',
                maxWidth: '420px',
                pointerEvents: 'auto',
                transform: isVisible ? 'translateX(0)' : 'translateX(120%)',
                opacity: isVisible ? 1 : 0,
                transition: 'all 0.2s cubic-bezier(0.16, 1, 0.3, 1)',
            }}
        >
            <Icon size={20} style={{ color: colors.icon, flexShrink: 0 }} />
            <span style={{
                flex: 1,
                fontSize: '0.875rem',
                color: 'var(--color-text-main)',
                fontWeight: 500
            }}>
                {message}
            </span>
            <button
                onClick={handleClose}
                style={{
                    background: 'transparent',
                    border: 'none',
                    color: 'var(--color-text-muted)',
                    cursor: 'pointer',
                    padding: '0.25rem',
                    borderRadius: 'var(--radius-sm)',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    transition: 'color 0.15s'
                }}
                onMouseOver={e => e.currentTarget.style.color = 'var(--color-text-main)'}
                onMouseOut={e => e.currentTarget.style.color = 'var(--color-text-muted)'}
            >
                <X size={16} />
            </button>
        </div>
    );
}

export default Toast;
