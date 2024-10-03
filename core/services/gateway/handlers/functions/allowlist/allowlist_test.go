package allowlist_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/allowlist"
	amocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/allowlist/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const (
	addr1           = "9ed925d8206a4f88a2f643b28b3035b315753cd6"
	addr2           = "ea6721ac65bced841b8ec3fc5fedea6141a0ade4"
	addr3           = "84689acc87ff22841b8ec378300da5e141a99911"
	ToSContractV100 = "Functions Terms of Service Allow List v1.0.0"
	ToSContractV110 = "Functions Terms of Service Allow List v1.1.0"
)

func TestUpdateAndCheck(t *testing.T) {
	t.Parallel()

	t.Run("OK-with_ToS_V1.0.0", func(t *testing.T) {
		client := mocks.NewClient(t)
		client.On("LatestBlockHeight", mock.Anything).Return(big.NewInt(42), nil)

		addr := common.HexToAddress("0x0000000000000000000000000000000000000020")
		typeAndVersionResponse, err := encodeTypeAndVersionResponse(ToSContractV100)
		require.NoError(t, err)

		client.On("CallContract", mock.Anything, ethereum.CallMsg{ // typeAndVersion
			To:   &addr,
			Data: hexutil.MustDecode("0x181f5a77"),
		}, mock.Anything).Return(typeAndVersionResponse, nil)

		client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(sampleEncodedAllowlist(t), nil)

		config := allowlist.OnchainAllowlistConfig{
			ContractVersion:    1,
			ContractAddress:    common.Address{},
			BlockConfirmations: 1,
		}

		orm := amocks.NewORM(t)
		orm.On("PurgeAllowedSenders", mock.Anything).Times(1).Return(nil)
		orm.On("CreateAllowedSenders", mock.Anything, []common.Address{common.HexToAddress(addr1), common.HexToAddress(addr2)}).Times(1).Return(nil)

		allowlist, err := allowlist.NewOnchainAllowlist(client, config, orm, logger.TestLogger(t))
		require.NoError(t, err)

		err = allowlist.Start(testutils.Context(t))
		require.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, allowlist.Close())
		})

		require.NoError(t, allowlist.UpdateFromContract(testutils.Context(t)))
		require.False(t, allowlist.Allow(common.Address{}))
		require.True(t, allowlist.Allow(common.HexToAddress(addr1)))
		require.True(t, allowlist.Allow(common.HexToAddress(addr2)))
		require.False(t, allowlist.Allow(common.HexToAddress(addr3)))
	})

	t.Run("OK-with_ToS_V1.1.0", func(t *testing.T) {
		client := mocks.NewClient(t)
		client.On("LatestBlockHeight", mock.Anything).Return(big.NewInt(42), nil)

		typeAndVersionResponse, err := encodeTypeAndVersionResponse(ToSContractV110)
		require.NoError(t, err)

		addr := common.HexToAddress("0x0000000000000000000000000000000000000020")
		client.On("CallContract", mock.Anything, ethereum.CallMsg{ // typeAndVersion
			To:   &addr,
			Data: hexutil.MustDecode("0x181f5a77"),
		}, mock.Anything).Return(typeAndVersionResponse, nil)

		client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(sampleEncodedAllowlist(t), nil)

		config := allowlist.OnchainAllowlistConfig{
			ContractVersion:    1,
			ContractAddress:    common.Address{},
			BlockConfirmations: 1,
		}

		orm := amocks.NewORM(t)
		orm.On("DeleteAllowedSenders", mock.Anything, []common.Address{common.HexToAddress(addr1), common.HexToAddress(addr2)}).Times(1).Return(nil)
		orm.On("CreateAllowedSenders", mock.Anything, []common.Address{common.HexToAddress(addr1), common.HexToAddress(addr2)}).Times(1).Return(nil)

		allowlist, err := allowlist.NewOnchainAllowlist(client, config, orm, logger.TestLogger(t))
		require.NoError(t, err)

		err = allowlist.Start(testutils.Context(t))
		require.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, allowlist.Close())
		})

		require.NoError(t, allowlist.UpdateFromContract(testutils.Context(t)))
		require.False(t, allowlist.Allow(common.Address{}))
		require.True(t, allowlist.Allow(common.HexToAddress(addr1)))
		require.True(t, allowlist.Allow(common.HexToAddress(addr2)))
		require.False(t, allowlist.Allow(common.HexToAddress(addr3)))
	})
}

func TestUnsupportedVersion(t *testing.T) {
	t.Parallel()

	client := mocks.NewClient(t)
	config := allowlist.OnchainAllowlistConfig{
		ContractVersion:    0,
		ContractAddress:    common.Address{},
		BlockConfirmations: 1,
	}

	orm := amocks.NewORM(t)
	_, err := allowlist.NewOnchainAllowlist(client, config, orm, logger.TestLogger(t))
	require.Error(t, err)
}

func TestUpdatePeriodically(t *testing.T) {
	t.Parallel()

	t.Run("OK-with_ToS_V1.0.0", func(t *testing.T) {
		ctx, cancel := context.WithCancel(testutils.Context(t))
		client := mocks.NewClient(t)
		client.On("LatestBlockHeight", mock.Anything).Return(big.NewInt(42), nil)

		addr := common.HexToAddress("0x0000000000000000000000000000000000000020")
		typeAndVersionResponse, err := encodeTypeAndVersionResponse(ToSContractV100)
		require.NoError(t, err)

		client.On("CallContract", mock.Anything, ethereum.CallMsg{ // typeAndVersion
			To:   &addr,
			Data: hexutil.MustDecode("0x181f5a77"),
		}, mock.Anything).Return(typeAndVersionResponse, nil)

		client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			cancel()
		}).Return(sampleEncodedAllowlist(t), nil)
		config := allowlist.OnchainAllowlistConfig{
			ContractAddress:    common.Address{},
			ContractVersion:    1,
			BlockConfirmations: 1,
			UpdateFrequencySec: 2,
			UpdateTimeoutSec:   1,
		}

		orm := amocks.NewORM(t)
		orm.On("PurgeAllowedSenders", mock.Anything).Times(1).Return(nil)
		orm.On("GetAllowedSenders", mock.Anything, uint(0), uint(1000)).Return([]common.Address{}, nil)
		orm.On("CreateAllowedSenders", mock.Anything, []common.Address{common.HexToAddress(addr1), common.HexToAddress(addr2)}).Times(1).Return(nil)

		allowlist, err := allowlist.NewOnchainAllowlist(client, config, orm, logger.TestLogger(t))
		require.NoError(t, err)

		err = allowlist.Start(ctx)
		require.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, allowlist.Close())
		})

		gomega.NewGomegaWithT(t).Eventually(func() bool {
			return allowlist.Allow(common.HexToAddress(addr1)) && !allowlist.Allow(common.HexToAddress(addr3))
		}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
	})

	t.Run("OK-with_ToS_V1.1.0", func(t *testing.T) {
		ctx, cancel := context.WithCancel(testutils.Context(t))
		client := mocks.NewClient(t)
		client.On("LatestBlockHeight", mock.Anything).Return(big.NewInt(42), nil)

		addr := common.HexToAddress("0x0000000000000000000000000000000000000020")
		typeAndVersionResponse, err := encodeTypeAndVersionResponse(ToSContractV110)
		require.NoError(t, err)

		client.On("CallContract", mock.Anything, ethereum.CallMsg{ // typeAndVersion
			To:   &addr,
			Data: hexutil.MustDecode("0x181f5a77"),
		}, mock.Anything).Return(typeAndVersionResponse, nil)

		client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			cancel()
		}).Return(sampleEncodedAllowlist(t), nil)
		config := allowlist.OnchainAllowlistConfig{
			ContractAddress:    common.Address{},
			ContractVersion:    1,
			BlockConfirmations: 1,
			UpdateFrequencySec: 2,
			UpdateTimeoutSec:   1,
		}

		orm := amocks.NewORM(t)
		orm.On("DeleteAllowedSenders", mock.Anything, []common.Address{common.HexToAddress(addr1), common.HexToAddress(addr2)}).Times(1).Return(nil)
		orm.On("GetAllowedSenders", mock.Anything, uint(0), uint(1000)).Return([]common.Address{}, nil)
		orm.On("CreateAllowedSenders", mock.Anything, []common.Address{common.HexToAddress(addr1), common.HexToAddress(addr2)}).Times(1).Return(nil)

		allowlist, err := allowlist.NewOnchainAllowlist(client, config, orm, logger.TestLogger(t))
		require.NoError(t, err)

		err = allowlist.Start(ctx)
		require.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, allowlist.Close())
		})

		gomega.NewGomegaWithT(t).Eventually(func() bool {
			return allowlist.Allow(common.HexToAddress(addr1)) && !allowlist.Allow(common.HexToAddress(addr3))
		}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
	})
}

func TestUpdateFromContract(t *testing.T) {
	t.Parallel()

	t.Run("OK-fetch_complete_list_of_allowed_senders", func(t *testing.T) {
		ctx, cancel := context.WithCancel(testutils.Context(t))
		client := mocks.NewClient(t)
		client.On("LatestBlockHeight", mock.Anything).Return(big.NewInt(42), nil)

		addr := common.HexToAddress("0x0000000000000000000000000000000000000020")
		typeAndVersionResponse, err := encodeTypeAndVersionResponse(ToSContractV100)
		require.NoError(t, err)

		client.On("CallContract", mock.Anything, ethereum.CallMsg{ // typeAndVersion
			To:   &addr,
			Data: hexutil.MustDecode("0x181f5a77"),
		}, mock.Anything).Return(typeAndVersionResponse, nil)

		client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			cancel()
		}).Return(sampleEncodedAllowlist(t), nil)
		config := allowlist.OnchainAllowlistConfig{
			ContractAddress:           common.HexToAddress(addr3),
			ContractVersion:           1,
			BlockConfirmations:        1,
			UpdateFrequencySec:        2,
			UpdateTimeoutSec:          1,
			StoredAllowlistBatchSize:  2,
			OnchainAllowlistBatchSize: 16,
			FetchingDelayInRangeSec:   0,
		}

		orm := amocks.NewORM(t)
		orm.On("PurgeAllowedSenders", mock.Anything).Times(1).Return(nil)
		orm.On("CreateAllowedSenders", mock.Anything, []common.Address{common.HexToAddress(addr1), common.HexToAddress(addr2)}).Times(1).Return(nil)

		allowlist, err := allowlist.NewOnchainAllowlist(client, config, orm, logger.TestLogger(t))
		require.NoError(t, err)

		err = allowlist.UpdateFromContract(ctx)
		require.NoError(t, err)

		gomega.NewGomegaWithT(t).Eventually(func() bool {
			return allowlist.Allow(common.HexToAddress(addr1)) && !allowlist.Allow(common.HexToAddress(addr3))
		}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
	})

	t.Run("OK-iterate_over_list_of_allowed_senders", func(t *testing.T) {
		ctx, cancel := context.WithCancel(testutils.Context(t))
		client := mocks.NewClient(t)
		client.On("LatestBlockHeight", mock.Anything).Return(big.NewInt(42), nil)

		addr := common.HexToAddress("0x0000000000000000000000000000000000000020")
		typeAndVersionResponse, err := encodeTypeAndVersionResponse(ToSContractV110)
		require.NoError(t, err)

		client.On("CallContract", mock.Anything, ethereum.CallMsg{ // typeAndVersion
			To:   &addr,
			Data: hexutil.MustDecode("0x181f5a77"),
		}, mock.Anything).Return(typeAndVersionResponse, nil)

		client.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			cancel()
		}).Return(sampleEncodedAllowlist(t), nil)
		config := allowlist.OnchainAllowlistConfig{
			ContractAddress:           common.HexToAddress(addr3),
			ContractVersion:           1,
			BlockConfirmations:        1,
			UpdateFrequencySec:        2,
			UpdateTimeoutSec:          1,
			StoredAllowlistBatchSize:  2,
			OnchainAllowlistBatchSize: 16,
			FetchingDelayInRangeSec:   0,
		}

		orm := amocks.NewORM(t)
		orm.On("DeleteAllowedSenders", mock.Anything, []common.Address{common.HexToAddress(addr1), common.HexToAddress(addr2)}).Times(2).Return(nil)
		orm.On("CreateAllowedSenders", mock.Anything, []common.Address{common.HexToAddress(addr1), common.HexToAddress(addr2)}).Times(2).Return(nil)

		allowlist, err := allowlist.NewOnchainAllowlist(client, config, orm, logger.TestLogger(t))
		require.NoError(t, err)

		err = allowlist.UpdateFromContract(ctx)
		require.NoError(t, err)

		gomega.NewGomegaWithT(t).Eventually(func() bool {
			return allowlist.Allow(common.HexToAddress(addr1)) && !allowlist.Allow(common.HexToAddress(addr3))
		}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
	})
}

func TestExtractContractVersion(t *testing.T) {
	type tc struct {
		name           string
		versionStr     string
		expectedResult string
		expectedError  *string
	}

	var errInvalidVersion = func(v string) *string {
		ev := fmt.Sprintf("version not found in string: %s", v)
		return &ev
	}

	tcs := []tc{
		{
			name:           "OK-Tos_type_and_version",
			versionStr:     "Functions Terms of Service Allow List v1.1.0",
			expectedResult: "v1.1.0",
			expectedError:  nil,
		},
		{
			name:           "OK-double_digits_minor",
			versionStr:     "Functions Terms of Service Allow List v1.20.0",
			expectedResult: "v1.20.0",
			expectedError:  nil,
		},
		{
			name:           "NOK-invalid_version",
			versionStr:     "invalid_version",
			expectedResult: "",
			expectedError:  errInvalidVersion("invalid_version"),
		},
		{
			name:           "NOK-incomplete_version",
			versionStr:     "v2.0",
			expectedResult: "",
			expectedError:  errInvalidVersion("v2.0"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actualResult, actualError := allowlist.ExtractContractVersion(tc.versionStr)
			require.Equal(t, tc.expectedResult, actualResult)

			if tc.expectedError != nil {
				require.EqualError(t, actualError, *tc.expectedError)
			} else {
				require.NoError(t, actualError)
			}
		})
	}
}

func encodeTypeAndVersionResponse(typeAndVersion string) ([]byte, error) {
	codecName := "my_codec"
	evmEncoderConfig := `[{"Name":"typeAndVersion","Type":"string"}]`
	codecConfig := types.CodecConfig{Configs: map[string]types.ChainCodecConfig{
		codecName: {TypeABI: evmEncoderConfig},
	}}
	encoder, err := codec.NewCodec(codecConfig)
	if err != nil {
		return nil, err
	}

	input := map[string]any{
		"typeAndVersion": typeAndVersion,
	}
	typeAndVersionResponse, err := encoder.Encode(context.Background(), input, codecName)
	if err != nil {
		return nil, err
	}

	return typeAndVersionResponse, nil
}

func sampleEncodedAllowlist(t *testing.T) []byte {
	abiEncodedAddresses :=
		"0000000000000000000000000000000000000000000000000000000000000020" +
			"0000000000000000000000000000000000000000000000000000000000000002" +
			"000000000000000000000000" + addr1 +
			"000000000000000000000000" + addr2
	rawData, err := hex.DecodeString(abiEncodedAddresses)
	require.NoError(t, err)
	return rawData
}
