package test

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ types.ErrorLog = (*StaticErrorLog)(nil)

type StaticErrorLog struct{}

func (s *StaticErrorLog) SaveError(ctx context.Context, msg string) error {
	if msg != errMsg {
		return fmt.Errorf("expected %q but got %q", errMsg, msg)
	}
	return nil
}
