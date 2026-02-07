DROP VIEW IF EXISTS public.vw_usr_auth_claims;
DROP VIEW IF EXISTS public.vw_usr_details;
DROP TABLE IF EXISTS public.usr_user_role;
DROP TABLE IF EXISTS public.usr_user;
DROP TABLE IF EXISTS public.usr_auth;
DROP TABLE IF EXISTS public.usr_role;

-- Drop sys_audit_logs table and indexes
DROP INDEX IF EXISTS idx_sys_audit_logs_created_at;
DROP INDEX IF EXISTS idx_sys_audit_logs_resource_id;
DROP INDEX IF EXISTS idx_sys_audit_logs_resource_entity;
DROP INDEX IF EXISTS idx_sys_audit_logs_action;
DROP INDEX IF EXISTS idx_sys_audit_logs_actor_id;
DROP TABLE IF EXISTS sys_audit_logs;

-- Drop audit log views
DROP VIEW IF EXISTS vw_audit_logs_detailed;
DROP VIEW IF EXISTS vw_audit_logs;

-- Drop pg_cron extension
SELECT cron.unschedule('delete_old_audits_logs');
DROP EXTENSION IF EXISTS "pg_cron";

-- Drop other extensions
DROP EXTENSION IF EXISTS "unaccent";
DROP EXTENSION IF EXISTS "citext";
DROP EXTENSION IF EXISTS "pgcrypto";
