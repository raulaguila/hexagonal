import React, { useEffect, useState, useCallback } from 'react';
import { Plus, Shield, Pen, Trash2, Search } from 'lucide-react';
import { useToast } from '../components/feedback/ToastProvider';
import { useRoles, usePermissions } from '../hooks';
import { roleService } from '../services/roleService';
import Button from '../components/common/Button';
import Modal from '../components/common/Modal';
import Badge from '../components/common/Badge';
import EmptyState from '../components/common/EmptyState';
import { ConfirmDialog, SkeletonTableRow } from '../components/feedback';
import { Table, Thead, Tbody, Tr, Th, Td } from '../components/common/Table';
import Pagination from '../components/common/Pagination';
import { usePreferences } from '../context/PreferencesContext';

// Styled search input component
const SearchInput = ({ value, onChange, onKeyDown, placeholder }) => (
    <input
        type="text"
        value={value}
        onChange={onChange}
        onKeyDown={onKeyDown}
        placeholder={placeholder}
        style={{
            flex: 1,
            padding: '0.625rem 1rem',
            backgroundColor: 'var(--color-surface)',
            border: '1px solid var(--color-border)',
            borderRadius: 'var(--radius-md)',
            fontSize: '0.875rem',
            color: 'var(--color-text-main)',
            outline: 'none',
            transition: 'border-color 0.2s, box-shadow 0.2s'
        }}
        onFocus={(e) => {
            e.target.style.borderColor = 'var(--color-primary)';
            e.target.style.boxShadow = '0 0 0 3px rgba(99, 102, 241, 0.1)';
        }}
        onBlur={(e) => {
            e.target.style.borderColor = 'var(--color-border)';
            e.target.style.boxShadow = 'none';
        }}
    />
);

// Permission modules and actions
const PERMISSION_MODULES = [
    { key: 'users', label: 'Users' },
    { key: 'roles', label: 'Roles' }
];

const PERMISSION_ACTIONS = [
    { key: 'view', label: 'View' },
    { key: 'create', label: 'Create' },
    { key: 'edit', label: 'Edit' },
    { key: 'delete', label: 'Delete' }
];

const RolesPage = () => {
    const toast = useToast();
    const { t } = usePreferences();
    const { hasPermission, isRoot } = usePermissions();

    // Data hooks
    const { roles, loading, fetchRoles, deleteRole: removeRole } = useRoles();

    // Local state
    const [searchTerm, setSearchTerm] = useState('');
    const [appliedSearch, setAppliedSearch] = useState('');
    const [currentPage, setCurrentPage] = useState(1);
    const [modalOpen, setModalOpen] = useState(false);
    const [editingId, setEditingId] = useState(null);
    const [saving, setSaving] = useState(false);
    const [formData, setFormData] = useState({
        name: '',
        permissions: [],
        enabled: true
    });
    const itemsPerPage = 10;

    // Confirm dialog state
    const [confirmOpen, setConfirmOpen] = useState(false);
    const [deleteTarget, setDeleteTarget] = useState(null);
    const [deleting, setDeleting] = useState(false);

    // Fetch data on mount
    useEffect(() => {
        fetchRoles();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    // Filter roles based on applied search only
    const filteredRoles = React.useMemo(() => {
        if (!appliedSearch.trim()) {
            return roles;
        }
        return roles.filter(role =>
            role.name.toLowerCase().includes(appliedSearch.toLowerCase())
        );
    }, [roles, appliedSearch]);

    // Calculate pagination
    const totalPages = Math.ceil(filteredRoles.length / itemsPerPage);
    const paginatedRoles = filteredRoles.slice(
        (currentPage - 1) * itemsPerPage,
        currentPage * itemsPerPage
    );

    // Search handler - only on button click
    const handleSearch = useCallback(() => {
        setAppliedSearch(searchTerm);
        setCurrentPage(1);
    }, [searchTerm]);

    // Handle Enter key in search
    const handleSearchKeyDown = (e) => {
        if (e.key === 'Enter') {
            handleSearch();
        }
    };

    // Toggle permission
    const togglePermission = (permission) => {
        setFormData(prev => ({
            ...prev,
            permissions: prev.permissions.includes(permission)
                ? prev.permissions.filter(p => p !== permission)
                : [...prev.permissions, permission]
        }));
    };

    // Toggle all permissions for a module
    const toggleModule = (moduleKey) => {
        const modulePermissions = PERMISSION_ACTIONS.map(a => `${moduleKey}:${a.key}`);
        const allChecked = modulePermissions.every(p => formData.permissions.includes(p));

        setFormData(prev => ({
            ...prev,
            permissions: allChecked
                ? prev.permissions.filter(p => !modulePermissions.includes(p))
                : [...new Set([...prev.permissions, ...modulePermissions])]
        }));
    };

    // Open edit modal
    const handleEdit = (role) => {
        setFormData({
            name: role.name,
            permissions: role.permissions || [],
            enabled: role.enabled !== false
        });
        setEditingId(role.id);
        setModalOpen(true);
    };

    // Open create modal
    const handleCreate = () => {
        setFormData({ name: '', permissions: [], enabled: true });
        setEditingId(null);
        setModalOpen(true);
    };

    // Submit form
    const handleSubmit = async (e) => {
        e.preventDefault();
        setSaving(true);

        try {
            if (editingId) {
                await roleService.updateRole(editingId, formData);
                toast.success('Role updated successfully');
            } else {
                await roleService.createRole(formData);
                toast.success('Role created successfully');
            }

            setModalOpen(false);
            fetchRoles();
        } catch (error) {
            toast.error(error.response?.data?.message || 'Failed to save role');
        } finally {
            setSaving(false);
        }
    };

    // Delete confirmation
    const handleDeleteClick = (role) => {
        setDeleteTarget(role);
        setConfirmOpen(true);
    };

    // Confirm delete
    const handleConfirmDelete = async () => {
        if (!deleteTarget) return;

        setDeleting(true);
        try {
            await removeRole(deleteTarget.id);
            toast.success('Role deleted successfully');
            setConfirmOpen(false);
            setDeleteTarget(null);
        } catch {
            toast.error('Failed to delete role');
        } finally {
            setDeleting(false);
        }
    };

    // Check permissions
    const canCreate = hasPermission('roles:create') || isRoot();
    const canEdit = hasPermission('roles:edit') || isRoot();
    const canDelete = hasPermission('roles:delete') || isRoot();

    return (
        <div>
            {/* Header */}
            <div style={{ marginBottom: '1.5rem' }}>
                <h1 style={{
                    fontSize: '1.5rem',
                    fontWeight: 700,
                    margin: 0,
                    color: 'var(--color-text-main)'
                }}>
                    {t('roles.title') || 'Roles & Permissions'}
                </h1>
                <p style={{
                    color: 'var(--color-text-secondary)',
                    marginTop: '0.25rem',
                    fontSize: '0.875rem'
                }}>
                    {t('roles.subtitle') || 'Manage access control and permissions'}
                </p>
            </div>

            {/* Search and Actions Row */}
            <div style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                marginBottom: '1rem',
                gap: '1rem'
            }}>
                <div style={{ display: 'flex', gap: '0.5rem', flex: 1, maxWidth: '400px' }}>
                    <SearchInput
                        placeholder={t('roles.search_placeholder') || 'Search roles...'}
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        onKeyDown={handleSearchKeyDown}
                    />
                    <Button variant="secondary" onClick={handleSearch}>
                        <Search size={18} />
                    </Button>
                </div>

                {canCreate && (
                    <Button onClick={handleCreate} variant="primary">
                        <Plus size={18} />
                        {t('roles.add') || 'Add Role'}
                    </Button>
                )}
            </div>

            {/* Table */}
            <div style={{
                backgroundColor: 'var(--color-surface)',
                borderRadius: 'var(--radius-lg)',
                border: '1px solid var(--color-border)',
                overflow: 'hidden'
            }}>
                <Table>
                    <Thead>
                        <Tr>
                            <Th>{t('roles.table.name') || 'Role Name'}</Th>
                            <Th style={{ textAlign: 'center' }}>{t('roles.table.status') || 'Status'}</Th>
                            <Th>{t('roles.table.permissions') || 'Permissions'}</Th>
                            <Th style={{ textAlign: 'center' }}>{t('roles.table.actions') || 'Actions'}</Th>
                        </Tr>
                    </Thead>
                    <Tbody>
                        {loading ? (
                            Array.from({ length: 3 }).map((_, i) => (
                                <SkeletonTableRow key={i} columns={4} />
                            ))
                        ) : paginatedRoles.length === 0 ? (
                            <Tr>
                                <Td colSpan={4}>
                                    <EmptyState
                                        icon={Shield}
                                        title={t('roles.empty') || 'No roles found'}
                                        description={appliedSearch ? t('roles.empty_search') : 'Get started by creating a new role.'}
                                        actionLabel={t('roles.add') || 'Add Role'}
                                        onAction={canCreate ? handleCreate : null}
                                    />
                                </Td>
                            </Tr>
                        ) : (
                            paginatedRoles.map(role => (
                                <Tr key={role.id} className="table-row">
                                    {/* Role Name */}
                                    <Td>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                                            <div style={{
                                                width: '36px',
                                                height: '36px',
                                                borderRadius: 'var(--radius-md)',
                                                backgroundColor: 'var(--color-primary-light)',
                                                color: 'var(--color-primary)',
                                                display: 'flex',
                                                alignItems: 'center',
                                                justifyContent: 'center'
                                            }}>
                                                <Shield size={18} />
                                            </div>
                                            <span style={{ fontWeight: 500, color: 'var(--color-text-main)' }}>
                                                {role.name}
                                            </span>
                                        </div>
                                    </Td>

                                    {/* Status */}
                                    <Td style={{ textAlign: 'center' }}>
                                        <Badge
                                            variant={role.enabled !== false ? 'success' : 'error'}
                                            size="sm"
                                        >
                                            {role.enabled !== false ? (t('common.active') || 'Active') : (t('common.inactive') || 'Inactive')}
                                        </Badge>
                                    </Td>

                                    {/* Permissions Preview */}
                                    <Td>
                                        <div style={{
                                            display: 'flex',
                                            gap: '0.25rem',
                                            flexWrap: 'wrap',
                                            maxWidth: '400px'
                                        }}>
                                            {role.permissions?.slice(0, 4).map(p => (
                                                <Badge key={p} variant="default" size="sm">
                                                    {p}
                                                </Badge>
                                            ))}
                                            {role.permissions?.length > 4 && (
                                                <Badge variant="primary" size="sm">
                                                    +{role.permissions.length - 4} more
                                                </Badge>
                                            )}
                                            {(!role.permissions || role.permissions.length === 0) && (
                                                <span style={{
                                                    fontSize: '0.75rem',
                                                    color: 'var(--color-text-muted)',
                                                    fontStyle: 'italic'
                                                }}>
                                                    No permissions assigned
                                                </span>
                                            )}
                                        </div>
                                    </Td>

                                    {/* Actions - Centered */}
                                    <Td style={{ textAlign: 'center' }}>
                                        <div style={{ display: 'inline-flex', gap: '0.25rem', justifyContent: 'center' }}>
                                            {canEdit && (
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => handleEdit(role)}
                                                    style={{ padding: '0.375rem' }}
                                                >
                                                    <Pen size={16} />
                                                </Button>
                                            )}
                                            {canDelete && role.name !== 'ROOT' && (
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => handleDeleteClick(role)}
                                                    style={{ padding: '0.375rem', color: 'var(--color-error)' }}
                                                >
                                                    <Trash2 size={16} />
                                                </Button>
                                            )}
                                        </div>
                                    </Td>
                                </Tr>
                            ))
                        )}
                    </Tbody>
                </Table>

                {/* Pagination */}
                {!loading && filteredRoles.length > 0 && (
                    <Pagination
                        currentPage={currentPage}
                        totalPages={totalPages}
                        totalItems={filteredRoles.length}
                        itemsPerPage={itemsPerPage}
                        onPageChange={setCurrentPage}
                        t={t}
                    />
                )}
            </div>

            {/* Create/Edit Modal */}
            <Modal
                isOpen={modalOpen}
                onClose={() => setModalOpen(false)}
                title={editingId ? t('role.modal.edit_title') : t('role.modal.create_title')}
                size="md"
            >
                <form onSubmit={handleSubmit}>
                    <div>
                        <label style={{
                            display: 'block',
                            fontSize: '0.875rem',
                            fontWeight: 500,
                            marginBottom: '0.5rem',
                            color: 'var(--color-text-main)'
                        }}>
                            {t('role.form.name') || 'Role Name'}
                        </label>
                        <input
                            type="text"
                            value={formData.name}
                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                            required
                            placeholder="e.g. EDITOR"
                            style={{
                                width: '100%',
                                padding: '0.625rem 1rem',
                                backgroundColor: 'var(--color-background)',
                                border: '1px solid var(--color-border)',
                                borderRadius: 'var(--radius-md)',
                                fontSize: '0.875rem',
                                color: 'var(--color-text-main)',
                                outline: 'none',
                                boxSizing: 'border-box'
                            }}
                        />
                    </div>

                    {/* Enabled Status */}
                    <div style={{ marginTop: '1rem' }}>
                        <label style={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: '0.75rem',
                            cursor: 'pointer'
                        }}>
                            <div
                                onClick={() => setFormData({ ...formData, enabled: !formData.enabled })}
                                style={{
                                    width: '44px',
                                    height: '24px',
                                    borderRadius: '12px',
                                    backgroundColor: formData.enabled ? 'var(--color-success)' : 'var(--color-border)',
                                    position: 'relative',
                                    cursor: 'pointer',
                                    transition: 'background-color 0.2s'
                                }}
                            >
                                <div style={{
                                    width: '20px',
                                    height: '20px',
                                    borderRadius: '50%',
                                    backgroundColor: 'white',
                                    position: 'absolute',
                                    top: '2px',
                                    left: formData.enabled ? '22px' : '2px',
                                    transition: 'left 0.2s',
                                    boxShadow: '0 1px 3px rgba(0,0,0,0.2)'
                                }} />
                            </div>
                            <span style={{
                                fontSize: '0.875rem',
                                fontWeight: 500,
                                color: 'var(--color-text-main)'
                            }}>
                                {t('role.form.enabled') || 'Role Enabled'}
                            </span>
                            <span style={{
                                fontSize: '0.75rem',
                                color: 'var(--color-text-muted)'
                            }}>
                                {formData.enabled
                                    ? (t('role.form.enabled_hint') || 'Users with this role can access the system')
                                    : (t('role.form.disabled_hint') || 'Users with only this role will be denied access')}
                            </span>
                        </label>
                    </div>

                    {/* Permissions Matrix */}
                    <div style={{ marginTop: '1.5rem' }}>
                        <label style={{
                            display: 'block',
                            fontSize: '0.875rem',
                            fontWeight: 500,
                            marginBottom: '0.75rem',
                            color: 'var(--color-text-main)'
                        }}>
                            {t('role.form.permissions') || 'Permissions'}
                        </label>

                        <div style={{
                            border: '1px solid var(--color-border)',
                            borderRadius: 'var(--radius-md)',
                            overflow: 'hidden'
                        }}>
                            {/* Header */}
                            <div style={{
                                display: 'grid',
                                gridTemplateColumns: '1fr repeat(4, 80px)',
                                gap: '0.5rem',
                                padding: '0.75rem 1rem',
                                backgroundColor: 'var(--color-background)',
                                borderBottom: '1px solid var(--color-border)',
                                fontSize: '0.75rem',
                                fontWeight: 600,
                                color: 'var(--color-text-secondary)',
                                textTransform: 'uppercase'
                            }}>
                                <span>Module</span>
                                {PERMISSION_ACTIONS.map(action => (
                                    <span key={action.key} style={{ textAlign: 'center' }}>
                                        {action.label}
                                    </span>
                                ))}
                            </div>

                            {/* Rows */}
                            {PERMISSION_MODULES.map(module => {
                                const modulePerms = PERMISSION_ACTIONS.map(a => `${module.key}:${a.key}`);
                                const allChecked = modulePerms.every(p => formData.permissions.includes(p));

                                return (
                                    <div
                                        key={module.key}
                                        style={{
                                            display: 'grid',
                                            gridTemplateColumns: '1fr repeat(4, 80px)',
                                            gap: '0.5rem',
                                            padding: '0.75rem 1rem',
                                            borderBottom: '1px solid var(--color-border)',
                                            alignItems: 'center'
                                        }}
                                    >
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                            <input
                                                type="checkbox"
                                                checked={allChecked}
                                                onChange={() => toggleModule(module.key)}
                                                style={{ accentColor: 'var(--color-primary)' }}
                                            />
                                            <span style={{
                                                fontSize: '0.875rem',
                                                fontWeight: 500,
                                                color: 'var(--color-text-main)'
                                            }}>
                                                {module.label}
                                            </span>
                                        </div>

                                        {PERMISSION_ACTIONS.map(action => {
                                            const perm = `${module.key}:${action.key}`;
                                            return (
                                                <div key={perm} style={{ textAlign: 'center' }}>
                                                    <input
                                                        type="checkbox"
                                                        checked={formData.permissions.includes(perm)}
                                                        onChange={() => togglePermission(perm)}
                                                        style={{
                                                            accentColor: 'var(--color-primary)',
                                                            width: '16px',
                                                            height: '16px'
                                                        }}
                                                    />
                                                </div>
                                            );
                                        })}
                                    </div>
                                );
                            })}
                        </div>
                    </div>

                    {/* Actions */}
                    <div style={{
                        display: 'flex',
                        justifyContent: 'flex-end',
                        gap: '0.75rem',
                        marginTop: '1.5rem',
                        paddingTop: '1rem',
                        borderTop: '1px solid var(--color-border)'
                    }}>
                        <Button type="button" variant="secondary" onClick={() => setModalOpen(false)}>
                            {t('user.form.cancel') || 'Cancel'}
                        </Button>
                        <Button type="submit" variant="primary" loading={saving}>
                            {editingId ? (t('user.form.update') || 'Save Changes') : (t('role.form.save') || 'Create Role')}
                        </Button>
                    </div>
                </form>
            </Modal>

            {/* Delete Confirmation */}
            <ConfirmDialog
                isOpen={confirmOpen}
                onClose={() => {
                    setConfirmOpen(false);
                    setDeleteTarget(null);
                }}
                onConfirm={handleConfirmDelete}
                title="Delete Role"
                message={t('roles.delete_confirm') || `Are you sure you want to delete "${deleteTarget?.name}"?`}
                confirmText="Delete"
                variant="danger"
                loading={deleting}
            />
        </div>
    );
};

export default RolesPage;
