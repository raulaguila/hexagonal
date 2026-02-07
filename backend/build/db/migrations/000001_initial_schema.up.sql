-- Habilita extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "unaccent";
CREATE EXTENSION IF NOT EXISTS "citext";
CREATE EXTENSION IF NOT EXISTS "pg_cron";

-- 1. ROLE
CREATE TABLE if not exists public.usr_role (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    "name" CITEXT NOT NULL,
    permissions text [ ] NOT NULL,
    "enabled" boolean DEFAULT true NOT NULL,
    CONSTRAINT uni_usr_role UNIQUE ("name")
);

INSERT INTO public.usr_role ("name", permissions, "enabled")
VALUES
    ('ROOT', ARRAY [ '*' ], true),
    ('ADMIN', ARRAY [ 
        'users:view', 'users:create', 'users:edit', 'users:delete', 
        'roles:view', 'roles:create', 'roles:edit', 'roles:delete' 
    ], true)
ON CONFLICT DO NOTHING;

-- 2. AUTH
CREATE TABLE if not exists public.usr_auth (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    "status" bool NOT NULL DEFAULT true,
    token CITEXT NULL,
    "password" CITEXT NULL,
    CONSTRAINT uni_usr_auth UNIQUE (token)
);

CREATE INDEX if not exists idx_usr_auth_token ON public.usr_auth USING btree (token);

INSERT INTO public.usr_auth ("status", token, "password")
VALUES
    (
        true,
        'd048aee9-dd65-4ca0-aee7-230c1bf19d8c',
        '$2a$10$vqkyIvgHRU2sl2FGtlbkNeGFeTsJHQYz18abMJiLlGyJt.Ge99zYy'
    )
ON CONFLICT DO NOTHING;

-- 3. USER
CREATE TABLE if not exists public.usr_user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    "name" CITEXT NOT NULL,
    username CITEXT NOT NULL,
    email CITEXT NOT NULL,
    auth_id UUID NOT NULL,
    CONSTRAINT fk_usr_user_auth FOREIGN KEY (auth_id) REFERENCES public.usr_auth (id) ON DELETE CASCADE,
    CONSTRAINT uni_usr_user_email UNIQUE (email),
    CONSTRAINT uni_usr_user_username UNIQUE (username),
    CONSTRAINT uni_usr_user_auth UNIQUE (auth_id)
);

INSERT INTO public.usr_user (auth_id, "name", email, username)
SELECT
    id,
    'Administrator',
    'admin@admin.com',
    'admin'
FROM
    public.usr_auth
WHERE
    token = 'd048aee9-dd65-4ca0-aee7-230c1bf19d8c'
ON CONFLICT DO NOTHING;

-- 4. USER ROLES
CREATE TABLE if not exists public.usr_user_role (
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    CONSTRAINT pkey_usr_user_role PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_uur_user FOREIGN KEY (user_id) REFERENCES public.usr_user (id) ON DELETE CASCADE,
    CONSTRAINT fk_uur_role FOREIGN KEY (role_id) REFERENCES public.usr_role (id) ON DELETE CASCADE
);

INSERT INTO public.usr_user_role (user_id, role_id)
SELECT
    u.id,
    r.id
FROM
    public.usr_user u
    CROSS JOIN public.usr_role r
WHERE
    u.username = 'admin'
ON CONFLICT DO NOTHING;

-- 5. AUDIT LOGS
CREATE TABLE IF NOT EXISTS sys_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID REFERENCES usr_user(id) ON DELETE SET NULL,
    "action" TEXT NOT NULL,
    resource_entity TEXT NOT NULL,
    resource_id TEXT NOT NULL,
    metadata JSONB,
    ip_address TEXT,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sys_audit_logs_actor_id ON sys_audit_logs(actor_id);
CREATE INDEX IF NOT EXISTS idx_sys_audit_logs_action ON sys_audit_logs("action");
CREATE INDEX IF NOT EXISTS idx_sys_audit_logs_resource_entity ON sys_audit_logs(resource_entity);
CREATE INDEX IF NOT EXISTS idx_sys_audit_logs_resource_id ON sys_audit_logs(resource_id);
CREATE INDEX IF NOT EXISTS idx_sys_audit_logs_created_at ON sys_audit_logs(created_at);

-- 6. CRON
SELECT cron.schedule('delete_old_audits_logs', '30 3 * * *', $$
    DELETE FROM sys_audit_logs WHERE created_at < now() - INTERVAL '390 days'
$$);

-- VIEWS
-- Create view for simplified user details
CREATE OR REPLACE VIEW public.vw_usr_details AS
SELECT
    u.id AS user_id,
    u.name,
    u.username,
    u.email,
    u.created_at,
    a.status AS is_active,
    array_agg(r.name) AS roles
FROM
    public.usr_user u
    JOIN public.usr_auth a ON u.auth_id = a.id
    LEFT JOIN public.usr_user_role uur ON u.id = uur.user_id
    LEFT JOIN public.usr_role r ON uur.role_id = r.id
GROUP BY
    u.id,
    a.id;

-- Create view for simplified user authentication claims
CREATE OR REPLACE VIEW public.vw_usr_auth_claims AS
SELECT
    u.id AS user_id,
    u.username,
    u.email,
    a.id AS auth_id,
    a.password AS password_hash,
    a.token,
    a.status,
    array_agg(DISTINCT r.name) AS roles,
    COALESCE(
        (
            SELECT
                array_agg(DISTINCT p)
            FROM
                public.usr_user_role uur_sub
                JOIN public.usr_role r_sub ON uur_sub.role_id = r_sub.id
                CROSS JOIN unnest(r_sub.permissions) AS p
            WHERE
                uur_sub.user_id = u.id
        ),
        ARRAY [ ]:: text [ ]
    ) AS all_permissions
FROM
    public.usr_user u
    JOIN public.usr_auth a ON u.auth_id = a.id
    LEFT JOIN public.usr_user_role uur ON u.id = uur.user_id
    LEFT JOIN public.usr_role r ON uur.role_id = r.id
GROUP BY
    u.id,
    a.id;

-- Create view for simplified audit log viewing with user information
CREATE OR REPLACE VIEW vw_audit_logs AS
SELECT
    al.id,
    al.created_at AS "timestamp",
    al.action,
    al.resource_entity AS resource_type,
    al.resource_id,
    al.actor_id,
    u.name AS actor_name,
    u.username AS actor_username,
    al.ip_address,
    al.user_agent,
    al.metadata
FROM sys_audit_logs al
LEFT JOIN usr_user u ON al.actor_id = u.id
ORDER BY al.created_at DESC;

-- Create a more detailed view with readable metadata
CREATE OR REPLACE VIEW vw_audit_logs_detailed AS
SELECT
    al.id,
    TO_CHAR(al.created_at AT TIME ZONE 'America/Manaus', 'DD/MM/YYYY HH24:MI:SS') AS "timestamp_br",
    al.created_at AS timestamp_utc,
    CASE al.action
        WHEN 'create' THEN 'Criação'
        WHEN 'update' THEN 'Atualização'
        WHEN 'delete' THEN 'Exclusão'
        ELSE al.action
    END AS action_label,
    al.action,
    CASE al.resource_entity
        WHEN 'user' THEN 'Usuário'
        WHEN 'role' THEN 'Perfil'
        ELSE al.resource_entity
    END AS resource_type_label,
    al.resource_entity AS resource_type,
    al.resource_id,
    al.actor_id,
    COALESCE(u.name, 'Sistema') AS actor_name,
    COALESCE(u.username, 'system') AS actor_username,
    al.ip_address,
    al.user_agent,
    al.metadata,
    al.metadata->>'input' AS input_data
FROM sys_audit_logs al
LEFT JOIN usr_user u ON al.actor_id = u.id
ORDER BY al.created_at DESC;

COMMENT ON VIEW vw_audit_logs IS 'Simplified view of audit logs with actor information';
COMMENT ON VIEW vw_audit_logs_detailed IS 'Detailed view of audit logs with localized labels and formatted timestamps';
