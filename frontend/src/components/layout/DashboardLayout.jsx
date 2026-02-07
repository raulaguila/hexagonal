import React, { useState } from 'react';
import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { usePreferences } from '../../context/PreferencesContext';
import { LayoutDashboard, Users, Shield, LogOut, Settings as SettingsIcon, ChevronLeft, ChevronRight, Moon, Sun, Globe } from 'lucide-react';
import { ConfirmDialog } from '../../components/feedback';
import styles from './DashboardLayout.module.css';

const SidebarLink = ({ to, icon: Icon, children, isCollapsed }) => {
    return (
        <NavLink
            to={to}
            className={({ isActive }) => `${styles.link} ${isActive ? styles.active : ''}`}
            title={isCollapsed ? children : ''}
        >
            <Icon size={20} className={styles.linkIcon} />
            <span className={styles.linkText}>
                {children}
            </span>
        </NavLink>
    );
};

const SidebarSubmenu = ({ icon: Icon, label, children, isCollapsed }) => {
    const [isOpen, setIsOpen] = React.useState(true);

    return (
        <div className={styles.submenuContainer}>
            <div
                className={styles.submenuHeader}
                onClick={() => setIsOpen(!isOpen)}
                title={isCollapsed ? label : ''}
            >
                <Icon size={20} className={styles.linkIcon} />
                {!isCollapsed && <span style={{ flex: 1, whiteSpace: 'nowrap', overflow: 'hidden' }}>{label}</span>}
                {!isCollapsed && <span className={`${styles.submenuArrow} ${isOpen ? styles.rotate180 : ''}`}>▲</span>}
            </div>
            {isOpen && (
                <div className={styles.submenuContent}>
                    {React.Children.toArray(children).filter(Boolean).map(child =>
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
    const [showLogoutConfirm, setShowLogoutConfirm] = useState(false);

    // Helper to get all user permissions (from top-level or aggregated from roles)
    const getUserPermissions = () => {
        // If top-level permissions exists, use it
        if (user?.permissions && user.permissions.length > 0) {
            return user.permissions;
        }
        // Otherwise aggregate from roles
        if (user?.roles && user.roles.length > 0) {
            const perms = new Set();
            user.roles.forEach(role => {
                if (role.permissions) {
                    role.permissions.forEach(p => perms.add(p));
                }
            });
            return Array.from(perms);
        }
        return [];
    };

    // Helper to check if user has any permission for a resource
    const hasAnyPermission = (resource) => {
        const permissions = getUserPermissions();
        if (permissions.length === 0) return true; // If no permissions defined, allow access (admin)
        const resourcePerms = [`${resource}:view`, `${resource}:create`, `${resource}:edit`, `${resource}:delete`];
        return resourcePerms.some(p => permissions.includes(p));
    };

    const canViewUsers = hasAnyPermission('users');
    const canViewRoles = hasAnyPermission('roles');
    const hasConfigsAccess = canViewUsers || canViewRoles;

    const handleLogout = async () => {
        await logout();
        navigate('/login');
    };

    const toggleLanguage = () => {
        setLanguage(prev => prev === 'en-US' ? 'pt-BR' : 'en-US');
    };

    return (
        <div className={styles.container}>
            {/* Sidebar */}
            <aside className={`${styles.sidebar} ${isCollapsed ? styles.collapsed : ''}`}>
                {/* Logo Area */}
                <div className={styles.logoArea}>
                    <div className={styles.logoIcon}>
                        <LayoutDashboard size={18} />
                    </div>
                    {!isCollapsed && (
                        <h1 className={styles.logoText}>
                            Hex<span className={styles.highlight}>Admin</span>
                        </h1>
                    )}
                </div>

                {/* Collapse Toggle */}
                <button
                    onClick={() => setIsCollapsed(!isCollapsed)}
                    className={styles.toggleButton}
                >
                    {isCollapsed ? <ChevronRight size={14} /> : <ChevronLeft size={14} />}
                </button>

                {/* Navigation */}
                <nav className={styles.nav}>
                    <SidebarLink to="/dashboard" icon={LayoutDashboard} isCollapsed={isCollapsed}>{t('sidebar.dashboard')}</SidebarLink>

                    {hasConfigsAccess && (
                        <div style={{ marginTop: '1.5rem' }}>
                            {!isCollapsed && (
                                <div className={styles.sectionLabel}>
                                    {t('sidebar.dashboard') === 'Dashboard' ? 'Settings' : 'Configurações'}
                                </div>
                            )}
                            <SidebarSubmenu icon={SettingsIcon} label={t('sidebar.dashboard') === 'Dashboard' ? 'Configs' : 'Sistema'} isCollapsed={isCollapsed}>
                                {canViewUsers && <SidebarLink to="/users" icon={Users} isCollapsed={isCollapsed}>{t('sidebar.users')}</SidebarLink>}
                                {canViewRoles && <SidebarLink to="/roles" icon={Shield} isCollapsed={isCollapsed}>{t('sidebar.roles')}</SidebarLink>}
                            </SidebarSubmenu>
                        </div>
                    )}
                </nav>

                {/* Footer Controls (Theme/Lang/User) */}
                <div className={styles.footer}>

                    {/* Preferences Toggles */}
                    <div className={styles.controlsRow}>
                        <button
                            onClick={toggleTheme}
                            className={styles.controlButton}
                            title="Toggle Theme"
                        >
                            {theme === 'dark' ? <Moon size={18} /> : <Sun size={18} />}
                        </button>

                        <button
                            onClick={toggleLanguage}
                            className={`${styles.controlButton} ${styles.langButton}`}
                            title="Switch Language"
                        >
                            {isCollapsed ? <Globe size={18} /> : (language === 'en-US' ? 'EN' : 'PT')}
                        </button>
                    </div>

                    {/* User Profile */}
                    <div className={styles.userProfile}>
                        <div className={styles.userAvatar}>
                            {user?.name?.[0] || 'U'}
                        </div>
                        {!isCollapsed && (
                            <div className={styles.userInfo}>
                                <p className={styles.userName}>{user?.name || 'User'}</p>
                                <p className={styles.userHandle}>@{user?.username}</p>
                            </div>
                        )}
                        {!isCollapsed && (
                            <button
                                onClick={() => setShowLogoutConfirm(true)}
                                className={styles.logoutButton}
                                title="Sign Out"
                            >
                                <LogOut size={16} />
                            </button>
                        )}
                    </div>
                </div>
            </aside>

            {/* Main Content Area */}
            <main className={styles.main}>
                <Outlet />
            </main>

            {/* Logout Confirmation Dialog */}
            <ConfirmDialog
                isOpen={showLogoutConfirm}
                onClose={() => setShowLogoutConfirm(false)}
                onConfirm={handleLogout}
                title={language === 'pt-BR' ? 'Confirmar Saída' : 'Confirm Logout'}
                message={language === 'pt-BR' ? 'Tem certeza que deseja sair do sistema?' : 'Are you sure you want to sign out?'}
                confirmText={language === 'pt-BR' ? 'Sair' : 'Sign Out'}
                variant="danger"
            />
        </div>
    );
};

export default DashboardLayout;
