package v2

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"

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
	decrypter     func(shares []string) (string, error)

	metrics *monitoring.WorkflowsMetricLabeler
}

func NewSecretsFetcher(
	metrics *monitoring.WorkflowsMetricLabeler,
	capRegistry core.CapabilitiesRegistry,
	lggr logger.Logger,
	semaphore *semaphore[[]*sdkpb.SecretResponse],
	workflowOwner string,
	workflowName string,
	decrypter func(shares []string) (string, error),
) *secretsFetcher {
	return &secretsFetcher{
		capRegistry:   capRegistry,
		lggr:          logger.Named(lggr, "SecretsFetcher"),
		semaphore:     semaphore,
		workflowOwner: workflowOwner,
		workflowName:  workflowName,
		decrypter:     decrypter,
		metrics:       metrics,
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

	logKeys := []string{}
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

	sdkResp := []*sdkpb.SecretResponse{}
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
						Error:     resp.GetError(),
					},
				},
			})
			continue
		}

		if len(resp.GetData().GetEncryptedDecryptionKeyShares()) != 1 {
			lggr.Errorw("unexpected number of decryption shares received", "key", key, "len", len(resp.GetData().GetEncryptedDecryptionKeyShares()))
			sdkResp = append(sdkResp, &sdkpb.SecretResponse{
				Response: &sdkpb.SecretResponse_Error{
					Error: &sdkpb.SecretError{
						Id:        r.Id,
						Namespace: r.Namespace,
						Owner:     s.workflowOwner,
						Error:     "unexpected error when getting secret for " + key,
					},
				},
			})
			continue
		}

		shares := resp.GetData().EncryptedDecryptionKeyShares[0].Shares
		secret, err := s.decrypter(shares)
		if err != nil {
			lggr.Errorw("failed to combine decryption shares", "key", key, "err", err)
			sdkResp = append(sdkResp, &sdkpb.SecretResponse{
				Response: &sdkpb.SecretResponse_Error{
					Error: &sdkpb.SecretError{
						Id:        r.Id,
						Namespace: r.Namespace,
						Owner:     s.workflowOwner,
						Error:     "unexpected error when getting secret for " + key,
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
		encryptionKeys = append(encryptionKeys, string(peerNode.EncryptionPublicKey[:]))
	}
	// Sort the encryption keys to ensure consistent ordering across all nodes.
	sort.Strings(encryptionKeys)
	return encryptionKeys, nil
}
