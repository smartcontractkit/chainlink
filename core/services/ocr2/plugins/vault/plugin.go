package vault

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/quorumhelper"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
)

const (
	defaultBatchSize = 20
	defaultNamespace = "main"
	keySeparator     = ":"
)

func NewReportingPluginFactory(lggr logger.Logger, store *requests.Store[*Request], publicKey *tdh2easy.PublicKey, privateKeyShare *tdh2easy.PrivateShare) *ReportingPluginFactory {
	return &ReportingPluginFactory{
		lggr:            lggr,
		store:           store,
		batchSize:       defaultBatchSize, // TODO fetch from onchain config
		publicKey:       publicKey,
		privateKeyShare: privateKeyShare,
	}
}

type ReportingPluginFactory struct {
	lggr            logger.Logger
	store           *requests.Store[*Request]
	batchSize       int
	publicKey       *tdh2easy.PublicKey
	privateKeyShare *tdh2easy.PrivateShare
}

func (r *ReportingPluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, fetcher ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	return &ReportingPlugin{
		lggr:            r.lggr.Named("ReportingPlugin"),
		store:           r.store,
		batchSize:       r.batchSize,
		config:          config,
		publicKey:       r.publicKey,
		privateKeyShare: r.privateKeyShare,
	}, ocr3_1types.ReportingPluginInfo{}, nil
}

type ReportingPlugin struct {
	lggr                  logger.Logger
	store                 *requests.Store[*Request]
	batchSize             int
	config                ocr3types.ReportingPluginConfig
	publicKey             *tdh2easy.PublicKey
	privateKeyShare       *tdh2easy.PrivateShare
	maxSecretsPerOwner    int
	maxCiphertextLenBytes int
	maxIdentifierLenBytes int
}

func (r *ReportingPlugin) Query(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Query, error) {
	return types.Query{}, nil
}

func (r *ReportingPlugin) Observation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, keyValueReader ocr3_1types.KeyValueReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Observation, error) {
	batch, err := r.store.FirstN(r.batchSize)
	if err != nil {
		return nil, fmt.Errorf("could not fetch batch of requests: %w", err)
	}

	ids := []string{}
	obs := []*vault.Observation{}
	for _, req := range batch {
		o := &vault.Observation{
			Id: req.ID(),
		}
		ids = append(ids, req.ID())

		switch tp := req.Payload.(type) {
		case *vault.GetSecretsRequest:
			o.RequestType = vault.RequestType_GET_SECRETS
			o.Request = &vault.Observation_GetSecretsRequest{
				GetSecretsRequest: tp,
			}

			resps := []*vault.SecretResponse{}
			for _, secretRequest := range tp.Requests {
				resp, err := r.handleGetSecretRequest(ctx, newReadStore(keyValueReader), secretRequest)
				if err != nil {
					r.lggr.Errorw("failed to handle get secret request", "id", secretRequest.Id, "error", err)
					errorMsg := "failed to handle get secret request"
					if errors.Is(err, &userError{}) {
						errorMsg = err.Error()
					}
					resps = append(resps, &vault.SecretResponse{
						Id: secretRequest.Id,
						Result: &vault.SecretResponse_Error{
							Error: errorMsg,
						},
					})
				} else {
					resps = append(resps, resp)
				}
			}

			o.Response = &vault.Observation_GetSecretsResponse{
				GetSecretsResponse: &vault.GetSecretsResponse{
					Responses: resps,
				},
			}
		case *vault.CreateSecretsRequest:
			o.RequestType = vault.RequestType_CREATE_SECRETS
			o.Request = &vault.Observation_CreateSecretsRequest{
				CreateSecretsRequest: tp,
			}

			resps := []*vault.CreateSecretResponse{}
			for _, sr := range tp.EncryptedSecrets {
				resp, err := r.handleCreateSecretRequest(ctx, newReadStore(keyValueReader), sr)
				if err != nil {
					r.lggr.Errorw("failed to handle create secret request", "id", sr.Id, "error", err)
					errorMsg := "failed to handle create secret request"
					if errors.Is(err, &userError{}) {
						errorMsg = err.Error()
					}
					resps = append(resps, &vault.CreateSecretResponse{
						Id:      sr.Id,
						Success: false,
						Error:   errorMsg,
					})
				} else {
					resps = append(resps, resp)
				}
			}

			o.Response = &vault.Observation_CreateSecretsResponse{
				CreateSecretsResponse: &vault.CreateSecretsResponse{
					Responses: resps,
				},
			}
		default:
			r.lggr.Debugw("unknown request type, skipping...", "requestType", fmt.Sprintf("%T", req.Payload), "id", req.ID())
			continue
		}

		obs = append(obs, o)
	}

	obsb, err := proto.MarshalOptions{Deterministic: true}.Marshal(&vault.Observations{
		Observations: obs,
	})
	if err != nil {
		return nil, fmt.Errorf("could not marshal observations: %w", err)
	}

	r.lggr.Debugw("Observation complete", "ids", ids, "batchSize", len(batch))
	return types.Observation(obsb), nil
}

func (r *ReportingPlugin) validateSecretIdentifier(id *vault.SecretIdentifier) (*vault.SecretIdentifier, error) {
	if id == nil {
		return nil, newUserError("invalid secret identifier: cannot be nil")
	}

	if id.Key == "" {
		return nil, newUserError("invalid secret identifier: key cannot be empty")
	}

	if id.Owner == "" {
		return nil, newUserError("invalid secret identifier: owner cannot be empty")
	}

	namespace := id.Namespace
	if namespace == "" {
		id.Namespace = defaultNamespace
	}

	if strings.Contains(id.Key, keySeparator) || strings.Contains(id.Owner, keySeparator) || strings.Contains(namespace, keySeparator) {
		return nil, newUserError(fmt.Sprintf("invalid secret identifier: id cannot contain `%s`", keySeparator))
	}

	newId := &vault.SecretIdentifier{
		Key:       id.Key,
		Owner:     id.Owner,
		Namespace: namespace,
	}

	if len(keyFor(newId)) > r.maxIdentifierLenBytes {
		return nil, newUserError(fmt.Sprintf("invalid secret identifier: id exceeds maximum length of %d bytes", r.maxIdentifierLenBytes))
	}

	return newId, nil
}

func newUserError(msg string) *userError {
	return &userError{msg: msg}
}

type userError struct {
	msg string
}

func (u *userError) Error() string {
	return u.msg
}

func (u *userError) Is(target error) bool {
	_, ok := target.(*userError)
	return ok
}

func keyFor(id *vault.SecretIdentifier) string {
	return fmt.Sprintf("%s::%s::%s", id.Owner, id.Namespace, id.Key)
}

func (r *ReportingPlugin) handleGetSecretRequest(ctx context.Context, reader readKVStore, secretRequest *vault.SecretRequest) (*vault.SecretResponse, error) {
	id, err := r.validateSecretIdentifier(secretRequest.Id)
	if err != nil {
		return nil, err
	}

	secret, err := reader.getSecret(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from key-value store: %w", err)
	}

	if secret == nil {
		return nil, newUserError("key does not exist")
	}

	ct := &tdh2easy.Ciphertext{}
	err = ct.UnmarshalVerify(secret.EncryptedSecret, r.publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ciphertext: %w", err)
	}

	share, err := tdh2easy.Decrypt(ct, r.privateKeyShare)
	if err != nil {
		return nil, fmt.Errorf("could not generate decryption share: %w", err)
	}

	shareb, err := share.Marshal()
	if err != nil {
		return nil, errors.New("could not marshal decryption share")
	}

	shares := []*vault.EncryptedShares{}
	for _, pk := range secretRequest.EncryptionKeys {
		publicKey, err := base64.StdEncoding.DecodeString(pk)
		if err != nil {
			return nil, newUserError("failed to convert public key to bytes: " + err.Error())
		}

		if len(publicKey) != curve25519.PointSize {
			return nil, newUserError(fmt.Sprintf("invalid public key size: expected %d bytes, got %d bytes", curve25519.PointSize, len(publicKey)))
		}

		publicKeyLength := [curve25519.PointSize]byte(publicKey)
		encrypted, err := box.SealAnonymous(nil, shareb, &publicKeyLength, rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt decryption share: %w", err)
		}

		shares = append(shares, &vault.EncryptedShares{
			Shares: []string{
				base64.StdEncoding.EncodeToString(encrypted),
			},
		})
	}

	return &vault.SecretResponse{
		Id: id,
		Result: &vault.SecretResponse_Data{
			Data: &vault.SecretData{
				EncryptedValue:               base64.StdEncoding.EncodeToString(secret.EncryptedSecret),
				EncryptedDecryptionKeyShares: shares,
			},
		},
	}, nil
}

func (r *ReportingPlugin) handleCreateSecretRequest(ctx context.Context, reader readKVStore, secretRequest *vault.EncryptedSecret) (*vault.CreateSecretResponse, error) {
	id, err := r.validateSecretIdentifier(secretRequest.Id)
	if err != nil {
		return nil, err
	}

	rawCiphertext := secretRequest.EncryptedValue
	rawCiphertextB, err := base64.StdEncoding.DecodeString(rawCiphertext)
	if err != nil {
		return nil, newUserError(fmt.Sprintf("invalid base64 encoding for ciphertext: %s", err.Error()))
	}

	if len(rawCiphertextB) > r.maxCiphertextLenBytes {
		return nil, newUserError(fmt.Sprintf("ciphertext size exceeds maximum allowed size: %d bytes", r.maxCiphertextLenBytes))
	}

	ct := &tdh2easy.Ciphertext{}
	err = ct.UnmarshalVerify(rawCiphertextB, r.publicKey)
	if err != nil {
		return nil, newUserError(fmt.Sprintf("failed to verify ciphertext: %s", err.Error()))
	}

	md, err := reader.getMetadata(id.Owner)
	if err != nil {
		return nil, err
	}

	// If the metadata record doesn't exist, we can assume this is the first time
	// creating a secret for this user and check against the default limit.
	count := 0
	if md != nil {
		count = len(md.Keys)
	}

	if count+1 > r.maxSecretsPerOwner {
		return nil, newUserError(fmt.Sprintf("maximum number of secrets per owner reached: %d", r.maxSecretsPerOwner))
	}

	secret, err := reader.getSecret(id)
	if err != nil {
		return nil, err
	}

	if secret != nil {
		return nil, newUserError("key already exists")
	}

	// Return an initialized response,
	// This will get filled in during the outcome step.
	return &vault.CreateSecretResponse{
		Id:      id,
		Success: false,
		Error:   "",
	}, nil
}

func (r *ReportingPlugin) ValidateObservation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, ao types.AttributedObservation, keyValueReader ocr3_1types.KeyValueReader, blobFetcher ocr3_1types.BlobFetcher) error {
	return nil
}

func (r *ReportingPlugin) ObservationQuorum(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReader ocr3_1types.KeyValueReader, blobFetcher ocr3_1types.BlobFetcher) (quorumReached bool, err error) {
	return quorumhelper.ObservationCountReachesObservationQuorum(quorumhelper.QuorumTwoFPlusOne, r.config.N, r.config.F, aos), nil
}

func shaForProto(msg proto.Message) (string, error) {
	protoBytes, err := proto.MarshalOptions{Deterministic: true}.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("could not generate sha for proto message: failed to marshal proto: %w", err)
	}

	return fmt.Sprintf("%x", sha256.Sum256(protoBytes)), nil
}

func shaForObservation(o *vault.Observation) (string, error) {
	switch o.RequestType {
	case vault.RequestType_GET_SECRETS:
		cloned := proto.CloneOf(o)
		for _, r := range cloned.GetGetSecretsResponse().Responses {
			if r.GetData() != nil {
				// Exclude the encrypted shares from the sha, as these need to be aggregated later.
				r.GetData().EncryptedDecryptionKeyShares = nil
			}
		}

		return shaForProto(cloned)
	default:
		return shaForProto(o)
	}
}

func validateObservation(o *vault.Observation) error {
	if o.Id == "" {
		return fmt.Errorf("observation id cannot be empty")
	}

	switch o.RequestType {
	case vault.RequestType_GET_SECRETS:
		if o.GetGetSecretsRequest() == nil || o.GetGetSecretsResponse() == nil {
			return errors.New("GetSecrets observation must have both request and response")
		}

		if len(o.GetGetSecretsRequest().Requests) != len(o.GetGetSecretsResponse().Responses) {
			return errors.New("GetSecrets request and response must have the same number of items")
		}
	case vault.RequestType_CREATE_SECRETS:
		if o.GetCreateSecretsRequest() == nil || o.GetCreateSecretsResponse() == nil {
			return errors.New("GetSecrets observation must have both request and response")
		}

		if len(o.GetCreateSecretsRequest().EncryptedSecrets) != len(o.GetCreateSecretsResponse().Responses) {
			return errors.New("GetSecrets request and response must have the same number of items")
		}
	default:
		return errors.New("invalid observation type: " + o.RequestType.String())
	}

	return nil
}

func (r *ReportingPlugin) StateTransition(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReadWriter ocr3_1types.KeyValueReadWriter, blobFetcher ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	store := newWriteStore(keyValueReadWriter)

	obsMap := map[string][]*vault.Observation{}
	for _, ao := range aos {
		obs := &vault.Observations{}
		if err := proto.Unmarshal([]byte(ao.Observation), obs); err != nil {
			r.lggr.Errorw("failed to unmarshal observations", "error", err, "observation", ao.Observation)
			continue
		}

		for _, o := range obs.Observations {
			err := validateObservation(o)
			if err != nil {
				r.lggr.Errorw("invalid observation", "error", err, "observation", o)
				continue
			}

			if _, ok := obsMap[o.Id]; !ok {
				obsMap[o.Id] = []*vault.Observation{}
			}
			obsMap[o.Id] = append(obsMap[o.Id], o)
		}
	}

	os := &vault.Outcomes{
		Outcomes: []*vault.Outcome{},
	}
	for id, obs := range obsMap {
		shaToObs := map[string][]*vault.Observation{}
		for _, ob := range obs {
			sha, err := shaForObservation(ob)
			if err != nil {
				r.lggr.Errorw("failed to compute sha for observation", "error", err, "observation", ob)
				continue
			}
			shaToObs[sha] = append(shaToObs[sha], ob)
		}

		chosen := []*vault.Observation{}
		threshold := 2*r.config.F + 1
		for sha, obs := range shaToObs {
			if len(obs) >= threshold {
				r.lggr.Debugw("sufficient observations for sha", "sha", sha, "count", len(obs), "threshold", threshold, "id", id)
				chosen = shaToObs[sha]
				break
			}
		}

		if len(chosen) == 0 {
			r.lggr.Warnw("insufficient observations found for id", "id", id, "threshold", threshold)
			continue
		}

		// The shas are the same so the requests will have
		// the same Id and Type.
		first := chosen[0]
		o := &vault.Outcome{
			Id:          first.Id,
			RequestType: first.RequestType,
		}
		switch first.RequestType {
		case vault.RequestType_GET_SECRETS:
			// First, let's generate the aggregated request.
			// We've validated that all requests with the same sha have the same
			// contents, so we can just sort the SecretRequests by their ID
			// and use that as the aggregated request.
			reqs := first.GetGetSecretsRequest().Requests
			idToReqs := map[string]*vault.SecretRequest{}
			for _, req := range reqs {
				idToReqs[keyFor(req.Id)] = req
			}

			newReqs := []*vault.SecretRequest{}
			for _, sreq := range slices.Sorted(maps.Keys(idToReqs)) {
				newReqs = append(newReqs, idToReqs[sreq])
			}

			o.Request = &vault.Outcome_GetSecretsRequest{
				GetSecretsRequest: &vault.GetSecretsRequest{
					Requests: newReqs,
				},
			}

			// Next, we deal with the responses.
			// For each request, we take the Id of the first observation
			// then aggregate the encrypted shares across all observations.
			// Like with the requests, we sort these by Id and use the result as the response.
			idToAggResponse := map[string]*vault.SecretResponse{}
			for _, resp := range chosen {
				getSecretsResp := resp.GetGetSecretsResponse()
				for _, r := range getSecretsResp.Responses {
					key := keyFor(r.Id)
					mergedResp, ok := idToAggResponse[key]
					if !ok {
						resp := &vault.SecretResponse{
							Id:     r.Id,
							Result: r.Result,
						}
						idToAggResponse[key] = resp
						continue
					}

					if r.GetData() != nil {
						data := mergedResp.GetData()

						if len(data.EncryptedDecryptionKeyShares) == 0 {
							data.EncryptedDecryptionKeyShares = []*vault.EncryptedShares{}
						}

						keyToShares := map[string]*vault.EncryptedShares{}
						for _, s := range data.EncryptedDecryptionKeyShares {
							keyToShares[s.EncryptionKey] = s
						}

						for _, existing := range r.GetData().EncryptedDecryptionKeyShares {
							if shares, ok := keyToShares[existing.EncryptionKey]; ok {
								shares.Shares = append(shares.Shares, existing.Shares...)
							} else {
								data.EncryptedDecryptionKeyShares = append(
									data.EncryptedDecryptionKeyShares,
									existing,
								)
							}
						}
					}
				}
			}

			sortedResponses := []*vault.SecretResponse{}
			for _, k := range slices.Sorted(maps.Keys(idToAggResponse)) {
				sortedResponses = append(sortedResponses, idToAggResponse[k])
			}

			o.Response = &vault.Outcome_GetSecretsResponse{
				GetSecretsResponse: &vault.GetSecretsResponse{
					Responses: sortedResponses,
				},
			}
			os.Outcomes = append(os.Outcomes, o)
		case vault.RequestType_CREATE_SECRETS:
			// First we'll aggregate the requests.
			// Since the shas for all requests match, we can just take the first entry
			// and sort the requests contained within it.
			req := first.GetCreateSecretsRequest().EncryptedSecrets
			idToReqs := map[string]*vault.EncryptedSecret{}
			for _, r := range req {
				idToReqs[keyFor(r.Id)] = r
			}

			newReqs := []*vault.EncryptedSecret{}
			for _, sreq := range slices.Sorted(maps.Keys(idToReqs)) {
				newReqs = append(newReqs, idToReqs[sreq])
			}

			o.Request = &vault.Outcome_CreateSecretsRequest{
				CreateSecretsRequest: &vault.CreateSecretsRequest{
					EncryptedSecrets: newReqs,
				},
			}

			// Next let's aggregate the responses.
			// We do this by taking the first response, and determine if
			// there was a validation error. If not, we write it to the key value store.
			// The responses are sorted by Id.
			resp := first.GetCreateSecretsResponse()
			idToResps := map[string]*vault.CreateSecretResponse{}
			for _, r := range resp.Responses {
				idToResps[keyFor(r.Id)] = r
			}

			sortedResps := []*vault.CreateSecretResponse{}
			for _, id := range slices.Sorted(maps.Keys(idToResps)) {
				resp := idToResps[id]
				req := idToReqs[id]
				if resp.GetError() != "" {
					sortedResps = append(sortedResps, resp)
					continue
				}

				encryptedSecret, err := base64.StdEncoding.DecodeString(req.EncryptedValue)
				if err != nil {
					sortedResps = append(sortedResps, &vault.CreateSecretResponse{
						Id:      resp.Id,
						Success: false,
						Error:   "could not decode secret value: invalid base64",
					})
					continue
				}

				err = store.writeSecret(req.Id, &vault.StoredSecret{
					EncryptedSecret: encryptedSecret,
				})
				if err != nil {
					r.lggr.Errorw("failed to write secret to key value store", "error", err, "id", resp.Id)
					sortedResps = append(sortedResps, &vault.CreateSecretResponse{
						Id:      resp.Id,
						Success: false,
						Error:   "failed to write secret to key value store",
					})
					continue
				}

				err = store.addKeyToMetadata(req.Id)
				if err != nil {
					r.lggr.Errorw("failed to add key to metadata", "error", err, "id", resp.Id)
					sortedResps = append(sortedResps, &vault.CreateSecretResponse{
						Id:      resp.Id,
						Success: false,
						Error:   "failed to write secret to key value store",
					})
					continue
				}

				sortedResps = append(sortedResps, &vault.CreateSecretResponse{
					Id:      resp.Id,
					Success: true,
					Error:   "",
				})
			}

			o.Response = &vault.Outcome_CreateSecretsResponse{
				CreateSecretsResponse: &vault.CreateSecretsResponse{
					Responses: sortedResps,
				},
			}
			os.Outcomes = append(os.Outcomes, o)
		default:
			r.lggr.Debugw("unknown request type, skipping...", "requestType", first.RequestType, "id", id)
			continue
		}
	}

	ospb, err := proto.MarshalOptions{Deterministic: true}.Marshal(os)
	r.lggr.Debugw("State transition complete", "count", len(os.Outcomes), "err", err)
	if err != nil {
		return ocr3_1types.ReportsPlusPrecursor{}, fmt.Errorf("could not marshal outcomes: %w", err)
	}

	return ocr3_1types.ReportsPlusPrecursor(ospb), nil
}

func (r *ReportingPlugin) Committed(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueReader) error {
	// Not currently used by the protocol, so we noop here.
	return nil
}

func (r *ReportingPlugin) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[[]byte], error) {
	outcomes := &vault.Outcomes{}
	err := proto.Unmarshal([]byte(reportsPlusPrecursor), outcomes)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal outcomes: %w", err)
	}

	reports := []ocr3types.ReportPlus[[]byte]{}
	for _, o := range outcomes.Outcomes {
		switch o.RequestType {
		case vault.RequestType_GET_SECRETS:
			rep, err := r.generateProtoReport(o.Id, o.RequestType, o.GetGetSecretsResponse())
			if err != nil {
				r.lggr.Errorw("failed to generate Proto report", "error", err, "id", o.Id)
				continue
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: rep,
			})
		case vault.RequestType_CREATE_SECRETS:
			rep, err := r.generateJSONReport(o.Id, o.RequestType, o.GetCreateSecretsResponse())
			if err != nil {
				r.lggr.Errorw("failed to generate JSON report", "error", err, "id", o.Id)
				continue
			}

			reports = append(reports, ocr3types.ReportPlus[[]byte]{
				ReportWithInfo: rep,
			})
		default:
		}
	}

	r.lggr.Debugw("Reports complete", "count", len(reports))
	return reports, nil
}

func (r *ReportingPlugin) generateProtoReport(id string, requestType vault.RequestType, msg proto.Message) (ocr3types.ReportWithInfo[[]byte], error) {
	if msg == nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("invalid report: response cannot be nil")
	}

	rpb, err := proto.MarshalOptions{Deterministic: true}.Marshal(msg)
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("failed to marshal response to proto: %w", err)
	}

	rip, err := proto.MarshalOptions{Deterministic: true}.Marshal(&vault.ReportInfo{
		Id:          id,
		RequestType: requestType,
		Format:      vault.ReportFormat_REPORT_FORMAT_PROTOBUF,
	})
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("failed to marshal report info: %w", err)
	}

	return ocr3types.ReportWithInfo[[]byte]{
		Report: rpb,
		Info:   rip,
	}, nil
}

func (r *ReportingPlugin) generateJSONReport(id string, requestType vault.RequestType, msg proto.Message) (ocr3types.ReportWithInfo[[]byte], error) {
	if msg == nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("invalid report: response cannot be nil")
	}

	jsonb, err := ToCanonicalJSON(msg)
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("failed to convert proto to canonical JSON: %w", err)
	}

	rip, err := proto.MarshalOptions{Deterministic: true}.Marshal(&vault.ReportInfo{
		Id:          id,
		RequestType: requestType,
		Format:      vault.ReportFormat_REPORT_FORMAT_JSON,
	})
	if err != nil {
		return ocr3types.ReportWithInfo[[]byte]{}, fmt.Errorf("failed to marshal report info: %w", err)
	}

	return ocr3types.ReportWithInfo[[]byte]{
		Report: jsonb,
		Info:   rip,
	}, nil
}

func (r *ReportingPlugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return true, nil
}

func (r *ReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return true, nil
}

func (r *ReportingPlugin) Close() error {
	return nil
}
