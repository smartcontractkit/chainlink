// Package confidentialcompute holds the pieces shared by confidential compute
// capabilities (e.g. confidential-workflows, confidential-http): reading the
// enclave list from configuration, sealing the capability's API key to each
// node, and building the on-chain registry entry that publishes the enclaves.
//
// The confidential relay handler reads that registry entry to discover which
// enclaves it may route to, so the list must be supplied before the environment
// starts. Features under cre/features consume these.
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

	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/box"
	"google.golang.org/protobuf/proto"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	cctypes "github.com/smartcontractkit/chainlink-confidential-compute/types"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// apiKey is the API key the capability presents to the enclaves. Enclaves in
// tests are started with a matching key; a real key is not needed, but using a
// non-empty one keeps the encrypt/decrypt path exercised.
const apiKey = "foobar"

// workflowEncryptionKey reads a node's workflow public encryption key, which is
// used to seal the capability's API key so it is not stored in plaintext.
func workflowEncryptionKey(workerNode *cre.Node) ([32]byte, error) {
	var publicKey [32]byte

	apiClient := workerNode.Clients.RestClient.APIClient
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiClient.BaseURL+"/v2/keys/workflow", nil)
	if err != nil {
		return publicKey, errors.Wrap(err, "failed to create request to get workflow keys")
	}
	if len(apiClient.Cookies) == 0 {
		return publicKey, errors.New("no session cookie available for get workflow keys request")
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
	if err = json.Unmarshal(body, &workflowKeysResp); err != nil {
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

// EnclavesConfigKey is the capability config value holding a JSON array of
// enclaves, letting a topology declare them instead of passing Go values.
const EnclavesConfigKey = "enclaves"

// EncryptedAPIKeys seals the capability's API key to each worker node's workflow
// public key, so it is never stored in plaintext. The capability is configured
// with every node's sealed copy; each node decrypts only its own.
func EncryptedAPIKeys(workerNodes []*cre.Node) ([]string, error) {
	encrypted := make([]string, 0, len(workerNodes))
	for _, workerNode := range workerNodes {
		publicKey, kErr := workflowEncryptionKey(workerNode)
		if kErr != nil {
			return nil, kErr
		}

		ctxt, sErr := box.SealAnonymous(nil, []byte(apiKey), &publicKey, rand.Reader)
		if sErr != nil {
			return nil, errors.Wrap(sErr, "failed to seal API key")
		}
		encrypted = append(encrypted, hex.EncodeToString(ctxt))
	}

	return encrypted, nil
}

// JobConfigJSON builds the capability job's config. Liveness detection is kept
// aggressive so failover traffic starts only once each node has had a chance to
// observe a dead enclave.
func JobConfigJSON(encryptedAPIKeys []string) (string, error) {
	config := map[string]any{
		"InsecureSkipTLSVerify":  true,
		"EncryptedAPIKeys":       strings.Join(encryptedAPIKeys, ","),
		"EnableCache":            true,
		"EnableProactiveRefresh": true,
		"MaxRetries":             3,
		"RetryBackoffSeconds":    5,
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal capability config")
	}

	return string(configBytes), nil
}

// RegistryCapabilityConfig builds the on-chain registry entry that publishes the
// enclave list, which is how the confidential relay handler discovers where it
// may route requests.
func RegistryCapabilityConfig(name, version string, enclaves []cctypes.Enclave) (keystone_changeset.DONCapabilityWithConfig, error) {
	wrappedConfig, err := values.WrapMap(cctypes.EnclavesList{Enclaves: enclaves})
	if err != nil {
		return keystone_changeset.DONCapabilityWithConfig{}, errors.Wrap(err, "failed to wrap enclave list config")
	}

	return keystone_changeset.DONCapabilityWithConfig{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   name,
			Version:        version,
			CapabilityType: 1, // ACTION
		},
		Config: &capabilitiespb.CapabilityConfig{
			DefaultConfig: values.Proto(wrappedConfig).GetMapValue(),
			LocalOnly:     true,
		},
	}, nil
}

// MarshalRegistryConfig encodes an enclave list as the capability's on-chain
// registry config, for publishing enclaves to a DON that is already running.
//
// A fresh CapabilityConfig is built rather than decoding and re-encoding what is
// already on-chain, so config left by a broken earlier run cannot corrupt this one.
func MarshalRegistryConfig(enclaves []cctypes.Enclave) ([]byte, error) {
	wrappedConfig, err := values.WrapMap(cctypes.EnclavesList{Enclaves: enclaves})
	if err != nil {
		return nil, errors.Wrap(err, "failed to wrap enclave list config")
	}

	// LocalOnly must match what RegistryCapabilityConfig writes at startup, or
	// updating the enclave list would silently change the capability's scope.
	encoded, err := proto.Marshal(&capabilitiespb.CapabilityConfig{
		DefaultConfig: values.Proto(wrappedConfig).GetMapValue(),
		LocalOnly:     true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal capability config")
	}

	return encoded, nil
}

// MarshalEnclaves encodes an enclave list for EnclavesConfigKey, for callers
// that discover their enclaves at runtime and hand them to the environment as
// configuration.
func MarshalEnclaves(enclaves []cctypes.Enclave) (string, error) {
	encoded, err := json.Marshal(enclaves)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal enclave list")
	}

	return string(encoded), nil
}

// EnclavesFromConfig reads the enclave list a topology declared for the named
// capability. Returns nil when the capability has no config or no enclaves key.
func EnclavesFromConfig(nodeSet *cre.NodeSet, name string) ([]cctypes.Enclave, error) {
	if nodeSet == nil {
		return nil, nil
	}

	capConfig, ok := nodeSet.CapabilityConfigs[name]
	if !ok || capConfig.Values == nil {
		return nil, nil
	}

	raw, ok := capConfig.Values[EnclavesConfigKey]
	if !ok {
		return nil, nil
	}

	encoded, ok := raw.(string)
	if !ok {
		return nil, errors.Errorf("capability %q: %q must be a JSON string, got %T", name, EnclavesConfigKey, raw)
	}

	var enclaves []cctypes.Enclave
	if err := json.Unmarshal([]byte(encoded), &enclaves); err != nil {
		return nil, errors.Wrapf(err, "capability %q: failed to parse %q", name, EnclavesConfigKey)
	}

	return enclaves, nil
}
