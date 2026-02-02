import React, { useEffect, useState } from 'react';
import api from '../utils/api';
import { Table, Thead, Tbody, Tr, Th, Td } from '../components/common/Table';
import Button from '../components/common/Button';
import Modal from '../components/common/Modal';
import Input from '../components/common/Input';
import { Plus, Trash2, Edit2, Shield } from 'lucide-react';

const RolesPage = () => {
    const [roles, setRoles] = useState([]);
    const [loading, setLoading] = useState(true);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [formData, setFormData] = useState({ name: '', permissions: [] });
    const [creating, setCreating] = useState(false);

    const MODULES = [
        { key: "users", label: "Users" },
        { key: "roles", label: "Roles" }
    ];

    const ACTIONS = [
        { key: "view", label: "View" },
        { key: "create", label: "Create" },
        { key: "edit", label: "Edit" },
        { key: "delete", label: "Delete" }
    ];

    const togglePermission = (perm) => {
        setFormData(prev => {
            const currentPerms = prev.permissions || [];
            if (currentPerms.includes(perm)) {
                return { ...prev, permissions: currentPerms.filter(p => p !== perm) };
            } else {
                return { ...prev, permissions: [...currentPerms, perm] };
            }
        });
    };

    const fetchRoles = async () => {
        try {
            const { data } = await api.get('/role');
            setRoles(Array.isArray(data) ? data : data.items || []);
        } catch (error) {
            console.error("Failed to fetch roles", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchRoles();
    }, []);

    const [editingId, setEditingId] = useState(null);

    const handleEdit = (role) => {
        setFormData({
            name: role.name,
            permissions: role.permissions || []
        });
        setEditingId(role.id);
        setIsModalOpen(true);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setCreating(true);
        try {
            if (editingId) {
                // Update role
                const { data } = await api.put(`/role/${editingId}`, {
                    name: formData.name,
                    permissions: formData.permissions
                });
                setRoles(roles.map(r => r.id === editingId ? data : r));
                alert("Role updated successfully");
            } else {
                // Create role
                const { data } = await api.post('/role', {
                    name: formData.name,
                    permissions: []
                });
                fetchRoles();
                alert("Role created successfully");
            }
            setIsModalOpen(false);
            setFormData({ name: '', permissions: [] });
            setEditingId(null);
        } catch (error) {
            alert("Failed to save role: " + (error.response?.data?.message || error.message));
        } finally {
            setCreating(false);
        }
    };

    const openCreateModal = () => {
        setFormData({ name: '', permissions: [] });
        setEditingId(null);
        setIsModalOpen(true);
    };

    const handleDelete = async (id) => {
        if (!window.confirm("Delete this role?")) return;
        try {
            await api.delete('/role', { data: { ids: [id] } });
            setRoles(roles.filter(r => r.id !== id));
        } catch (error) {
            alert("Failed to delete role");
        }
    };

    return (
        <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
                <div>
                    <h1 style={{ fontSize: '1.875rem', fontWeight: 700, margin: 0 }}>Roles</h1>
                </div>
                <Button onClick={openCreateModal}>
                    <Plus size={18} />
                    Add Role
                </Button>
            </div>

            <Table>
                <Thead>
                    <Tr>
                        <Th>Role Name</Th>
                        <Th>Permissions</Th>
                        <Th>Actions</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {roles.map(role => (
                        <Tr key={role.id}>
                            <Td>
                                <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                                    <Shield size={18} color="var(--color-primary)" />
                                    <span style={{ fontWeight: 500 }}>{role.name}</span>
                                </div>
                            </Td>
                            <Td>
                                <div style={{ fontSize: '0.75rem', color: 'var(--color-text-secondary)', maxWidth: '400px', overflow: 'hidden', whiteSpace: 'nowrap', textOverflow: 'ellipsis' }}>
                                    {role.permissions && role.permissions.join(', ')}
                                </div>
                            </Td>
                            <Td>
                                <div style={{ display: 'flex', gap: '0.5rem' }}>
                                    <button
                                        onClick={() => handleEdit(role)}
                                        style={{ color: 'var(--color-text-secondary)', background: 'none', cursor: 'pointer' }}
                                        title="Edit"
                                    >
                                        <Edit2 size={18} />
                                    </button>
                                    <button
                                        onClick={() => handleDelete(role.id)}
                                        style={{ color: 'var(--color-error)', background: 'none', cursor: 'pointer' }}
                                        title="Delete"
                                    >
                                        <Trash2 size={18} />
                                    </button>
                                </div>
                            </Td>
                        </Tr>
                    ))}
                </Tbody>
            </Table>

            <Modal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                title={editingId ? "Edit Role" : "Create New Role"}
            >
                <form onSubmit={handleSubmit}>
                    <Input
                        label="Role Name"
                        value={formData.name}
                        onChange={e => setFormData({ ...formData, name: e.target.value })}
                        required
                        placeholder="e.g. EDITOR"
                    />
                    <div style={{ marginTop: '1rem' }}>
                        <label style={{ display: 'block', fontSize: '0.875rem', fontWeight: 500, marginBottom: '0.5rem', color: 'var(--color-text-main)' }}>
                            Permissions
                        </label>
                        <div style={{
                            padding: '0.5rem',
                            border: '1px solid var(--color-border)',
                            borderRadius: 'var(--radius-md)',
                            maxHeight: '300px',
                            overflowY: 'auto'
                        }}>
                            {MODULES.map(module => (
                                <div key={module.key} style={{ marginBottom: '1rem' }}>
                                    <div style={{
                                        fontSize: '0.875rem',
                                        fontWeight: 600,
                                        color: 'var(--color-text-main)',
                                        marginBottom: '0.5rem',
                                        textTransform: 'capitalize'
                                    }}>
                                        {module.label} Module
                                    </div>
                                    <div style={{
                                        display: 'grid',
                                        gridTemplateColumns: 'repeat(2, 1fr)',
                                        gap: '0.5rem'
                                    }}>
                                        {ACTIONS.map(action => {
                                            const permString = `${module.key}:${action.key}`;
                                            return (
                                                <label key={permString} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.875rem', cursor: 'pointer' }}>
                                                    <input
                                                        type="checkbox"
                                                        checked={(formData.permissions || []).includes(permString)}
                                                        onChange={() => togglePermission(permString)}
                                                        style={{ accentColor: 'var(--color-primary)' }}
                                                    />
                                                    {action.label}
                                                </label>
                                            );
                                        })}
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '1rem', marginTop: '1.5rem' }}>
                        <Button type="button" variant="ghost" onClick={() => setIsModalOpen(false)}>Cancel</Button>
                        <Button type="submit" disabled={creating}>
                            {creating ? 'Saving...' : (editingId ? 'Update Role' : 'Create Role')}
                        </Button>
                    </div>
                </form>
            </Modal>
        </div >
    );
};

export default RolesPage;
