package v2

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

type SecretsFetcher interface {
	GetSecrets(ctx context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error)
}

type secretsFetcher struct {
	capRegistry core.CapabilitiesRegistry
	lggr        logger.Logger

	semaphore limits.ResourcePoolLimiter[int]

	workflowOwner         string
	workflowName          string
	workflowID            string
	phaseID               string
	workflowEncryptionKey workflowkey.Key

	metrics *monitoring.WorkflowsMetricLabeler
}

func NewSecretsFetcher(
	metrics *monitoring.WorkflowsMetricLabeler,
	capRegistry core.CapabilitiesRegistry,
	lggr logger.Logger,
	semaphore limits.ResourcePoolLimiter[int],
	workflowOwner string,
	workflowName string,
	workflowID string,
	phaseID string,
	workflowEncryptionKey workflowkey.Key,
) *secretsFetcher {
	lggr = logger.Named(lggr, "WorkflowEngine.SecretsFetcher")
	lggr = logger.With(lggr, "workflowID", workflowID, "workflowName", workflowName, "workflowOwner", workflowOwner, "phaseID", phaseID)
	return &secretsFetcher{
		capRegistry:           capRegistry,
		lggr:                  lggr,
		semaphore:             semaphore,
		workflowOwner:         workflowOwner,
		workflowName:          workflowName,
		phaseID:               phaseID,
		workflowEncryptionKey: workflowEncryptionKey,
		metrics:               metrics,
	}
}

func keyFor(owner, namespace, id string) string {
	return fmt.Sprintf("%s::%s::%s", owner, namespace, id)
}

func (s *secretsFetcher) GetSecrets(ctx context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	ctx = contexts.WithCRE(ctx, contexts.CRE{
		Owner:    s.workflowOwner,
		Workflow: s.workflowName,
	})
	start := time.Now()
	resp, err := func() ([]*sdkpb.SecretResponse, error) {
		free, err := s.semaphore.Wait(ctx, 1)
		if err != nil {
			return nil, err
		}
		defer free()
		return s.getSecretsForBatch(ctx, request)
	}()

	s.metrics.With(
		"workflowOwner", s.workflowOwner,
		"workflowName", s.workflowName,
		"success", strconv.FormatBool(err == nil),
	).RecordGetSecretsDuration(ctx, time.Since(start).Milliseconds())

	return resp, err
}

func sha(strs ...string) string {
	h := sha256.New()
	for _, s := range strs {
		h.Write([]byte(s))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func normalizeOwner(owner string) (string, error) {
	if len(owner) < 40 {
		return "", errors.New("invalid owner address: too short")
	}

	if owner[:2] != "0x" {
		owner = "0x" + owner
	}

	return common.HexToAddress(owner).Hex(), nil
}

func (s *secretsFetcher) getSecretsForBatch(ctx context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	vaultCap, err := s.capRegistry.GetExecutable(ctx, vault.CapabilityID)
	if err != nil {
		return nil, errors.New("failed to get vault capability: " + err.Error())
	}

	vaultCapInfo, err := vaultCap.Info(ctx)
	if err != nil {
		return nil, errors.New("failed to get vault capability Info: " + err.Error())
	}

	var donID uint32
	if vaultCapInfo.IsLocal {
		// If the capability is local, we can use the local node's DON ID.
		localNode, err2 := s.capRegistry.LocalNode(ctx)
		if err2 != nil {
			return nil, errors.New("failed to get local node from registry: " + err2.Error())
		}
		donID = localNode.WorkflowDON.ID
	} else {
		don := vaultCapInfo.DON
		if don == nil {
			return nil, errors.New("vault capability is not associated with any DON")
		}
		donID = don.ID
	}
	vaultCapConfig, err := s.capRegistry.ConfigForCapability(ctx, vault.CapabilityID, donID)
	if err != nil {
		return nil, errors.New("failed to get vault capability config for donID: " + strconv.FormatUint(uint64(donID), 10) + ". Error: " + err.Error())
	}

	cfg, err := unmarshalConfig(vaultCapConfig)
	if err != nil {
		return nil, errors.New("failed to extract vault public key from capability config: " + err.Error())
	}

	encryptionKeys, err := s.getEncryptionKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption keys: %w", err)
	}
	vp := &vault.GetSecretsRequest{
		Requests: make([]*vault.SecretRequest, 0),
	}

	owner, err := normalizeOwner(s.workflowOwner)
	if err != nil {
		return nil, fmt.Errorf("could not normalize workflowOwner: %w", err)
	}

	logKeys := make([]string, 0, len(request.Requests))
	for _, r := range request.Requests {
		logKeys = append(logKeys, keyFor(owner, r.Namespace, r.Id))
		vp.Requests = append(vp.Requests, &vault.SecretRequest{
			Id: &vault.SecretIdentifier{
				Key:       r.Id,
				Namespace: r.Namespace,
				Owner:     owner,
			},
			EncryptionKeys: encryptionKeys,
		})
	}

	anypbReq, err := anypb.New(vp)
	if err != nil {
		return nil, fmt.Errorf("failed to convert vault request to any: %w", err)
	}

	lggr := logger.With(s.lggr, "requestedKeys", logKeys)
	lggr.Debug("fetching secrets...")

	capabilityResponse, err := vaultCap.Execute(ctx, capabilities.CapabilityRequest{
		Payload:      anypbReq,
		Method:       vault.MethodGetSecrets,
		CapabilityId: vault.CapabilityID,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          s.workflowID,
			WorkflowOwner:       s.workflowOwner,
			WorkflowName:        s.workflowName,
			WorkflowExecutionID: sha(s.phaseID, strconv.FormatInt(int64(request.CallbackId), 10)),
			ReferenceID:         strconv.FormatInt(int64(request.CallbackId), 10),
		},
	})
	if err != nil {
		lggr.Errorw("failed to fetch secrets", "err", err)
		return nil, fmt.Errorf("failed to execute vault.GetSecrets: %w", err)
	}

	lggr.Debug("successfully fetched secrets from vault capability")

	batchedVaultResponse := &vault.GetSecretsResponse{}
	err = capabilityResponse.Payload.UnmarshalTo(batchedVaultResponse)
	if err != nil {
		lggr.Errorw("failed to unmarshal vault payload to GetSecretsResponse", "err", err)
		return nil, fmt.Errorf("failed to unmarshal vault payload to GetSecretsResponse: %w", err)
	}

	m := map[string]*vault.SecretResponse{}
	for _, secretResponse := range batchedVaultResponse.Responses {
		key := keyFor(secretResponse.Id.Owner, secretResponse.Id.Namespace, secretResponse.Id.Key)
		m[key] = secretResponse
	}

	sdkResp := make([]*sdkpb.SecretResponse, 0, len(request.Requests))
	for _, r := range request.Requests {
		key := keyFor(owner, r.Namespace, r.Id)
		resp, ok := m[key]
		if !ok {
			errorMessage := "could not find response for the request: " + key
			errorResponse := s.wrapErrorResponse(lggr, r.Id, r.Namespace, owner, errorMessage)
			sdkResp = append(sdkResp, &errorResponse)
			continue
		}
		response := s.getSecretForSingleRequest(logger.With(lggr, "key", key), r.Id, owner, r.Namespace, cfg, resp)
		sdkResp = append(sdkResp, &response)
	}
	return sdkResp, nil
}

func (s *secretsFetcher) getSecretForSingleRequest(lggr logger.Logger, id, owner, namespace string, cfg *vaultConfig, response *vault.SecretResponse) sdkpb.SecretResponse {
	if response.GetId() != nil {
		if response.GetId().GetKey() != "" {
			id = response.GetId().GetKey()
		}
		if response.GetId().GetNamespace() != "" {
			namespace = response.GetId().GetNamespace()
		}
		if response.GetId().GetOwner() != "" {
			owner = response.GetId().GetOwner()
		}
	}
	if response.GetError() != "" {
		errorMessage := "secret request returned an error: " + response.GetError()
		return s.wrapErrorResponse(lggr, id, namespace, owner, errorMessage)
	}

	var localNodeShares []string
	workflowNodeEncryptionPublicKeyStr := s.workflowEncryptionKey.PublicKeyString()
	for _, share := range response.GetData().GetEncryptedDecryptionKeyShares() {
		if share.EncryptionKey == workflowNodeEncryptionPublicKeyStr {
			localNodeShares = share.Shares
		}
	}
	if len(localNodeShares) == 0 {
		errorMessage := "no shares found for this node's encryption key: " + workflowNodeEncryptionPublicKeyStr
		return s.wrapErrorResponse(lggr, id, namespace, owner, errorMessage)
	}

	encryptedSecretBytes, err := hex.DecodeString(response.GetData().GetEncryptedValue())
	if err != nil {
		errorMessage := "failed to decode the secret string: " + err.Error()
		return s.wrapErrorResponse(lggr, id, namespace, owner, errorMessage)
	}

	secret, err := s.decryptSecret(lggr, encryptedSecretBytes, localNodeShares, cfg)
	if err != nil {
		errorMessage := "failed to decrypt secret: " + err.Error()
		return s.wrapErrorResponse(lggr, id, namespace, owner, errorMessage)
	}

	return sdkpb.SecretResponse{
		Response: &sdkpb.SecretResponse_Secret{
			Secret: &sdkpb.Secret{
				Id:        response.GetId().GetKey(),
				Namespace: response.GetId().GetNamespace(),
				Owner:     response.GetId().GetOwner(),
				Value:     secret,
			},
		},
	}
}

func (s *secretsFetcher) wrapErrorResponse(lggr logger.Logger, id, namespace, owner, errorMessage string) sdkpb.SecretResponse {
	lggr.Debugw(errorMessage)
	return sdkpb.SecretResponse{
		Response: &sdkpb.SecretResponse_Error{
			Error: &sdkpb.SecretError{
				Id:        id,
				Namespace: namespace,
				Owner:     owner,
				Error:     errorMessage,
			},
		},
	}
}

func (s *secretsFetcher) decryptSecret(lggr logger.Logger, encryptedSecretBytes []byte, encryptedDecryptionShares []string, cfg *vaultConfig) (string, error) {
	lggr.Debug("decrypting secret...")

	cipherText := &tdh2easy.Ciphertext{}
	errOuter := cipherText.UnmarshalVerify(encryptedSecretBytes, cfg.VaultPublicKey)
	if errOuter != nil {
		return "", errors.New("failed to unmarshal encrypted secret: " + errOuter.Error())
	}

	decryptionShares := make([]*tdh2easy.DecryptionShare, 0, len(encryptedDecryptionShares))
	for i, encryptedDecryptionShare := range encryptedDecryptionShares {
		encryptedDecryptionShareBytes, err := hex.DecodeString(encryptedDecryptionShare)
		if err != nil {
			lggr.Debugw("failed to hex decode the encryptedDecryptionShare", "index", i)
			continue
		}
		decryptionShareBytes, err := s.workflowEncryptionKey.Decrypt(encryptedDecryptionShareBytes)
		if err != nil {
			lggr.Debugw("failed to decrypt the encryptedDecryptionShare", "index", i)
			continue
		}
		decryptionShare := &tdh2easy.DecryptionShare{}
		err = decryptionShare.Unmarshal(decryptionShareBytes)
		if err != nil {
			lggr.Debugw("failed to unmarshal decryption share", "index", i)
			continue
		}
		err = tdh2easy.VerifyShare(cipherText, cfg.VaultPublicKey, decryptionShare)
		if err != nil {
			lggr.Debugw("failed to verify decryption share", "index", i)
			continue
		}
		decryptionShares = append(decryptionShares, decryptionShare)
	}
	lggr.Debugw("decryption shares collected", "count", len(decryptionShares), "expected", len(encryptedDecryptionShares), "threshold", cfg.Threshold)

	if len(decryptionShares) < cfg.Threshold {
		return "", fmt.Errorf("not enough decryption shares to decrypt the secret: have %d, need at least %d", len(encryptedDecryptionShares), cfg.Threshold)
	}

	// Note that the last parameter 'n' to tdh2easy.Aggregate() isn't verified by the library at all.
	// Thus, the len(encryptedDecryptionShares) set below is just an optional hint for memory allocation.
	decryptedSecret, err := tdh2easy.Aggregate(cipherText, decryptionShares, len(encryptedDecryptionShares))
	if err != nil {
		return "", errors.New("failed to aggregate decryption shares: " + err.Error())
	}
	return string(decryptedSecret), nil
}

func (s *secretsFetcher) getEncryptionKeys(ctx context.Context) ([]string, error) {
	s.lggr.Debug("Fetching encryption keys...")
	myNode, err := s.capRegistry.LocalNode(ctx)
	if err != nil {
		return nil, errors.New("failed to get local node from registry" + err.Error())
	}

	encryptionKeys := make([]string, 0, len(myNode.WorkflowDON.Members))
	for _, peerID := range myNode.WorkflowDON.Members {
		peerNode, err := s.capRegistry.NodeByPeerID(ctx, peerID)
		if err != nil {
			return nil, errors.New("failed to get node info for peerID: " + peerID.String() + " - " + err.Error())
		}
		encryptionKeys = append(encryptionKeys, hex.EncodeToString(peerNode.EncryptionPublicKey[:]))
	}
	// Sort the encryption keys to ensure consistent ordering across all nodes.
	sort.Strings(encryptionKeys)
	return encryptionKeys, nil
}

type VaultCapabilityRegistryConfig struct {
	VaultPublicKey string
	Threshold      int
}

type vaultConfig struct {
	VaultPublicKey *tdh2easy.PublicKey
	Threshold      int
}

func unmarshalConfig(config capabilities.CapabilityConfiguration) (*vaultConfig, error) {
	cfg := &VaultCapabilityRegistryConfig{}
	err := config.DefaultConfig.UnwrapTo(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unwrap capability config: %w", err)
	}

	if cfg.Threshold <= 0 {
		return nil, errors.New("invalid Threshold in the capability config")
	}

	if cfg.VaultPublicKey == "" {
		return nil, errors.New("VaultPublicKey is not provided in the capability config")
	}

	pkBytes, err := hex.DecodeString(cfg.VaultPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode vault public key from registry: %w", err)
	}

	pk := tdh2easy.PublicKey{}
	err = pk.Unmarshal(pkBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to construct vault public key from raw bytes: %w", err)
	}

	return &vaultConfig{
		Threshold:      cfg.Threshold,
		VaultPublicKey: &pk,
	}, nil
}
