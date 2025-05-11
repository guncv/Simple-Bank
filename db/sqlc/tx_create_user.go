package db

import (
	"context"
)

// CreateUserTxParams contains the input parameters of the create user transaction
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user Users) error
}

// CreateUserTxResult is the result of the create user transaction
type CreateUserTxResult struct {
	User Users `json:"user"`
}

// CreateUserTx performs a create user transaction
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create user record
		result.User, err = q.CreateUser(ctx, CreateUserParams{
			Username:       arg.Username,
			HashedPassword: arg.HashedPassword,
			FullName:       arg.FullName,
			Email:          arg.Email,
		})
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
