package oraclecreator

import (
	"fmt"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"go.uber.org/multierr"
)

// wrappedOracle is a wrapper for cctypes.CCIPOracle that allows custom actions on Oracle shutdown.
type wrappedOracle struct {
	baseOracle cctypes.CCIPOracle

	// afterCloseFunc will run after calling baseOracle.Close()
	afterCloseFunc func() error
}

func newWrappedOracle(baseOracle cctypes.CCIPOracle, afterCloseFunc func() error) cctypes.CCIPOracle {
	return &wrappedOracle{
		baseOracle:     baseOracle,
		afterCloseFunc: afterCloseFunc,
	}
}

func (o *wrappedOracle) Start() error {
	return o.baseOracle.Start()
}

func (o *wrappedOracle) Close() error {
	errs := make([]error, 0)

	if err := o.baseOracle.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close base oracle: %w", err))
	}

	if o.afterCloseFunc != nil {
		if err := o.afterCloseFunc(); err != nil {
			errs = append(errs, fmt.Errorf("after close func: %w", err))
		}
	}

	return multierr.Combine(errs...)
}
