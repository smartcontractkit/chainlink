package keeper

import "github.com/pkg/errors"

var (
	ErrContractCallFailure = errors.New("failure in calling contract")
)
