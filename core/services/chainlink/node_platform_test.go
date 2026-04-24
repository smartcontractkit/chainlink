package chainlink_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/csakey"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	commoncfg "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	commonv1 "github.com/smartcontractkit/chainlink-protos/node-platform/common/v1"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	keystoremocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
)

type fakeNodePlatformJobReader struct {
	jobs []job.Job
	err  error
}

func (f fakeNodePlatformJobReader) FindJobs(context.Context, int, int) ([]job.Job, int, error) {
	return f.jobs, len(f.jobs), f.err
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

func TestNodePlatformJobInfo_EmitsTransmittersFromOCR2Jobs(t *testing.T) {
	obs := beholdertest.NewObserver(t)

	servicetest.Run(t, chainlink.NewNodePlatformJobInfoService(chainlink.NodePlatformJobInfoConfig{
		Beat:         10 * time.Millisecond,
		Lggr:         logger.TestLogger(t),
		CSAPublicKey: "csa-public-key",
		JobReader: fakeNodePlatformJobReader{
			jobs: []job.Job{
				{
					OCR2OracleSpec: &job.OCR2OracleSpec{
						Relay:         "evm",
						ChainID:       "1",
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
					OCR2OracleSpec: &job.OCR2OracleSpec{
						Relay:         "evm",
						ChainID:       "10",
						TransmitterID: null.StringFrom("0x4444444444444444444444444444444444444444"),
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
	require.Equal(t, "csa-public-key", payload.CsaPublicKey)
	require.Len(t, payload.Transmitters, 2)
	require.Equal(t, "1", payload.Transmitters[0].ChainId)
	require.Equal(t, []string{
		"0x1111111111111111111111111111111111111111",
		"0x3333333333333333333333333333333333333333",
	}, payload.Transmitters[0].Addresses["ocr2"].Values)
	require.Equal(t, []string{"0x2222222222222222222222222222222222222222"}, payload.Transmitters[0].Addresses["ocr2_dual_transmission"].Values)
	require.Equal(t, "10", payload.Transmitters[1].ChainId)
	require.Equal(t, []string{"0x4444444444444444444444444444444444444444"}, payload.Transmitters[1].Addresses["ocr2"].Values)
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
