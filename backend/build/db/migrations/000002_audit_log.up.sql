-- Create sys_audit_logs table for tracking CRUD operations
CREATE TABLE IF NOT EXISTS sys_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID REFERENCES usr_user(id) ON DELETE SET NULL,
    action TEXT NOT NULL,
    resource_entity TEXT NOT NULL,
    resource_id TEXT NOT NULL,
    metadata JSONB,
    ip_address TEXT,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_sys_audit_logs_actor_id ON sys_audit_logs(actor_id);
CREATE INDEX IF NOT EXISTS idx_sys_audit_logs_action ON sys_audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_sys_audit_logs_resource_entity ON sys_audit_logs(resource_entity);
CREATE INDEX IF NOT EXISTS idx_sys_audit_logs_resource_id ON sys_audit_logs(resource_id);
CREATE INDEX IF NOT EXISTS idx_sys_audit_logs_created_at ON sys_audit_logs(created_at);
