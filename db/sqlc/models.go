// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Todo struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	Group     string             `json:"group"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	// oneof(ongoing, suspended, completed)
	Status   string             `json:"status"`
	Deadline pgtype.Timestamptz `json:"deadline"`
}

type User struct {
	ID                int64              `json:"id"`
	Username          string             `json:"username"`
	HashedPassword    string             `json:"hashed_password"`
	FullName          string             `json:"full_name"`
	Email             string             `json:"email"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	CreatedAt         pgtype.Timestamptz `json:"created_at"`
}

type UsersTodo struct {
	ID             int64 `json:"id"`
	UserID         int64 `json:"user_id"`
	TodosID        int64 `json:"todos_id"`
	HasPermissions bool  `json:"has_permissions"`
}
