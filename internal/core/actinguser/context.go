package actinguser

import (
	"context"
	"fmt"

	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
)

type Principal struct {
	UserID  int
	IsAdmin bool
}

type ctxKey struct{}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, ctxKey{}, principal)
}

func FromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(ctxKey{}).(Principal)
	return principal, ok
}

func IsAdmin(ctx context.Context) bool {
	principal, ok := FromContext(ctx)
	return ok && principal.IsAdmin
}

func WithContext(ctx context.Context, userID int) context.Context {
	return WithPrincipal(ctx, Principal{UserID: userID})
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	principal, ok := FromContext(ctx)
	if !ok || principal.IsAdmin {
		return 0, false
	}
	return principal.UserID, true
}

func Require(ctx context.Context) (int, error) {
	userID, ok := UserIDFromContext(ctx)
	if !ok {
		if IsAdmin(ctx) {
			return 0, fmt.Errorf("acting user is required for this operation: %w", core_errors.ErrForbidden)
		}
		return 0, fmt.Errorf("acting user is required: %w", core_errors.ErrUnauthorized)
	}
	return userID, nil
}

func RequirePrincipal(ctx context.Context) (Principal, error) {
	principal, ok := FromContext(ctx)
	if !ok {
		return Principal{}, fmt.Errorf("acting principal is required: %w", core_errors.ErrUnauthorized)
	}
	return principal, nil
}

func EnsureSelf(ctx context.Context, resourceOwnerID int) error {
	if IsAdmin(ctx) {
		return nil
	}

	actingID, err := Require(ctx)
	if err != nil {
		return err
	}
	if actingID != resourceOwnerID {
		return fmt.Errorf("access to another user's resource is forbidden: %w", core_errors.ErrForbidden)
	}
	return nil
}
