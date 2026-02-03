import React, { useState } from 'react';
import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { usePreferences } from '../../context/PreferencesContext';
import { LayoutDashboard, Users, Shield, LogOut, Settings as SettingsIcon, ChevronLeft, ChevronRight, Moon, Sun, Globe } from 'lucide-react';

const SidebarLink = ({ to, icon: Icon, children, isCollapsed }) => {
    return (
        <NavLink
            to={to}
            className="sidebar-link"
            style={({ isActive }) => ({
                display: 'flex',
                alignItems: 'center',
                justifyContent: isCollapsed ? 'center' : 'flex-start',
                gap: isCollapsed ? 0 : '0.875rem',
                padding: '0.875rem 1rem',
                borderRadius: 'var(--radius-md)',
                color: isActive ? 'var(--color-primary)' : 'var(--color-text-secondary)',
                backgroundColor: isActive ? 'var(--color-primary-light)' : 'transparent',
                fontWeight: isActive ? 600 : 500,
                transition: 'all var(--transition-fast)',
                marginBottom: '0.375rem',
                borderLeft: isActive && !isCollapsed ? '3px solid var(--color-primary)' : '3px solid transparent',
                height: '48px',
                width: '100%',
                overflow: 'hidden'
            })}
            title={isCollapsed ? children : ''}
        >
            <Icon size={20} style={{ minWidth: '20px' }} />
            <span style={{
                opacity: isCollapsed ? 0 : 1,
                width: isCollapsed ? 0 : 'auto',
                whiteSpace: 'nowrap',
                transition: 'opacity 0.2s',
                overflow: 'hidden'
            }}>
                {children}
            </span>
        </NavLink>
    );
};

const SidebarSubmenu = ({ icon: Icon, label, children, isCollapsed }) => {
    const [isOpen, setIsOpen] = React.useState(true);

    return (
        <div style={{ marginBottom: '0.5rem' }}>
            <div
                style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: isCollapsed ? 'center' : 'flex-start',
                    gap: isCollapsed ? 0 : '0.875rem',
                    padding: '0.75rem 1rem',
                    color: 'var(--color-text-main)',
                    fontWeight: 600,
                    cursor: 'pointer',
                    borderRadius: 'var(--radius-md)',
                    transition: 'background-color 0.2s',
                    height: '48px'
                }}
                onClick={() => setIsOpen(!isOpen)}
                onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'var(--color-surface-hover)'}
                onMouseOut={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                title={isCollapsed ? label : ''}
            >
                <Icon size={20} style={{ minWidth: '20px' }} />
                {!isCollapsed && <span style={{ flex: 1, whiteSpace: 'nowrap', overflow: 'hidden' }}>{label}</span>}
                {!isCollapsed && <span style={{ fontSize: '0.625rem', transform: isOpen ? 'rotate(180deg)' : 'rotate(0deg)', transition: 'transform 0.2s' }}>▲</span>}
            </div>
            {isOpen && (
                <div style={{
                    marginLeft: isCollapsed ? 0 : '1.25rem',
                    borderLeft: isCollapsed ? 'none' : '2px solid var(--color-border)',
                    paddingLeft: isCollapsed ? 0 : '0.5rem',
                    marginTop: '0.25rem',
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: isCollapsed ? 'center' : 'stretch'
                }}>
                    {React.Children.map(children, child =>
                        React.cloneElement(child, { isCollapsed })
                    )}
                </div>
            )}
        </div>
    );
};

const DashboardLayout = () => {
    const { user, logout } = useAuth();
    const { theme, toggleTheme, language, setLanguage, t } = usePreferences();
    const navigate = useNavigate();
    const [isCollapsed, setIsCollapsed] = useState(false);

    const handleLogout = async () => {
        await logout();
        navigate('/login');
    };

    const toggleLanguage = () => {
        setLanguage(prev => prev === 'en-US' ? 'pt-BR' : 'en-US');
    };

    return (
        <div style={{ display: 'flex', height: '100vh', backgroundColor: 'var(--color-background)', overflow: 'hidden' }}>
            {/* Sidebar */}
            <aside style={{
                width: isCollapsed ? '80px' : '280px',
                backgroundColor: 'var(--color-surface)',
                borderRight: '1px solid var(--color-border)',
                padding: '1.5rem 1rem',
                display: 'flex',
                flexDirection: 'column',
                transition: 'width 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
                position: 'relative',
                zIndex: 20
            }}>
                {/* Logo Area */}
                <div style={{
                    marginBottom: '2rem',
                    paddingLeft: isCollapsed ? 0 : '0.5rem',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: isCollapsed ? 'center' : 'flex-start',
                    height: '40px',
                    overflow: 'hidden'
                }}>
                    <div style={{
                        width: '32px', height: '32px',
                        background: 'linear-gradient(135deg, var(--color-primary), var(--color-primary-hover))',
                        borderRadius: '8px',
                        display: 'flex', alignItems: 'center', justifyContent: 'center',
                        color: 'white', flexShrink: 0,
                        boxShadow: '0 4px 6px -1px rgba(99, 102, 241, 0.4)'
                    }}>
                        <LayoutDashboard size={18} />
                    </div>
                    {!isCollapsed && (
                        <h1 style={{
                            fontSize: '1.25rem', fontWeight: 800, color: 'var(--color-text-main)',
                            margin: '0 0 0 0.75rem', letterSpacing: '-0.025em', whiteSpace: 'nowrap'
                        }}>
                            Hex<span style={{ color: 'var(--color-primary)' }}>Admin</span>
                        </h1>
                    )}
                </div>

                {/* Collapse Toggle */}
                <button
                    onClick={() => setIsCollapsed(!isCollapsed)}
                    style={{
                        position: 'absolute',
                        top: '24px',
                        right: '-12px',
                        width: '24px',
                        height: '24px',
                        borderRadius: '50%',
                        backgroundColor: 'var(--color-surface)',
                        border: '1px solid var(--color-border)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        cursor: 'pointer',
                        color: 'var(--color-text-secondary)',
                        boxShadow: 'var(--shadow-md)',
                        zIndex: 30,
                    }}
                >
                    {isCollapsed ? <ChevronRight size={14} /> : <ChevronLeft size={14} />}
                </button>

                {/* Navigation */}
                <nav style={{ flex: 1, overflowY: 'auto', scrollbarWidth: 'none' }}>
                    <SidebarLink to="/dashboard" icon={LayoutDashboard} isCollapsed={isCollapsed}>{t('sidebar.dashboard')}</SidebarLink>

                    <div style={{ marginTop: '1.5rem' }}>
                        {!isCollapsed && (
                            <div style={{
                                fontSize: '0.75rem', fontWeight: 700, textTransform: 'uppercase',
                                color: 'var(--color-text-muted)', marginBottom: '0.75rem', paddingLeft: '1rem',
                                letterSpacing: '0.05em'
                            }}>
                                {t('sidebar.dashboard') === 'Dashboard' ? 'Settings' : 'Configurações'}
                            </div>
                        )}
                        <SidebarSubmenu icon={SettingsIcon} label={t('sidebar.dashboard') === 'Dashboard' ? 'Configs' : 'Sistema'} isCollapsed={isCollapsed}>
                            <SidebarLink to="/users" icon={Users} isCollapsed={isCollapsed}>{t('sidebar.users')}</SidebarLink>
                            <SidebarLink to="/roles" icon={Shield} isCollapsed={isCollapsed}>{t('sidebar.roles')}</SidebarLink>
                        </SidebarSubmenu>
                    </div>
                </nav>

                {/* Footer Controls (Theme/Lang/User) */}
                <div style={{
                    borderTop: '1px solid var(--color-border)',
                    paddingTop: '1rem',
                    display: 'flex',
                    flexDirection: 'column',
                    gap: '0.5rem'
                }}>

                    {/* Preferences Toggles */}
                    <div style={{ display: 'flex', padding: isCollapsed ? '0 0.25rem' : '0 1rem', gap: '0.5rem' }}>
                        <button
                            onClick={toggleTheme}
                            style={{
                                flex: 1,
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                padding: '0.5rem',
                                borderRadius: 'var(--radius-md)',
                                backgroundColor: 'var(--color-background)',
                                color: 'var(--color-text-secondary)',
                                cursor: 'pointer',
                                transition: 'all 0.2s'
                            }}
                            title="Toggle Theme"
                        >
                            {theme === 'dark' ? <Moon size={18} /> : <Sun size={18} />}
                        </button>

                        <button
                            onClick={toggleLanguage}
                            style={{
                                flex: 1,
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                padding: '0.5rem',
                                borderRadius: 'var(--radius-md)',
                                backgroundColor: 'var(--color-background)',
                                color: 'var(--color-text-secondary)',
                                cursor: 'pointer',
                                transition: 'all 0.2s',
                                fontWeight: 600,
                                fontSize: '0.8rem'
                            }}
                            title="Switch Language"
                        >
                            {isCollapsed ? <Globe size={18} /> : (language === 'en-US' ? 'EN' : 'PT')}
                        </button>
                    </div>

                    {/* User Profile */}
                    <div style={{
                        marginTop: '0.5rem',
                        padding: isCollapsed ? '0.5rem' : '0.75rem',
                        backgroundColor: 'var(--color-surface-hover)',
                        borderRadius: 'var(--radius-md)',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.75rem',
                        cursor: 'default'
                    }}>
                        <div style={{
                            width: '32px', height: '32px', borderRadius: '50%',
                            background: 'var(--color-primary-light)', color: 'var(--color-primary)',
                            display: 'flex', alignItems: 'center', justifyContent: 'center',
                            fontWeight: 'bold', flexShrink: 0
                        }}>
                            {user?.name?.[0] || 'U'}
                        </div>
                        {!isCollapsed && (
                            <div style={{ overflow: 'hidden', flex: 1 }}>
                                <p style={{ margin: 0, fontWeight: 600, color: 'var(--color-text-main)', fontSize: '0.875rem' }}>{user?.name || 'User'}</p>
                                <p style={{ margin: 0, fontSize: '0.75rem', color: 'var(--color-text-muted)', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>@{user?.username}</p>
                            </div>
                        )}
                        {!isCollapsed && (
                            <button
                                onClick={handleLogout}
                                style={{
                                    background: 'transparent',
                                    color: 'var(--color-text-muted)',
                                    cursor: 'pointer',
                                    padding: '4px',
                                    borderRadius: '4px'
                                }}
                                onMouseOver={e => e.currentTarget.style.color = 'var(--color-error)'}
                                onMouseOut={e => e.currentTarget.style.color = 'var(--color-text-muted)'}
                                title="Sign Out"
                            >
                                <LogOut size={16} />
                            </button>
                        )}
                    </div>
                </div>
            </aside>

            {/* Main Content Area */}
            <main style={{
                flex: 1,
                padding: '2rem',
                overflowY: 'auto',
                backgroundColor: 'var(--color-background)'
            }}>
                <Outlet />
            </main>
        </div>
    );
};

export default DashboardLayout;
