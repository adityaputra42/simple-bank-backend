package db

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	UniwueViolation     = "23505"
)

var ErrRecordNotFound = pgx.ErrNoRows
var ErrUniqueViolation = &pgconn.PgError{
	Code: UniwueViolation,
}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""

}
