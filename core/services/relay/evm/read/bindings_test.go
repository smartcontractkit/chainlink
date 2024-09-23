package read_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/read"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/read/mocks"
)

var (
	contractName1  = "Contract1"
	methodName1    = "Method1"
	methodName2    = "Method2"
	eventName0     = "EventName0"
	filterWithSigs = logpoller.Filter{
		Name:      eventName0,
		EventSigs: []common.Hash{common.HexToHash("0x25"), common.HexToHash("0x29")},
	}
)

func TestBindingsRegistry(t *testing.T) {
	t.Parallel()

	t.Run("readers are addable, bindable, and retrievable", func(t *testing.T) {
		t.Parallel()

		mBatch := new(mocks.BatchCaller)
		mRdr := new(mocks.Reader)
		mReg := new(mocks.Registrar)

		mRdr.EXPECT().Bind(mock.Anything, mock.Anything).Return(nil)

		named := read.NewBindingsRegistry()
		named.SetBatchCaller(mBatch)

		require.NoError(t, named.AddReader(contractName1, methodName1, mRdr))

		bindings := []commontypes.BoundContract{{Address: "0x24", Name: contractName1}}
		_ = named.Bind(context.Background(), mReg, bindings)

		rdr, _, err := named.GetReader(bindings[0].ReadIdentifier(methodName1))

		require.NoError(t, err)
		require.NotNil(t, rdr)
	})

	t.Run("register all before bind", func(t *testing.T) {
		t.Parallel()

		mBatch := new(mocks.BatchCaller)
		mRdr0 := new(mocks.Reader)
		mRdr1 := new(mocks.Reader)
		mReg := new(mocks.Registrar)

		named := read.NewBindingsRegistry()
		named.SetBatchCaller(mBatch)

		// register is called once through RegisterAll and again in Bind
		mRdr0.EXPECT().Register(mock.Anything).Return(nil)
		mRdr1.EXPECT().Register(mock.Anything).Return(nil)

		mRdr0.EXPECT().Bind(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mRdr1.EXPECT().Bind(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		mReg.EXPECT().HasFilter(mock.Anything).Return(false)
		mReg.EXPECT().RegisterFilter(mock.Anything, mock.Anything).Return(nil)
		mRdr0.EXPECT().GetLatestValue(mock.Anything, common.HexToAddress("0x25"), mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mRdr0.EXPECT().GetLatestValue(mock.Anything, common.HexToAddress("0x24"), mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mRdr1.EXPECT().GetLatestValue(mock.Anything, common.HexToAddress("0x26"), mock.Anything, mock.Anything, mock.Anything).Return(nil)

		// part of the init phase of chain reader
		require.NoError(t, named.AddReader(contractName1, methodName1, mRdr0))
		require.NoError(t, named.AddReader(contractName1, methodName2, mRdr1))
		_ = named.SetFilter(contractName1, filterWithSigs)

		// run within the start phase of chain reader
		require.NoError(t, named.RegisterAll(context.Background(), mReg))

		bindings := []commontypes.BoundContract{
			{Address: "0x24", Name: contractName1},
			{Address: "0x25", Name: contractName1},
			{Address: "0x26", Name: contractName1},
		}

		// calling bind will also call register for each reader
		_ = named.Bind(context.Background(), mReg, bindings)

		rdr1, _, err := named.GetReader(bindings[0].ReadIdentifier(methodName1))
		require.NoError(t, err)

		rdr2, _, err := named.GetReader(bindings[0].ReadIdentifier(methodName2))
		require.NoError(t, err)

		require.NoError(t, rdr1.GetLatestValue(context.Background(), common.HexToAddress("0x25"), primitives.Finalized, nil, nil))
		require.NoError(t, rdr1.GetLatestValue(context.Background(), common.HexToAddress("0x24"), primitives.Finalized, nil, nil))
		require.NoError(t, rdr2.GetLatestValue(context.Background(), common.HexToAddress("0x26"), primitives.Finalized, nil, nil))

		mBatch.AssertExpectations(t)
		mRdr0.AssertExpectations(t)
		mRdr1.AssertExpectations(t)
		mReg.AssertExpectations(t)
	})
}
