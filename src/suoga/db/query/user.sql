
-- name: CreateUser :one
INSERT INTO users (
    username,
    pwd,
    email
) VALUES (
    $1,$2,$3
) RETURNING *;


-- name: GetUserForUpdate :one
SELECT * FROM users
    WHERE id = $1 LIMIT 1
    FOR NO KEY UPDATE;


-- name: ListUsers :many
SELECT * FROM users
    ORDER BY id
    LIMIT $1
    OFFSET $2;


-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;