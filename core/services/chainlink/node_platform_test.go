package chainlink_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	commonv1 "github.com/smartcontractkit/chainlink-protos/node-platform/common/v1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

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
