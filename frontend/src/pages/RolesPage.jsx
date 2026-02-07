import React, { useEffect, useState, useCallback } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
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
import { roleSchema } from '../utils/schemas';
import styles from './RolesPage.module.css';

// Styled search input component
const SearchInput = ({ value, onChange, onKeyDown, placeholder }) => (
    <input
        type="text"
        value={value}
        onChange={onChange}
        onKeyDown={onKeyDown}
        placeholder={placeholder}
        className={styles.searchInput}
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
    const [currentPage, setCurrentPage] = useState(1);
    const [sortKey, setSortKey] = useState('name');
    const [sortOrder, setSortOrder] = useState('asc');
    const [modalOpen, setModalOpen] = useState(false);
    const [editingId, setEditingId] = useState(null);
    const [saving, setSaving] = useState(false);
    const itemsPerPage = 10;

    // React Hook Form
    const {
        register,
        handleSubmit: formSubmit,
        reset,
        watch,
        setValue,
        formState: { errors },
    } = useForm({
        resolver: zodResolver(roleSchema),
        defaultValues: {
            name: '',
            permissions: [],
            enabled: true
        },
    });

    const watchPermissions = watch('permissions', []);
    const watchEnabled = watch('enabled', true);

    // Confirm dialog state
    const [confirmOpen, setConfirmOpen] = useState(false);
    const [deleteTarget, setDeleteTarget] = useState(null);
    const [deleting, setDeleting] = useState(false);

    // Fetch data function
    // Note: Backend expects sort=columnName, order=asc|desc
    const loadRoles = useCallback((params = {}) => {
        fetchRoles({
            page: params.page ?? currentPage,
            limit: itemsPerPage,
            search: params.search ?? searchTerm,
            sort: params.sortKey ?? sortKey,
            order: params.sortOrder ?? sortOrder
        });
    }, [fetchRoles, currentPage, searchTerm, sortKey, sortOrder]);

    // Fetch data on mount
    useEffect(() => {
        loadRoles({ page: 1, search: '' });
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    // Calculate pagination from data
    const totalPages = Math.ceil(roles.length / itemsPerPage);
    const paginatedRoles = roles.slice(
        (currentPage - 1) * itemsPerPage,
        currentPage * itemsPerPage
    );

    // Search handler - sends to backend
    const handleSearch = useCallback(() => {
        setCurrentPage(1);
        loadRoles({ page: 1, search: searchTerm });
    }, [searchTerm, loadRoles]);

    // Handle Enter key in search
    const handleSearchKeyDown = (e) => {
        if (e.key === 'Enter') {
            handleSearch();
        }
    };

    // Sort handler - sends to backend
    const handleSort = useCallback((key, order) => {
        setSortKey(key);
        setSortOrder(order);
        loadRoles({ sortKey: key, sortOrder: order });
    }, [loadRoles]);

    // Page change handler
    const handlePageChange = (page) => {
        setCurrentPage(page);
        loadRoles({ page });
    };

    // Toggle permission
    const togglePermission = (permission) => {
        const newPermissions = watchPermissions.includes(permission)
            ? watchPermissions.filter(p => p !== permission)
            : [...watchPermissions, permission];
        setValue('permissions', newPermissions);
    };

    // Toggle all permissions for a module
    const toggleModule = (moduleKey) => {
        const modulePermissions = PERMISSION_ACTIONS.map(a => `${moduleKey}:${a.key}`);
        const allChecked = modulePermissions.every(p => watchPermissions.includes(p));

        const newPermissions = allChecked
            ? watchPermissions.filter(p => !modulePermissions.includes(p))
            : [...new Set([...watchPermissions, ...modulePermissions])];
        setValue('permissions', newPermissions);
    };

    // Open edit modal
    const handleEdit = (role) => {
        reset({
            name: role.name,
            permissions: role.permissions || [],
            enabled: role.enabled !== false
        });
        setEditingId(role.id);
        setModalOpen(true);
    };

    // Open create modal
    const handleCreate = () => {
        reset({ name: '', permissions: [], enabled: true });
        setEditingId(null);
        setModalOpen(true);
    };

    // Submit form
    const onSubmit = async (data) => {
        setSaving(true);

        try {
            const payload = {
                name: data.name,
                permissions: data.permissions,
                enabled: data.enabled
            };

            if (editingId) {
                await roleService.updateRole(editingId, payload);
                toast.success(t('roles.updated') || 'Role updated successfully');
            } else {
                await roleService.createRole(payload);
                toast.success(t('roles.created') || 'Role created successfully');
            }

            setModalOpen(false);
            loadRoles();
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
        <div className={styles.pageContainer}>
            {/* Header */}
            <div className={styles.header}>
                <h1 className={styles.title}>
                    {t('roles.title') || 'Roles & Permissions'}
                </h1>
                <p className={styles.subtitle}>
                    {t('roles.subtitle') || 'Manage access control and permissions'}
                </p>
            </div>

            {/* Search and Actions Row */}
            <div className={styles.actionsRow}>
                <div className={styles.searchContainer}>
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
                        <span className={styles.btnText}>{t('roles.add') || 'Add Role'}</span>
                    </Button>
                )}
            </div>

            {/* Table */}
            <div className={styles.tableContainer}>
                <Table>
                    <Thead>
                        <Tr>
                            <Th
                                sortKey="name"
                                currentSort={sortKey}
                                sortOrder={sortOrder}
                                onSort={handleSort}
                            >
                                {t('roles.table.name') || 'Role Name'}
                            </Th>
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
                                <Td colSpan={4} className={styles.emptyStateWrapper}>
                                    <EmptyState
                                        icon={Shield}
                                        title={t('roles.empty') || 'No roles found'}
                                        description={searchTerm ? t('roles.empty_search') : 'Get started by creating a new role.'}
                                        actionLabel={t('roles.add') || 'Add Role'}
                                        onAction={canCreate ? handleCreate : null}
                                    />
                                </Td>
                            </Tr>
                        ) : (
                            paginatedRoles.map(role => (
                                <Tr key={role.id}>
                                    {/* Role Name */}
                                    <Td>
                                        <div className={styles.roleInfo}>
                                            <div className={styles.roleIcon}>
                                                <Shield size={18} />
                                            </div>
                                            <span className={styles.roleName}>
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
                                        <div className={styles.permissionsWrapper}>
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
                                                <span className={styles.noPermissions}>
                                                    No permissions assigned
                                                </span>
                                            )}
                                        </div>
                                    </Td>

                                    {/* Actions - Centered */}
                                    <Td className={styles.actionsCell}>
                                        <div className={styles.actionsWrapper}>
                                            {canEdit && (
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => handleEdit(role)}
                                                    className={styles.actionBtn}
                                                >
                                                    <Pen size={16} />
                                                </Button>
                                            )}
                                            {canDelete && role.name !== 'ROOT' && (
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => handleDeleteClick(role)}
                                                    className={`${styles.actionBtn} ${styles.deleteBtn}`}
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
                {!loading && roles.length > 0 && (
                    <Pagination
                        currentPage={currentPage}
                        totalPages={totalPages}
                        totalItems={roles.length}
                        itemsPerPage={itemsPerPage}
                        onPageChange={handlePageChange}
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
                <form onSubmit={formSubmit(onSubmit)}>
                    <div>
                        <label className={styles.label}>
                            {t('role.form.name') || 'Role Name'}
                        </label>
                        <input
                            type="text"
                            {...register('name')}
                            placeholder="e.g. EDITOR"
                            className={`${styles.input} ${errors.name ? styles.inputError : ''}`}
                        />
                        {errors.name && (
                            <span className={styles.errorMsg}>
                                {errors.name.message}
                            </span>
                        )}
                    </div>

                    {/* Enabled Status */}
                    <div style={{ marginTop: '1rem' }}>
                        <label className={styles.enabledToggle}>
                            <div
                                onClick={() => setValue('enabled', !watchEnabled)}
                                className={`${styles.toggleSwitch} ${watchEnabled ? styles.toggleSwitchActive : styles.toggleSwitchInactive}`}
                            >
                                <div className={`${styles.toggleKnob} ${watchEnabled ? styles.toggleKnobActive : styles.toggleKnobInactive}`} />
                            </div>
                            <span className={styles.toggleLabel}>
                                {t('role.form.enabled') || 'Role Enabled'}
                            </span>
                            <span className={styles.toggleHint}>
                                {watchEnabled
                                    ? (t('role.form.enabled_hint') || 'Users with this role can access the system')
                                    : (t('role.form.disabled_hint') || 'Users with only this role will be denied access')}
                            </span>
                        </label>
                    </div>

                    {/* Permissions Matrix */}
                    <div className={styles.permissionsSection}>
                        <label className={styles.permissionsLabel}>
                            {t('role.form.permissions') || 'Permissions'}
                        </label>

                        <div className={styles.matrixContainer}>
                            {/* Header */}
                            <div className={styles.matrixHeader}>
                                <span>Module</span>
                                {PERMISSION_ACTIONS.map(action => (
                                    <span key={action.key} className={styles.headerAction}>
                                        {action.label}
                                    </span>
                                ))}
                            </div>

                            {/* Rows */}
                            {PERMISSION_MODULES.map(module => {
                                const modulePerms = PERMISSION_ACTIONS.map(a => `${module.key}:${a.key}`);
                                const allChecked = modulePerms.every(p => watchPermissions.includes(p));

                                return (
                                    <div
                                        key={module.key}
                                        className={styles.matrixRow}
                                    >
                                        <div className={styles.moduleCell}>
                                            <input
                                                type="checkbox"
                                                checked={allChecked}
                                                onChange={() => toggleModule(module.key)}
                                                className={styles.moduleCheckbox}
                                            />
                                            <span className={styles.moduleName}>
                                                {module.label}
                                            </span>
                                        </div>

                                        {PERMISSION_ACTIONS.map(action => {
                                            const perm = `${module.key}:${action.key}`;
                                            return (
                                                <div key={perm} className={styles.permissionCell}>
                                                    <input
                                                        type="checkbox"
                                                        checked={watchPermissions.includes(perm)}
                                                        onChange={() => togglePermission(perm)}
                                                        className={styles.permissionCheckbox}
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
                    <div className={styles.modalActions}>
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
