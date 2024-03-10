package blockhashstore

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type TestCoordinator struct {
	RequestEvents     []Event
	FulfillmentEvents []Event
}

func (t *TestCoordinator) Addresses() []common.Address {
	return []common.Address{}
}

func (t *TestCoordinator) Requests(_ context.Context, fromBlock uint64, toBlock uint64) ([]Event, error) {
	var result []Event
	for _, req := range t.RequestEvents {
		if req.Block >= fromBlock && req.Block <= toBlock {
			result = append(result, req)
		}
	}
	return result, nil
}

func (t *TestCoordinator) Fulfillments(_ context.Context, fromBlock uint64) ([]Event, error) {
	var result []Event
	for _, ful := range t.FulfillmentEvents {
		if ful.Block >= fromBlock {
			result = append(result, ful)
		}
	}
	return result, nil
}

type TestBHS struct {
	Stored []uint64

	StoredEarliest bool

	// errorsStore defines which block numbers should return errors on Store.
	ErrorsStore []uint64

	// errorsIsStored defines which block numbers should return errors on IsStored.
	ErrorsIsStored []uint64
}

func (t *TestBHS) Store(_ context.Context, blockNum uint64) error {
	for _, e := range t.ErrorsStore {
		if e == blockNum {
			return errors.New("error storing")
		}
	}

	t.Stored = append(t.Stored, blockNum)
	return nil
}

func (t *TestBHS) IsTrusted() bool {
	return false
}

func (t *TestBHS) StoreTrusted(
	ctx context.Context, blockNums []uint64, blockhashes []common.Hash, recentBlock uint64, recentBlockhash common.Hash,
) error {
	return errors.New("not implemented")
}

func (t *TestBHS) IsStored(_ context.Context, blockNum uint64) (bool, error) {
	for _, e := range t.ErrorsIsStored {
		if e == blockNum {
			return false, errors.New("error checking if stored")
		}
	}

	for _, s := range t.Stored {
		if s == blockNum {
			return true, nil
		}
	}
	return false, nil
}

func (t *TestBHS) StoreEarliest(ctx context.Context) error {
	t.StoredEarliest = true
	return nil
}

type TestBatchBHS struct {
	Stored                       []uint64
	GetBlockhashesCallCounter    uint16
	StoreVerifyHeaderCallCounter uint16
	GetBlockhashesError          error
	StoreVerifyHeadersError      error
}

func (t *TestBatchBHS) GetBlockhashes(_ context.Context, blockNumbers []*big.Int) ([][32]byte, error) {
	t.GetBlockhashesCallCounter++
	if t.GetBlockhashesError != nil {
		return nil, t.GetBlockhashesError
	}
	var blockhashes [][32]byte
	for _, b := range blockNumbers {
		for _, stored := range t.Stored {
			var randomBlockhash [32]byte
			if stored == b.Uint64() {
				_, err := rand.Read(randomBlockhash[:])
				if err != nil {
					return nil, err
				}
			}
			blockhashes = append(blockhashes, randomBlockhash)
		}
	}
	return blockhashes, nil
}

func (t *TestBatchBHS) StoreVerifyHeader(ctx context.Context, blockNumbers []*big.Int, blockHeaders [][]byte, fromAddress common.Address) error {
	t.StoreVerifyHeaderCallCounter++
	if t.StoreVerifyHeadersError != nil {
		return t.StoreVerifyHeadersError
	}
	if len(blockNumbers) != len(blockHeaders) {
		return errors.Errorf("input length did not match. blockNumbers length: %d, blockHeaders length: %d", len(blockNumbers), len(blockHeaders))
	}
	for _, blockNumber := range blockNumbers {
		t.Stored = append(t.Stored, blockNumber.Uint64())
	}
	return nil
}

type TestBlockHeaderProvider struct {
}

func (p *TestBlockHeaderProvider) RlpHeadersBatch(ctx context.Context, blockRange []*big.Int) ([][]byte, error) {
	var headers [][]byte
	for range blockRange {
		var randomBytes [30]byte //random length
		_, err := rand.Read(randomBytes[:])
		if err != nil {
			return nil, err
		}
		headers = append(headers, randomBytes[:])
	}
	return headers, nil
}
