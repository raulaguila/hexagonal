import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { ChevronRight, Home } from 'lucide-react';
import { usePreferences } from '../../context/PreferencesContext';

/**
 * Breadcrumbs navigation component
 * Displays the current navigation path with clickable links
 */
const Breadcrumbs = ({ customLabels = {} }) => {
    const location = useLocation();
    const { t } = usePreferences();

    // Get path segments
    const pathnames = location.pathname.split('/').filter(x => x);

    // If at root or dashboard, don't show breadcrumbs
    if (pathnames.length <= 1 && pathnames[0] === 'dashboard') {
        return null;
    }

    // Default labels for routes
    const defaultLabels = {
        dashboard: t('sidebar.dashboard') || 'Dashboard',
        users: t('sidebar.users') || 'Users',
        roles: t('sidebar.roles') || 'Roles',
        settings: 'Settings',
        ...customLabels
    };

    const getLabel = (segment) => {
        return defaultLabels[segment] || segment.charAt(0).toUpperCase() + segment.slice(1);
    };

    return (
        <nav
            aria-label="Breadcrumb"
            style={{
                display: 'flex',
                alignItems: 'center',
                gap: '0.5rem',
                fontSize: '0.875rem',
                marginBottom: '1rem',
                color: 'var(--color-text-secondary)'
            }}
        >
            {/* Home link */}
            <Link
                to="/dashboard"
                style={{
                    display: 'flex',
                    alignItems: 'center',
                    color: 'var(--color-text-muted)',
                    transition: 'color 0.15s'
                }}
                onMouseOver={(e) => e.currentTarget.style.color = 'var(--color-primary)'}
                onMouseOut={(e) => e.currentTarget.style.color = 'var(--color-text-muted)'}
            >
                <Home size={16} />
            </Link>

            {/* Path segments */}
            {pathnames.map((segment, index) => {
                const routeTo = `/${pathnames.slice(0, index + 1).join('/')}`;
                const isLast = index === pathnames.length - 1;

                return (
                    <React.Fragment key={routeTo}>
                        <ChevronRight size={14} style={{ color: 'var(--color-text-muted)' }} />
                        {isLast ? (
                            <span style={{
                                color: 'var(--color-text-main)',
                                fontWeight: 500
                            }}>
                                {getLabel(segment)}
                            </span>
                        ) : (
                            <Link
                                to={routeTo}
                                style={{
                                    color: 'var(--color-text-secondary)',
                                    transition: 'color 0.15s'
                                }}
                                onMouseOver={(e) => e.currentTarget.style.color = 'var(--color-primary)'}
                                onMouseOut={(e) => e.currentTarget.style.color = 'var(--color-text-secondary)'}
                            >
                                {getLabel(segment)}
                            </Link>
                        )}
                    </React.Fragment>
                );
            })}
        </nav>
    );
};

export default Breadcrumbs;
