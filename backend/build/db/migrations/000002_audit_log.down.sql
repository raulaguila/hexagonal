-- Drop sys_audit_logs table and indexes
DROP INDEX IF EXISTS idx_sys_audit_logs_created_at;
DROP INDEX IF EXISTS idx_sys_audit_logs_resource_id;
DROP INDEX IF EXISTS idx_sys_audit_logs_resource_entity;
DROP INDEX IF EXISTS idx_sys_audit_logs_action;
DROP INDEX IF EXISTS idx_sys_audit_logs_actor_id;
DROP TABLE IF EXISTS sys_audit_logs;
