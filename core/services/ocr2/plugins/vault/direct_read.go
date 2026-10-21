package vault

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

var _ vaultcap.DirectSecretsReader = (*getSecretsHandler)(nil)

// getSecretsHandler reads secrets from replicated state and generates
// GetSecrets responses with the key material of a single DKG instance
type getSecretsHandler struct {
	lggr            logger.SugaredLogger
	metrics         *pluginMetrics
	validator       *vaultcap.RequestValidator
	publicKey       *tdh2easy.PublicKey
	privateKeyShare *tdh2easy.PrivateShare

	// Nil unless this instance serves direct GetSecrets reads.
	readOnlyKV ocr3_1types.ReadOnlyKeyValueState
	registry   *vaultcap.LazyDirectSecretsReader
}

type directSecretRead struct {
	id     *vaultcommon.SecretIdentifier
	secret *vaultcommon.StoredSecret
	err    error
}

// GetSecretsDirect serves a direct GetSecrets request from local replicated state.
// Per-item errors match the OCR path's so the aggregator can group identical ones.
func (d *getSecretsHandler) GetSecretsDirect(ctx context.Context, req *vaultcommon.GetSecretsRequest) (*vaultcommon.GetSecretsResponse, error) {
	if d.readOnlyKV == nil {
		return nil, vaultcap.ErrDirectReadNotReady
	}

	publicKey, err := d.directReadPublicKey()
	if err != nil {
		return nil, err
	}

	reads, err := d.readSecretsDirect(ctx, req.Requests)
	if err != nil {
		return nil, err
	}

	// Shares are generated after the transaction is discarded; it should only be held for reads.
	resps := make([]*vaultcommon.SecretResponse, 0, len(req.Requests))
	for i, secretRequest := range req.Requests {
		read := reads[i]
		resp, ierr := read.response(d, secretRequest)
		if ierr != nil {
			logUserErrorAware(d.lggr, "failed to serve direct get secret request item", ierr, "id", secretRequest.Id)
			resps = append(resps, &vaultcommon.SecretResponse{
				Id: secretRequest.Id,
				Result: &vaultcommon.SecretResponse_Error{
					Error: userFacingError(ierr, vaulttypes.SecretGetSystemErrorFallback),
				},
			})
			continue
		}
		resps = append(resps, resp)
	}

	return &vaultcommon.GetSecretsResponse{Responses: resps, RawVaultPublicKey: publicKey}, nil
}

func (d *getSecretsHandler) directReadPublicKey() (string, error) {
	if d.publicKey == nil {
		return "", vaultcap.ErrDirectReadNotReady
	}
	pkb, err := d.publicKey.Marshal()
	if err != nil {
		return "", fmt.Errorf("failed to marshal vault public key: %w", err)
	}
	return hex.EncodeToString(pkb), nil
}

// buildGetSecretsResponse generates this instance's share of the requested
// secret, encrypted to each requested key; shared with the plugin's OCR
// observation path.
func (d *getSecretsHandler) buildGetSecretsResponse(id *vaultcommon.SecretIdentifier, secretRequest *vaultcommon.SecretRequest, secret *vaultcommon.StoredSecret) (*vaultcommon.SecretResponse, error) {
	sh, err := generatePlaintextShare(d.publicKey, d.privateKeyShare, secret.EncryptedSecret, id.Owner)
	if err != nil {
		return nil, err
	}

	shares := []*vaultcommon.EncryptedShares{}
	for _, pk := range secretRequest.EncryptionKeys {
		encShare, err := sh.encryptWithKeyBinary(pk)
		if err != nil {
			return nil, err
		}

		shares = append(shares, &vaultcommon.EncryptedShares{
			EncryptionKey: pk,
			BinaryShares:  [][]byte{encShare},
		})
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

func (d directSecretRead) response(h *getSecretsHandler, secretRequest *vaultcommon.SecretRequest) (*vaultcommon.SecretResponse, error) {
	if d.err != nil {
		return nil, d.err
	}
	return h.buildGetSecretsResponse(d.id, secretRequest, d.secret)
}

func (d *getSecretsHandler) readSecretsDirect(ctx context.Context, requests []*vaultcommon.SecretRequest) ([]directSecretRead, error) {
	reads := make([]directSecretRead, len(requests))
	counts := buildGetSecretsRequestIdentifierCounts(requests)
	for i, sr := range requests {
		reads[i].id, reads[i].err = validateGetSecretsRequestItem(ctx, d.validator, sr, counts)
	}

	txn, err := d.readOnlyKV.NewReadTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open read transaction: %w", err)
	}
	defer txn.Discard()

	store := NewReadStore(txn, d.metrics)
	for i := range reads {
		if reads[i].err != nil {
			continue
		}
		secret, err := store.GetSecret(ctx, reads[i].id)
		switch {
		case err != nil:
			reads[i].err = fmt.Errorf("failed to read secret from key-value store: %w", err)
		case secret == nil:
			reads[i].err = vaulttypes.NewUserError("key does not exist")
		default:
			reads[i].secret = secret
		}
	}

	d.lggr.Debugw("read secrets for direct get secrets request", "seqNr", txn.SeqNr(), "count", len(reads))
	return reads, nil
}

func (d *getSecretsHandler) Close() {
	if d.registry != nil {
		d.registry.Clear(d)
	}
}
