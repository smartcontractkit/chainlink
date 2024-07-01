package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

type NoContractExistsError struct {
	Address common.Address
}

func (e NoContractExistsError) Error() string {
	return fmt.Sprintf("contract does not exist at address: %s", e.Address)
}

type methodBinding struct {
	lp                   logpoller.LogPoller
	address              common.Address
	contractName         string
	method               string
	client               evmclient.Client
	codec                commontypes.Codec
	bound                bool
	confirmationsMapping map[primitives.ConfidenceLevel]evmtypes.Confirmations
}

var _ readBinding = &methodBinding{}

func (m *methodBinding) SetCodec(codec commontypes.RemoteCodec) {
	m.codec = codec
}

func (m *methodBinding) Bind(ctx context.Context, binding commontypes.BoundContract) error {
	addr := common.HexToAddress(binding.Address)

	// check for contract byte code at the latest block and provided address
	byteCode, err := m.client.CodeAt(ctx, addr, nil)
	if err != nil {
		return err
	}

	if len(byteCode) == 0 {
		return NoContractExistsError{Address: addr}
	}

	m.address = addr
	m.bound = true

	return nil
}

func (m *methodBinding) Register(_ context.Context) error {
	return nil
}

func (m *methodBinding) Unregister(_ context.Context) error {
	return nil
}

func (m *methodBinding) GetLatestValue(ctx context.Context, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	if !m.bound {
		return fmt.Errorf("%w: method not bound", commontypes.ErrInvalidType)
	}

	data, err := m.codec.Encode(ctx, params, wrapItemType(m.contractName, m.method, true))
	if err != nil {
		return err
	}

	callMsg := ethereum.CallMsg{
		To:   &m.address,
		From: m.address,
		Data: data,
	}

	block, err := m.blockNumberFromConfidence(ctx, confidenceLevel)
	if err != nil {
		return err
	}

	bytes, err := m.client.CallContract(ctx, callMsg, block)
	if err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}

	return m.codec.Decode(ctx, bytes, returnVal, wrapItemType(m.contractName, m.method, false))
}

func (m *methodBinding) QueryKey(_ context.Context, _ query.KeyFilter, _ query.LimitAndSort, _ any) ([]commontypes.Sequence, error) {
	return nil, nil
}

func (m *methodBinding) blockNumberFromConfidence(ctx context.Context, confidenceLevel primitives.ConfidenceLevel) (*big.Int, error) {
	value, ok := m.confirmationsMapping[confidenceLevel]
	if !ok {
		// TODO is this ok? Maybe some things have to always be faster than Finalized?
		value = evmtypes.Finalized
	}

	lpBlock, err := m.lp.LatestBlock(ctx)
	if err != nil {
		return nil, err
	}

	if value == evmtypes.Finalized {
		return big.NewInt(lpBlock.FinalizedBlockNumber), nil
	} else if value == evmtypes.Unconfirmed {
		// this is the latest block
		return big.NewInt(lpBlock.BlockNumber), nil
	}

	return nil, fmt.Errorf("unknown confidence level: %v", confidenceLevel)
}
