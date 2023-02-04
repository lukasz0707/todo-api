// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: groups.sql

package db

import (
	"context"
)

const creategroup = `-- name: Creategroup :one
INSERT INTO groups (
  group_name,
  owner_id
) VALUES (
    $1, $2
)RETURNING id, group_name, owner_id
`

type CreategroupParams struct {
	GroupName string `json:"group_name"`
	OwnerID   int64  `json:"owner_id"`
}

func (q *Queries) Creategroup(ctx context.Context, arg CreategroupParams) (Group, error) {
	row := q.db.QueryRowContext(ctx, creategroup, arg.GroupName, arg.OwnerID)
	var i Group
	err := row.Scan(&i.ID, &i.GroupName, &i.OwnerID)
	return i, err
}
