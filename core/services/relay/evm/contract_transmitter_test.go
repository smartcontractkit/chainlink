package evm

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"

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
	lggr := logger.TestLogger(t)
	c := new(evmmocks.Client)
	// scanLogs = false
	digestAndEpochDontScanLogs, _ := hex.DecodeString(
		"0000000000000000000000000000000000000000000000000000000000000000" + // false
			"000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc776" + // config digest
			"0000000000000000000000000000000000000000000000000000000000000002") // epoch
	c.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(digestAndEpochDontScanLogs, nil).Once()
	contractABI, _ := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	ot := NewOCRContractTransmitter(gethcommon.Address{}, c, contractABI, nil, lggr)
	digest, epoch, err := ot.LatestConfigDigestAndEpoch(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc776", hex.EncodeToString(digest[:]))
	assert.Equal(t, uint32(2), epoch)

	// scanLogs = true
	digestAndEpochScanLogs, _ := hex.DecodeString(
		"0000000000000000000000000000000000000000000000000000000000000001" + // false
			"000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc776" + // config digest
			"0000000000000000000000000000000000000000000000000000000000000002") // epoch
	c.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(digestAndEpochScanLogs, nil).Once()
	// We expect it will call for latest config details to get a lower bound on the log search.
	latestConfigDetails, _ := hex.DecodeString(
		"0000000000000000000000000000000000000000000000000000000000000001" + // config count
			"0000000000000000000000000000000000000000000000000000000000000002" + //  block num
			"000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc776") // digest
	c.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(latestConfigDetails, nil)

	transmitted1, _ := hex.DecodeString(
		"000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc776" + // config digest
			"0000000000000000000000000000000000000000000000000000000000000001") // epoch
	transmitted2, _ := hex.DecodeString(
		"000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc777" + // config digest
			"0000000000000000000000000000000000000000000000000000000000000002") // epoch
	c.On("FilterLogs", mock.Anything, mock.Anything).Return(
		[]types.Log{
			{
				Data: transmitted1,
			},
			{
				Data: transmitted2,
			},
		}, nil)
	digest, epoch, err = ot.LatestConfigDigestAndEpoch(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "000130da6b9315bd59af6b0a3f5463c0d0a39e92eaa34cbcbdbace7b3bfcc777", hex.EncodeToString(digest[:]))
	assert.Equal(t, uint32(2), epoch)
	c.AssertExpectations(t)
}
