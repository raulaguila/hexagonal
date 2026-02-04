import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useAuth } from '../context/AuthContext';
import { usePreferences } from '../context/PreferencesContext';
import { useNavigate } from 'react-router-dom';
import { Hexagon, Lock, User, Sparkles, Eye, EyeOff, Moon, Sun } from 'lucide-react';
import { loginSchema } from '../utils/schemas';

const LoginPage = () => {
    const [showPassword, setShowPassword] = useState(false);
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const { login } = useAuth();
    const { t, theme, toggleTheme } = usePreferences();
    const navigate = useNavigate();

    const isDark = theme === 'dark';

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm({
        resolver: zodResolver(loginSchema),
        defaultValues: {
            login: '',
            password: '',
        },
    });

    const onSubmit = async (data) => {
        setError('');
        setIsLoading(true);

        try {
            await login(data.login, data.password);
            navigate('/dashboard');
        } catch (err) {
            console.error(err);
            setError(t('login.error') || 'Invalid username or password');
        } finally {
            setIsLoading(false);
        }
    };

    // Theme-aware colors - clean solid design for both modes
    const colors = {
        background: isDark ? '#0f172a' : '#f1f5f9',
        cardBg: isDark ? '#1e293b' : '#ffffff',
        cardBorder: isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.08)',
        cardShadow: isDark
            ? '0 10px 40px rgba(0, 0, 0, 0.4)'
            : '0 10px 40px rgba(0, 0, 0, 0.08)',
        title: isDark ? '#ffffff' : '#0f172a',
        subtitle: isDark ? '#94a3b8' : '#64748b',
        label: isDark ? '#94a3b8' : '#475569',
        inputBg: isDark ? '#0f172a' : '#f8fafc',
        inputBorder: isDark ? '#334155' : '#e2e8f0',
        inputText: isDark ? '#ffffff' : '#0f172a',
        inputIcon: isDark ? '#64748b' : '#94a3b8',
        footerText: '#64748b',
    };

    return (
        <div style={{
            minHeight: '100vh',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            background: colors.background,
            transition: 'background 0.3s ease'
        }}>
            {/* Theme Toggle Button */}
            <button
                onClick={toggleTheme}
                style={{
                    position: 'absolute',
                    top: '1.5rem',
                    right: '1.5rem',
                    width: '44px',
                    height: '44px',
                    borderRadius: '12px',
                    border: 'none',
                    backgroundColor: isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.05)',
                    color: isDark ? '#fff' : '#64748b',
                    cursor: 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    transition: 'all 0.2s ease',
                    zIndex: 10
                }}
                onMouseOver={(e) => {
                    e.currentTarget.style.transform = 'scale(1.05)';
                    e.currentTarget.style.backgroundColor = isDark ? 'rgba(255, 255, 255, 0.15)' : 'rgba(0, 0, 0, 0.1)';
                }}
                onMouseOut={(e) => {
                    e.currentTarget.style.transform = 'scale(1)';
                    e.currentTarget.style.backgroundColor = isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.05)';
                }}
                title={isDark ? 'Switch to Light Mode' : 'Switch to Dark Mode'}
            >
                {isDark ? <Sun size={20} /> : <Moon size={20} />}
            </button>

            {/* Login Card */}
            <div style={{
                position: 'relative',
                width: '100%',
                maxWidth: '420px',
                margin: '1rem',
                backgroundColor: colors.cardBg,
                backdropFilter: 'blur(20px)',
                borderRadius: '24px',
                boxShadow: colors.cardShadow,
                border: `1px solid ${colors.cardBorder}`,
                padding: '2.5rem',
                overflow: 'hidden',
                transition: 'all 0.3s ease'
            }}>
                {/* Logo */}
                <div style={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    marginBottom: '2rem'
                }}>
                    <div style={{
                        width: '72px',
                        height: '72px',
                        borderRadius: '20px',
                        background: 'linear-gradient(135deg, var(--color-primary) 0%, #818cf8 100%)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        marginBottom: '1.25rem',
                        boxShadow: '0 8px 32px rgba(99, 102, 241, 0.3)'
                    }}>
                        <Hexagon size={36} strokeWidth={2} color="#fff" />
                    </div>
                    <h1 style={{
                        fontSize: '1.5rem',
                        fontWeight: 700,
                        margin: 0,
                        color: colors.title,
                        letterSpacing: '-0.02em',
                        transition: 'color 0.3s ease'
                    }}>
                        {t('login.title') || 'HexAdmin'}
                    </h1>
                    <p style={{
                        fontSize: '0.875rem',
                        color: colors.subtitle,
                        marginTop: '0.5rem',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.5rem',
                        transition: 'color 0.3s ease'
                    }}>
                        <Sparkles size={14} />
                        {t('login.subtitle') || 'Administrative Portal'}
                    </p>
                </div>

                {/* Error Message */}
                {error && (
                    <div style={{
                        backgroundColor: 'rgba(239, 68, 68, 0.1)',
                        color: '#f87171',
                        padding: '0.875rem 1rem',
                        borderRadius: '12px',
                        fontSize: '0.875rem',
                        marginBottom: '1.5rem',
                        border: '1px solid rgba(239, 68, 68, 0.2)',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.5rem'
                    }}>
                        <Lock size={16} />
                        {error}
                    </div>
                )}

                {/* Form */}
                <form onSubmit={handleSubmit(onSubmit)}>
                    <div style={{ marginBottom: '1.25rem' }}>
                        <label style={{
                            display: 'block',
                            fontSize: '0.875rem',
                            fontWeight: 500,
                            color: colors.label,
                            marginBottom: '0.5rem',
                            transition: 'color 0.3s ease'
                        }}>
                            {t('login.username') || 'Username'}
                        </label>
                        <div style={{ position: 'relative' }}>
                            <User size={18} style={{
                                position: 'absolute',
                                left: '1rem',
                                top: '50%',
                                transform: 'translateY(-50%)',
                                color: colors.inputIcon,
                                transition: 'color 0.3s ease'
                            }} />
                            <input
                                type="text"
                                placeholder={t('login.username_placeholder') || 'Enter your username'}
                                {...register('login')}
                                style={{
                                    width: '100%',
                                    padding: '0.875rem 1rem 0.875rem 2.75rem',
                                    backgroundColor: colors.inputBg,
                                    border: errors.login ? '1px solid #ef4444' : `1px solid ${colors.inputBorder}`,
                                    borderRadius: '12px',
                                    fontSize: '0.9375rem',
                                    color: colors.inputText,
                                    outline: 'none',
                                    transition: 'all 0.2s',
                                    boxSizing: 'border-box'
                                }}
                                onFocus={(e) => {
                                    if (!errors.login) {
                                        e.target.style.borderColor = 'var(--color-primary)';
                                        e.target.style.boxShadow = '0 0 0 3px rgba(99, 102, 241, 0.15)';
                                    }
                                }}
                                onBlur={(e) => {
                                    if (!errors.login) {
                                        e.target.style.borderColor = colors.inputBorder;
                                        e.target.style.boxShadow = 'none';
                                    }
                                }}
                            />
                        </div>
                        {errors.login && (
                            <span style={{ color: '#f87171', fontSize: '0.75rem', marginTop: '0.25rem', display: 'block' }}>
                                {errors.login.message}
                            </span>
                        )}
                    </div>

                    <div style={{ marginBottom: '1.5rem' }}>
                        <label style={{
                            display: 'block',
                            fontSize: '0.875rem',
                            fontWeight: 500,
                            color: colors.label,
                            marginBottom: '0.5rem',
                            transition: 'color 0.3s ease'
                        }}>
                            {t('login.password') || 'Password'}
                        </label>
                        <div style={{ position: 'relative' }}>
                            <Lock size={18} style={{
                                position: 'absolute',
                                left: '1rem',
                                top: '50%',
                                transform: 'translateY(-50%)',
                                color: colors.inputIcon,
                                transition: 'color 0.3s ease'
                            }} />
                            <input
                                type={showPassword ? 'text' : 'password'}
                                placeholder="••••••••"
                                {...register('password')}
                                style={{
                                    width: '100%',
                                    padding: '0.875rem 3rem 0.875rem 2.75rem',
                                    backgroundColor: colors.inputBg,
                                    border: errors.password ? '1px solid #ef4444' : `1px solid ${colors.inputBorder}`,
                                    borderRadius: '12px',
                                    fontSize: '0.9375rem',
                                    color: colors.inputText,
                                    outline: 'none',
                                    transition: 'all 0.2s',
                                    boxSizing: 'border-box'
                                }}
                                onFocus={(e) => {
                                    if (!errors.password) {
                                        e.target.style.borderColor = 'var(--color-primary)';
                                        e.target.style.boxShadow = '0 0 0 3px rgba(99, 102, 241, 0.15)';
                                    }
                                }}
                                onBlur={(e) => {
                                    if (!errors.password) {
                                        e.target.style.borderColor = colors.inputBorder;
                                        e.target.style.boxShadow = 'none';
                                    }
                                }}
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                style={{
                                    position: 'absolute',
                                    right: '1rem',
                                    top: '50%',
                                    transform: 'translateY(-50%)',
                                    background: 'none',
                                    border: 'none',
                                    padding: '0.25rem',
                                    cursor: 'pointer',
                                    color: colors.inputIcon,
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'center',
                                    transition: 'color 0.2s'
                                }}
                                onMouseOver={(e) => e.currentTarget.style.color = 'var(--color-primary)'}
                                onMouseOut={(e) => e.currentTarget.style.color = colors.inputIcon}
                            >
                                {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                            </button>
                        </div>
                        {errors.password && (
                            <span style={{ color: '#f87171', fontSize: '0.75rem', marginTop: '0.25rem', display: 'block' }}>
                                {errors.password.message}
                            </span>
                        )}
                    </div>

                    <button
                        type="submit"
                        disabled={isLoading}
                        style={{
                            width: '100%',
                            padding: '0.875rem',
                            background: 'linear-gradient(135deg, var(--color-primary) 0%, #818cf8 100%)',
                            border: 'none',
                            borderRadius: '12px',
                            color: '#fff',
                            fontSize: '0.9375rem',
                            fontWeight: 600,
                            cursor: isLoading ? 'not-allowed' : 'pointer',
                            opacity: isLoading ? 0.7 : 1,
                            transition: 'all 0.2s',
                            boxShadow: '0 4px 15px rgba(99, 102, 241, 0.35)',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            gap: '0.5rem'
                        }}
                        onMouseOver={(e) => {
                            if (!isLoading) e.target.style.transform = 'translateY(-1px)';
                        }}
                        onMouseOut={(e) => {
                            e.target.style.transform = 'translateY(0)';
                        }}
                    >
                        {isLoading ? (
                            <>
                                <div style={{
                                    width: '18px',
                                    height: '18px',
                                    border: '2px solid rgba(255,255,255,0.3)',
                                    borderTopColor: '#fff',
                                    borderRadius: '50%',
                                    animation: 'spin 0.8s linear infinite'
                                }} />
                                {t('login.signing_in') || 'Signing in...'}
                            </>
                        ) : (
                            t('login.sign_in') || 'Sign In'
                        )}
                    </button>
                </form>

                {/* Footer */}
                <div style={{
                    marginTop: '2rem',
                    textAlign: 'center',
                    fontSize: '0.8125rem',
                    color: colors.footerText,
                    transition: 'color 0.3s ease'
                }}>
                    {t('login.no_account') || "Don't have an account?"}{' '}
                    <span style={{
                        color: 'var(--color-primary)',
                        cursor: 'pointer',
                        fontWeight: 500
                    }}>
                        {t('login.contact_admin') || 'Contact Admin'}
                    </span>
                </div>
            </div>

            {/* Global Styles */}
            <style>{`
                @keyframes spin {
                    to { transform: rotate(360deg); }
                }
                @keyframes pulse {
                    0%, 100% { transform: scale(1); opacity: 1; }
                    50% { transform: scale(1.05); opacity: 0.8; }
                }
            `}</style>
        </div>
    );
};

export default LoginPage;
