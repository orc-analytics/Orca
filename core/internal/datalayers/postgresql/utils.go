package postgresql

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Tool to check for pg error
func errorToPgError(err error) *pgconn.PgError {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return pgErr
		}
	}
	return nil
}
