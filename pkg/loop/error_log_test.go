package loop_test

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
)

var _ loop.ErrorLog = (*staticErrorLog)(nil)

type staticErrorLog struct{}

func (s *staticErrorLog) SaveError(ctx context.Context, msg string) error {
	if msg != errMsg {
		return fmt.Errorf("expected %q but got %q", errMsg, msg)
	}
	return nil
}
