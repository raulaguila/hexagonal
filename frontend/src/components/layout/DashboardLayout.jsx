import React, { useState } from 'react';
import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { LayoutDashboard, Users, Shield, LogOut, Menu, Settings as SettingsIcon, ChevronLeft, ChevronRight } from 'lucide-react';

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
                transition: 'all 0.2s ease',
                marginBottom: '0.375rem',
                borderLeft: isActive && !isCollapsed ? '3px solid var(--color-primary)' : '3px solid transparent',
                height: '48px', // Fixed height for consistency
                width: '100%'
            })}
            title={isCollapsed ? children : ''}
        >
            <Icon size={20} style={{ minWidth: '20px' }} />
            {!isCollapsed && <span style={{ whiteSpace: 'nowrap', overflow: 'hidden' }}>{children}</span>}
        </NavLink>
    );
};

const SidebarSubmenu = ({ icon: Icon, label, children, isCollapsed }) => {
    const [isOpen, setIsOpen] = React.useState(true);

    // If collapsed, we disable the accordion effect for now or handle it differently.
    // Simplifying: If collapsed, we force show children but just as icons? 
    // Or we keep the toggle behavior but remove text?
    // Let's keep toggle but adjust styles.

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
                    userSelect: 'none',
                    borderRadius: 'var(--radius-md)',
                    transition: 'background-color 0.2s',
                    height: '48px'
                }}
                onClick={() => setIsOpen(!isOpen)}
                onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'var(--color-background)'}
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
                    {/* Pass isCollapsed down to children */}
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
    const navigate = useNavigate();
    const [isCollapsed, setIsCollapsed] = useState(false);

    const handleLogout = async () => {
        await logout();
        navigate('/login');
    };

    return (
        <div style={{ display: 'flex', minHeight: '100vh', backgroundColor: 'var(--color-background)' }}>
            {/* Sidebar */}
            <aside style={{
                width: isCollapsed ? '80px' : '280px',
                backgroundColor: 'var(--color-surface)',
                borderRight: '1px solid var(--color-border)',
                padding: '2rem 1rem', // Reduced padding horizontal
                display: 'flex',
                flexDirection: 'column',
                boxShadow: '4px 0 24px rgba(0,0,0,0.02)',
                zIndex: 10,
                transition: 'width 0.3s ease',
            }}>
                {/* Header / Toggle */}
                <div style={{
                    marginBottom: '3rem',
                    paddingLeft: isCollapsed ? 0 : '0.5rem',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: isCollapsed ? 'center' : 'space-between',
                    gap: '0.75rem'
                }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                        <div style={{ width: '32px', height: '32px', background: 'var(--color-primary)', borderRadius: '8px', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'white', flexShrink: 0 }}>
                            <LayoutDashboard size={18} />
                        </div>
                        {!isCollapsed && (
                            <h1 style={{ fontSize: '1.5rem', fontWeight: 800, color: 'var(--color-text-main)', margin: 0, letterSpacing: '-0.5px', whiteSpace: 'nowrap' }}>
                                Hex<span style={{ color: 'var(--color-primary)' }}>Admin</span>
                            </h1>
                        )}
                    </div>
                </div>

                {/* Toggle Button - floated or inline? Let's put it top right or bottom */}
                <button
                    onClick={() => setIsCollapsed(!isCollapsed)}
                    style={{
                        position: 'absolute',
                        top: '50px',
                        left: isCollapsed ? '65px' : '265px',
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
                        boxShadow: '0 2px 5px rgba(0,0,0,0.05)',
                        zIndex: 20,
                        transition: 'left 0.3s ease'
                    }}
                >
                    {isCollapsed ? <ChevronRight size={14} /> : <ChevronLeft size={14} />}
                </button>


                <nav style={{ flex: 1 }}>
                    {!isCollapsed && (
                        <div style={{ fontSize: '0.75rem', fontWeight: 700, textTransform: 'uppercase', color: 'var(--color-text-muted)', marginBottom: '1rem', paddingLeft: '1rem', letterSpacing: '0.05em' }}>
                            Menu
                        </div>
                    )}
                    <SidebarLink to="/dashboard" icon={LayoutDashboard} isCollapsed={isCollapsed}>Overview</SidebarLink>

                    <div style={{ marginTop: '1rem' }}>
                        {!isCollapsed && (
                            <div style={{ fontSize: '0.75rem', fontWeight: 700, textTransform: 'uppercase', color: 'var(--color-text-muted)', marginBottom: '0.5rem', paddingLeft: '1rem', letterSpacing: '0.05em' }}>
                                Settings
                            </div>
                        )}
                        <SidebarSubmenu icon={SettingsIcon} label="Configs" isCollapsed={isCollapsed}>
                            <SidebarLink to="/users" icon={Users} isCollapsed={isCollapsed}>Users</SidebarLink>
                            <SidebarLink to="/roles" icon={Shield} isCollapsed={isCollapsed}>Roles & Permissions</SidebarLink>
                        </SidebarSubmenu>
                    </div>
                </nav>

                <div style={{ borderTop: '1px solid var(--color-border)', paddingTop: '1.5rem', marginTop: '1.5rem' }}>
                    <div style={{
                        marginBottom: '1rem',
                        padding: isCollapsed ? '0.5rem' : '1rem',
                        backgroundColor: 'var(--color-background)',
                        borderRadius: 'var(--radius-md)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: isCollapsed ? 'center' : 'flex-start',
                        gap: '0.75rem'
                    }}>
                        <div style={{ width: '36px', height: '36px', borderRadius: '50%', background: 'var(--color-primary-light)', color: 'var(--color-primary)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 'bold', flexShrink: 0 }}>
                            {user?.name?.[0] || 'U'}
                        </div>
                        {!isCollapsed && (
                            <div style={{ overflow: 'hidden' }}>
                                <p style={{ margin: 0, fontWeight: 600, color: 'var(--color-text-main)', fontSize: '0.875rem' }}>{user?.name || 'User'}</p>
                                <p style={{ margin: 0, fontSize: '0.75rem', color: 'var(--color-text-secondary)', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>@{user?.corp_id || user?.username}</p>
                            </div>
                        )}
                    </div>
                    <button
                        onClick={handleLogout}
                        style={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: isCollapsed ? 'center' : 'flex-start',
                            gap: '0.75rem',
                            width: '100%',
                            padding: '0.875rem 1rem',
                            background: 'transparent',
                            color: 'var(--color-error)',
                            borderRadius: 'var(--radius-md)',
                            transition: 'all 0.2s',
                            cursor: 'pointer',
                            fontWeight: 500,
                            border: '1px solid transparent'
                        }}
                        onMouseOver={(e) => {
                            e.currentTarget.style.backgroundColor = '#FEF2F2';
                            e.currentTarget.style.borderColor = '#FECACA';
                        }}
                        onMouseOut={(e) => {
                            e.currentTarget.style.backgroundColor = 'transparent';
                            e.currentTarget.style.borderColor = 'transparent';
                        }}
                        title={isCollapsed ? "Sign Out" : ""}
                    >
                        <LogOut size={20} />
                        {!isCollapsed && <span>Sign Out</span>}
                    </button>
                </div>
            </aside>

            {/* Main Content */}
            <main style={{ flex: 1, padding: '2rem', overflowY: 'auto' }}>
                <Outlet />
            </main>
        </div>
    );
};

export default DashboardLayout;
