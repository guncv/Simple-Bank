package db

import (
	"context"
	"database/sql"
)

// CreateUserTxParams contains the input parameters of the create user transaction
type VerifyEmailTxParams struct {
	EmailID    int64  `json:"email_id"`
	SecretCode string `json:"secret_code"`
}

// CreateUserTxResult is the result of the create user transaction
type VerifyEmailTxResult struct {
	User        Users        `json:"user"`
	VerifyEmail VerifyEmails `json:"verify_email"`
}

// CreateUserTx performs a create user transaction
func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create verify email record
		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         int32(arg.EmailID),
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
