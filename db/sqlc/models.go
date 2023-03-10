// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID        int64  `json:"id"`
	GroupName string `json:"group_name"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Todo struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	TodoName  string    `json:"todo_name"`
	CreatedAt time.Time `json:"created_at"`
	// oneof(todo, suspended, completed)
	Status   string    `json:"status"`
	Deadline time.Time `json:"deadline"`
}

type User struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	// oneof(user, moderator, admin)
	Role              string    `json:"role"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	IsBlocked         bool      `json:"is_blocked"`
}

type UsersGroup struct {
	ID      int64 `json:"id"`
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
	// oneof(user, owner)
	Role string `json:"role"`
}
