-- name: CreateItem :one
INSERT INTO items (
    content,
    user_id
) VALUES (
    $1, $2
) RETURNING *;


-- name: ListItemByUserId :many
SELECT * FROM items 
    WHERE user_id = $1;
