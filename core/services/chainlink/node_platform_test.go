package chainlink_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

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
	keystoremocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
)

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
