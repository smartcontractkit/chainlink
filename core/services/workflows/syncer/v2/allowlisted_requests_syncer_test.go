package v2

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/testutils"
)

func Test_AllowlistedRequestsSyncer(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	backendTH := testutils.NewEVMBackendTH(t)
	donFamily := "A"

	wfRegistryAddr, _, wfRegistryC, err := workflow_registry_wrapper_v2.DeployWorkflowRegistry(backendTH.ContractsOwner, backendTH.Backend.Client())
	backendTH.Backend.Commit()
	require.NoError(t, err)

	updateAllowedDONsV2(t, backendTH, wfRegistryC, []string{donFamily})
	updateAuthorizedAddressV2(t, backendTH, wfRegistryC, backendTH.ContractsOwner.From, donFamily)

	activeAllowlistedRequestsCount := int(MaxResultsPerQuery + 1)
	expiryTimestamp := time.Now().Add(24 * time.Hour)
	for i := 0; i < activeAllowlistedRequestsCount; i++ {
		createSecretsRequestParams, marshalErr := json.Marshal(vaultcommon.CreateSecretsRequest{
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{
					Id: &vaultcommon.SecretIdentifier{
						Key:       strconv.Itoa(i),
						Namespace: "active",
					},
					EncryptedValue: "encrypted-value",
				},
			},
		})
		require.NoError(t, marshalErr)

		allowlistRequest(t, backendTH, wfRegistryC, allowlistRequestParams{
			Request: jsonrpc.Request[json.RawMessage]{
				Method: vaulttypes.MethodSecretsCreate,
				Params: (*json.RawMessage)(&createSecretsRequestParams),
			},
			Owner:           backendTH.ContractsOwner.From,
			ExpiryTimestamp: expiryTimestamp,
		})
	}

	contractReaderFn := func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
		return backendTH.NewContractReader(ctx, t, bytes)
	}
	allowlistedSyncer := NewAllowlistedRequestsSyncer(lggr, contractReaderFn, wfRegistryAddr.Hex())

	servicetest.Run(t, allowlistedSyncer)

	require.Eventually(t, func() bool {
		return len(allowlistedSyncer.GetAllowlistedRequests(t.Context())) == activeAllowlistedRequestsCount
	}, tests.WaitTimeout(t), time.Second, "synced allowlisted requests do not match expectations")
}
