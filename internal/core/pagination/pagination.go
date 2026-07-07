package core_pagination

import (
	"fmt"

	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
)

const (
	DefaultLimit = 50
	MaxLimit     = 100
)

// Resolve применяет дефолтный limit и проверяет границы пагинации.
func Resolve(limit, offset *int) (resolvedLimit int, resolvedOffset int, err error) {
	resolvedLimit = DefaultLimit
	if limit != nil {
		if *limit <= 0 {
			return 0, 0, fmt.Errorf(
				"limit must be greater than 0: %w",
				core_errors.ErrInvalidArgument,
			)
		}
		if *limit > MaxLimit {
			return 0, 0, fmt.Errorf(
				"limit must not exceed %d: %w",
				MaxLimit,
				core_errors.ErrInvalidArgument,
			)
		}
		resolvedLimit = *limit
	}

	if offset != nil {
		if *offset < 0 {
			return 0, 0, fmt.Errorf(
				"offset must be non-negative: %w",
				core_errors.ErrInvalidArgument,
			)
		}
		resolvedOffset = *offset
	}

	return resolvedLimit, resolvedOffset, nil
}
