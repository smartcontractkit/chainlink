package scenarios

import (
	"testing"
	"time"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/integration-tests/smoke/ccip/canonical/types"
	"github.com/stretchr/testify/require"
)

// ValidationType is the type of validation to perform in the scenario.
type ValidationType int

const (
	ValidationTypeNone ValidationType = iota
	ValidationTypeCommitOnly
	ValidationTypeExec
)

func NewMessagingToCCIPReceiverScenario(t *testing.T, validationType ValidationType) types.Scenario {
	return &messagingToCCIPReceiverScenario{validationType: validationType}
}

type messagingToCCIPReceiverScenario struct {
	validationType ValidationType
}

func (s *messagingToCCIPReceiverScenario) Run(t *testing.T, ctx types.ExecContext) {
	// sanity checks
	require.GreaterOrEqual(t, len(ctx.Sources()), 1)
	require.NotNil(t, ctx.Dest())

	// determine the source and chain families so we can appropriately build the message
	// note that we select the first source, this scenario doesn't support multiple sources
	// anyways.
	source := ctx.Sources()[0]
	dest := ctx.Dest()

	receiver := dest.CCIPReceiver()
	nativeFeeToken := source.NativeFeeToken()

	extraArgs, err := dest.GetExtraArgs(receiver, source.ChainFamily())
	require.NoError(t, err)

	components := types.MessageComponents{
		Receiver:     receiver,
		Data:         []byte("hello ccip receiver testing scenario"),
		FeeToken:     nativeFeeToken,
		ExtraArgs:    extraArgs,
		TokenAmounts: []types.TokenAmount{},
	}

	t.Logf("receiver length: %d, receiver: %x", len(receiver), receiver)

	msg, err := source.BuildMessage(components)
	require.NoError(t, err)

	t.Logf("source: %d, dest: %d", source.ChainSelector(), dest.ChainSelector())

	// send the message on the source chain to the destination chain.
	sendEvent := testhelpers.TestSendRequest(
		t,
		ctx.Env(),
		ctx.OnchainState(),
		source.ChainSelector(),
		dest.ChainSelector(),
		false,
		msg,
	)

	time.Sleep(30 * time.Second)
	ctx.ReplayLogs(t, map[uint64]uint64{
		source.ChainSelector(): 0,
	})

	switch s.validationType {
	case ValidationTypeNone:
		return
	case ValidationTypeCommitOnly:
		dest.ValidateCommit(t, source.ChainSelector(), ccipocr3.SeqNumRange{
			ccipocr3.SeqNum(sendEvent.SequenceNumber),
			ccipocr3.SeqNum(sendEvent.SequenceNumber),
		})
	case ValidationTypeExec:
		dest.ValidateCommit(t, source.ChainSelector(), ccipocr3.SeqNumRange{
			ccipocr3.SeqNum(sendEvent.SequenceNumber),
			ccipocr3.SeqNum(sendEvent.SequenceNumber),
		})
		dest.ValidateExec(t, source.ChainSelector(), []uint64{sendEvent.SequenceNumber})
	}
}

func NewMessagingToEOAScenario(t *testing.T, validationType ValidationType) types.Scenario {
	return &messagingToEOAScenario{validationType: validationType}
}

type messagingToEOAScenario struct {
	validationType ValidationType
}

func (s *messagingToEOAScenario) Run(t *testing.T, ctx types.ExecContext) {
	// sanity checks
	require.GreaterOrEqual(t, len(ctx.Sources()), 1)
	require.NotNil(t, ctx.Dest())
	// Solana doesn't support sending messages to EOAs
	require.NotEqual(t, chain_selectors.FamilySolana, ctx.Dest().ChainFamily())
	// Aptos doesn't support sending messages to EOAs
	require.NotEqual(t, chain_selectors.FamilyAptos, ctx.Dest().ChainFamily())

	// determine the source and chain families so we can appropriately build the message
	// note that we select the first source, this scenario doesn't support multiple sources
	// anyways.
	source := ctx.Sources()[0]
	dest := ctx.Dest()

	receiver := dest.RandomReceiver()
	nativeFeeToken := source.NativeFeeToken()

	components := types.MessageComponents{
		Receiver:     receiver,
		Data:         []byte("hello eoa testing scenario"),
		FeeToken:     nativeFeeToken,
		ExtraArgs:    []byte{},
		TokenAmounts: []types.TokenAmount{},
	}

	msg, err := source.BuildMessage(components)
	require.NoError(t, err)

	// send the message on the source chain to the destination chain.
	sendEvent := testhelpers.TestSendRequest(
		t,
		ctx.Env(),
		ctx.OnchainState(),
		source.ChainSelector(),
		dest.ChainSelector(),
		false,
		msg,
	)

	time.Sleep(30 * time.Second)
	ctx.ReplayLogs(t, map[uint64]uint64{
		source.ChainSelector(): 0,
	})

	switch s.validationType {
	case ValidationTypeNone:
		return
	case ValidationTypeCommitOnly:
		dest.ValidateCommit(t, source.ChainSelector(), ccipocr3.SeqNumRange{
			ccipocr3.SeqNum(sendEvent.SequenceNumber),
			ccipocr3.SeqNum(sendEvent.SequenceNumber),
		})
	case ValidationTypeExec:
		dest.ValidateCommit(t, source.ChainSelector(), ccipocr3.SeqNumRange{
			ccipocr3.SeqNum(sendEvent.SequenceNumber),
			ccipocr3.SeqNum(sendEvent.SequenceNumber),
		})
		dest.ValidateExec(t, source.ChainSelector(), []uint64{sendEvent.SequenceNumber})
	}
}
