import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useAuth } from '../context/AuthContext';
import { usePreferences } from '../context/PreferencesContext';
import { useNavigate } from 'react-router-dom';
import Button from '../components/common/Button';
import Input from '../components/common/Input';
import { Hexagon, Lock, User, Sparkles, Eye, EyeOff } from 'lucide-react';
import { loginSchema } from '../utils/schemas';

const LoginPage = () => {
    const [showPassword, setShowPassword] = useState(false);
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const { login } = useAuth();
    const { t } = usePreferences();
    const navigate = useNavigate();

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

    return (
        <div style={{
            minHeight: '100vh',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            background: 'linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%)',
            position: 'relative',
            overflow: 'hidden'
        }}>
            {/* Background Effects */}
            <div style={{
                position: 'absolute',
                top: '-50%',
                left: '-50%',
                width: '200%',
                height: '200%',
                background: 'radial-gradient(circle at 50% 50%, rgba(99, 102, 241, 0.08) 0%, transparent 50%)',
                animation: 'pulse 15s ease-in-out infinite'
            }} />

            {/* Floating Particles */}
            <div style={{
                position: 'absolute',
                top: '20%',
                left: '10%',
                width: '300px',
                height: '300px',
                borderRadius: '50%',
                background: 'linear-gradient(135deg, rgba(99, 102, 241, 0.15) 0%, transparent 70%)',
                filter: 'blur(60px)'
            }} />
            <div style={{
                position: 'absolute',
                bottom: '10%',
                right: '15%',
                width: '400px',
                height: '400px',
                borderRadius: '50%',
                background: 'linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, transparent 70%)',
                filter: 'blur(80px)'
            }} />

            {/* Login Card */}
            <div style={{
                position: 'relative',
                width: '100%',
                maxWidth: '420px',
                margin: '1rem',
                backgroundColor: 'rgba(30, 41, 59, 0.95)',
                backdropFilter: 'blur(20px)',
                borderRadius: '24px',
                boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.5), inset 0 1px 0 rgba(255, 255, 255, 0.1)',
                border: '1px solid rgba(99, 102, 241, 0.2)',
                padding: '2.5rem',
                overflow: 'hidden'
            }}>
                {/* Glow Effect */}
                <div style={{
                    position: 'absolute',
                    top: '-2px',
                    left: '50%',
                    transform: 'translateX(-50%)',
                    width: '40%',
                    height: '4px',
                    background: 'linear-gradient(90deg, transparent, var(--color-primary), transparent)',
                    borderRadius: '2px'
                }} />

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
                        color: '#fff',
                        letterSpacing: '-0.02em'
                    }}>
                        {t('login.title') || 'HexAdmin'}
                    </h1>
                    <p style={{
                        fontSize: '0.875rem',
                        color: 'rgba(148, 163, 184, 1)',
                        marginTop: '0.5rem',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.5rem'
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
                            color: 'rgba(148, 163, 184, 1)',
                            marginBottom: '0.5rem'
                        }}>
                            {t('login.username') || 'Username'}
                        </label>
                        <div style={{ position: 'relative' }}>
                            <User size={18} style={{
                                position: 'absolute',
                                left: '1rem',
                                top: '50%',
                                transform: 'translateY(-50%)',
                                color: 'rgba(100, 116, 139, 1)'
                            }} />
                            <input
                                type="text"
                                placeholder={t('login.username_placeholder') || 'Enter your username'}
                                {...register('login')}
                                style={{
                                    width: '100%',
                                    padding: '0.875rem 1rem 0.875rem 2.75rem',
                                    backgroundColor: 'rgba(15, 23, 42, 0.6)',
                                    border: errors.login ? '1px solid #ef4444' : '1px solid rgba(71, 85, 105, 0.5)',
                                    borderRadius: '12px',
                                    fontSize: '0.9375rem',
                                    color: '#fff',
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
                                        e.target.style.borderColor = 'rgba(71, 85, 105, 0.5)';
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
                            color: 'rgba(148, 163, 184, 1)',
                            marginBottom: '0.5rem'
                        }}>
                            {t('login.password') || 'Password'}
                        </label>
                        <div style={{ position: 'relative' }}>
                            <Lock size={18} style={{
                                position: 'absolute',
                                left: '1rem',
                                top: '50%',
                                transform: 'translateY(-50%)',
                                color: 'rgba(100, 116, 139, 1)'
                            }} />
                            <input
                                type={showPassword ? 'text' : 'password'}
                                placeholder="••••••••"
                                {...register('password')}
                                style={{
                                    width: '100%',
                                    padding: '0.875rem 3rem 0.875rem 2.75rem',
                                    backgroundColor: 'rgba(15, 23, 42, 0.6)',
                                    border: errors.password ? '1px solid #ef4444' : '1px solid rgba(71, 85, 105, 0.5)',
                                    borderRadius: '12px',
                                    fontSize: '0.9375rem',
                                    color: '#fff',
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
                                        e.target.style.borderColor = 'rgba(71, 85, 105, 0.5)';
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
                                    color: 'rgba(100, 116, 139, 1)',
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'center',
                                    transition: 'color 0.2s'
                                }}
                                onMouseOver={(e) => e.currentTarget.style.color = 'var(--color-primary)'}
                                onMouseOut={(e) => e.currentTarget.style.color = 'rgba(100, 116, 139, 1)'}
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
                    color: 'rgba(100, 116, 139, 1)'
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
