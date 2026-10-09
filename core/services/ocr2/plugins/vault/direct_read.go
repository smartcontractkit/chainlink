package vault

import (
	"context"
	"encoding/hex"
	"fmt"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

var _ vaultcap.DirectSecretsReader = (*ReportingPlugin)(nil)

type directSecretRead struct {
	id     *vaultcommon.SecretIdentifier
	secret *vaultcommon.StoredSecret
	err    error
}

// Per-item errors match the OCR path's so the aggregator can group identical ones.
func (r *ReportingPlugin) GetSecretsDirect(ctx context.Context, req *vaultcommon.GetSecretsRequest) (*vaultcommon.GetSecretsResponse, error) {
	if r.readOnlyKV == nil {
		return nil, vaultcap.ErrDirectReadNotReady
	}

	publicKey, err := r.directReadPublicKey()
	if err != nil {
		return nil, err
	}

	reads, err := r.readSecretsDirect(ctx, req.Requests)
	if err != nil {
		return nil, err
	}

	// Shares are generated after the transaction is discarded; it should only be held for reads.
	resps := make([]*vaultcommon.SecretResponse, 0, len(req.Requests))
	for i, secretRequest := range req.Requests {
		read := reads[i]
		resp, ierr := read.response(r, secretRequest)
		if ierr != nil {
			logUserErrorAware(r.lggr, "failed to serve direct get secret request item", ierr, "id", secretRequest.Id)
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

func (r *ReportingPlugin) directReadPublicKey() (string, error) {
	if r.cfg.PublicKey == nil {
		return "", vaultcap.ErrDirectReadNotReady
	}
	pkb, err := r.cfg.PublicKey.Marshal()
	if err != nil {
		return "", fmt.Errorf("failed to marshal vault public key: %w", err)
	}
	return hex.EncodeToString(pkb), nil
}

func (d directSecretRead) response(r *ReportingPlugin, secretRequest *vaultcommon.SecretRequest) (*vaultcommon.SecretResponse, error) {
	if d.err != nil {
		return nil, d.err
	}
	return r.buildGetSecretsResponse(d.id, secretRequest, d.secret)
}

func (r *ReportingPlugin) readSecretsDirect(ctx context.Context, requests []*vaultcommon.SecretRequest) ([]directSecretRead, error) {
	reads := make([]directSecretRead, len(requests))
	counts := buildGetSecretsRequestIdentifierCounts(requests)
	for i, sr := range requests {
		reads[i].id, reads[i].err = r.validateGetSecretsRequestItem(ctx, sr, counts)
	}

	txn, err := r.readOnlyKV.NewReadTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open read transaction: %w", err)
	}
	defer txn.Discard()

	store := NewReadStore(txn, r.metrics)
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

	r.lggr.Debugw("read secrets for direct get secrets request", "seqNr", txn.SeqNr(), "count", len(reads))
	return reads, nil
}
