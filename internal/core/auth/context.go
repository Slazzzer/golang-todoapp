package core_auth

import (
	"context"
	"fmt"

	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
)

type contextKey struct{}

func ContextWithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, contextKey{}, userID)
}

func UserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(contextKey{}).(int)
	if !ok || userID <= 0 {
		return 0, fmt.Errorf("user id not found in context: %w", core_errors.ErrUnauthorized)
	}

	return userID, nil
}

func EnsureSelfAccess(requesterID, targetUserID int) error {
	if requesterID != targetUserID {
		return fmt.Errorf("access to user is forbidden: %w", core_errors.ErrForbidden)
	}

	return nil
}

func EnsureTaskOwner(requesterID, authorUserID int) error {
	if requesterID != authorUserID {
		return fmt.Errorf("access to task is forbidden: %w", core_errors.ErrForbidden)
	}

	return nil
}
