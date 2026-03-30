package aptos

import (
	"fmt"
	"strings"

	aptossdk "github.com/aptos-labs/aptos-go-sdk"
	pkgerrors "github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	aptoschain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/aptos"
)

func normalizeTransmitter(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", pkgerrors.New("empty Aptos transmitter")
	}

	var addr aptossdk.AccountAddress
	if err := addr.ParseStringRelaxed(s); err != nil {
		return "", err
	}
	return addr.StringLong(), nil
}

func normalizeForwarderAddress(raw string) (string, error) {
	var addr aptossdk.AccountAddress
	if err := addr.ParseStringRelaxed(strings.TrimSpace(raw)); err != nil {
		return "", err
	}
	return addr.StringLong(), nil
}

func findAptosChainByChainID(chains []creblockchains.Blockchain, chainID uint64) (*aptoschain.Blockchain, error) {
	for _, bc := range chains {
		if bc.IsFamily(chainselectors.FamilyAptos) && bc.ChainID() == chainID {
			aptosBlockchain, ok := bc.(*aptoschain.Blockchain)
			if !ok {
				return nil, fmt.Errorf("Aptos blockchain for chain id %d has unexpected type %T", chainID, bc)
			}
			return aptosBlockchain, nil
		}
	}
	return nil, fmt.Errorf("Aptos blockchain for chain id %d not found", chainID)
}
