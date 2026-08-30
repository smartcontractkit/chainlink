package vault

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/nacl/box"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	coreCapabilities "github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type stubOwnerResolver struct {
	mu     sync.Mutex
	owners map[string]string
}

func (s *stubOwnerResolver) ResolveOwner(_ context.Context, donID uint32, workflowID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := fmt.Sprintf("%d:%s", donID, workflowID)
	owner, ok := s.owners[key]
	if !ok {
		return "", fmt.Errorf("no workflow found for DON %d, workflow %q", donID, workflowID)
	}
	return owner, nil
}

func TestIntegration_CrossDONOwnerSpoof(t *testing.T) {
	const (
		attackerOwner = "0x000000000000000000000000000000000000a11c"
		victimOwner   = "0x00000000000000000000000000000000deadbeef"
		attackerDonID = uint32(1)
	)

	attackerPub, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	attackerEncKey := hex.EncodeToString(attackerPub[:])

	ownerResolver := &stubOwnerResolver{
		owners: map[string]string{
			fmt.Sprintf("%d:attacker-workflow", attackerDonID): attackerOwner,
		},
	}

	newVaultCapability := func(t *testing.T) *Capability {
		t.Helper()
		lggr := logger.TestLogger(t)
		clock := clockwork.NewFakeClock()
		store := requests.NewStore[*vaulttypes.Request]()
		handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](
			lggr, store, clock, 10*time.Second,
		)
		reg := coreCapabilities.NewRegistry(lggr)
		lf := limits.Factory{Settings: cresettings.DefaultGetter}
		cap, err := NewCapability(lggr, clock, 10*time.Second, handler, reg, nil, lf, newTestRequestLifecycleTracker(t))
		require.NoError(t, err)
		cap.SetOwnerResolver(ownerResolver)
		servicetest.Run(t, cap)
		return cap
	}

	buildGetSecretsRequest := func(t *testing.T, requestedOwner string) *anypb.Any {
		t.Helper()
		payload, err := anypb.New(&vault.GetSecretsRequest{
			Requests: []*vault.SecretRequest{{
				Id: &vault.SecretIdentifier{
					Owner:     requestedOwner,
					Namespace: "main",
					Key:       "api_key",
				},
				EncryptionKeys: []string{attackerEncKey},
			}},
		})
		require.NoError(t, err)
		return payload
	}

	buildCapabilityRequest := func(payload *anypb.Any, claimedOwner, execID string) capabilities.CapabilityRequest {
		return capabilities.CapabilityRequest{
			Method:  vaulttypes.MethodSecretsGet,
			Payload: payload,
			Metadata: capabilities.RequestMetadata{
				WorkflowOwner:       claimedOwner,
				WorkflowID:          "attacker-workflow",
				WorkflowExecutionID: execID,
				ReferenceID:         "ref",
				WorkflowDonID:       attackerDonID,
			},
		}
	}

	t.Run("honest_caller_rejected", func(t *testing.T) {
		cap := newVaultCapability(t)
		req := buildCapabilityRequest(
			buildGetSecretsRequest(t, victimOwner),
			attackerOwner,
			"exec-honest",
		)
		_, err := cap.Execute(t.Context(), req)
		require.ErrorContains(t, err, "does not match workflow owner")
	})

	t.Run("spoofed_owner_rejected", func(t *testing.T) {
		cap := newVaultCapability(t)
		req := buildCapabilityRequest(
			buildGetSecretsRequest(t, victimOwner),
			victimOwner,
			"exec-spoof",
		)
		_, err := cap.Execute(t.Context(), req)
		require.ErrorContains(t, err, "does not match claimed workflow owner")
	})
}
