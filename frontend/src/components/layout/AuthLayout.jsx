import React from 'react';
import { Outlet } from 'react-router-dom';

const AuthLayout = () => {
    return (
        <div style={{
            display: 'flex',
            minHeight: '100vh',
            alignItems: 'center',
            justifyContent: 'center',
            backgroundColor: 'var(--color-background)',
            padding: '1rem'
        }}>
            <div style={{
                width: '100%',
                maxWidth: '400px'
            }}>
                <Outlet />
            </div>
        </div>
    );
};

export default AuthLayout;
