package v1_6

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
)

type OffRampParams struct {
	GasForCallExactCheck                    uint16
	PermissionLessExecutionThresholdSeconds uint32
	MessageInterceptor                      common.Address
}

func (c OffRampParams) Validate(ignoreGasForCallExactCheck bool) error {
	if !ignoreGasForCallExactCheck && c.GasForCallExactCheck == 0 {
		return errors.New("GasForCallExactCheck is 0")
	}
	if c.PermissionLessExecutionThresholdSeconds == 0 {
		return errors.New("PermissionLessExecutionThresholdSeconds is 0")
	}
	return nil
}

func DefaultOffRampParams() OffRampParams {
	return OffRampParams{
		GasForCallExactCheck:                    uint16(5000),
		PermissionLessExecutionThresholdSeconds: uint32(globals.PermissionLessExecutionThreshold.Seconds()),
	}
}
