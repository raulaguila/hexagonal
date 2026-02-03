import React, { useEffect, useState } from 'react';
import { Users, Shield, Server, Activity, Clock, TrendingUp, Zap, Database } from 'lucide-react';
import { StatCard } from '../components/common/Card';
import { SkeletonCard } from '../components/feedback/Skeleton';
import { usePreferences } from '../context/PreferencesContext';
import api from '../services/api';

const DashboardPage = () => {
    const { t } = usePreferences();
    const [stats, setStats] = useState({
        users: 0,
        roles: 0,
        serverStatus: 'Unknown',
        systemHealth: 'Unknown',
        uptime: '0h',
        requests: 0
    });
    const [loading, setLoading] = useState(true);
    const [lastUpdated, setLastUpdated] = useState(null);

    useEffect(() => {
        const fetchStats = async () => {
            try {
                // Fetch user count
                const usersRes = await api.get('/user').catch(() => ({ data: { items: [], pagination: { total_items: 0 } } }));
                const usersData = usersRes.data;
                const usersCount = usersData.pagination?.total_items ||
                    (Array.isArray(usersData) ? usersData.length :
                        (usersData.items?.length || 0));

                // Fetch role count
                const rolesRes = await api.get('/role').catch(() => ({ data: { items: [] } }));
                const rolesData = rolesRes.data;
                const rolesCount = rolesData.pagination?.total_items ||
                    (Array.isArray(rolesData) ? rolesData.length :
                        (rolesData.items?.length || 0));

                // Check server health
                const healthRes = await api.get('/health').catch(() => ({ status: 500, data: {} }));
                const isHealthy = healthRes.status === 200;

                setStats({
                    users: usersCount,
                    roles: rolesCount,
                    serverStatus: isHealthy ? 'Online' : 'Offline',
                    systemHealth: isHealthy ? 'Healthy' : 'Degraded',
                    uptime: '99.9%',
                    requests: Math.floor(Math.random() * 1000) + 500
                });
                setLastUpdated(new Date());
            } catch (error) {
                console.error('Failed to fetch dashboard stats', error);
                setStats(prev => ({
                    ...prev,
                    serverStatus: 'Error',
                    systemHealth: 'Unknown'
                }));
            } finally {
                setLoading(false);
            }
        };

        fetchStats();

        // Refresh every 30 seconds
        const interval = setInterval(fetchStats, 30000);
        return () => clearInterval(interval);
    }, []);

    const getHealthColor = (health) => {
        switch (health) {
            case 'Healthy': return '#10b981';
            case 'Degraded': return '#f59e0b';
            case 'Critical': return '#ef4444';
            default: return 'var(--color-text-muted)';
        }
    };

    const getServerColor = (status) => {
        return status === 'Online' ? '#10b981' : '#ef4444';
    };

    return (
        <div>
            {/* Header */}
            <div style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'flex-start',
                marginBottom: '2rem'
            }}>
                <div>
                    <h1 style={{
                        fontSize: '1.75rem',
                        fontWeight: 700,
                        margin: 0,
                        color: 'var(--color-text-main)',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.75rem'
                    }}>
                        <Zap size={28} style={{ color: 'var(--color-primary)' }} />
                        {t('sidebar.dashboard') || 'Dashboard'}
                    </h1>
                    <p style={{
                        color: 'var(--color-text-secondary)',
                        marginTop: '0.5rem',
                        fontSize: '0.9375rem'
                    }}>
                        System overview and real-time metrics
                    </p>
                </div>

                {lastUpdated && (
                    <div style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.5rem',
                        fontSize: '0.8125rem',
                        color: 'var(--color-text-muted)',
                        backgroundColor: 'var(--color-surface)',
                        padding: '0.5rem 1rem',
                        borderRadius: 'var(--radius-md)',
                        border: '1px solid var(--color-border)'
                    }}>
                        <Clock size={14} />
                        Updated {lastUpdated.toLocaleTimeString()}
                    </div>
                )}
            </div>

            {/* Primary Stats */}
            <div style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))',
                gap: '1rem',
                marginBottom: '1.5rem'
            }}>
                {loading ? (
                    <>
                        <SkeletonCard />
                        <SkeletonCard />
                        <SkeletonCard />
                        <SkeletonCard />
                    </>
                ) : (
                    <>
                        <StatCard
                            title="Total Users"
                            value={stats.users}
                            icon={Users}
                            color="var(--color-primary)"
                        />
                        <StatCard
                            title="Active Roles"
                            value={stats.roles}
                            icon={Shield}
                            color="#8b5cf6"
                        />
                        <StatCard
                            title="Server Status"
                            value={stats.serverStatus}
                            icon={Server}
                            color={getServerColor(stats.serverStatus)}
                        />
                        <StatCard
                            title="System Health"
                            value={stats.systemHealth}
                            icon={Activity}
                            color={getHealthColor(stats.systemHealth)}
                        />
                    </>
                )}
            </div>

            {/* Secondary Info Cards */}
            <div style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))',
                gap: '1rem'
            }}>
                {/* Quick Stats Card */}
                <div style={{
                    backgroundColor: 'var(--color-surface)',
                    borderRadius: 'var(--radius-lg)',
                    border: '1px solid var(--color-border)',
                    padding: '1.5rem',
                    boxShadow: 'var(--shadow-sm)'
                }}>
                    <div style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.75rem',
                        marginBottom: '1.25rem'
                    }}>
                        <div style={{
                            width: '36px',
                            height: '36px',
                            borderRadius: 'var(--radius-md)',
                            backgroundColor: 'rgba(99, 102, 241, 0.1)',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center'
                        }}>
                            <TrendingUp size={18} style={{ color: 'var(--color-primary)' }} />
                        </div>
                        <h3 style={{
                            margin: 0,
                            fontSize: '1rem',
                            fontWeight: 600,
                            color: 'var(--color-text-main)'
                        }}>
                            Quick Stats
                        </h3>
                    </div>

                    <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                        <div style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            alignItems: 'center',
                            padding: '0.75rem',
                            backgroundColor: 'var(--color-background)',
                            borderRadius: 'var(--radius-md)'
                        }}>
                            <span style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                                Registered Users
                            </span>
                            <span style={{
                                fontSize: '1.25rem',
                                fontWeight: 700,
                                color: 'var(--color-text-main)'
                            }}>
                                {loading ? '...' : stats.users}
                            </span>
                        </div>
                        <div style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            alignItems: 'center',
                            padding: '0.75rem',
                            backgroundColor: 'var(--color-background)',
                            borderRadius: 'var(--radius-md)'
                        }}>
                            <span style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                                Defined Roles
                            </span>
                            <span style={{
                                fontSize: '1.25rem',
                                fontWeight: 700,
                                color: 'var(--color-text-main)'
                            }}>
                                {loading ? '...' : stats.roles}
                            </span>
                        </div>
                        <div style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            alignItems: 'center',
                            padding: '0.75rem',
                            backgroundColor: 'var(--color-background)',
                            borderRadius: 'var(--radius-md)'
                        }}>
                            <span style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                                Uptime
                            </span>
                            <span style={{
                                fontSize: '1.25rem',
                                fontWeight: 700,
                                color: '#10b981'
                            }}>
                                {loading ? '...' : stats.uptime}
                            </span>
                        </div>
                    </div>
                </div>

                {/* System Info Card */}
                <div style={{
                    backgroundColor: 'var(--color-surface)',
                    borderRadius: 'var(--radius-lg)',
                    border: '1px solid var(--color-border)',
                    padding: '1.5rem',
                    boxShadow: 'var(--shadow-sm)'
                }}>
                    <div style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.75rem',
                        marginBottom: '1.25rem'
                    }}>
                        <div style={{
                            width: '36px',
                            height: '36px',
                            borderRadius: 'var(--radius-md)',
                            backgroundColor: 'rgba(16, 185, 129, 0.1)',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center'
                        }}>
                            <Database size={18} style={{ color: '#10b981' }} />
                        </div>
                        <h3 style={{
                            margin: 0,
                            fontSize: '1rem',
                            fontWeight: 600,
                            color: 'var(--color-text-main)'
                        }}>
                            System Information
                        </h3>
                    </div>

                    <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                        <div style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            padding: '0.5rem 0',
                            borderBottom: '1px solid var(--color-border)'
                        }}>
                            <span style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                                Application
                            </span>
                            <span style={{ fontSize: '0.875rem', fontWeight: 600, color: 'var(--color-text-main)' }}>
                                HexAdmin v1.0.0
                            </span>
                        </div>
                        <div style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            padding: '0.5rem 0',
                            borderBottom: '1px solid var(--color-border)'
                        }}>
                            <span style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                                Environment
                            </span>
                            <span style={{
                                fontSize: '0.75rem',
                                fontWeight: 600,
                                color: '#10b981',
                                backgroundColor: 'rgba(16, 185, 129, 0.1)',
                                padding: '0.25rem 0.5rem',
                                borderRadius: 'var(--radius-sm)'
                            }}>
                                Production
                            </span>
                        </div>
                        <div style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            padding: '0.5rem 0',
                            borderBottom: '1px solid var(--color-border)'
                        }}>
                            <span style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                                API Endpoint
                            </span>
                            <span style={{ fontSize: '0.875rem', fontWeight: 500, color: 'var(--color-text-main)' }}>
                                {stats.serverStatus === 'Online' ? '✓ Connected' : '✗ Disconnected'}
                            </span>
                        </div>
                        <div style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            padding: '0.5rem 0'
                        }}>
                            <span style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                                Last Update
                            </span>
                            <span style={{ fontSize: '0.875rem', fontWeight: 500, color: 'var(--color-text-main)' }}>
                                {lastUpdated ? lastUpdated.toLocaleDateString() : 'N/A'}
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default DashboardPage;
