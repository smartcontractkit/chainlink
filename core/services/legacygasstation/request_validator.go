package legacygasstation

import (
	"math/big"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func validateSendTransactionRequest(req types.SendTransactionRequest) error {
	if utils.IsEmptyAddress(req.From) {
		return errors.New("empty from address")
	}
	if utils.IsEmptyAddress(req.Target) {
		return errors.New("empty target address")
	}
	if utils.IsEmptyAddress(req.Receiver) {
		return errors.New("empty receiver address")
	}
	if req.TargetName == "" {
		return errors.New("empty target_name")
	}
	if req.Version == "" {
		return errors.New("empty version")
	}
	if req.Nonce == nil {
		return errors.New("empty nonce")
	}
	if req.Amount == nil || req.Amount.Cmp(big.NewInt(0)) == 0 {
		return errors.New("empty amount")
	}
	if req.SourceChainID == 0 {
		return errors.New("invalid chain_id")
	}
	if req.DestinationChainID == 0 {
		return errors.New("invalid destination_chain_id")
	}
	if req.Signature == nil {
		return errors.New("empty signature")
	}
	if req.ValidUntilTime == nil {
		return errors.New("empty valid_until_time")
	}
	return nil
}
