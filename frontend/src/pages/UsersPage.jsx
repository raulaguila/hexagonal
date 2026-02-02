import React, { useEffect, useState } from 'react';
import api from '../utils/api';
import { Table, Thead, Tbody, Tr, Th, Td } from '../components/common/Table';
import Button from '../components/common/Button';
import Modal from '../components/common/Modal';
import Input from '../components/common/Input';
import { Plus, Trash2, Edit2, User } from 'lucide-react';

const UsersPage = () => {
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [formData, setFormData] = useState({ name: '', username: '', mail: '', role_ids: [], status: true });
    const [creating, setCreating] = useState(false);

    // Fetch Users
    const fetchUsers = async () => {
        try {
            const { data } = await api.get('/user');
            // Handle both paginated ({ items: [] }) and array responses
            setUsers(Array.isArray(data) ? data : data.items || []);
        } catch (error) {
            console.error("Failed to fetch users", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchUsers();
        fetchRoles();
    }, []);

    const [availableRoles, setAvailableRoles] = useState([]);

    const fetchRoles = async () => {
        try {
            const { data } = await api.get('/role');
            setAvailableRoles(Array.isArray(data) ? data : data.items || []);
        } catch (error) {
            console.error("Failed to fetch roles", error);
        }
    };

    const [editingId, setEditingId] = useState(null);

    const handleEdit = (user) => {
        setFormData({
            name: user.name,
            username: user.corp_id, // backend uses corp_id for username output
            mail: user.email,
            role_ids: user.roles ? user.roles.map(r => r.id) : [],
            status: user.status !== undefined ? user.status : true
        });
        setEditingId(user.id);
        setIsModalOpen(true);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setCreating(true);
        try {
            if (editingId) {
                // Update user
                const { data } = await api.put(`/user/${editingId}`, {
                    name: formData.name,
                    username: formData.username,
                    email: formData.mail,
                    role_ids: formData.role_ids,
                    status: formData.status
                });
                setUsers(users.map(u => u.id === editingId ? data : u));
                alert("User updated successfully");
            } else {
                // Create user
                const { data } = await api.post('/user', {
                    name: formData.name,
                    username: formData.username,
                    email: formData.mail,
                    role_ids: formData.role_ids,
                    status: formData.status
                });
                setUsers([...users, data]);
                alert("User created successfully");
            }
            setIsModalOpen(false);
            setFormData({ name: '', username: '', mail: '' });
            setEditingId(null);
        } catch (error) {
            alert("Failed to save user: " + (error.response?.data?.message || error.message));
        } finally {
            setCreating(false);
        }
    };

    const openCreateModal = () => {
        setFormData({ name: '', username: '', mail: '', role_ids: [], status: true });
        setEditingId(null);
        setIsModalOpen(true);
    };

    const handleDelete = async (id) => {
        if (!window.confirm("Are you sure you want to delete this user?")) return;
        try {
            // API expects { ids: [id] } body for delete?
            // Handler: deleteUser -> idsBodyDTO -> IDsInput { IDs []string }
            await api.delete('/user', { data: { ids: [id] } });
            setUsers(users.filter(u => u.id !== id));
        } catch (error) {
            alert("Failed to delete user");
        }
    };

    return (
        <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
                <div>
                    <h1 style={{ fontSize: '1.875rem', fontWeight: 700, margin: 0 }}>Users</h1>
                    <p style={{ color: 'var(--color-text-secondary)' }}>Manage system access and profiles</p>
                </div>
                <Button onClick={openCreateModal}>
                    <Plus size={18} />
                    Add User
                </Button>
            </div>

            <Table>
                <Thead>
                    <Tr>
                        <Th>User</Th>
                        <Th>Details</Th>
                        <Th>Roles</Th>
                        <Th>Actions</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {users.map(user => (
                        <Tr key={user.id}>
                            <Td>
                                <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                                    <div style={{
                                        width: '40px', height: '40px', borderRadius: '50%',
                                        backgroundColor: 'var(--color-primary-light)',
                                        color: 'var(--color-primary)',
                                        display: 'flex', alignItems: 'center', justifyContent: 'center'
                                    }}>
                                        <User size={20} />
                                    </div>
                                    <div>
                                        <div style={{ fontWeight: 500 }}>{user.name}</div>
                                        <div style={{ fontSize: '0.75rem', color: 'var(--color-text-secondary)' }}>@{user.corp_id}</div>
                                    </div>
                                </div>
                            </Td>
                            <Td>{user.email}</Td>
                            <Td>
                                <div style={{ display: 'flex', gap: '0.25rem', flexWrap: 'wrap' }}>
                                    {user.roles && user.roles.map(role => (
                                        <span key={role.id} style={{
                                            fontSize: '0.75rem',
                                            padding: '0.125rem 0.5rem',
                                            borderRadius: '999px',
                                            backgroundColor: 'var(--color-background)',
                                            border: '1px solid var(--color-border)',
                                            color: 'var(--color-text-secondary)'
                                        }}>
                                            {role.name}
                                        </span>
                                    ))}
                                </div>
                            </Td>
                            <Td>
                                <div style={{ display: 'flex', gap: '0.5rem' }}>
                                    <button
                                        onClick={() => handleEdit(user)}
                                        style={{ color: 'var(--color-text-secondary)', background: 'none', cursor: 'pointer' }}
                                        title="Edit"
                                    >
                                        <Edit2 size={18} />
                                    </button>
                                    <button
                                        onClick={() => handleDelete(user.id)}
                                        style={{ color: 'var(--color-error)', background: 'none', cursor: 'pointer' }}
                                        title="Delete"
                                    >
                                        <Trash2 size={18} />
                                    </button>
                                </div>
                            </Td>
                        </Tr>
                    ))}
                    {!loading && users.length === 0 && (
                        <Tr>
                            <Td colSpan={4} style={{ textAlign: 'center', padding: '3rem' }}>
                                No users found.
                            </Td>
                        </Tr>
                    )}
                </Tbody>
            </Table>

            <Modal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                title={editingId ? "Edit User" : "Create New User"}
            >
                <form onSubmit={handleSubmit}>
                    <Input
                        label="Full Name"
                        value={formData.name}
                        onChange={e => setFormData({ ...formData, name: e.target.value })}
                        required
                        placeholder="John Doe"
                    />
                    <Input
                        label="Username"
                        value={formData.username}
                        onChange={e => setFormData({ ...formData, username: e.target.value })}
                        required
                        placeholder="jdoe"
                    />
                    <Input
                        label="Email"
                        type="email"
                        value={formData.mail}
                        onChange={e => setFormData({ ...formData, mail: e.target.value })}
                        required
                        placeholder="john@example.com"
                    />

                    {/* Role Selection */}
                    <div style={{ marginTop: '1rem' }}>
                        <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: 500, marginBottom: '0.5rem', color: 'var(--color-text-main)' }}>
                            Roles
                        </label>
                        <div style={{
                            display: 'grid',
                            gridTemplateColumns: 'repeat(2, 1fr)',
                            gap: '0.5rem',
                            padding: '0.5rem',
                            border: '1px solid var(--color-border)',
                            borderRadius: 'var(--radius-md)',
                            maxHeight: '150px',
                            overflowY: 'auto'
                        }}>
                            {availableRoles.map(role => (
                                <label key={role.id} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.875rem', cursor: 'pointer' }}>
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
                                    {role.name}
                                </label>
                            ))}
                        </div>
                    </div>

                    {/* Status Toggle (Only for editing or if supported during create) */}
                    <div style={{ marginTop: '1rem' }}>
                        <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.875rem', fontWeight: 500, cursor: 'pointer', color: 'var(--color-text-main)' }}>
                            <input
                                type="checkbox"
                                checked={formData.status}
                                onChange={e => setFormData({ ...formData, status: e.target.checked })}
                                style={{ accentColor: 'var(--color-primary)' }}
                            />
                            Active User
                        </label>
                    </div>

                    <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '1rem', marginTop: '1.5rem' }}>
                        <Button type="button" variant="ghost" onClick={() => setIsModalOpen(false)}>Cancel</Button>
                        <Button type="submit" disabled={creating}>
                            {creating ? 'Saving...' : (editingId ? 'Update User' : 'Create User')}
                        </Button>
                    </div>
                </form>
            </Modal>
        </div>
    );
};

export default UsersPage;
