package logpoller

import (
	"context"
	"math/big"

	avaclient "github.com/ava-labs/avalanche-rosetta/client"
	"github.com/ava-labs/coreth/interfaces"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
)

type avaClient struct {
	aec *avaclient.EthClient
}

func (a *avaClient) HeaderByNumber(ctx context.Context, number *big.Int) (*Header, error) {
	header, err := a.aec.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return &Header{
		Hash:       header.Hash(),
		ParentHash: header.ParentHash,
		Number:     header.Number,
	}, nil
}

func (a *avaClient) HeaderByHash(ctx context.Context, hash common.Hash) (*Header, error) {
	header, err := a.aec.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return &Header{
		Hash:       header.Hash(),
		ParentHash: header.ParentHash,
		Number:     header.Number,
	}, nil
}

func (a *avaClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	lgs, err := a.aec.FilterLogs(ctx, interfaces.FilterQuery{
		BlockHash: q.BlockHash,
		FromBlock: q.FromBlock,
		ToBlock:   q.ToBlock,
		Addresses: q.Addresses,
		Topics:    q.Topics,
	})
	if err != nil {
		return nil, err
	}
	// Cast to geth type
	var elgs []types.Log
	for _, lg := range lgs {
		lgs = append(lgs, lg)
	}
	return elgs, err
}

func (a *avaClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	//TODO implement me
	panic("implement me")
}

func (a *avaClient) ChainID() *big.Int {
	return big.NewInt(43113)
}

func NewAvaClient(rpcURL string) *avaClient {
	aec, err := avaclient.NewEthClient(context.Background(), rpcURL)
	if err != nil {
		panic(err)
	}
	return &avaClient{
		aec: aec,
	}
}

type ethClient struct {
	ec      *ethclient.Client
	rc      *rpc.Client
	chainID *big.Int
}

func (e *ethClient) HeaderByNumber(ctx context.Context, number *big.Int) (*Header, error) {
	h, err := e.ec.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return &Header{
		Hash:       h.Hash(),
		ParentHash: h.ParentHash,
		Number:     h.Number,
	}, nil
}

func (e *ethClient) HeaderByHash(ctx context.Context, hash common.Hash) (*Header, error) {
	h, err := e.ec.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return &Header{
		Hash:       h.Hash(),
		ParentHash: h.ParentHash,
		Number:     h.Number,
	}, nil
}

func (e *ethClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return e.ec.FilterLogs(ctx, q)
}

func (e *ethClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return e.rc.BatchCallContext(ctx, b)
}

func (e *ethClient) ChainID() *big.Int {
	return e.chainID
}

func NewEthClient(rpcURL string, chainID *big.Int) *ethClient {
	rc, err := rpc.DialHTTP(rpcURL)
	if err != nil {
		panic(err)
	}
	return &ethClient{
		ec:      ethclient.NewClient(rc),
		rc:      rc,
		chainID: chainID,
	}
}

type simClient struct {
	b *evmclient.SimulatedBackendClient
}

func (s simClient) HeaderByNumber(ctx context.Context, number *big.Int) (*Header, error) {
	h, err := s.b.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return &Header{
		Hash:       h.Hash(),
		ParentHash: h.ParentHash,
		Number:     h.Number,
	}, nil
}

func (s simClient) HeaderByHash(ctx context.Context, hash common.Hash) (*Header, error) {
	h, err := s.b.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return &Header{
		Hash:       h.Hash(),
		ParentHash: h.ParentHash,
		Number:     h.Number,
	}, nil
}

func (s simClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return s.b.FilterLogs(ctx, q)
}

func (s simClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return s.b.BatchCallContext(ctx, b)
}

func (s simClient) ChainID() *big.Int {
	return s.b.ChainID()
}

func NewEthClientFromSim(b *evmclient.SimulatedBackendClient) *simClient {
	return &simClient{
		b: b,
	}
}
