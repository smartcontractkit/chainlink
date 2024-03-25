package test

import (
	"context"
	"fmt"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var ErrorLog = StaticErrorLog{errMsg: "an error"}

type StaticErrorLog struct {
	errMsg string
}

var _ testtypes.ErrorLogEvaluator = StaticErrorLog{}

func (s StaticErrorLog) SaveError(ctx context.Context, msg string) error {
	if msg != s.errMsg {
		return fmt.Errorf("expected %q but got %q", s.errMsg, msg)
	}
	return nil
}

func (s StaticErrorLog) Evaluate(ctx context.Context, other types.ErrorLog) error {
	return s.SaveError(ctx, s.errMsg)
}
