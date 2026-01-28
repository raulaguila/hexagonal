-- name: CreateProfile :one
INSERT INTO usr_profile (
  created_at, updated_at, name, permissions
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetProfile :one
SELECT * FROM usr_profile
WHERE id = $1 LIMIT 1;

-- name: ListProfiles :many
SELECT * FROM usr_profile
ORDER BY id;

-- name: CreateUser :one
INSERT INTO usr_user (
  created_at, updated_at, name, username, mail, auth_id
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: CreateUsersBatch :many
INSERT INTO usr_user (name, username, mail)
SELECT unnest($1::text[]), unnest($2::text[]), unnest($3::text[])
RETURNING *;

-- name: GetUser :one
SELECT * FROM usr_user
WHERE id = $1 LIMIT 1;
