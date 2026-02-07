import React, { useEffect, useState } from 'react';
import { Users, Shield, Server, Activity, Clock, TrendingUp, Zap, Database } from 'lucide-react';
import { StatCard } from '../components/common/Card';
import { SkeletonCard } from '../components/feedback/Skeleton';
import { usePreferences } from '../context/PreferencesContext';
import api from '../services/api';
import styles from './DashboardPage.module.css';

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
            <div className={styles.header}>
                <div className={styles.titleWrapper}>
                    <h1 className={styles.title}>
                        <Zap size={28} className={styles.titleIcon} />
                        {t('dashboard.title') || 'Dashboard'}
                    </h1>
                    <p className={styles.subtitle}>
                        {t('dashboard.subtitle') || 'System overview and real-time metrics'}
                    </p>
                </div>

                {lastUpdated && (
                    <div className={styles.updatedBadge}>
                        <Clock size={14} />
                        {t('dashboard.updated') || 'Updated'} {lastUpdated.toLocaleTimeString()}
                    </div>
                )}
            </div>

            {/* Primary Stats */}
            <div className={styles.statsGrid}>
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
            <div className={styles.infoGrid}>
                {/* Quick Stats Card */}
                <div className={styles.infoCard}>
                    <div className={styles.infoCardHeader}>
                        <div className={`${styles.infoIconWrapper} ${styles.quickStatsIcon}`}>
                            <TrendingUp size={18} />
                        </div>
                        <h3 className={styles.infoCardTitle}>
                            {t('dashboard.quick_stats') || 'Quick Stats'}
                        </h3>
                    </div>

                    <div className={styles.infoList}>
                        <div className={styles.listItem}>
                            <span className={styles.label}>
                                {t('dashboard.registered_users') || 'Registered Users'}
                            </span>
                            <span className={styles.value}>
                                {loading ? '...' : stats.users}
                            </span>
                        </div>
                        <div className={styles.listItem}>
                            <span className={styles.label}>
                                {t('dashboard.defined_roles') || 'Defined Roles'}
                            </span>
                            <span className={styles.value}>
                                {loading ? '...' : stats.roles}
                            </span>
                        </div>
                        <div className={styles.listItem}>
                            <span className={styles.label}>
                                {t('dashboard.uptime') || 'Uptime'}
                            </span>
                            <span className={styles.valueSuccess}>
                                {loading ? '...' : stats.uptime}
                            </span>
                        </div>
                    </div>
                </div>

                {/* System Info Card */}
                <div className={styles.infoCard}>
                    <div className={styles.infoCardHeader}>
                        <div className={`${styles.infoIconWrapper} ${styles.systemIcon}`}>
                            <Database size={18} />
                        </div>
                        <h3 className={styles.infoCardTitle}>
                            {t('dashboard.system_info') || 'System Information'}
                        </h3>
                    </div>

                    <div className={styles.systemList}>
                        <div className={styles.listItemBordered}>
                            <span className={styles.label}>
                                {t('dashboard.application') || 'Application'}
                            </span>
                            <span className={styles.valueSmall}>
                                HexAdmin {stats.version || 'v1.0.0'}
                            </span>
                        </div>
                        <div className={styles.listItemBordered}>
                            <span className={styles.label}>
                                {t('dashboard.environment') || 'Environment'}
                            </span>
                            <span className={styles.tag}>
                                Production
                            </span>
                        </div>
                        <div className={styles.listItemBordered}>
                            <span className={styles.label}>
                                {t('dashboard.api_endpoint') || 'API Endpoint'}
                            </span>
                            <span className={styles.valueStatus}>
                                {stats.serverStatus === 'Online' ? `✓ ${t('dashboard.connected') || 'Connected'}` : `✗ ${t('dashboard.disconnected') || 'Disconnected'}`}
                            </span>
                        </div>
                        <div className={styles.listItemBordered}>
                            <span className={styles.label}>
                                {t('dashboard.last_update') || 'Last Update'}
                            </span>
                            <span className={styles.valueStatus}>
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
