package oraclecreator

import (
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

type wrappedOracle struct {
	baseOracle cctypes.CCIPOracle

	// beforeCloseFunc will run before calling baseOracle.Close()
	beforeCloseFunc func() error
}

func newWrappedOracle(baseOracle cctypes.CCIPOracle, beforeCloseFunc func() error) cctypes.CCIPOracle {
	return &wrappedOracle{
		baseOracle:      baseOracle,
		beforeCloseFunc: beforeCloseFunc,
	}
}

func (o *wrappedOracle) withBeforeCloseFunc(fn func() error) *wrappedOracle {
	o.beforeCloseFunc = fn
	return o
}

func (o *wrappedOracle) Start() error {
	return o.baseOracle.Start()
}

func (o *wrappedOracle) Close() error {
	if o.beforeCloseFunc != nil {
		if err := o.beforeCloseFunc(); err != nil {
			return err
		}
	}
	return o.baseOracle.Close()
}
