package evm

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestContractTransmitter(t *testing.T) {
	c := new(evmmocks.Client)
	// scanLogs = false
	digestAndEpoch, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc7760000000000000000000000000000000000000000000000000000000000000002")
	c.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(digestAndEpoch, nil)
	lggr := logger.TestLogger(t)
	contractABI, _ := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	ot := NewOCRContractTransmitter(gethcommon.Address{}, c, contractABI, nil, lggr)
	digest, epoch, err := ot.LatestConfigDigestAndEpoch(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc776", hex.EncodeToString(digest[:]))
	assert.Equal(t, uint32(2), epoch)

	// TODO test scanLogs = true with ccip
}
