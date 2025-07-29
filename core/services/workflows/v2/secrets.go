package v2

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

type SecretsFetcher interface {
	GetSecrets(ctx context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error)
}

type secretsFetcher struct {
	capRegistry core.CapabilitiesRegistry
	lggr        logger.Logger

	semaphore *semaphore[[]*sdkpb.SecretResponse]

	workflowOwner string
	workflowName  string
	workflowKey   workflowkey.Key

	vaultPublicKey *tdh2easy.PublicKey

	metrics *monitoring.WorkflowsMetricLabeler
}

func NewSecretsFetcher(
	metrics *monitoring.WorkflowsMetricLabeler,
	capRegistry core.CapabilitiesRegistry,
	lggr logger.Logger,
	semaphore *semaphore[[]*sdkpb.SecretResponse],
	workflowOwner string,
	workflowName string,
	workflowKey workflowkey.Key,
	vaultPublicKey *tdh2easy.PublicKey,
) *secretsFetcher {
	return &secretsFetcher{
		capRegistry:    capRegistry,
		lggr:           logger.Named(lggr, "SecretsFetcher"),
		semaphore:      semaphore,
		workflowOwner:  workflowOwner,
		workflowName:   workflowName,
		workflowKey:    workflowKey,
		vaultPublicKey: vaultPublicKey,
		metrics:        metrics,
	}
}

func keyFor(owner, namespace, id string) string {
	return fmt.Sprintf("%s::%s::%s", owner, namespace, id)
}

func (s *secretsFetcher) GetSecrets(ctx context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	start := time.Now()
	resp, err := s.semaphore.WhenAcquired(ctx, func() ([]*sdkpb.SecretResponse, error) {
		return s.getSecrets(ctx, request)
	})

	s.metrics.With(
		"workflowOwner", s.workflowOwner,
		"workflowName", s.workflowName,
		"success", strconv.FormatBool(err == nil),
	).RecordGetSecretsDuration(ctx, time.Since(start).Milliseconds())

	return resp, err
}

func (s *secretsFetcher) getSecrets(ctx context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	vaultCap, err := s.capRegistry.GetExecutable(ctx, vault.CapabilityID)
	if err != nil {
		return nil, errors.New("failed to get vault capability: " + err.Error())
	}

	encryptionKeys, err := s.getEncryptionKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption keys: %w", err)
	}
	vp := &vault.GetSecretsRequest{
		Requests: make([]*vault.SecretRequest, 0),
	}

	logKeys := make([]string, len(request.Requests))
	for _, r := range request.Requests {
		logKeys = append(logKeys, keyFor(s.workflowOwner, r.Namespace, r.Id))
		vp.Requests = append(vp.Requests, &vault.SecretRequest{
			Id: &vault.SecretIdentifier{
				Key:       r.Id,
				Namespace: r.Namespace,
				Owner:     s.workflowOwner,
			},
			EncryptionKeys: encryptionKeys,
		})
	}

	anypbReq, err := anypb.New(vp)
	if err != nil {
		return nil, fmt.Errorf("failed to convert vault request to any: %w", err)
	}

	lggr := logger.With(s.lggr, "requestedKeys", logKeys, "owner", s.workflowOwner, "workflow", s.workflowName)
	lggr.Debug("fetching secrets...")

	resp, err := vaultCap.Execute(ctx, capabilities.CapabilityRequest{
		Payload:      anypbReq,
		Method:       vault.MethodGetSecrets,
		CapabilityId: vault.CapabilityID,
		Metadata: capabilities.RequestMetadata{
			WorkflowOwner: s.workflowOwner,
			WorkflowName:  s.workflowName,
		},
	})
	if err != nil {
		lggr.Errorw("failed to fetch secrets", "err", err)
		return nil, fmt.Errorf("failed to execute vault.GetSecrets: %w", err)
	}

	lggr.Debug("successfully fetched secrets")

	respPayload := &vault.GetSecretsResponse{}
	err = resp.Payload.UnmarshalTo(respPayload)
	if err != nil {
		lggr.Errorw("failed to unmarshal vault payload to GetSecretsResponse", "err", err)
		return nil, fmt.Errorf("failed to unmarshal vault payload to GetSecretsResponse: %w", err)
	}

	m := map[string]*vault.SecretResponse{}
	for _, s := range respPayload.Responses {
		key := keyFor(s.Id.Owner, s.Id.Namespace, s.Id.Key)
		m[key] = s
	}

	localNode, err := s.capRegistry.LocalNode(ctx)
	if err != nil {
		lggr.Errorw("failed to get local node from registry" + err.Error())
		return nil, errors.New("failed to get local node from registry: " + err.Error())
	}

	sdkResp := make([]*sdkpb.SecretResponse, 0, len(request.Requests))
	for _, r := range request.Requests {
		key := keyFor(s.workflowOwner, r.Namespace, r.Id)
		resp, ok := m[key]
		if !ok {
			lggr.Debugw("could not find secret in response map", "key", key)
			sdkResp = append(sdkResp, &sdkpb.SecretResponse{
				Response: &sdkpb.SecretResponse_Error{
					Error: &sdkpb.SecretError{
						Id:        r.Id,
						Namespace: r.Namespace,
						Owner:     s.workflowOwner,
						Error:     "could not find secret for " + key,
					},
				},
			})
			continue
		}

		if resp.GetError() != "" {
			lggr.Debugw("secret request returned an error", "key", key, "err", resp.GetError())
			sdkResp = append(sdkResp, &sdkpb.SecretResponse{
				Response: &sdkpb.SecretResponse_Error{
					Error: &sdkpb.SecretError{
						Id:        r.Id,
						Namespace: r.Namespace,
						Owner:     s.workflowOwner,
						Error:     "secret request returned an error for key " + key + ". Error: " + resp.GetError(),
					},
				},
			})
			continue
		}

		var localNodeShares []string
		for _, share := range resp.GetData().GetEncryptedDecryptionKeyShares() {
			if share.EncryptionKey == base64.StdEncoding.EncodeToString(localNode.EncryptionPublicKey[:]) {
				localNodeShares = share.Shares
			}
		}
		if len(localNodeShares) == 0 {
			lggr.Errorw("no shares found for this node's encryption key", "key", key, "encryptionkey", string(localNode.EncryptionPublicKey[:]))
			sdkResp = append(sdkResp, &sdkpb.SecretResponse{
				Response: &sdkpb.SecretResponse_Error{
					Error: &sdkpb.SecretError{
						Id:        r.Id,
						Namespace: r.Namespace,
						Owner:     s.workflowOwner,
						Error:     "unexpected error when getting secret for " + key + ", no shares found for this node's encryption key(base64 encoded): " + base64.StdEncoding.EncodeToString(localNode.EncryptionPublicKey[:]),
					},
				},
			})
			continue
		}

		encryptedSecretBytes, err := base64.StdEncoding.DecodeString(resp.GetData().GetEncryptedValue())
		if err != nil {
			lggr.Debugw("failed to base64 decode the secret", "key", key, "err", resp.GetError())
			sdkResp = append(sdkResp, &sdkpb.SecretResponse{
				Response: &sdkpb.SecretResponse_Error{
					Error: &sdkpb.SecretError{
						Id:        r.Id,
						Namespace: r.Namespace,
						Owner:     s.workflowOwner,
						Error:     "failed to base64 decode the secret for key " + key + ". Error: " + err.Error(),
					},
				},
			})
			continue
		}

		secret, err := s.decryptSecret(encryptedSecretBytes, localNodeShares)
		if err != nil {
			lggr.Errorw("failed to decrypt secret", "key", key, "err", err)
			sdkResp = append(sdkResp, &sdkpb.SecretResponse{
				Response: &sdkpb.SecretResponse_Error{
					Error: &sdkpb.SecretError{
						Id:        r.Id,
						Namespace: r.Namespace,
						Owner:     s.workflowOwner,
						Error:     "failed to decrypt secret for key " + key + ". Error: " + err.Error(),
					},
				},
			})
			continue
		}

		sdkResp = append(sdkResp, &sdkpb.SecretResponse{
			Response: &sdkpb.SecretResponse_Secret{
				Secret: &sdkpb.Secret{
					Id:        resp.GetId().GetKey(),
					Namespace: resp.GetId().GetNamespace(),
					Owner:     resp.GetId().GetOwner(),
					Value:     secret,
				},
			},
		})
	}

	return sdkResp, nil
}

func (s *secretsFetcher) decryptSecret(encryptedSecretBytes []byte, encodedShares []string) (string, error) {
	cipher := &tdh2easy.Ciphertext{}
	err := cipher.UnmarshalVerify(encryptedSecretBytes, s.vaultPublicKey)
	if err != nil {
		return "", errors.New("failed to unmarshal encrypted secret: " + err.Error())
	}

	decryptionShares := make([]*tdh2easy.DecryptionShare, 0, len(encodedShares))
	for _, encodedShare := range encodedShares {
		encodedShareBytes, err := base64.StdEncoding.DecodeString(encodedShare)
		if err != nil {
			return "", errors.New("Failed to base64 decode the encodedShare: " + err.Error())
		}
		privateShareBytes, err := s.workflowKey.Decrypt(encodedShareBytes)
		if err != nil {
			return "", errors.New("failed to decrypt the encodedShare: " + err.Error())
		}
		var privateShare tdh2easy.PrivateShare
		err = privateShare.Unmarshal(privateShareBytes)
		if err != nil {
			return "", errors.New("failed to unmarshal privateShare: " + err.Error())
		}
		decryptionShare, err := tdh2easy.Decrypt(cipher, &privateShare)
		if err != nil {
			return "", errors.New("failed to decrypt privateShare: " + err.Error())
		}
		err = tdh2easy.VerifyShare(cipher, s.vaultPublicKey, decryptionShare)
		if err != nil {
			return "", errors.New("failed to verifyshare the decryptionshare: " + err.Error())
		}
		decryptionShares = append(decryptionShares, decryptionShare)
	}
	decryptedSecret, err := tdh2easy.Aggregate(cipher, decryptionShares, len(encodedShares))
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
		encryptionKeys = append(encryptionKeys, base64.StdEncoding.EncodeToString(peerNode.EncryptionPublicKey[:]))
	}
	// Sort the encryption keys to ensure consistent ordering across all nodes.
	sort.Strings(encryptionKeys)
	return encryptionKeys, nil
}
