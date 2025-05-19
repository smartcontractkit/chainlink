package oraclecreator

import (
	"errors"
	"fmt"
	"io"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

// wrappedOracle is a wrapper for cctypes.CCIPOracle that allows custom actions on Oracle shutdown.
type wrappedOracle struct {
	baseOracle cctypes.CCIPOracle

	// closableResources will be closed after calling baseOracle.Close()
	closableResources []io.Closer
}

func newWrappedOracle(baseOracle cctypes.CCIPOracle, closableResources []io.Closer) cctypes.CCIPOracle {
	return &wrappedOracle{
		baseOracle:        baseOracle,
		closableResources: closableResources,
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

	errs = append(errs, services.CloseAll(o.closableResources...))

	return errors.Join(errs...)
}
