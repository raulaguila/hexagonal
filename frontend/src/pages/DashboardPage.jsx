import React, { useEffect, useState } from 'react';
import api from '../utils/api';
import { Users, Shield, Server, Activity } from 'lucide-react';

const StatCard = ({ title, value, icon: Icon, color }) => {
    return (
        <div style={{
            backgroundColor: 'var(--color-surface)',
            padding: '1.5rem',
            borderRadius: 'var(--radius-lg)',
            boxShadow: 'var(--shadow-sm)',
            display: 'flex',
            alignItems: 'center',
            gap: '1rem',
            border: '1px solid var(--color-border)'
        }}>
            <div style={{
                padding: '0.75rem',
                borderRadius: 'var(--radius-md)',
                backgroundColor: color + '20', // 20% opacity
                color: color
            }}>
                <Icon size={24} />
            </div>
            <div>
                <h3 style={{ margin: 0, fontSize: '0.875rem', color: 'var(--color-text-secondary)', fontWeight: 500 }}>{title}</h3>
                <p style={{ margin: '0.25rem 0 0 0', fontSize: '1.5rem', fontWeight: 700, color: 'var(--color-text-main)' }}>
                    {value}
                </p>
            </div>
        </div>
    );
};

const DashboardPage = () => {
    const [stats, setStats] = useState({
        users: 0,
        roles: 0,
        serverStatus: 'Unknown'
    });
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchStats = async () => {
            try {
                // Fetch Users
                const usersRes = await api.get('/user').catch(() => ({ data: { total: 0, items: [] } }));
                // Adapt based on API response structure. 
                // If it returns { total: 10, items: [...] } or just [...]
                const usersCount = Array.isArray(usersRes.data) ? usersRes.data.length : (usersRes.data.total || 0);

                // Fetch Roles
                const rolesRes = await api.get('/role').catch(() => ({ data: [] }));
                const rolesCount = Array.isArray(rolesRes.data) ? rolesRes.data.length : (rolesRes.data.total || 0);

                // Check Health
                const healthRes = await api.get('/health').catch(() => ({ status: 500 }));
                const isHealthy = healthRes.status === 200;

                setStats({
                    users: usersCount,
                    roles: rolesCount,
                    serverStatus: isHealthy ? 'Online' : 'Offline'
                });
            } catch (error) {
                console.error("Failed to fetch dashboard stats", error);
            } finally {
                setLoading(false);
            }
        };

        fetchStats();
    }, []);

    return (
        <div>
            <div style={{ marginBottom: '2rem' }}>
                <h1 style={{ fontSize: '1.875rem', fontWeight: 700, margin: 0, color: 'var(--color-text-main)' }}>Dashboard</h1>
                <p style={{ color: 'var(--color-text-secondary)', marginTop: '0.5rem' }}>System overview and statistics</p>
            </div>

            <div style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))',
                gap: '1.5rem'
            }}>
                <StatCard
                    title="Total Users"
                    value={loading ? '...' : stats.users}
                    icon={Users}
                    color="var(--color-primary)"
                />
                <StatCard
                    title="Active Roles"
                    value={loading ? '...' : stats.roles}
                    icon={Shield}
                    color="#10B981" // Emerald
                />
                <StatCard
                    title="Server Status"
                    value={loading ? '...' : stats.serverStatus}
                    icon={Server}
                    color={stats.serverStatus === 'Online' ? '#10B981' : '#EF4444'}
                />
                <StatCard
                    title="System Health"
                    value="Good"
                    icon={Activity}
                    color="#F59E0B" // Amber
                />
            </div>
        </div>
    );
};

export default DashboardPage;
