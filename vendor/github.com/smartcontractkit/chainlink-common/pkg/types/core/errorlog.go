package core

import "context"

type ErrorLog interface {
	SaveError(ctx context.Context, msg string) error
}
