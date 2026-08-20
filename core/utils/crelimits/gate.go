package crelimits

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

type ErrorLogger interface {
	Errorw(msg string, keysAndValues ...any)
}

// GateOpen reports whether gate is open.
//
// Prefer this over calling [limits.GateLimiter.Limit] directly: Limit exists to hand
// out a raw value and only records the limit gauge, whereas AllowErr (used here) is the
// enforcement method and records the limiter's usage/denied metrics. GateOpen keeps
// Limit's ergonomics — a closed gate is (false, nil), so only a genuine evaluation
// failure returns a non-nil error and callers that fail closed can keep doing so.
func GateOpen(ctx context.Context, gate limits.GateLimiter) (bool, error) {
	err := gate.AllowErr(ctx)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, limits.ErrorNotAllowed{}):
		return false, nil
	default:
		return false, err
	}
}

// GateAllows is [GateOpen] for callers that treat an unevaluatable gate as closed.
// The evaluation failure is logged at error level and reported as false.
//
// Only use this where "closed" is the safe outcome. If a gate read failure must abort
// the operation instead (a fail-closed gate whose closed state *permits* something),
// use GateOpen and handle the error explicitly.
func GateAllows(ctx context.Context, lggr ErrorLogger, gate limits.GateLimiter, gateName string) bool {
	open, err := GateOpen(ctx, gate)
	if err != nil {
		lggr.Errorw("unexpected error evaluating CRE gate", "gate", gateName, "error", err)
		return false
	}
	return open
}
