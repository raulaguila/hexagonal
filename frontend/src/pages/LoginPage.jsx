import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useAuth } from '../context/AuthContext';
import { usePreferences } from '../context/PreferencesContext';
import { useNavigate } from 'react-router-dom';
import { Hexagon, Lock, User, Sparkles, Eye, EyeOff, Moon, Sun } from 'lucide-react';
import { loginSchema } from '../utils/schemas';
import styles from './LoginPage.module.css';

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

    return (
        <div className={styles.container}>
            {/* Theme Toggle Button */}
            <button
                onClick={toggleTheme}
                className={styles.themeToggle}
                title={isDark ? 'Switch to Light Mode' : 'Switch to Dark Mode'}
            >
                {isDark ? <Sun size={20} /> : <Moon size={20} />}
            </button>

            {/* Login Card */}
            <div className={styles.card}>
                {/* Logo */}
                <div className={styles.logoSection}>
                    <div className={styles.logoWrapper}>
                        <Hexagon size={36} strokeWidth={2} color="#fff" />
                    </div>
                    <h1 className={styles.title}>
                        {t('login.title') || 'HexAdmin'}
                    </h1>
                    <p className={styles.subtitle}>
                        <Sparkles size={14} />
                        {t('login.subtitle') || 'Administrative Portal'}
                    </p>
                </div>

                {/* Error Message */}
                {error && (
                    <div className={styles.errorMessage}>
                        <Lock size={16} />
                        {error}
                    </div>
                )}

                {/* Form */}
                <form onSubmit={handleSubmit(onSubmit)}>
                    <div className={styles.formGroup}>
                        <label className={styles.label}>
                            {t('login.username') || 'Username'}
                        </label>
                        <div className={styles.inputWrapper}>
                            <User size={18} className={styles.inputIcon} />
                            <input
                                type="text"
                                placeholder={t('login.username_placeholder') || 'Enter your username'}
                                {...register('login')}
                                className={`${styles.input} ${errors.login ? styles.inputError : ''}`}
                            />
                        </div>
                        {errors.login && (
                            <span className={styles.errorText}>
                                {errors.login.message}
                            </span>
                        )}
                    </div>

                    <div className={`${styles.formGroup} ${styles.last}`}>
                        <label className={styles.label}>
                            {t('login.password') || 'Password'}
                        </label>
                        <div className={styles.inputWrapper}>
                            <Lock size={18} className={styles.inputIcon} />
                            <input
                                type={showPassword ? 'text' : 'password'}
                                placeholder="••••••••"
                                {...register('password')}
                                className={`${styles.input} ${errors.password ? styles.inputError : ''}`}
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className={styles.passwordToggle}
                            >
                                {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                            </button>
                        </div>
                        {errors.password && (
                            <span className={styles.errorText}>
                                {errors.password.message}
                            </span>
                        )}
                    </div>

                    <button
                        type="submit"
                        disabled={isLoading}
                        className={styles.submitButton}
                    >
                        {isLoading ? (
                            <>
                                <div className={styles.spinner} />
                                {t('login.signing_in') || 'Signing in...'}
                            </>
                        ) : (
                            t('login.sign_in') || 'Sign In'
                        )}
                    </button>
                </form>

                {/* Footer */}
                <div className={styles.footer}>
                    {t('login.no_account') || "Don't have an account?"}{' '}
                    <span className={styles.contactAdmin}>
                        {t('login.contact_admin') || 'Contact Admin'}
                    </span>
                </div>
            </div>
        </div>
    );
};

export default LoginPage;
