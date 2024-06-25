package ocr3impls_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/ocr3impls"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

func TestMultichainTransmitter(t *testing.T) {
	t.Parallel()
	// create many separate transmitters and separate chains
	numChains := 4
	unis := make(map[int]testUniverse[multichainMeta])
	for i := 0; i < numChains; i++ {
		unis[i] = newTestUniverse[multichainMeta](t, nil)
	}

	mct, err := ocr3impls.NewMultichainTransmitterOCR3[multichainMeta](
		map[commontypes.RelayID]ocr3types.ContractTransmitter[multichainMeta]{
			commontypes.NewRelayID(relay.NetworkEVM, "0"): unis[0].ocr3Transmitter,
			commontypes.NewRelayID(relay.NetworkEVM, "1"): unis[1].ocr3Transmitter,
			commontypes.NewRelayID(relay.NetworkEVM, "2"): unis[2].ocr3Transmitter,
			commontypes.NewRelayID(relay.NetworkEVM, "3"): unis[3].ocr3Transmitter,
		},
		logger.TestLogger(t),
	)
	require.NoError(t, err)

	expectedTransmitters := []string{
		ocr3impls.EncodeTransmitter(commontypes.NewRelayID(relay.NetworkEVM, "0"), ocrtypes.Account(unis[0].transmitters[0].From.String())),
		ocr3impls.EncodeTransmitter(commontypes.NewRelayID(relay.NetworkEVM, "1"), ocrtypes.Account(unis[1].transmitters[0].From.String())),
		ocr3impls.EncodeTransmitter(commontypes.NewRelayID(relay.NetworkEVM, "2"), ocrtypes.Account(unis[2].transmitters[0].From.String())),
		ocr3impls.EncodeTransmitter(commontypes.NewRelayID(relay.NetworkEVM, "3"), ocrtypes.Account(unis[3].transmitters[0].From.String())),
	}
	slices.Sort(expectedTransmitters)
	expectedFromAccount := strings.Join(expectedTransmitters, ",")
	fromAccount, err := mct.FromAccount()
	require.NoError(t, err)
	require.Equal(t, expectedFromAccount, string(fromAccount))

	// generate a report for each chain and sign it
	// note that in this test each chain has a different set of signers
	// this is okay because it's just a test
	// in actuality the same signers will be used across all chains
	var reports []ocr3types.ReportWithInfo[multichainMeta]
	for i := 0; i < numChains; i++ {
		c, err2 := unis[i].wrapper.LatestConfigDetails(nil)
		require.NoError(t, err2, "failed to get latest config details")
		report := ocr3types.ReportWithInfo[multichainMeta]{
			Info:   multichainMeta{destChainIndex: i, configDigest: c.ConfigDigest},
			Report: []byte{},
		}
		reports = append(reports, report)
	}
	seqNum := uint64(1)
	for i := range reports {
		c, err2 := unis[i].wrapper.LatestConfigDetails(nil)
		require.NoError(t, err2, "failed to get latest config details")
		attributedSigs := unis[i].SignReport(t, c.ConfigDigest, reports[i], seqNum)
		err = mct.Transmit(testutils.Context(t), c.ConfigDigest, seqNum, reports[i], attributedSigs)
		require.NoError(t, err)
		events := unis[i].TransmittedEvents(t)
		require.Len(t, events, 1)
		require.Equal(t, c.ConfigDigest, events[0].ConfigDigest, "config digest mismatch")
		require.Equal(t, seqNum, events[0].SequenceNumber, "sequence number mismatch")
		// increment sequence number so that each chain gets a unique one for this test
		seqNum++
	}
}

type multichainMeta struct {
	destChainIndex int
	configDigest   ocrtypes.ConfigDigest
}

func (m multichainMeta) GetDestinationChain() commontypes.RelayID {
	return commontypes.NewRelayID(relay.NetworkEVM, fmt.Sprintf("%d", m.destChainIndex))
}

func (m multichainMeta) GetDestinationConfigDigest() ocrtypes.ConfigDigest {
	return m.configDigest
}
