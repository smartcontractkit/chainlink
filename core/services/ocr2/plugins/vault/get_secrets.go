package vault

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"slices"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

type share struct {
	data []byte
}

func (s *share) encryptWithKeyBinary(pk string) ([]byte, error) {
	publicKey, err := hex.DecodeString(pk)
	if err != nil {
		return nil, newUserError("failed to convert public key to bytes: " + err.Error())
	}

	if len(publicKey) != curve25519.PointSize {
		return nil, newUserError(fmt.Sprintf("invalid public key size: expected %d bytes, got %d bytes", curve25519.PointSize, len(publicKey)))
	}

	publicKeyLength := [curve25519.PointSize]byte(publicKey)
	encrypted, err := box.SealAnonymous(nil, s.data, &publicKeyLength, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt decryption share: %w", err)
	}

	return encrypted, nil
}

func generatePlaintextShare(publicKey *tdh2easy.PublicKey, privateKeyShare *tdh2easy.PrivateShare, encryptedSecret []byte, workflowOwner string) (*share, error) {
	ct := &tdh2easy.Ciphertext{}
	err := ct.UnmarshalVerify(encryptedSecret, publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ciphertext: %w", err)
	}

	es := hex.EncodeToString(encryptedSecret)
	err = vaultcap.EnsureRightLabelOnSecret(publicKey, es, workflowOwner)
	if err != nil {
		return nil, errors.New("failed to verify label on secret. error: " + err.Error())
	}

	s, err := tdh2easy.Decrypt(ct, privateKeyShare)
	if err != nil {
		return nil, fmt.Errorf("could not generate decryption share: %w", err)
	}

	sb, err := s.Marshal()
	if err != nil {
		return nil, errors.New("could not marshal decryption share")
	}

	return &share{data: sb}, nil
}

func (r *ReportingPlugin) observeGetSecrets(ctx context.Context, reader ReadKVStore, req proto.Message, o *vaultcommon.Observation) {
	tp := req.(*vaultcommon.GetSecretsRequest)
	o.RequestType = vaultcommon.RequestType_GET_SECRETS
	if !r.optimizationsEnabled(ctx) {
		o.Request = &vaultcommon.Observation_GetSecretsRequest{
			GetSecretsRequest: tp,
		}
	}
	resps := []*vaultcommon.SecretResponse{}
	for _, secretRequest := range tp.Requests {
		resp, ierr := r.observeGetSecretsRequest(ctx, reader, secretRequest)
		if ierr != nil {
			logUserErrorAware(r.lggr, "failed to observe get secret request item", ierr, "id", secretRequest.Id)
			errorMsg := userFacingError(ierr, "failed to handle get secret request")
			resps = append(resps, &vaultcommon.SecretResponse{
				Id: secretRequest.Id,
				Result: &vaultcommon.SecretResponse_Error{
					Error: errorMsg,
				},
			})
		} else {
			r.lggr.Debugw("observed get secret request item", "id", resp.Id)
			resps = append(resps, resp)
		}
	}

	o.Response = &vaultcommon.Observation_GetSecretsResponse{
		GetSecretsResponse: &vaultcommon.GetSecretsResponse{
			Responses: resps,
		},
	}
}

func (r *ReportingPlugin) observeGetSecretsRequest(ctx context.Context, reader ReadKVStore, secretRequest *vaultcommon.SecretRequest) (*vaultcommon.SecretResponse, error) {
	id, err := r.validateSecretIdentifier(ctx, secretRequest.Id)
	if err != nil {
		return nil, err
	}

	secret, err := reader.GetSecret(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}
	if secret == nil {
		return nil, newUserError("key does not exist")
	}

	sh, err := generatePlaintextShare(r.cfg.PublicKey, r.cfg.PrivateKeyShare, secret.EncryptedSecret, id.Owner)
	if err != nil {
		return nil, err
	}

	shares := []*vaultcommon.EncryptedShares{}
	useBinaryShares := r.optimizationsEnabled(ctx)
	for _, pk := range secretRequest.EncryptionKeys {
		encShare, err := sh.encryptWithKeyBinary(pk)
		if err != nil {
			return nil, err
		}

		if useBinaryShares {
			shares = append(shares, &vaultcommon.EncryptedShares{
				EncryptionKey: pk,
				BinaryShares:  [][]byte{encShare},
			})
		} else {
			shares = append(shares, &vaultcommon.EncryptedShares{
				EncryptionKey: pk,
				Shares: []string{
					hex.EncodeToString(encShare),
				},
			})
		}
	}

	return &vaultcommon.SecretResponse{
		Id: id,
		Result: &vaultcommon.SecretResponse_Data{
			Data: &vaultcommon.SecretData{
				EncryptedValue:               hex.EncodeToString(secret.EncryptedSecret),
				EncryptedDecryptionKeyShares: shares,
			},
		},
	}, nil
}

func (r *ReportingPlugin) stateTransitionGetSecrets(ctx context.Context, chosen []*vaultcommon.Observation, o *vaultcommon.Outcome) {
	if !r.optimizationsEnabled(ctx) {
		first := chosen[0]
		reqs := first.GetGetSecretsRequest().Requests
		idToReqs := map[string]*vaultcommon.SecretRequest{}
		for _, req := range reqs {
			idToReqs[vaulttypes.KeyFor(req.Id)] = req
		}

		newReqs := make([]*vaultcommon.SecretRequest, 0, len(idToReqs))
		for _, sreq := range slices.Sorted(maps.Keys(idToReqs)) {
			newReqs = append(newReqs, idToReqs[sreq])
		}

		o.Request = &vaultcommon.Outcome_GetSecretsRequest{
			GetSecretsRequest: &vaultcommon.GetSecretsRequest{
				Requests: newReqs,
			},
		}
	}

	// Next, we deal with the responses.
	// For each request, we take the Id of the first observation
	// then aggregate the encrypted shares across all observations.
	// We sort these by Id and use the result as the response.
	idToAggResponse := map[string]*vaultcommon.SecretResponse{}
	for _, resp := range chosen {
		getSecretsResp := resp.GetGetSecretsResponse()
		for _, rsp := range getSecretsResp.Responses {
			key := vaulttypes.KeyFor(rsp.Id)
			mergedResp, ok := idToAggResponse[key]
			if !ok {
				resp := &vaultcommon.SecretResponse{
					Id:     rsp.Id,
					Result: rsp.Result,
				}
				idToAggResponse[key] = resp
				continue
			}

			if rsp.GetData() != nil {
				data := mergedResp.GetData()

				if len(data.EncryptedDecryptionKeyShares) == 0 {
					data.EncryptedDecryptionKeyShares = []*vaultcommon.EncryptedShares{}
				}

				keyToShares := map[string]*vaultcommon.EncryptedShares{}
				for _, s := range data.EncryptedDecryptionKeyShares {
					keyToShares[s.EncryptionKey] = s
				}

				innerCtx := contexts.WithCRE(ctx, contexts.CRE{Owner: rsp.Id.Owner})
				for _, existing := range rsp.GetData().EncryptedDecryptionKeyShares {
					if err := validateEncryptedSharesEntry(existing); err != nil {
						// This should not happen because we validate against this in ValidateObservation.
						r.lggr.Errorw("exactly 1 share must be provided in the response, skipping", "id", rsp.Id)
						continue
					}
					shareSize, err := encryptedShareSizeForLimit(existing)
					if err != nil {
						r.lggr.Errorw("could not measure share size, skipping", "id", rsp.Id, "encryptionKey", existing.EncryptionKey, "err", err)
						continue
					}
					if err := r.cfg.MaxShareLengthBytes.Check(innerCtx, pkgconfig.Size(shareSize)*pkgconfig.Byte); err != nil {
						var errBoundLimited limits.ErrorBoundLimited[pkgconfig.Size]
						if errors.As(err, &errBoundLimited) {
							r.lggr.Errorw("share exceeds max allowed size, skipping...", "id", rsp.Id, "encryptionKey", existing.EncryptionKey, "err", err)
						} else {
							r.lggr.Errorw("could not check max allowed share size, skipping...", "id", rsp.Id, "encryptionKey", existing.EncryptionKey, "err", err)
						}
						continue
					}

					if shares, ok := keyToShares[existing.EncryptionKey]; ok {
						appendEncryptedShareEntry(shares, existing)
					} else {
						// This shouldn't happen -- this is because we're aggregating
						// requests that have a matching sha (excluding the decryption share).
						// Accordingly, we can assume that the request has been made with the same
						// set of encryption keys.
						r.lggr.Errorw("unexpected encryption key in response", "id", rsp.Id, "encryptionKey", existing.EncryptionKey)
					}
				}
			}
		}
	}

	sortedResponses := []*vaultcommon.SecretResponse{}
	for _, k := range slices.Sorted(maps.Keys(idToAggResponse)) {
		sortedResponses = append(sortedResponses, idToAggResponse[k])
	}

	o.Response = &vaultcommon.Outcome_GetSecretsResponse{
		GetSecretsResponse: &vaultcommon.GetSecretsResponse{
			Responses: sortedResponses,
		},
	}
}
