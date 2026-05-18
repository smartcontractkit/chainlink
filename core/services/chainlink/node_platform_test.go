package chainlink_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/csakey"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	commoncfg "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	commonv1 "github.com/smartcontractkit/chainlink-protos/node-platform/common/v1"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	keystoremocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

type fakeNodePlatformJobReader struct {
	jobs []job.Job
	err  error
}

func (f fakeNodePlatformJobReader) FindJobs(_ context.Context, offset, limit int) ([]job.Job, int, error) {
	if f.err != nil {
		return nil, 0, f.err
	}

	if offset >= len(f.jobs) {
		return nil, len(f.jobs), nil
	}

	end := offset + limit
	if end > len(f.jobs) {
		end = len(f.jobs)
	}
	return f.jobs[offset:end], len(f.jobs), nil
}

func eip55Address(raw string) evmtypes.EIP55Address {
	return evmtypes.MustEIP55Address(raw)
}

func eip55AddressPtr(raw string) *evmtypes.EIP55Address {
	address := eip55Address(raw)
	return &address
}

func TestNewNodePlatformBuildInfoConfig_UsesThreeMinuteBeat(t *testing.T) {
	csaStore := &keystoremocks.CSA{}
	keyStore := &keystoremocks.Master{}
	keyStore.EXPECT().CSA().Return(csaStore).Once()

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, _ *chainlink.Secrets) {
		c.Telemetry.HeartbeatInterval = commoncfg.MustNewDuration(5 * time.Second)
	})

	buildInfoCfg := chainlink.NewNodePlatformBuildInfoConfig(chainlink.ApplicationOpts{
		Config:   cfg,
		Logger:   logger.TestLogger(t),
		KeyStore: keyStore,
	})

	require.Equal(t, 3*time.Minute, buildInfoCfg.Beat)
	require.Same(t, csaStore, buildInfoCfg.CSAKeyStore)
}

func TestNodePlatformBuildInfo_EmitsNodeBuildInfo(t *testing.T) {
	obs := beholdertest.NewObserver(t)

	servicetest.Run(t, chainlink.NewNodePlatformBuildInfoService(chainlink.NodePlatformBuildInfoConfig{
		Beat:         10 * time.Millisecond,
		Lggr:         logger.TestLogger(t),
		CSAPublicKey: "csa-public-key",
		CommitSHA:    "commit-sha",
		DockerTag:    "docker-tag",
		VersionTag:   "version-tag",
		Version:      "1.2.3",
	}))

	require.Eventually(t, func() bool {
		return obs.Len(t, beholder.AttrKeyEntity, "common.v1.NodeBuildInfo") > 0
	}, time.Second, 10*time.Millisecond)

	msgs := obs.Messages(t, beholder.AttrKeyEntity, "common.v1.NodeBuildInfo")
	require.NotEmpty(t, msgs)

	msg := msgs[0]
	require.Equal(t, "node-platform", msg.Attrs[beholder.AttrKeyDomain])
	require.Equal(t, "/node-platform/common/v1", msg.Attrs[beholder.AttrKeyDataSchema])

	var payload commonv1.NodeBuildInfo
	require.NoError(t, proto.Unmarshal(msg.Body, &payload))
	require.Equal(t, "csa-public-key", payload.CsaPublicKey)
	require.Equal(t, "commit-sha", payload.CommitSha)
	require.Equal(t, "docker-tag", payload.DockerTag)
	require.Equal(t, "version-tag", payload.VersionTag)
	require.Equal(t, "1.2.3", payload.Version)
}

func TestNodePlatformJobInfo_EmitsSubmitterAddressesFromJobFields(t *testing.T) {
	obs := beholdertest.NewObserver(t)

	servicetest.Run(t, chainlink.NewNodePlatformJobInfoService(chainlink.NodePlatformJobInfoConfig{
		Beat:         10 * time.Millisecond,
		Lggr:         logger.TestLogger(t),
		CSAPublicKey: "csa-public-key",
		JobReader: fakeNodePlatformJobReader{
			jobs: []job.Job{
				{
					Type: job.OffchainReporting,
					OCROracleSpec: &job.OCROracleSpec{
						TransmitterAddress: eip55AddressPtr("0x9999999999999999999999999999999999999999"),
						EVMChainID:         sqlutil.NewI(1),
					},
				},
				{
					Type: job.OffchainReporting2,
					OCR2OracleSpec: &job.OCR2OracleSpec{
						Relay:         "evm",
						ChainID:       "1",
						PluginType:    commontypes.Median,
						TransmitterID: null.StringFrom("0x1111111111111111111111111111111111111111"),
						RelayConfig: job.JSONConfig{
							"chainID":                "1",
							"sendingKeys":            []any{"0x1111111111111111111111111111111111111111", "0x3333333333333333333333333333333333333333"},
							"enableDualTransmission": true,
							"dualTransmission": map[string]any{
								"transmitterAddress": "0x2222222222222222222222222222222222222222",
							},
						},
					},
				},
				{
					Type: job.OffchainReporting2,
					OCR2OracleSpec: &job.OCR2OracleSpec{
						Relay:         "evm",
						ChainID:       "2",
						PluginType:    commontypes.Mercury,
						TransmitterID: null.StringFrom("0x4444444444444444444444444444444444444444"),
					},
				},
				{
					Type: job.VRF,
					VRFSpec: &job.VRFSpec{
						FromAddresses: []evmtypes.EIP55Address{
							eip55Address("0x6666666666666666666666666666666666666666"),
							eip55Address("0x7777777777777777777777777777777777777777"),
						},
						EVMChainID: sqlutil.NewI(4),
					},
				},
				{
					Type: job.BlockhashStore,
					BlockhashStoreSpec: &job.BlockhashStoreSpec{
						FromAddresses: []evmtypes.EIP55Address{eip55Address("0x8888888888888888888888888888888888888888")},
						EVMChainID:    sqlutil.NewI(5),
					},
				},
				{
					Type: job.BlockHeaderFeeder,
					BlockHeaderFeederSpec: &job.BlockHeaderFeederSpec{
						FromAddresses: []evmtypes.EIP55Address{eip55Address("0x9999999999999999999999999999999999999999")},
						EVMChainID:    sqlutil.NewI(6),
					},
				},
				{
					Type: job.LegacyGasStationServer,
					LegacyGasStationServerSpec: &job.LegacyGasStationServerSpec{
						FromAddresses: []evmtypes.EIP55Address{eip55Address("0x1010101010101010101010101010101010101010")},
						EVMChainID:    sqlutil.NewI(7),
					},
				},
				{
					Type: job.StandardCapabilities,
					StandardCapabilitiesSpec: &job.StandardCapabilitiesSpec{
						OracleFactory: job.OracleFactoryConfig{
							Enabled:       true,
							ChainID:       "8",
							TransmitterID: "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
						},
					},
				},
			},
		},
	}))

	require.Eventually(t, func() bool {
		return obs.Len(t, beholder.AttrKeyEntity, "common.v1.NodeJobInfo") > 0
	}, time.Second, 10*time.Millisecond)

	msg := obs.Messages(t, beholder.AttrKeyEntity, "common.v1.NodeJobInfo")[0]
	require.Equal(t, "node-platform", msg.Attrs[beholder.AttrKeyDomain])
	require.Equal(t, "/node-platform/common/v1", msg.Attrs[beholder.AttrKeyDataSchema])

	var payload commonv1.NodeJobInfo
	require.NoError(t, proto.Unmarshal(msg.Body, &payload))
	expected := &commonv1.NodeJobInfo{
		CsaPublicKey: "csa-public-key",
		SubmitterAddresses: []*commonv1.NodeSubmitterAddress{
			{
				ChainId:   "1",
				JobType:   "offchainreporting",
				FieldPath: "transmitterAddress",
				Addresses: []string{"0x9999999999999999999999999999999999999999"},
			},
			{
				ChainId:    "1",
				JobType:    "offchainreporting2",
				PluginType: "median",
				FieldPath:  "relayConfig.dualTransmission.transmitterAddress",
				Addresses:  []string{"0x2222222222222222222222222222222222222222"},
			},
			{
				ChainId:    "1",
				JobType:    "offchainreporting2",
				PluginType: "median",
				FieldPath:  "relayConfig.sendingKeys",
				Addresses: []string{
					"0x1111111111111111111111111111111111111111",
					"0x3333333333333333333333333333333333333333",
				},
			},
			{
				ChainId:    "1",
				JobType:    "offchainreporting2",
				PluginType: "median",
				FieldPath:  "transmitterID",
				Addresses:  []string{"0x1111111111111111111111111111111111111111"},
			},
			{
				ChainId:   "4",
				JobType:   "vrf",
				FieldPath: "fromAddresses",
				Addresses: []string{
					"0x6666666666666666666666666666666666666666",
					"0x7777777777777777777777777777777777777777",
				},
			},
			{
				ChainId:   "5",
				JobType:   "blockhashstore",
				FieldPath: "fromAddresses",
				Addresses: []string{"0x8888888888888888888888888888888888888888"},
			},
			{
				ChainId:   "6",
				JobType:   "blockheaderfeeder",
				FieldPath: "fromAddresses",
				Addresses: []string{"0x9999999999999999999999999999999999999999"},
			},
			{
				ChainId:   "7",
				JobType:   "legacygasstationserver",
				FieldPath: "fromAddresses",
				Addresses: []string{"0x1010101010101010101010101010101010101010"},
			},
			{
				ChainId:   "8",
				JobType:   "standardcapabilities",
				FieldPath: "oracle_factory.transmitter_id",
				Addresses: []string{"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
			},
		},
	}
	require.Truef(t, proto.Equal(expected, &payload), "expected:\n%sgot:\n%s", prototext.Format(expected), prototext.Format(&payload))
}

func TestNodePlatformJobInfo_PaginatesSubmitterAddressJobs(t *testing.T) {
	obs := beholdertest.NewObserver(t)

	jobs := make([]job.Job, 1001)
	jobs[1000] = job.Job{
		Type: job.VRF,
		VRFSpec: &job.VRFSpec{
			EVMChainID: sqlutil.NewI(10),
		},
		Pipeline: pipeline.Pipeline{Tasks: []pipeline.Task{
			&pipeline.ETHTxTask{
				From: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			},
		}},
	}

	servicetest.Run(t, chainlink.NewNodePlatformJobInfoService(chainlink.NodePlatformJobInfoConfig{
		Beat:         100 * time.Millisecond,
		Lggr:         logger.TestLogger(t),
		CSAPublicKey: "csa-public-key",
		JobReader: fakeNodePlatformJobReader{
			jobs: jobs,
		},
	}))

	require.Eventually(t, func() bool {
		msgs := obs.Messages(t, beholder.AttrKeyEntity, "common.v1.NodeJobInfo")
		for _, msg := range msgs {
			var payload commonv1.NodeJobInfo
			if proto.Unmarshal(msg.Body, &payload) != nil {
				continue
			}

			for _, submitterAddress := range payload.SubmitterAddresses {
				if submitterAddress.ChainId == "10" &&
					submitterAddress.JobType == "vrf" &&
					submitterAddress.FieldPath == "observationSource.ethtx.from" &&
					len(submitterAddress.Addresses) == 1 &&
					submitterAddress.Addresses[0] == "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
					return true
				}
			}
		}
		return false
	}, time.Second, 100*time.Millisecond)
}

func TestNodePlatformBuildInfo_ResolvesCSAKeyOnStart(t *testing.T) {
	obs := beholdertest.NewObserver(t)
	csaStore := &keystoremocks.CSA{}

	csaStore.EXPECT().EnsureKey(mock.Anything).Return(nil).Once()
	csaStore.EXPECT().GetAll().Return([]csakey.KeyV2{cltest.DefaultCSAKey}, nil).Once()

	servicetest.Run(t, chainlink.NewNodePlatformBuildInfoService(chainlink.NodePlatformBuildInfoConfig{
		Beat:        10 * time.Millisecond,
		Lggr:        logger.TestLogger(t),
		CSAKeyStore: csaStore,
		CommitSHA:   "commit-sha",
		DockerTag:   "docker-tag",
		VersionTag:  "version-tag",
		Version:     "1.2.3",
	}))

	require.Eventually(t, func() bool {
		return obs.Len(t, beholder.AttrKeyEntity, "common.v1.NodeBuildInfo") > 0
	}, time.Second, 10*time.Millisecond)

	msg := obs.Messages(t, beholder.AttrKeyEntity, "common.v1.NodeBuildInfo")[0]
	var payload commonv1.NodeBuildInfo
	require.NoError(t, proto.Unmarshal(msg.Body, &payload))
	require.Equal(t, cltest.DefaultCSAKey.PublicKeyString(), payload.CsaPublicKey)
}
