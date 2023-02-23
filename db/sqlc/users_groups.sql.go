// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: users_groups.sql

package db

import (
	"context"
)

const assignOwnerToGroup = `-- name: AssignOwnerToGroup :one
INSERT INTO users_groups (
  user_id,
  group_id,
  role
) VALUES (
    $1, $2, 'owner'
)RETURNING id, user_id, group_id, role
`

type AssignOwnerToGroupParams struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
}

func (q *Queries) AssignOwnerToGroup(ctx context.Context, arg AssignOwnerToGroupParams) (UsersGroup, error) {
	row := q.db.QueryRowContext(ctx, assignOwnerToGroup, arg.UserID, arg.GroupID)
	var i UsersGroup
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.GroupID,
		&i.Role,
	)
	return i, err
}

const assignUserToGroup = `-- name: AssignUserToGroup :one
INSERT INTO users_groups (
  user_id,
  group_id
) VALUES (
    $1, $2
)RETURNING id, user_id, group_id, role
`

type AssignUserToGroupParams struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
}

func (q *Queries) AssignUserToGroup(ctx context.Context, arg AssignUserToGroupParams) (UsersGroup, error) {
	row := q.db.QueryRowContext(ctx, assignUserToGroup, arg.UserID, arg.GroupID)
	var i UsersGroup
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.GroupID,
		&i.Role,
	)
	return i, err
}

const getFromUsersGroups = `-- name: GetFromUsersGroups :one
SELECT id, user_id, group_id, role FROM users_groups WHERE user_id = $1 AND group_id = $2
`

type GetFromUsersGroupsParams struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
}

func (q *Queries) GetFromUsersGroups(ctx context.Context, arg GetFromUsersGroupsParams) (UsersGroup, error) {
	row := q.db.QueryRowContext(ctx, getFromUsersGroups, arg.UserID, arg.GroupID)
	var i UsersGroup
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.GroupID,
		&i.Role,
	)
	return i, err
}
