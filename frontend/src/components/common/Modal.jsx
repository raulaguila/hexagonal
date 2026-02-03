import React from 'react';
import { X } from 'lucide-react';
import { createPortal } from 'react-dom';

const Modal = ({ isOpen, onClose, title, children, size = 'md' }) => {
    // Simple approach: always render if isOpen is true
    if (!isOpen) return null;

    const maxWidth = {
        sm: '400px',
        md: '500px',
        lg: '800px',
        xl: '1100px'
    }[size];

    return createPortal(
        <div
            style={{
                position: 'fixed',
                inset: 0,
                zIndex: 50,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                backgroundColor: 'rgba(0, 0, 0, 0.4)',
                backdropFilter: 'blur(4px)',
            }}
            onClick={onClose}
        >
            {/* Content */}
            <div style={{
                position: 'relative',
                backgroundColor: 'var(--color-surface)',
                borderRadius: 'var(--radius-lg)',
                boxShadow: 'var(--shadow-xl)',
                width: '100%',
                maxWidth: maxWidth,
                margin: '1.5rem',
                maxHeight: 'calc(100vh - 3rem)',
                display: 'flex',
                flexDirection: 'column',
                border: '1px solid var(--color-border)',
                animation: 'modalFadeIn 0.2s ease-out'
            }} onClick={e => e.stopPropagation()}>

                {/* Header */}
                <div style={{
                    padding: '1.25rem 1.5rem',
                    borderBottom: '1px solid var(--color-border)',
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    backgroundColor: 'var(--color-surface)'
                }}>
                    <h3 style={{ margin: 0, fontSize: '1.125rem', fontWeight: 600, color: 'var(--color-text-main)' }}>{title}</h3>
                    <button
                        onClick={onClose}
                        style={{
                            background: 'transparent',
                            color: 'var(--color-text-muted)',
                            cursor: 'pointer',
                            padding: '0.25rem',
                            borderRadius: 'var(--radius-sm)',
                            transition: 'color 0.2s, background-color 0.2s'
                        }}
                        onMouseOver={e => {
                            e.currentTarget.style.color = 'var(--color-text-main)';
                            e.currentTarget.style.backgroundColor = 'var(--color-background)';
                        }}
                        onMouseOut={e => {
                            e.currentTarget.style.color = 'var(--color-text-muted)';
                            e.currentTarget.style.backgroundColor = 'transparent';
                        }}
                    >
                        <X size={20} />
                    </button>
                </div>

                {/* Body */}
                <div style={{
                    padding: '1.5rem',
                    overflowY: 'auto',
                    flex: 1
                }}>
                    {children}
                </div>
            </div>
            <style>{`
                @keyframes modalFadeIn {
                    from { opacity: 0; transform: scale(0.95) translateY(10px); }
                    to { opacity: 1; transform: scale(1) translateY(0); }
                }
            `}</style>
        </div>,
        document.body
    );
};

export default Modal;
