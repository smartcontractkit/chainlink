package changeset_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestAddNodes(t *testing.T) {
	t.Parallel()

	type input struct {
		te                 test.TestEnv
		CreateNodeRequests map[string]changeset.CreateNodeRequest
		MCMSConfig         *changeset.MCMSConfig
	}
	type testCase struct {
		name     string
		input    input
		checkErr func(t *testing.T, useMCMS bool, err error)
	}
	// run the same tests for both mcms and non-mcms
	var mcmsConfigs = []*changeset.MCMSConfig{nil, {MinDuration: 0}}
	for _, mc := range mcmsConfigs {
		prefix := "no mcms"
		if mc != nil {
			prefix = "with mcms"
		}

		t.Run(prefix, func(t *testing.T) {
			te := test.SetupTestEnv(t, test.TestConfig{
				WFDonConfig:     test.DonConfig{N: 4},
				AssetDonConfig:  test.DonConfig{N: 4},
				WriterDonConfig: test.DonConfig{N: 4},
				NumChains:       1,
				UseMCMS:         mc != nil,
			})

			var cases = []testCase{
				{
					name: "error - unregistered nop",
					input: input{
						te: te,
						CreateNodeRequests: map[string]changeset.CreateNodeRequest{
							"test-node": {
								NOPIdentity: changeset.NOPIdentity{
									RegistrationID: math.MaxUint32,
								},
								Signer:              [32]byte{0: 4},
								EncryptionPublicKey: [32]byte{0: 23},
								P2PID:               testPeerID(t, "test-peer-id"),
							},
						},
						MCMSConfig: mc,
					},
					checkErr: func(t *testing.T, useMCMS bool, err error) {
						require.Error(t, err)
						if useMCMS {
							assert.ErrorContains(t, err, "underlying transaction reverted")
						} else {
							assert.ErrorContains(t, err, "NodeOperatorDoesNotExist")
						}
					},
				},
				{
					name: "error - unregistered capability",
					input: input{
						te: te,
						CreateNodeRequests: map[string]changeset.CreateNodeRequest{
							"test-node": {
								NOPIdentity: changeset.NOPIdentity{
									RegistrationID: 1,
								},
								Signer:              [32]byte{0: 4},
								EncryptionPublicKey: [32]byte{0: 23},
								P2PID:               testPeerID(t, "test-peer-id"),
								CapabilityIdentities: changeset.CapabilityIdentities{
									{RegistrationID: test32byte(t, "no exist")},
								},
							},
						},
						MCMSConfig: mc,
					},
					checkErr: func(t *testing.T, useMCMS bool, err error) {
						require.Error(t, err)
						if useMCMS {
							assert.ErrorContains(t, err, "underlying transaction reverted")
						} else {
							assert.ErrorContains(t, err, "InvalidNodeCapabilities")
						}
					},
				},
				{
					name: "add one node",
					input: input{
						te: te,
						CreateNodeRequests: map[string]changeset.CreateNodeRequest{
							"test-node": {
								NOPIdentity: changeset.NOPIdentity{
									RegistrationID: te.Nops()[0].NodeOperatorId,
								},
								Signer:              [32]byte{0: 4},
								EncryptionPublicKey: [32]byte{0: 23},
								P2PID:               testPeerID(t, "test-peer-id"),
								CapabilityIdentities: changeset.CapabilityIdentities{
									{RegistrationID: te.CapabilityInfos()[0].HashedId},
								},
							},
						},
						MCMSConfig: mc,
					},
				},
				{
					name: "add two nodes",
					input: input{
						te: te,
						CreateNodeRequests: map[string]changeset.CreateNodeRequest{
							"test-node": {
								NOPIdentity: changeset.NOPIdentity{
									RegistrationID: te.Nops()[0].NodeOperatorId,
								},
								Signer:              test32byte(t, "2 nodes signer 1"),
								EncryptionPublicKey: test32byte(t, "2 nodes enc key 1"),
								P2PID:               testPeerID(t, "2 nodes p2p 1"),
								CapabilityIdentities: changeset.CapabilityIdentities{
									{RegistrationID: te.CapabilityInfos()[0].HashedId},
								},
							},
							"test-node-2": {
								NOPIdentity: changeset.NOPIdentity{
									RegistrationID: te.Nops()[0].NodeOperatorId,
								},
								Signer:              test32byte(t, "2 nodes signer 2"),
								EncryptionPublicKey: test32byte(t, "2 nodes enc key 2"),
								P2PID:               testPeerID(t, "2 nodes p2p 2"),
								CapabilityIdentities: changeset.CapabilityIdentities{
									{RegistrationID: te.CapabilityInfos()[0].HashedId},
								},
							},
						},
						MCMSConfig: mc,
					},
				},

				{
					name: "error - deduplicate p2p",
					input: input{
						te: te,
						CreateNodeRequests: map[string]changeset.CreateNodeRequest{
							"test-node": {
								NOPIdentity: changeset.NOPIdentity{
									RegistrationID: te.Nops()[0].NodeOperatorId,
								},
								Signer:              test32byte(t, "signer p2p 1"),
								EncryptionPublicKey: test32byte(t, "enc key"),
								P2PID:               testPeerID(t, "dupe p2p key"),
								CapabilityIdentities: changeset.CapabilityIdentities{
									{RegistrationID: te.CapabilityInfos()[0].HashedId},
								},
							},
							"test-node-2": {
								NOPIdentity: changeset.NOPIdentity{
									RegistrationID: te.Nops()[0].NodeOperatorId,
								},
								Signer:              test32byte(t, "signer p2p 2"),
								EncryptionPublicKey: test32byte(t, "another enc key"),
								P2PID:               testPeerID(t, "dupe p2p key"),
								CapabilityIdentities: changeset.CapabilityIdentities{
									{RegistrationID: te.CapabilityInfos()[0].HashedId},
								},
							},
						},
						MCMSConfig: mc,
					},
					checkErr: func(t *testing.T, useMCMS bool, err error) {
						require.Error(t, err)
						// this error is the same independent of mcms; ie in creating the proposal not exec'ing it
						assert.ErrorContains(t, err, "duplicate p2pid")
					},
				},
				{
					name: "error - deduplicate signer",
					input: input{
						te: te,
						CreateNodeRequests: map[string]changeset.CreateNodeRequest{
							"test-node": {
								NOPIdentity: changeset.NOPIdentity{
									RegistrationID: te.Nops()[0].NodeOperatorId,
								},
								Signer:              test32byte(t, "dupe signer"),
								EncryptionPublicKey: test32byte(t, "enc key"),
								P2PID:               testPeerID(t, "p2p key a"),
								CapabilityIdentities: changeset.CapabilityIdentities{
									{RegistrationID: te.CapabilityInfos()[0].HashedId},
								},
							},
							"test-node-2": {
								NOPIdentity: changeset.NOPIdentity{
									RegistrationID: te.Nops()[0].NodeOperatorId,
								},
								Signer:              test32byte(t, "dupe signer"),
								EncryptionPublicKey: test32byte(t, "another enc key"),
								P2PID:               testPeerID(t, "p2p key b"),
								CapabilityIdentities: changeset.CapabilityIdentities{
									{RegistrationID: te.CapabilityInfos()[0].HashedId},
								},
							},
						},
						MCMSConfig: mc,
					},
					checkErr: func(t *testing.T, useMCMS bool, err error) {
						require.Error(t, err)
						// contract error if two nodes have the same signer
						if useMCMS {
							assert.ErrorContains(t, err, "underlying transaction reverted")
						} else {
							assert.ErrorContains(t, err, "InvalidNodeSigner")
						}
					},
				},
			}

			for _, tc := range cases {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					req := &changeset.AddNodesRequest{
						RegistryChainSel:   tc.input.te.RegistrySelector,
						CreateNodeRequests: tc.input.CreateNodeRequests,
						MCMSConfig:         tc.input.MCMSConfig,
					}
					r, err := changeset.AddNodes(tc.input.te.Env, req)
					if err != nil && tc.checkErr == nil {
						t.Errorf("non nil err from Add Node %v but no checkErr func defined", err)
					}
					useMCMS := req.MCMSConfig != nil
					if !useMCMS {
						if tc.checkErr != nil {
							tc.checkErr(t, useMCMS, err)
							return
						}
					} else {
						// when using mcms there are two kinds of errors:
						// those from creating the proposal and those executing the proposal
						// if we have a non-nil err here, its from creating the proposal
						// so check it and do not proceed to applying the proposal
						if err != nil {
							tc.checkErr(t, useMCMS, err)
							return
						}
						require.NotNil(t, r.Proposals) //nolint:staticcheck //SA1019 ignoring deprecated field for compatibility; we don't have tools to generate the new field
						require.Len(t, r.Proposals, 1) //nolint:staticcheck //SA1019 ignoring deprecated field for compatibility; we don't have tools to generate the new field
						applyErr := applyProposal(t, tc.input.te, []commonchangeset.ChangesetApplication{
							{
								Changeset: commonchangeset.WrapChangeSet(changeset.AddNodes),
								Config:    req,
							},
						})
						if tc.checkErr != nil {
							tc.checkErr(t, useMCMS, applyErr)
							return
						}
					}
					// whether or not we are using mcms, we expect new nodes to be add to the registry based on the input
					var expectedNewNodes []kcr.CapabilitiesRegistryNodeParams
					dedupedInputs := make(map[[32]byte]kcr.CapabilitiesRegistryNodeParams)
					for _, cr := range tc.input.CreateNodeRequests {
						p, err := cr.Resolve(tc.input.te.CapabilitiesRegistry())
						require.NoError(t, err)
						dedupedInputs[p.P2pId] = p
					}
					for _, p := range dedupedInputs {
						expectedNewNodes = append(expectedNewNodes, p)
					}
					assertNodesExist(t, tc.input.te.CapabilitiesRegistry(), expectedNewNodes...)
				})
			}
		})
	}
}

func assertNodesExist(t *testing.T, registry *kcr.CapabilitiesRegistry, nodes ...kcr.CapabilitiesRegistryNodeParams) {
	for i, node := range nodes {
		got, err := registry.GetNode(nil, node.P2pId)
		require.NoError(t, err) // careful here: the err is rpc, contract return empty info if it doesn't find the p2p as opposed to non-exist err.
		assert.Equal(t, node.EncryptionPublicKey, got.EncryptionPublicKey, "mismatch node encryption public key node %d", i)
		assert.Equal(t, node.Signer, got.Signer, "mismatch node signer node %d", i)
		assert.Equal(t, node.NodeOperatorId, got.NodeOperatorId, "mismatch node operator id node %d", i)
		assert.Equal(t, node.HashedCapabilityIds, got.HashedCapabilityIds, "mismatch node hashed capability ids node %d", i)
		assert.Equal(t, node.P2pId, got.P2pId, "mismatch node p2p id node %d", i)
	}
}

func applyProposal(t *testing.T, te test.TestEnv, applicable []commonchangeset.ChangesetApplication) error {
	// now apply the changeset such that the proposal is signed and execed
	contracts := te.ContractSets()[te.RegistrySelector]
	timelockContracts := map[uint64]*proposalutils.TimelockExecutionContracts{
		te.RegistrySelector: {
			Timelock:  contracts.Timelock,
			CallProxy: contracts.CallProxy,
		},
	}
	_, err := commonchangeset.ApplyChangesets(t, te.Env, timelockContracts, applicable)
	return err
}

func testPeerID(t *testing.T, s string) p2pkey.PeerID {
	var out [32]byte
	b := []byte(s)
	copy(out[:], b)
	return p2pkey.PeerID(out)
}

func test32byte(t *testing.T, s string) [32]byte {
	var out [32]byte
	b := []byte(s)
	copy(out[:], b)
	return out
}
