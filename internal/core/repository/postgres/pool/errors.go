package core_postgres_pool

import "errors"

var (
	ErrNoRows             = errors.New("no rows in result set")
	ErrViolatesForeignKey = errors.New("violates foreign key constraint")
)

func IsErrNoRows(err error) bool {
	return errors.Is(err, ErrNoRows)
}

func IsErrViolatesForeignKey(err error) bool {
	return errors.Is(err, ErrViolatesForeignKey)
}
