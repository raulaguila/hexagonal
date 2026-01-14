\connect api;

-- User Profile -------------------------------------------------------------------------------------------------------------------------------------
-- DROP SEQUENCE IF EXISTS public.seq_usr_profile_id;
CREATE SEQUENCE if not exists public.seq_usr_profile_id INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE;

-- DROP TABLE public.usr_profile;
CREATE TABLE if not exists public.usr_profile (
    id bigint PRIMARY KEY DEFAULT nextval('seq_usr_profile_id':: regclass) NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    "name" varchar(100) NOT NULL,
    permissions text [ ] NOT NULL,
    CONSTRAINT uni_usr_profile UNIQUE ("name")
);

INSERT INTO
    public.usr_profile (id, "name", permissions)
VALUES
    (1, 'ROOT', ARRAY [ 'users', 'profiles' ]);

ALTER SEQUENCE public.seq_usr_profile_id RESTART WITH 10;

-- User Auth ----------------------------------------------------------------------------------------------------------------------------------------
-- DROP sequence IF EXISTS public.seq_usr_auth_id;
CREATE SEQUENCE if not exists public.seq_usr_auth_id INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE;

-- DROP TABLE public.usr_auth;
CREATE TABLE if not exists public.usr_auth (
    id bigint PRIMARY KEY DEFAULT nextval('seq_usr_auth_id':: regclass) NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    "status" bool NOT NULL,
    profile_id bigint NOT NULL,
    token varchar(255) NULL,
    "password" varchar(255) NULL,
    CONSTRAINT uni_usr_auth UNIQUE (token),
    CONSTRAINT fk_usr_auth_profile FOREIGN KEY (profile_id) REFERENCES public.usr_profile (id)
);

CREATE INDEX if not exists idx_usr_auth_profile_id ON public.usr_auth USING btree (profile_id);

CREATE INDEX if not exists idx_usr_auth_token ON public.usr_auth USING btree (token);

-- Password: 12345678
INSERT INTO
    public.usr_auth (id, "status", profile_id, token, "password")
VALUES
    (1, true, 1, 'd048aee9-dd65-4ca0-aee7-230c1bf19d8c', '$2a$10$vqkyIvgHRU2sl2FGtlbkNeGFeTsJHQYz18abMJiLlGyJt.Ge99zYy');

ALTER SEQUENCE public.seq_usr_auth_id RESTART WITH 10;

-- User ---------------------------------------------------------------------------------------------------------------------------------------------
-- DROP SEQUENCE IF EXISTS public.seq_usr_user_id;
CREATE SEQUENCE if not exists public.seq_usr_user_id INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE;

-- DROP TABLE public.usr_user;
CREATE TABLE if not exists public.usr_user (
    id bigint PRIMARY KEY DEFAULT nextval('seq_usr_user_id':: regclass) NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    "name" varchar(255) NOT NULL,
    username varchar(255) NOT NULL,
    mail varchar(255) NOT NULL,
    auth_id bigint NOT NULL,
    CONSTRAINT fk_usr_user_auth FOREIGN KEY (auth_id) REFERENCES public.usr_auth (id) ON DELETE CASCADE,
    CONSTRAINT uni_usr_user UNIQUE (mail),
    CONSTRAINT uni_usr_user_username UNIQUE (username)
);

INSERT INTO
    public.usr_user (id, auth_id, "name", mail, username)
VALUES
    (1, 1, 'Administrator', 'admin@admin.com', 'admin');

ALTER SEQUENCE public.seq_usr_user_id RESTART WITH 10;