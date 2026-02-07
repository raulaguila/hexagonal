import React, { useEffect, useState, useCallback } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Plus, Search, Pen, Trash2, User, Mail } from 'lucide-react';
import { useToast } from '../components/feedback/ToastProvider';
import { useUsers, useRoles, usePermissions } from '../hooks';
import { userService } from '../services/userService';
import Button from '../components/common/Button';
import Modal from '../components/common/Modal';
import Badge from '../components/common/Badge';
import EmptyState from '../components/common/EmptyState';
import { ConfirmDialog, SkeletonTableRow } from '../components/feedback';
import { Table, Thead, Tbody, Tr, Th, Td } from '../components/common/Table';
import Pagination from '../components/common/Pagination';
import { usePreferences } from '../context/PreferencesContext';
import { userSchema } from '../utils/schemas';
import styles from './UsersPage.module.css';

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

const UsersPage = () => {
    const toast = useToast();
    const { t } = usePreferences();
    const { hasPermission, isRoot } = usePermissions();

    // Data hooks
    const { users, loading, pagination, fetchUsers, deleteUser: removeUser } = useUsers();
    const { roles, fetchRoles } = useRoles();

    // Local state
    const [search, setSearch] = useState('');
    const [sortKey, setSortKey] = useState('name');
    const [sortOrder, setSortOrder] = useState('asc');
    const [modalOpen, setModalOpen] = useState(false);
    const [editingId, setEditingId] = useState(null);
    const [saving, setSaving] = useState(false);

    // React Hook Form
    const {
        register,
        handleSubmit: formSubmit,
        reset,
        watch,
        setValue,
        formState: { errors },
    } = useForm({
        resolver: zodResolver(userSchema),
        defaultValues: {
            name: '',
            username: '',
            email: '',
            role_ids: [],
            status: true
        },
    });

    const watchRoleIds = watch('role_ids', []);
    const watchStatus = watch('status', true);

    // Confirm dialog state
    const [confirmOpen, setConfirmOpen] = useState(false);
    const [deleteTarget, setDeleteTarget] = useState(null);
    const [deleting, setDeleting] = useState(false);

    // Note: Backend expects sort=columnName, order=asc|desc
    const loadUsers = useCallback((params = {}) => {
        fetchUsers({
            page: params.page ?? pagination.page ?? 1,
            limit: pagination.limit || 10,
            search: params.search ?? search,
            sort: params.sortKey ?? sortKey,
            order: params.sortOrder ?? sortOrder
        });
    }, [fetchUsers, pagination, search, sortKey, sortOrder]);

    // Fetch data on mount
    useEffect(() => {
        fetchUsers({ page: 1, limit: 10, search: '' });
        fetchRoles();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    // Search handler - sends to backend
    const handleSearch = useCallback(() => {
        loadUsers({ page: 1, search: search });
    }, [search, loadUsers]);

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
        loadUsers({ sortKey: key, sortOrder: order });
    }, [loadUsers]);

    // Page change handler
    const handlePageChange = (page) => {
        loadUsers({ page });
    };

    // Open edit modal
    const handleEdit = (user) => {
        reset({
            name: user.name,
            username: user.username,
            email: user.email,
            role_ids: user.roles ? user.roles.map(r => r.id) : [],
            status: user.status !== undefined ? user.status : true
        });
        setEditingId(user.id);
        setModalOpen(true);
    };

    // Open create modal
    const handleCreate = () => {
        reset({
            name: '',
            username: '',
            email: '',
            role_ids: [],
            status: true
        });
        setEditingId(null);
        setModalOpen(true);
    };

    // Submit form
    const onSubmit = async (data) => {
        setSaving(true);

        try {
            const payload = {
                name: data.name,
                username: data.username,
                email: data.email,
                role_ids: data.role_ids,
                status: data.status
            };

            if (editingId) {
                await userService.updateUser(editingId, payload);
                toast.success(t('users.updated') || 'User updated successfully');
            } else {
                await userService.createUser(payload);
                toast.success(t('users.created') || 'User created successfully');
            }

            setModalOpen(false);
            loadUsers({ page: pagination.page || 1 });
        } catch (error) {
            toast.error(error.response?.data?.message || 'Failed to save user');
        } finally {
            setSaving(false);
        }
    };

    // Delete confirmation
    const handleDeleteClick = (user) => {
        setDeleteTarget(user);
        setConfirmOpen(true);
    };

    // Confirm delete
    const handleConfirmDelete = async () => {
        if (!deleteTarget) return;

        setDeleting(true);
        try {
            await removeUser(deleteTarget.id);
            toast.success(t('users.deleted') || 'User deleted successfully');
            setConfirmOpen(false);
            setDeleteTarget(null);
        } catch {
            toast.error('Failed to delete user');
        } finally {
            setDeleting(false);
        }
    };

    // Check permissions
    const canCreate = hasPermission('users:create') || isRoot();
    const canEdit = hasPermission('users:edit') || isRoot();
    const canDelete = hasPermission('users:delete') || isRoot();

    return (
        <div className={styles.pageContainer}>
            {/* Header */}
            <div className={styles.header}>
                <h1 className={styles.title}>
                    {t('users.title') || 'Users'}
                </h1>
                <p className={styles.subtitle}>
                    {t('users.subtitle') || 'Manage system access and profiles'}
                </p>
            </div>

            {/* Search and Actions Row */}
            <div className={styles.actionsRow}>
                <div className={styles.searchContainer}>
                    <SearchInput
                        placeholder={t('users.search_placeholder') || 'Search users...'}
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        onKeyDown={handleSearchKeyDown}
                    />
                    <Button variant="secondary" onClick={handleSearch}>
                        <Search size={18} />
                    </Button>
                </div>

                {canCreate && (
                    <Button onClick={handleCreate} variant="primary">
                        <Plus size={18} />
                        <span className={styles.btnText}>{t('users.add') || 'Add User'}</span>
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
                                {t('users.table.user') || 'User'}
                            </Th>
                            <Th
                                sortKey="email"
                                currentSort={sortKey}
                                sortOrder={sortOrder}
                                onSort={handleSort}
                            >
                                {t('users.table.contact') || 'Contact'}
                            </Th>
                            <Th>{t('users.table.roles') || 'Roles'}</Th>
                            <Th style={{ textAlign: 'center' }}>{t('users.table.status') || 'Status'}</Th>
                            <Th style={{ textAlign: 'center' }}>{t('users.table.actions') || 'Actions'}</Th>
                        </Tr>
                    </Thead>
                    <Tbody>
                        {loading ? (
                            Array.from({ length: 5 }).map((_, i) => (
                                <SkeletonTableRow key={i} columns={5} />
                            ))
                        ) : users.length === 0 ? (
                            <Tr>
                                <Td colSpan={5} className={styles.emptyStateWrapper}>
                                    <EmptyState
                                        icon={User}
                                        title={t('users.empty') || 'No users found'}
                                        description={search ? t('users.empty_search') : 'Get started by creating a new user.'}
                                        actionLabel={t('users.add') || 'Add User'}
                                        onAction={canCreate ? handleCreate : null}
                                    />
                                </Td>
                            </Tr>
                        ) : (
                            users.map(user => (
                                <Tr key={user.id}>
                                    {/* User Info */}
                                    <Td>
                                        <div className={styles.userInfo}>
                                            <div className={styles.avatar}>
                                                {user.name?.charAt(0) || '?'}
                                            </div>
                                            <div className={styles.userDetails}>
                                                <div className={styles.userName}>
                                                    {user.name}
                                                </div>
                                                <div className={styles.userUsername}>
                                                    @{user.username}
                                                </div>
                                            </div>
                                        </div>
                                    </Td>

                                    {/* Contact */}
                                    <Td>
                                        <div className={styles.contactInfo}>
                                            <Mail size={14} className={styles.contactIcon} />
                                            <span className={styles.contactEmail}>
                                                {user.email}
                                            </span>
                                        </div>
                                    </Td>

                                    {/* Roles */}
                                    <Td>
                                        <div className={styles.rolesWrapper}>
                                            {user.roles?.slice(0, 2).map(role => (
                                                <Badge key={role.id} variant="default" size="sm">
                                                    {role.name}
                                                </Badge>
                                            ))}
                                            {user.roles?.length > 2 && (
                                                <Badge variant="default" size="sm">
                                                    +{user.roles.length - 2}
                                                </Badge>
                                            )}
                                        </div>
                                    </Td>

                                    {/* Status */}
                                    <Td style={{ textAlign: 'center' }}>
                                        {user.status ? (
                                            <Badge variant="success" size="sm" dot>
                                                {t('common.active') || 'Active'}
                                            </Badge>
                                        ) : (
                                            <Badge variant="default" size="sm" dot>
                                                {t('common.inactive') || 'Inactive'}
                                            </Badge>
                                        )}
                                    </Td>

                                    {/* Actions - Centered */}
                                    <Td className={styles.actionsCell}>
                                        <div className={styles.actionsWrapper}>
                                            {canEdit && (
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => handleEdit(user)}
                                                    className={styles.actionBtn}
                                                >
                                                    <Pen size={16} />
                                                </Button>
                                            )}
                                            {canDelete && (
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => handleDeleteClick(user)}
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
                {!loading && users.length > 0 && (
                    <Pagination
                        currentPage={pagination.page || 1}
                        totalPages={pagination.totalPages || 1}
                        totalItems={pagination.total || users.length}
                        itemsPerPage={pagination.limit || 10}
                        onPageChange={handlePageChange}
                        t={t}
                    />
                )}
            </div>

            {/* Create/Edit Modal */}
            <Modal
                isOpen={modalOpen}
                onClose={() => setModalOpen(false)}
                title={editingId ? t('user.modal.edit_title') : t('user.modal.create_title')}
                size="md"
            >
                <form onSubmit={formSubmit(onSubmit)}>
                    <div className={styles.formGrid}>
                        <div className={styles.fullWidth}>
                            <label className={styles.label}>
                                {t('user.form.name') || 'Full Name'}
                            </label>
                            <input
                                type="text"
                                {...register('name')}
                                placeholder="e.g. John Doe"
                                className={`${styles.input} ${errors.name ? styles.inputError : ''}`}
                            />
                            {errors.name && (
                                <span className={styles.errorMsg}>
                                    {errors.name.message}
                                </span>
                            )}
                        </div>
                        <div>
                            <label className={styles.label}>
                                {t('user.form.username') || 'Username'}
                            </label>
                            <input
                                type="text"
                                {...register('username')}
                                placeholder="e.g. jdoe"
                                className={`${styles.input} ${errors.username ? styles.inputError : ''}`}
                            />
                            {errors.username && (
                                <span className={styles.errorMsg}>
                                    {errors.username.message}
                                </span>
                            )}
                        </div>
                        <div>
                            <label className={styles.label}>
                                {t('user.form.email') || 'Email'}
                            </label>
                            <input
                                type="email"
                                {...register('email')}
                                placeholder="john@example.com"
                                className={`${styles.input} ${errors.email ? styles.inputError : ''}`}
                            />
                            {errors.email && (
                                <span className={styles.errorMsg}>
                                    {errors.email.message}
                                </span>
                            )}
                        </div>
                    </div>

                    {/* Roles Selection */}
                    <div style={{ marginTop: '1rem' }}>
                        <label className={styles.label}>
                            {t('user.form.roles') || 'Assign Roles'}
                        </label>
                        <div className={styles.rolesGrid}>
                            {roles.map(role => {
                                const isSelected = watchRoleIds.includes(role.id);
                                return (
                                    <label
                                        key={role.id}
                                        className={`${styles.roleOption} ${isSelected ? styles.roleOptionSelected : ''}`}
                                    >
                                        <input
                                            type="checkbox"
                                            checked={isSelected}
                                            onChange={() => {
                                                const newRoles = isSelected
                                                    ? watchRoleIds.filter(id => id !== role.id)
                                                    : [...watchRoleIds, role.id];
                                                setValue('role_ids', newRoles);
                                            }}
                                            className={styles.roleCheckbox}
                                        />
                                        <span className={`${styles.roleName} ${isSelected ? styles.roleNameSelected : ''}`}>
                                            {role.name}
                                        </span>
                                    </label>
                                );
                            })}
                        </div>
                    </div>

                    {/* Status Toggle */}
                    <div style={{ marginTop: '1rem' }}>
                        <label className={styles.statusToggle}>
                            <input
                                type="checkbox"
                                checked={watchStatus}
                                onChange={(e) => setValue('status', e.target.checked)}
                                className={styles.statusCheckbox}
                            />
                            <span style={{ fontSize: '0.875rem', color: 'var(--color-text-main)' }}>
                                {t('user.form.active') || 'User is active'}
                            </span>
                        </label>
                    </div>

                    {/* Actions */}
                    <div className={styles.modalActions}>
                        <Button type="button" variant="secondary" onClick={() => setModalOpen(false)}>
                            {t('user.form.cancel') || 'Cancel'}
                        </Button>
                        <Button type="submit" variant="primary" loading={saving}>
                            {editingId ? (t('user.form.update') || 'Save Changes') : (t('user.form.save') || 'Create User')}
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
                title="Delete User"
                message={t('users.delete_confirm') || `Are you sure you want to delete "${deleteTarget?.name}"?`}
                confirmText="Delete"
                variant="danger"
                loading={deleting}
            />
        </div>
    );
};

export default UsersPage;
