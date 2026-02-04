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
        uptime: 'N/A',
        version: 'v1.0.0'
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

                // Check server health and get uptime
                const healthRes = await api.get('/health').catch(() => ({ status: 500, data: {} }));
                const isHealthy = healthRes.status === 200;
                // API returns { code, message, object } - data is in object
                const healthData = healthRes.data?.object || healthRes.data?.data || healthRes.data || {};

                setStats({
                    users: usersCount,
                    roles: rolesCount,
                    serverStatus: isHealthy ? 'Online' : 'Offline',
                    systemHealth: isHealthy ? 'Healthy' : 'Degraded',
                    uptime: healthData.uptime || 'N/A',
                    version: healthData.version || 'v1.0.0'
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
                        {t('dashboard.title') || 'Dashboard'}
                    </h1>
                    <p style={{
                        color: 'var(--color-text-secondary)',
                        marginTop: '0.5rem',
                        fontSize: '0.9375rem'
                    }}>
                        {t('dashboard.subtitle') || 'System overview and real-time metrics'}
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
                        {t('dashboard.updated') || 'Updated'} {lastUpdated.toLocaleTimeString()}
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
                            title={t('dashboard.total_users') || 'Total Users'}
                            value={stats.users}
                            icon={Users}
                            color="var(--color-primary)"
                        />
                        <StatCard
                            title={t('dashboard.active_roles') || 'Active Roles'}
                            value={stats.roles}
                            icon={Shield}
                            color="#8b5cf6"
                        />
                        <StatCard
                            title={t('dashboard.server_status') || 'Server Status'}
                            value={stats.serverStatus === 'Online' ? (t('dashboard.online') || 'Online') : (t('dashboard.offline') || 'Offline')}
                            icon={Server}
                            color={getServerColor(stats.serverStatus)}
                        />
                        <StatCard
                            title={t('dashboard.system_health') || 'System Health'}
                            value={stats.systemHealth === 'Healthy' ? (t('dashboard.healthy') || 'Healthy') : (t('dashboard.degraded') || 'Degraded')}
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
                            {t('dashboard.quick_stats') || 'Quick Stats'}
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
                                {t('dashboard.registered_users') || 'Registered Users'}
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
                                {t('dashboard.defined_roles') || 'Defined Roles'}
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
                                {t('dashboard.uptime') || 'Uptime'}
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
                            {t('dashboard.system_info') || 'System Information'}
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
                                {t('dashboard.application') || 'Application'}
                            </span>
                            <span style={{ fontSize: '0.875rem', fontWeight: 600, color: 'var(--color-text-main)' }}>
                                HexAdmin {stats.version || 'v1.0.0'}
                            </span>
                        </div>
                        <div style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            padding: '0.5rem 0',
                            borderBottom: '1px solid var(--color-border)'
                        }}>
                            <span style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                                {t('dashboard.environment') || 'Environment'}
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
                                {t('dashboard.api_endpoint') || 'API Endpoint'}
                            </span>
                            <span style={{ fontSize: '0.875rem', fontWeight: 500, color: 'var(--color-text-main)' }}>
                                {stats.serverStatus === 'Online' ? `✓ ${t('dashboard.connected') || 'Connected'}` : `✗ ${t('dashboard.disconnected') || 'Disconnected'}`}
                            </span>
                        </div>
                        <div style={{
                            display: 'flex',
                            justifyContent: 'space-between',
                            padding: '0.5rem 0'
                        }}>
                            <span style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                                {t('dashboard.last_update') || 'Last Update'}
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
