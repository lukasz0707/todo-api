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
	UserID    int64     `json:"user_id"`
	TodoName  string    `json:"todo_name"`
	GroupName string    `json:"group_name"`
	Deadline  time.Time `json:"deadline"`
}

type CreateTodoTxResult struct {
	Todo  Todo  `json:"todo"`
	Group Group `json:"group"`
}

func (store *SQLStore) CreateTodoTx(ctx context.Context, arg CreateTodoTxParams) (CreateTodoTxResult, error) {
	var result CreateTodoTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Group, err = q.CreateGroup(ctx, CreateGroupParams{
			GroupName: arg.GroupName,
			OwnerID:   arg.UserID,
		})
		if err != nil {
			return err
		}

		result.Todo, err = q.CreateTodo(ctx, CreateTodoParams{
			GroupID:  result.Group.ID,
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
