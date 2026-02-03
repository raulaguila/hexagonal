import React, { useEffect, useState, useCallback } from 'react';
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

const UsersPage = () => {
    const toast = useToast();
    const { t } = usePreferences();
    const { hasPermission, isRoot } = usePermissions();

    // Data hooks
    const { users, loading, pagination, fetchUsers, deleteUser: removeUser } = useUsers();
    const { roles, fetchRoles } = useRoles();

    // Local state
    const [search, setSearch] = useState('');
    const [modalOpen, setModalOpen] = useState(false);
    const [editingId, setEditingId] = useState(null);
    const [saving, setSaving] = useState(false);
    const [formData, setFormData] = useState({
        name: '',
        username: '',
        mail: '',
        role_ids: [],
        status: true
    });

    // Confirm dialog state
    const [confirmOpen, setConfirmOpen] = useState(false);
    const [deleteTarget, setDeleteTarget] = useState(null);
    const [deleting, setDeleting] = useState(false);

    // Fetch data on mount
    useEffect(() => {
        fetchUsers({ page: 1, limit: 10, search: '' });
        fetchRoles();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    // Search handler - only on button click
    const handleSearch = useCallback(() => {
        fetchUsers({ page: 1, limit: pagination.limit || 10, search });
    }, [fetchUsers, pagination.limit, search]);

    // Handle Enter key in search
    const handleSearchKeyDown = (e) => {
        if (e.key === 'Enter') {
            handleSearch();
        }
    };

    // Page change handler
    const handlePageChange = (page) => {
        fetchUsers({ page, limit: pagination.limit || 10, search });
    };

    // Open edit modal
    const handleEdit = (user) => {
        setFormData({
            name: user.name,
            username: user.username,
            mail: user.email,
            role_ids: user.roles ? user.roles.map(r => r.id) : [],
            status: user.status !== undefined ? user.status : true
        });
        setEditingId(user.id);
        setModalOpen(true);
    };

    // Open create modal
    const handleCreate = () => {
        setFormData({
            name: '',
            username: '',
            mail: '',
            role_ids: [],
            status: true
        });
        setEditingId(null);
        setModalOpen(true);
    };

    // Submit form
    const handleSubmit = async (e) => {
        e.preventDefault();
        setSaving(true);

        try {
            const payload = {
                name: formData.name,
                username: formData.username,
                email: formData.mail,
                role_ids: formData.role_ids,
                status: formData.status
            };

            if (editingId) {
                await userService.updateUser(editingId, payload);
                toast.success(t('users.updated') || 'User updated successfully');
            } else {
                await userService.createUser(payload);
                toast.success(t('users.created') || 'User created successfully');
            }

            setModalOpen(false);
            fetchUsers({ page: pagination.page || 1, limit: pagination.limit || 10, search });
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
        <div>
            {/* Header */}
            <div style={{ marginBottom: '1.5rem' }}>
                <h1 style={{
                    fontSize: '1.5rem',
                    fontWeight: 700,
                    margin: 0,
                    color: 'var(--color-text-main)'
                }}>
                    {t('users.title') || 'Users'}
                </h1>
                <p style={{
                    color: 'var(--color-text-secondary)',
                    marginTop: '0.25rem',
                    fontSize: '0.875rem'
                }}>
                    {t('users.subtitle') || 'Manage system access and profiles'}
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
                        {t('users.add') || 'Add User'}
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
                            <Th>{t('users.table.user') || 'User'}</Th>
                            <Th>{t('users.table.contact') || 'Contact'}</Th>
                            <Th>{t('users.table.roles') || 'Roles'}</Th>
                            <Th>{t('users.table.status') || 'Status'}</Th>
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
                                <Td colSpan={5}>
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
                                <Tr key={user.id} className="table-row">
                                    {/* User Info */}
                                    <Td>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                                            <div style={{
                                                width: '40px',
                                                height: '40px',
                                                borderRadius: '50%',
                                                background: 'linear-gradient(135deg, var(--color-primary-light), var(--color-surface))',
                                                color: 'var(--color-primary)',
                                                display: 'flex',
                                                alignItems: 'center',
                                                justifyContent: 'center',
                                                fontWeight: 600,
                                                fontSize: '1rem',
                                                border: '1px solid var(--color-border)'
                                            }}>
                                                {user.name?.charAt(0) || '?'}
                                            </div>
                                            <div>
                                                <div style={{ fontWeight: 500, color: 'var(--color-text-main)' }}>
                                                    {user.name}
                                                </div>
                                                <div style={{ fontSize: '0.75rem', color: 'var(--color-text-muted)' }}>
                                                    @{user.username}
                                                </div>
                                            </div>
                                        </div>
                                    </Td>

                                    {/* Contact */}
                                    <Td>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'var(--color-text-secondary)' }}>
                                            <Mail size={14} />
                                            <span style={{ fontSize: '0.875rem' }}>{user.email}</span>
                                        </div>
                                    </Td>

                                    {/* Roles */}
                                    <Td>
                                        <div style={{ display: 'flex', gap: '0.25rem', flexWrap: 'wrap' }}>
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
                                    <Td>
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
                                    <Td style={{ textAlign: 'center' }}>
                                        <div style={{ display: 'inline-flex', gap: '0.25rem', justifyContent: 'center' }}>
                                            {canEdit && (
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => handleEdit(user)}
                                                    style={{ padding: '0.375rem' }}
                                                >
                                                    <Pen size={16} />
                                                </Button>
                                            )}
                                            {canDelete && (
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => handleDeleteClick(user)}
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
                <form onSubmit={handleSubmit}>
                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                        <div style={{ gridColumn: 'span 2' }}>
                            <label style={{
                                display: 'block',
                                fontSize: '0.875rem',
                                fontWeight: 500,
                                marginBottom: '0.5rem',
                                color: 'var(--color-text-main)'
                            }}>
                                {t('user.form.name') || 'Full Name'}
                            </label>
                            <input
                                type="text"
                                value={formData.name}
                                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                required
                                placeholder="e.g. John Doe"
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
                        <div>
                            <label style={{
                                display: 'block',
                                fontSize: '0.875rem',
                                fontWeight: 500,
                                marginBottom: '0.5rem',
                                color: 'var(--color-text-main)'
                            }}>
                                {t('user.form.username') || 'Username'}
                            </label>
                            <input
                                type="text"
                                value={formData.username}
                                onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                                required
                                placeholder="e.g. jdoe"
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
                        <div>
                            <label style={{
                                display: 'block',
                                fontSize: '0.875rem',
                                fontWeight: 500,
                                marginBottom: '0.5rem',
                                color: 'var(--color-text-main)'
                            }}>
                                {t('user.form.email') || 'Email'}
                            </label>
                            <input
                                type="email"
                                value={formData.mail}
                                onChange={(e) => setFormData({ ...formData, mail: e.target.value })}
                                required
                                placeholder="john@example.com"
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
                    </div>

                    {/* Roles Selection */}
                    <div style={{ marginTop: '1rem' }}>
                        <label style={{
                            display: 'block',
                            fontSize: '0.875rem',
                            fontWeight: 500,
                            marginBottom: '0.5rem',
                            color: 'var(--color-text-main)'
                        }}>
                            {t('user.form.roles') || 'Assign Roles'}
                        </label>
                        <div style={{
                            display: 'grid',
                            gridTemplateColumns: 'repeat(auto-fill, minmax(140px, 1fr))',
                            gap: '0.5rem'
                        }}>
                            {roles.map(role => (
                                <label
                                    key={role.id}
                                    style={{
                                        display: 'flex',
                                        alignItems: 'center',
                                        gap: '0.5rem',
                                        padding: '0.5rem 0.75rem',
                                        borderRadius: 'var(--radius-md)',
                                        border: formData.role_ids.includes(role.id)
                                            ? '1px solid var(--color-primary)'
                                            : '1px solid var(--color-border)',
                                        backgroundColor: formData.role_ids.includes(role.id)
                                            ? 'var(--color-primary-light)'
                                            : 'var(--color-surface)',
                                        cursor: 'pointer',
                                        transition: 'all 0.15s'
                                    }}
                                >
                                    <input
                                        type="checkbox"
                                        checked={formData.role_ids.includes(role.id)}
                                        onChange={() => {
                                            const newRoles = formData.role_ids.includes(role.id)
                                                ? formData.role_ids.filter(id => id !== role.id)
                                                : [...formData.role_ids, role.id];
                                            setFormData({ ...formData, role_ids: newRoles });
                                        }}
                                        style={{ accentColor: 'var(--color-primary)' }}
                                    />
                                    <span style={{
                                        fontSize: '0.875rem',
                                        fontWeight: formData.role_ids.includes(role.id) ? 600 : 400,
                                        color: formData.role_ids.includes(role.id)
                                            ? 'var(--color-primary)'
                                            : 'var(--color-text-main)'
                                    }}>
                                        {role.name}
                                    </span>
                                </label>
                            ))}
                        </div>
                    </div>

                    {/* Status Toggle */}
                    <div style={{ marginTop: '1rem' }}>
                        <label style={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: '0.75rem',
                            cursor: 'pointer'
                        }}>
                            <input
                                type="checkbox"
                                checked={formData.status}
                                onChange={(e) => setFormData({ ...formData, status: e.target.checked })}
                                style={{
                                    width: '18px',
                                    height: '18px',
                                    accentColor: 'var(--color-primary)'
                                }}
                            />
                            <span style={{ fontSize: '0.875rem', color: 'var(--color-text-main)' }}>
                                {t('user.form.active') || 'User is active'}
                            </span>
                        </label>
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
