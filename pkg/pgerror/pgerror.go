package pgerror

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrDuplicatedKey           = errors.New("duplicated key not allowed")
	ErrForeignKeyViolated      = errors.New("violates foreign key constraint")
	ErrUndefinedColumn         = errors.New("undefined column or parameter name")
	ErrDatabaseAlreadyExists   = errors.New("database already exists")
	ErrCheckConstraintViolated = errors.New("check constraint violated")
)

// HandlerError converts PostgreSQL errors to domain errors
func HandlerError(err error) error {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch pgError.SQLState() {
		case "23505":
			return ErrDuplicatedKey
		case "23503":
			return ErrForeignKeyViolated
		case "42703":
			return ErrUndefinedColumn
		case "42P04":
			return ErrDatabaseAlreadyExists
		case "23514":
			return ErrCheckConstraintViolated
		default:
			fmt.Printf("PostgreSQL error not detected: %v\n", err)
		}
	}

	return err
}
