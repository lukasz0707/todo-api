package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	CreateTodoTx(ctx context.Context, arg CreateTodoTxParams) (CreateTodoTxResult, error)
	CreateGroupTx(ctx context.Context, arg CreateGroupTxParams) (CreateGroupTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// exectX executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type CreateTodoTxParams struct {
	UserID   int64     `json:"user_id"`
	TodoName string    `json:"todo_name"`
	GroupID  int64     `json:"group_id"`
	Deadline time.Time `json:"deadline"`
}

type CreateTodoTxResult struct {
	Todo      Todo       `json:"todo"`
	UserGroup UsersGroup `json:"user_group"`
}

func (store *SQLStore) CreateTodoTx(ctx context.Context, arg CreateTodoTxParams) (CreateTodoTxResult, error) {
	var result CreateTodoTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.UserGroup, err = q.GetFromUsersGroups(ctx, GetFromUsersGroupsParams{
			UserID:  arg.UserID,
			GroupID: arg.GroupID,
		})
		if err != nil {
			return err
		}

		result.Todo, err = q.CreateTodo(ctx, CreateTodoParams{
			GroupID:  arg.GroupID,
			TodoName: arg.TodoName,
			Deadline: arg.Deadline,
		})
		if err != nil {
			return err
		}
		return nil

	})
	return result, err
}

type CreateGroupTxParams struct {
	UserID    int64
	GroupName string
}

type CreateGroupTxResult struct {
	Group     Group      `json:"group"`
	UserGroup UsersGroup `json:"user_group"`
}

func (store *SQLStore) CreateGroupTx(ctx context.Context, arg CreateGroupTxParams) (CreateGroupTxResult, error) {
	var result CreateGroupTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Group, err = q.CreateGroup(ctx, arg.GroupName)
		if err != nil {
			return err
		}

		result.UserGroup, err = q.AssignOwnerToGroup(ctx, AssignOwnerToGroupParams{
			UserID:  arg.UserID,
			GroupID: result.Group.ID,
		})
		if err != nil {
			return err
		}

		return nil

	})
	return result, err
}
