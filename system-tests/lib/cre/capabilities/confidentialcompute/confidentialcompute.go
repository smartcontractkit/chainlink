// Package confidentialcompute registers a confidential compute capability (e.g.
// confidential-workflows, confidential-http) with a CRE DON: it proposes the standardcapabilities job
// that runs the capability binary on each worker node, and writes the enclave
// list into the capability's on-chain registry config.
//
// The confidential relay handler reads that registry config to discover which
// enclaves to route requests to, so the enclave list must be supplied before
// the environment starts.
package confidentialcompute

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/box"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	cctypes "github.com/smartcontractkit/chainlink-confidential-compute/types"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

// apiKey is the API key the capability presents to the enclaves. Enclaves in
// tests are started with a matching key; a real key is not needed, but using a
// non-empty one keeps the encrypt/decrypt path exercised.
const apiKey = "foobar"

var jobTemplate = `
type = "standardcapabilities"
schemaVersion = 1
externalJobID = "%s"
forwardingAllowed = false
command = "%s"
name = "%s"
config = %s
`

// jobsDelivered guards against the job spec function being invoked more than
// once per capability name for a single environment.
var jobsDelivered = make(map[string]bool)

// ResetDeliveryState clears the jobsDelivered guard so job specs can be
// re-delivered when a new CRE environment is created (e.g. across subtests).
func ResetDeliveryState() {
	jobsDelivered = make(map[string]bool)
}

// New returns an InstallableCapability for a confidential compute capability.
// name is both the DON capability flag and the registered LabelledName (e.g.
// "confidential-http"); binaryName is the capability binary the node runs.
// Pass a nil enclaves slice to register the capability with an empty enclave
// list, which is enough to satisfy config validation for capabilities the test
// does not exercise.
func New(name, version, binaryName string, enclaves []cctypes.Enclave) (*capabilities.Capability, error) {
	return capabilities.New( //nolint:staticcheck // SA1019 mirrors existing capability registrations
		name,
		capabilities.WithJobSpecFn(jobSpec(name, binaryName)),
		capabilities.WithCapabilityRegistryV2ConfigFn(registryConfigFn(name, version, enclaves)),
	)
}

func jobSpec(name string, binaryName string) cre.JobSpecFn {
	return func(input *cre.JobSpecInput) (cre.DonJobs, error) {
		if jobsDelivered[name] {
			return nil, nil
		}
		jobsDelivered[name] = true

		donJobs := make(cre.DonJobs, 0)
		for _, don := range input.Dons.List() {
			if !don.HasFlag(name) {
				continue
			}

			workerNodes, wErr := don.Workers()
			if wErr != nil {
				return nil, errors.Wrap(wErr, "failed to find worker nodes")
			}

			encryptedAPIKeys := make([]string, 0, len(workerNodes))
			for _, workerNode := range workerNodes {
				publicKey, kErr := workflowEncryptionKey(workerNode)
				if kErr != nil {
					return nil, kErr
				}

				ctxt, sErr := box.SealAnonymous(nil, []byte(apiKey), &publicKey, rand.Reader)
				if sErr != nil {
					return nil, errors.Wrap(sErr, "failed to seal API key")
				}
				encryptedAPIKeys = append(encryptedAPIKeys, hex.EncodeToString(ctxt))
			}

			for _, workerNode := range workerNodes {
				// Keep liveness detection aggressive in e2e so failover traffic starts
				// only after each node has had a chance to observe a dead enclave.
				config := map[string]any{
					"InsecureSkipTLSVerify":  true,
					"EncryptedAPIKeys":       strings.Join(encryptedAPIKeys, ","),
					"EnableCache":            true,
					"EnableProactiveRefresh": true,
					"MaxRetries":             3,
					"RetryBackoffSeconds":    5,
				}
				configBytes, mErr := json.Marshal(config)
				if mErr != nil {
					return nil, errors.Wrap(mErr, "failed to marshal capability config")
				}
				donJobs = append(donJobs, &jobv1.ProposeJobRequest{
					NodeId: workerNode.JobDistributorDetails.NodeID,
					Spec:   fmt.Sprintf(jobTemplate, uuid.NewString(), binaryName, name, fmt.Sprintf("'%s'", string(configBytes))),
				})
			}
		}

		return donJobs, nil
	}
}

// workflowEncryptionKey reads a node's workflow public encryption key, which is
// used to seal the capability's API key so it is not stored in plaintext.
func workflowEncryptionKey(workerNode *cre.Node) ([32]byte, error) {
	var publicKey [32]byte

	apiClient := workerNode.Clients.RestClient.APIClient
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiClient.BaseURL+"/v2/keys/workflow", nil)
	if err != nil {
		return publicKey, errors.Wrap(err, "failed to create request to get workflow keys")
	}
	req.AddCookie(apiClient.Cookies[0])

	resp, err := apiClient.GetClient().Do(req)
	if err != nil {
		return publicKey, errors.Wrap(err, "failed to send request to get workflow keys")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return publicKey, fmt.Errorf("expected 200 OK from get workflow keys request, got %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return publicKey, errors.Wrap(err, "failed to read response body from get workflow keys request")
	}

	var workflowKeysResp struct {
		Data []struct {
			Attributes struct {
				PublicKey string `json:"publicKey"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &workflowKeysResp); err != nil {
		return publicKey, errors.Wrap(err, "failed to unmarshal workflow keys response")
	}
	if len(workflowKeysResp.Data) == 0 {
		return publicKey, errors.New("no workflow keys found in response")
	}

	publicKeyBytes, err := hex.DecodeString(workflowKeysResp.Data[0].Attributes.PublicKey)
	if err != nil {
		return publicKey, errors.Wrap(err, "failed to decode public key hex")
	}
	if len(publicKeyBytes) != len(publicKey) {
		return publicKey, fmt.Errorf("expected public key to be %d bytes, got %d", len(publicKey), len(publicKeyBytes))
	}
	copy(publicKey[:], publicKeyBytes)

	return publicKey, nil
}

// registryConfigFn writes the enclave list into the capability's on-chain
// registry config, which is how the confidential relay handler discovers the
// enclaves it may route to.
func registryConfigFn(name string, version string, enclaves []cctypes.Enclave) cre.CapabilityRegistryConfigFn {
	return func(donFlags []string, _ *cre.NodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
		if !flags.HasFlag(donFlags, name) {
			return nil, nil
		}

		wrappedConfig, err := values.WrapMap(cctypes.EnclavesList{Enclaves: enclaves})
		if err != nil {
			return nil, errors.Wrap(err, "failed to wrap enclave list config")
		}

		return []keystone_changeset.DONCapabilityWithConfig{
			{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   name,
					Version:        version,
					CapabilityType: 1, // ACTION
				},
				Config: &capabilitiespb.CapabilityConfig{
					DefaultConfig: values.Proto(wrappedConfig).GetMapValue(),
					LocalOnly:     true,
				},
			},
		}, nil
	}
}
