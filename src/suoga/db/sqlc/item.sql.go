// Code generated by sqlc. DO NOT EDIT.
// source: item.sql

package db

import (
	"context"
	"database/sql"
)

const createItem = `-- name: CreateItem :one
INSERT INTO items (
    content,
    user_id
) VALUES (
    $1, $2
) RETURNING id, content, user_id
`

type CreateItemParams struct {
	Content sql.NullInt64 `json:"content"`
	UserID  sql.NullInt32 `json:"userID"`
}

func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (Items, error) {
	row := q.db.QueryRowContext(ctx, createItem, arg.Content, arg.UserID)
	var i Items
	err := row.Scan(&i.ID, &i.Content, &i.UserID)
	return i, err
}

const listItemByUserId = `-- name: ListItemByUserId :many
SELECT id, content, user_id FROM items 
    WHERE user_id = $1
`

func (q *Queries) ListItemByUserId(ctx context.Context, userID sql.NullInt32) ([]Items, error) {
	rows, err := q.db.QueryContext(ctx, listItemByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Items
	for rows.Next() {
		var i Items
		if err := rows.Scan(&i.ID, &i.Content, &i.UserID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}