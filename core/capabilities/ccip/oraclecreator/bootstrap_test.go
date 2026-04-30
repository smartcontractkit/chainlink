package oraclecreator

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"sort"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

func TestCalculateSyncActions(t *testing.T) {
	tests := []struct {
		name            string
		currentDigests  []cciptypes.Bytes32
		activeDigest    cciptypes.Bytes32
		candidateDigest cciptypes.Bytes32
		expectedActions []syncAction
	}{
		{
			name:            "no changes needed",
			currentDigests:  []cciptypes.Bytes32{{1}, {2}},
			activeDigest:    cciptypes.Bytes32{1},
			candidateDigest: cciptypes.Bytes32{2},
			expectedActions: nil,
		},
		{
			name:            "need to close candidate",
			currentDigests:  []cciptypes.Bytes32{{1}, {2}},
			activeDigest:    cciptypes.Bytes32{1},
			candidateDigest: cciptypes.Bytes32{}, // empty
			expectedActions: []syncAction{
				{actionType: ActionClose, endpointConfigDigest: cciptypes.Bytes32{2}},
			},
		},
		{
			name:            "need to create candidate",
			currentDigests:  []cciptypes.Bytes32{{1}},
			activeDigest:    cciptypes.Bytes32{1},
			candidateDigest: cciptypes.Bytes32{2},
			expectedActions: []syncAction{
				{actionType: ActionCreate, endpointConfigDigest: cciptypes.Bytes32{2}},
			},
		},
		{
			name:            "both configs empty",
			currentDigests:  []cciptypes.Bytes32{{1}, {2}},
			activeDigest:    cciptypes.Bytes32{},
			candidateDigest: cciptypes.Bytes32{},
			expectedActions: []syncAction{
				{actionType: ActionClose, endpointConfigDigest: cciptypes.Bytes32{1}},
				{actionType: ActionClose, endpointConfigDigest: cciptypes.Bytes32{2}},
			},
		},
		{
			name:            "replace both configs",
			currentDigests:  []cciptypes.Bytes32{{1}, {2}},
			activeDigest:    cciptypes.Bytes32{3},
			candidateDigest: cciptypes.Bytes32{4},
			expectedActions: []syncAction{
				{actionType: ActionClose, endpointConfigDigest: cciptypes.Bytes32{1}},
				{actionType: ActionClose, endpointConfigDigest: cciptypes.Bytes32{2}},
				{actionType: ActionCreate, endpointConfigDigest: cciptypes.Bytes32{3}},
				{actionType: ActionCreate, endpointConfigDigest: cciptypes.Bytes32{4}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commitConfigDigest := cciptypes.Bytes32{1}

			currentDigests := make([]cciptypes.Bytes32, len(tt.currentDigests))
			for i, d := range tt.currentDigests {
				currentDigests[i] = writePrefix(ocr2types.ConfigDigestPrefixCCIPMultiRoleRMNCombo,
					sha256.Sum256(append(commitConfigDigest[:], d[:]...)))
			}

			actions := calculateSyncActions(
				commitConfigDigest,
				currentDigests,
				tt.activeDigest,
				tt.candidateDigest,
			)

			require.Len(t, actions, len(tt.expectedActions))

			// Sort both slices to ensure consistent comparison
			sort.Slice(actions, func(i, j int) bool {
				if actions[i].actionType != actions[j].actionType {
					return actions[i].actionType < actions[j].actionType
				}
				return bytes.Compare(actions[i].endpointConfigDigest[:], actions[j].endpointConfigDigest[:]) < 0
			})
			sort.Slice(tt.expectedActions, func(i, j int) bool {
				if tt.expectedActions[i].actionType != tt.expectedActions[j].actionType {
					return tt.expectedActions[i].actionType < tt.expectedActions[j].actionType
				}
				return bytes.Compare(tt.expectedActions[i].endpointConfigDigest[:], tt.expectedActions[j].endpointConfigDigest[:]) < 0
			})

			for i := range actions {
				require.Equal(t, tt.expectedActions[i].actionType, actions[i].actionType)

				expEndpointConfigDigest := writePrefix(ocr2types.ConfigDigestPrefixCCIPMultiRoleRMNCombo,
					sha256.Sum256(append(commitConfigDigest[:], tt.expectedActions[i].endpointConfigDigest[:]...)))

				require.Equal(t, expEndpointConfigDigest, actions[i].endpointConfigDigest)
			}
		})
	}
}

func TestBootstrapOracleCreatorPopulateAddressCodecRegistryWithProviderCodecs(t *testing.T) {
	t.Parallel()

	tonSelector, ok := chainsel.TonChainIdToChainSelector()[-3]
	require.True(t, ok)

	relayID := commontypes.NewRelayID(chainsel.FamilyTon, "-3")
	provider := &fakeBootstrapCCIPProvider{
		codec: cciptypes.Codec{
			ChainSpecificAddressCodec: bootstrapTestAddressCodec{prefix: "provider"},
		},
	}
	relayer := &fakeBootstrapRelayer{provider: provider}
	addressCodec := ccipcommon.NewAddressCodec(nil)
	creator := &bootstrapOracleCreator{
		lggr:         logger.Test(t),
		relayers:     map[commontypes.RelayID]loop.Relayer{relayID: relayer},
		transmitters: map[commontypes.RelayID][]string{relayID: {"ton-transmitter"}},
	}
	ccipProviderSupported := map[string]bool{chainsel.FamilyTon: true}

	closers, err := creator.populateAddressCodecRegistryWithProviderCodecs(t.Context(), relayID, cctypes.OCR3ConfigWithMeta{
		Config: ccipreaderpkg.OCR3Config{
			PluginType:     uint8(cctypes.PluginTypeCCIPCommit),
			ChainSelector:  cciptypes.ChainSelector(tonSelector),
			OfframpAddress: []byte{0x01, 0x02},
		},
	}, &addressCodec, ccipProviderSupported)

	require.NoError(t, err)
	require.Len(t, closers, 1)
	require.Equal(t, 1, relayer.calls)
	require.Equal(t, cciptypes.PluginTypeCCIPCommit, relayer.args[0].PluginType)
	require.Equal(t, cciptypes.UnknownEncodedAddress("ton-transmitter"), relayer.args[0].TransmitterAddress)
	require.True(t, addressCodec.HasCodec(chainsel.FamilyTon))

	transmitter, err := addressCodec.TransmitterBytesToString([]byte{0x42}, cciptypes.ChainSelector(tonSelector))
	require.NoError(t, err)
	require.Equal(t, "provider-transmitter-42", transmitter)
}

func TestBootstrapOracleCreatorPopulateAddressCodecRegistryWithProviderCodecsSkipsWhenDestinationCodecExists(t *testing.T) {
	t.Parallel()

	evmSelector := cciptypes.ChainSelector(chainsel.ETHEREUM_MAINNET.Selector)
	relayer := &fakeBootstrapRelayer{}
	addressCodec := ccipcommon.NewAddressCodec(map[string]ccipcommon.ChainSpecificAddressCodec{
		chainsel.FamilyEVM: bootstrapTestAddressCodec{prefix: "evm"},
	})
	creator := &bootstrapOracleCreator{
		lggr:     logger.Test(t),
		relayers: map[commontypes.RelayID]loop.Relayer{commontypes.NewRelayID(chainsel.FamilySolana, "unused"): relayer},
	}
	ccipProviderSupported := map[string]bool{chainsel.FamilySolana: true}

	closers, err := creator.populateAddressCodecRegistryWithProviderCodecs(t.Context(), commontypes.NewRelayID(chainsel.FamilyEVM, "1"), cctypes.OCR3ConfigWithMeta{
		Config: ccipreaderpkg.OCR3Config{
			PluginType:    uint8(cctypes.PluginTypeCCIPCommit),
			ChainSelector: evmSelector,
		},
	}, &addressCodec, ccipProviderSupported)

	require.NoError(t, err)
	require.Empty(t, closers)
	require.Zero(t, relayer.calls)

	transmitter, err := addressCodec.TransmitterBytesToString([]byte{0x01}, evmSelector)
	require.NoError(t, err)
	require.Equal(t, "evm-transmitter-01", transmitter)
}

type bootstrapTestAddressCodec struct {
	prefix string
}

func (c bootstrapTestAddressCodec) AddressBytesToString(addr []byte) (string, error) {
	return fmt.Sprintf("%s-address-%x", c.prefix, addr), nil
}

func (c bootstrapTestAddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	return []byte(fmt.Sprintf("%s-%s", c.prefix, addr)), nil
}

func (c bootstrapTestAddressCodec) OracleIDAsAddressBytes(oracleID uint8) ([]byte, error) {
	return []byte{oracleID}, nil
}

func (c bootstrapTestAddressCodec) TransmitterBytesToString(addr []byte) (string, error) {
	return fmt.Sprintf("%s-transmitter-%x", c.prefix, addr), nil
}

type fakeBootstrapRelayer struct {
	loop.Relayer

	provider commontypes.CCIPProvider
	args     []commontypes.CCIPProviderArgs
	calls    int
	err      error
}

func (r *fakeBootstrapRelayer) NewCCIPProvider(ctx context.Context, args commontypes.CCIPProviderArgs) (commontypes.CCIPProvider, error) {
	r.calls++
	r.args = append(r.args, args)
	return r.provider, r.err
}

type fakeBootstrapCCIPProvider struct {
	commontypes.CCIPProvider

	codec cciptypes.Codec
}

func (p *fakeBootstrapCCIPProvider) Codec() cciptypes.Codec {
	return p.codec
}

func (p *fakeBootstrapCCIPProvider) Close() error {
	return nil
}
