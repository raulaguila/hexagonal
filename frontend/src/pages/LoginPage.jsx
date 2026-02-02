import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import Button from '../components/common/Button';
import Input from '../components/common/Input';
import { Hexagon } from 'lucide-react';

const LoginPage = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const { login } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setIsLoading(true);

        try {
            await login(username, password);
            // Valid login redirects to dashboard via ProtectedRoute logic or manual
            navigate('/dashboard');
        } catch (err) {
            console.error(err);
            setError('Invalid email or password');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div style={{
            backgroundColor: 'var(--color-surface)',
            padding: '2.5rem',
            borderRadius: 'var(--radius-lg)',
            boxShadow: 'var(--shadow-lg)',
            border: '1px solid var(--color-border)',
            textAlign: 'center'
        }}>
            <div style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                gap: '0.75rem',
                marginBottom: '2rem',
                color: 'var(--color-primary)'
            }}>
                <Hexagon size={40} strokeWidth={2.5} />
                <h1 style={{ fontSize: '1.75rem', fontWeight: 700, margin: 0 }}>HexAdmin</h1>
            </div>

            <h2 style={{ fontSize: '1.25rem', fontWeight: 600, marginBottom: '0.5rem' }}>Welcome back</h2>
            <p style={{ color: 'var(--color-text-seconday)', fontSize: '0.875rem', marginBottom: '2rem' }}>
                Please sign in to your account
            </p>

            {error && (
                <div style={{
                    backgroundColor: '#FEF2F2',
                    color: 'var(--color-error)',
                    padding: '0.75rem',
                    borderRadius: 'var(--radius-sm)',
                    fontSize: '0.875rem',
                    marginBottom: '1rem',
                    textAlign: 'left'
                }}>
                    {error}
                </div>
            )}

            <form onSubmit={handleSubmit} style={{ textAlign: 'left' }}>
                <Input
                    label="Username"
                    type="text"
                    placeholder="admin"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    required
                />
                <Input
                    label="Password"
                    type="password"
                    placeholder="••••••••"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                />

                <Button
                    type="submit"
                    variant="primary"
                    disabled={isLoading}
                    style={{
                        marginTop: '1.5rem',
                        height: '48px',
                        fontSize: '1rem',
                        fontWeight: 600,
                        boxShadow: '0 4px 12px rgba(99, 102, 241, 0.25)'
                    }}
                >
                    {isLoading ? 'Signing in...' : 'Sign In'}
                </Button>
            </form>

            <div style={{ marginTop: '1.5rem', fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                Don't have an account? <span style={{ color: 'var(--color-primary)', cursor: 'pointer' }}>Contact Admin</span>
            </div>
        </div>
    );
};

export default LoginPage;
