\connect api;

-- Habilita extensão para UUIDs (caso seja Postgres antigo, ou use pgcrypto para gen_random_uuid)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE EXTENSION IF NOT EXISTS "citext";

-- =================================================================================================
-- 1. ROLE (Perfis de Acesso)
-- =================================================================================================
CREATE TABLE if not exists public.usr_role (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    "name" CITEXT NOT NULL,
    permissions text [ ] NOT NULL,
    CONSTRAINT uni_usr_role UNIQUE ("name")
);

-- Inserindo Roles Padrão (IDs fixos ou gerados, aqui deixamos gerar automático)
-- DICA: Se precisar de IDs fixos para seeds, pode passar explicitamente: '...-...'::uuid
INSERT INTO
    public.usr_role ("name", permissions)
VALUES
    ('ROOT', ARRAY [ '*' ]),
    ('ADMIN', ARRAY [ 
        'users:view', 'users:create', 'users:edit', 'users:delete', 
        'roles:view', 'roles:create', 'roles:edit', 'roles:delete' 
    ]);

-- =================================================================================================
-- 2. AUTH (Credenciais)
-- =================================================================================================
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

INSERT INTO
    public.usr_auth ("status", token, "password")
VALUES
    (
        true,
        'd048aee9-dd65-4ca0-aee7-230c1bf19d8c',
        '$2a$10$vqkyIvgHRU2sl2FGtlbkNeGFeTsJHQYz18abMJiLlGyJt.Ge99zYy'
    );

-- =================================================================================================
-- 3. USER (Identidade)
-- =================================================================================================
CREATE TABLE if not exists public.usr_user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    "name" CITEXT NOT NULL,
    username CITEXT NOT NULL,
    mail CITEXT NOT NULL,
    auth_id UUID NOT NULL,
    CONSTRAINT fk_usr_user_auth FOREIGN KEY (auth_id) REFERENCES public.usr_auth (id) ON DELETE CASCADE,
    CONSTRAINT uni_usr_user_mail UNIQUE (mail),
    CONSTRAINT uni_usr_user_username UNIQUE (username),
    CONSTRAINT uni_usr_user_auth UNIQUE (auth_id)
);

-- Precisamos dos IDs gerados para vincular.
-- Em um script de migração real, você usaria variáveis ou DO block.
-- Aqui, assumindo inserção direta com subquery para pegar o auth criado acima:
INSERT INTO
    public.usr_user (auth_id, "name", mail, username)
SELECT
    id,
    'Administrator',
    'admin@admin.com',
    'admin'
FROM
    public.usr_auth
WHERE
    token = 'd048aee9-dd65-4ca0-aee7-230c1bf19d8c';

-- Usando token como âncora para este exemplo

-- =================================================================================================
-- 4. USER ROLES (Pivot)
-- =================================================================================================
CREATE TABLE if not exists public.usr_user_role (
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    CONSTRAINT pkey_usr_user_role PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_uur_user FOREIGN KEY (user_id) REFERENCES public.usr_user (id) ON DELETE CASCADE,
    CONSTRAINT fk_uur_role FOREIGN KEY (role_id) REFERENCES public.usr_role (id) ON DELETE CASCADE
);

-- Vinculando Admin às Roles (Todas as roles existentes)
INSERT INTO
    public.usr_user_role (user_id, role_id)
SELECT
    u.id,
    r.id
FROM
    public.usr_user u
    CROSS JOIN public.usr_role r
WHERE
    u.username = 'admin';

-- =================================================================================================
-- 5. VIEWS (Adaptadas para UUID)
-- =================================================================================================
-- VIEW 1: Detalhes
CREATE
OR REPLACE VIEW public.vw_usr_details AS
SELECT
    u.id AS user_id,
    u.name,
    u.username,
    u.mail,
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

-- VIEW 2: Claims Claims
CREATE
OR REPLACE VIEW public.vw_usr_auth_claims AS
SELECT
    u.id AS user_id,
    u.username,
    u.mail,
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