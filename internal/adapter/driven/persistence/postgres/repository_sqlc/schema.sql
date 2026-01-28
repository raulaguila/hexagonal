CREATE TABLE usr_profile (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    name VARCHAR(100) NOT NULL UNIQUE,
    permissions TEXT[] NOT NULL
);

CREATE TABLE usr_auth (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    status BOOLEAN NOT NULL,
    profile_id BIGINT NOT NULL REFERENCES usr_profile(id),
    token VARCHAR(255) UNIQUE,
    password VARCHAR(255)
);

CREATE TABLE usr_user (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    name TEXT,
    username TEXT,
    mail TEXT,
    auth_id BIGINT REFERENCES usr_auth(id) ON DELETE CASCADE
);
