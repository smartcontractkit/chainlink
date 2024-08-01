package decryptionplugin

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/tdh2/go/ocr2/decryptionplugin/config"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"google.golang.org/protobuf/proto"
)

type DecryptionReportingPluginFactory struct {
	DecryptionQueue  DecryptionQueuingService
	ConfigParser     config.ConfigParser
	PublicKey        *tdh2easy.PublicKey
	PrivKeyShare     *tdh2easy.PrivateShare
	OracleToKeyShare map[commontypes.OracleID]int
	Logger           commontypes.Logger
}

type decryptionPlugin struct {
	logger           commontypes.Logger
	decryptionQueue  DecryptionQueuingService
	publicKey        *tdh2easy.PublicKey
	privKeyShare     *tdh2easy.PrivateShare
	oracleToKeyShare map[commontypes.OracleID]int
	genericConfig    *types.ReportingPluginConfig
	specificConfig   *config.ReportingPluginConfigWrapper
}

// NewReportingPlugin complies with ReportingPluginFactory.
func (f DecryptionReportingPluginFactory) NewReportingPlugin(rpConfig types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	pluginConfig, err := f.ConfigParser.ParseConfig(rpConfig.OffchainConfig)
	if err != nil {
		return nil, types.ReportingPluginInfo{},
			fmt.Errorf("unable to decode reporting plugin config: %w", err)
	}

	// The number of decryption shares K needed to reconstruct the plaintext should satisfy F<K<=2F+1.
	// The lower bound ensure that no F parties can alone reconstruct the secret.
	// The upper bound ensures that there can be always enough decryption shares.
	// It depends on the minimum number of observations collected by the leader (2F+1).
	// Note that for configurations with K>F+1 liveness is not always satisfied as the leader might
	// include an observation from a malicious party, whose decryption share is invalid.
	// However, this configuration that favours safety over liveness might be desirable in certain use cases.
	if int(pluginConfig.Config.K) <= rpConfig.F || int(pluginConfig.Config.K) > 2*rpConfig.F+1 {
		return nil, types.ReportingPluginInfo{},
			fmt.Errorf("invalid configuration with K=%d and F=%d: decryption threshold K must satisfy F < K <= 2F+1", pluginConfig.Config.K, rpConfig.F)
	}

	info := types.ReportingPluginInfo{
		Name:          "ThresholdDecryption",
		UniqueReports: false, // Aggregating any k valid decryption shares results in the same plaintext. Must match setting in OCR2Base.sol.
		// TODO calculate limits based on the maximum size of the plaintext and ciphertextID
		Limits: types.ReportingPluginLimits{
			MaxQueryLength:       int(pluginConfig.Config.GetMaxQueryLengthBytes()),
			MaxObservationLength: int(pluginConfig.Config.GetMaxObservationLengthBytes()),
			MaxReportLength:      int(pluginConfig.Config.GetMaxReportLengthBytes()),
		},
	}

	plugin := decryptionPlugin{
		f.Logger,
		f.DecryptionQueue,
		f.PublicKey,
		f.PrivKeyShare,
		f.OracleToKeyShare,
		&rpConfig,
		pluginConfig,
	}

	return &plugin, info, nil
}

// Query creates a query with the oldest pending decryption requests.
func (dp *decryptionPlugin) Query(ctx context.Context, ts types.ReportTimestamp) (types.Query, error) {
	dp.logger.Debug("DecryptionReporting Query: start", commontypes.LogFields{
		"epoch": ts.Epoch,
		"round": ts.Round,
	})

	decryptionRequests := dp.decryptionQueue.GetRequests(
		int(dp.specificConfig.Config.RequestCountLimit),
		int(dp.specificConfig.Config.RequestTotalBytesLimit),
	)

	queryProto := Query{}
	ciphertextIDs := make(map[string]bool)
	allIDs := []string{}
	for _, request := range decryptionRequests {
		if _, ok := ciphertextIDs[string(request.CiphertextId)]; ok {
			dp.logger.Error("DecryptionReporting Query: duplicate request, skipping it", commontypes.LogFields{
				"ciphertextID": request.CiphertextId.String(),
			})
			continue
		}
		ciphertextIDs[string(request.CiphertextId)] = true

		ciphertext := &tdh2easy.Ciphertext{}
		if err := ciphertext.UnmarshalVerify(request.Ciphertext, dp.publicKey); err != nil {
			dp.decryptionQueue.SetResult(request.CiphertextId, nil, ErrUnmarshalling)
			dp.logger.Error("DecryptionReporting Query: cannot unmarshal the ciphertext, skipping it", commontypes.LogFields{
				"error":        err,
				"ciphertextID": request.CiphertextId.String(),
			})
			continue
		}
		queryProto.DecryptionRequests = append(queryProto.GetDecryptionRequests(), &CiphertextWithID{
			CiphertextId: request.CiphertextId,
			Ciphertext:   request.Ciphertext,
		})
		allIDs = append(allIDs, request.CiphertextId.String())
	}

	dp.logger.Debug("DecryptionReporting Query: end", commontypes.LogFields{
		"epoch":         ts.Epoch,
		"round":         ts.Round,
		"queryLen":      len(queryProto.DecryptionRequests),
		"ciphertextIDs": allIDs,
	})
	queryProtoBytes, err := proto.Marshal(&queryProto)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal query: %w", err)
	}
	return queryProtoBytes, nil
}

// Observation creates a decryption share for each request in the query.
// If dp.specificConfig.Config.LocalRequest is true, then the oracle
// only creates a decryption share for the decryption requests which it has locally.
func (dp *decryptionPlugin) Observation(ctx context.Context, ts types.ReportTimestamp, query types.Query) (types.Observation, error) {
	dp.logger.Debug("DecryptionReporting Observation: start", commontypes.LogFields{
		"epoch": ts.Epoch,
		"round": ts.Round,
	})

	queryProto := &Query{}
	if err := proto.Unmarshal(query, queryProto); err != nil {
		return nil, fmt.Errorf("cannot unmarshal query: %w", err)
	}

	observationProto := Observation{}
	ciphertextIDs := make(map[string]bool)
	decryptedIDs := []string{}
	for _, request := range queryProto.DecryptionRequests {
		ciphertextId := CiphertextId(request.CiphertextId)
		if _, ok := ciphertextIDs[string(ciphertextId)]; ok {
			dp.logger.Error("DecryptionReporting Observation: duplicate request in the same query, the leader is faulty", commontypes.LogFields{
				"ciphertextID": ciphertextId.String(),
			})
			return nil, fmt.Errorf("duplicate request in the same query")
		}
		ciphertextIDs[string(ciphertextId)] = true

		ciphertext := &tdh2easy.Ciphertext{}
		ciphertextBytes := request.Ciphertext
		if err := ciphertext.UnmarshalVerify(ciphertextBytes, dp.publicKey); err != nil {
			dp.logger.Error("DecryptionReporting Observation: cannot unmarshal and verify the ciphertext, the leader is faulty", commontypes.LogFields{
				"error":        err,
				"ciphertextID": ciphertextId.String(),
			})
			return nil, fmt.Errorf("cannot unmarshal and verify the ciphertext: %w", err)
		}
		if dp.specificConfig.Config.RequireLocalRequestCheck {
			queueCiphertextBytes, err := dp.decryptionQueue.GetCiphertext(ciphertextId)
			if err != nil && errors.Is(err, ErrNotFound) {
				dp.logger.Warn("DecryptionReporting Observation: cannot find ciphertext locally, skipping it", commontypes.LogFields{
					"error":        err,
					"ciphertextID": ciphertextId.String(),
				})
				continue
			} else if err != nil {
				dp.logger.Error("DecryptionReporting Observation: failed when looking for ciphertext locally, skipping it", commontypes.LogFields{
					"error":        err,
					"ciphertextID": ciphertextId.String(),
				})
				continue
			}
			if !bytes.Equal(queueCiphertextBytes, ciphertextBytes) {
				dp.logger.Error("DecryptionReporting Observation: local ciphertext does not match the query ciphertext, skipping it", commontypes.LogFields{
					"ciphertextID": ciphertextId.String(),
				})
				continue
			}
		}

		decryptionShare, err := tdh2easy.Decrypt(ciphertext, dp.privKeyShare)
		if err != nil {
			dp.decryptionQueue.SetResult(ciphertextId, nil, ErrDecryption)
			dp.logger.Error("DecryptionReporting Observation: cannot decrypt the ciphertext with the private key share", commontypes.LogFields{
				"error":        err,
				"ciphertextID": ciphertextId.String(),
			})
			continue
		}
		decryptionShareBytes, err := decryptionShare.Marshal()
		if err != nil {
			dp.logger.Error("DecryptionReporting Observation: cannot marshal the decryption share, skipping it", commontypes.LogFields{
				"error":        err,
				"ciphertextID": ciphertextId.String(),
			})
			continue
		}
		observationProto.DecryptionShares = append(observationProto.DecryptionShares, &DecryptionShareWithID{
			CiphertextId:    ciphertextId,
			DecryptionShare: decryptionShareBytes,
		})
		decryptedIDs = append(decryptedIDs, ciphertextId.String())
	}

	dp.logger.Debug("DecryptionReporting Observation: end", commontypes.LogFields{
		"epoch":             ts.Epoch,
		"round":             ts.Round,
		"decryptedRequests": len(observationProto.DecryptionShares),
		"totalRequests":     len(queryProto.DecryptionRequests),
		"ciphertextIDs":     decryptedIDs,
	})
	observationProtoBytes, err := proto.Marshal(&observationProto)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal observation: %w", err)
	}
	return observationProtoBytes, nil
}

// Report aggregates decryption shares from Observations to derive the plaintext.
func (dp *decryptionPlugin) Report(ctx context.Context, ts types.ReportTimestamp, query types.Query, obs []types.AttributedObservation) (bool, types.Report, error) {
	dp.logger.Debug("DecryptionReporting Report: start", commontypes.LogFields{
		"epoch":         ts.Epoch,
		"round":         ts.Round,
		"nObservations": len(obs),
	})

	queryProto := &Query{}
	if err := proto.Unmarshal(query, queryProto); err != nil {
		return false, nil, fmt.Errorf("cannot unmarshal query: %w ", err)
	}
	ciphertexts := make(map[string]*tdh2easy.Ciphertext)
	for _, request := range queryProto.DecryptionRequests {
		ciphertextId := CiphertextId(request.CiphertextId)
		ciphertext := &tdh2easy.Ciphertext{}
		if err := ciphertext.UnmarshalVerify(request.Ciphertext, dp.publicKey); err != nil {
			dp.logger.Error("DecryptionReporting Report: cannot unmarshal and verify the ciphertext, the leader is faulty", commontypes.LogFields{
				"error":        err,
				"ciphertextID": ciphertextId.String(),
			})
			return false, nil, fmt.Errorf("cannot unmarshal and verify the ciphertext: %w", err)
		}
		ciphertexts[string(ciphertextId)] = ciphertext
	}

	validDecryptionShares := make(map[string][]*tdh2easy.DecryptionShare)
	for _, ob := range obs {
		observationProto := &Observation{}
		if err := proto.Unmarshal(ob.Observation, observationProto); err != nil {
			dp.logger.Error("DecryptionReporting Report: cannot unmarshal observation, skipping it", commontypes.LogFields{
				"error":    err,
				"observer": ob.Observer,
			})
			continue
		}

		ciphertextIDs := make(map[string]bool)
		for _, decryptionShareWithID := range observationProto.DecryptionShares {
			ciphertextId := CiphertextId(decryptionShareWithID.CiphertextId)
			ciphertextIdRawStr := string(ciphertextId)
			if _, ok := ciphertextIDs[ciphertextIdRawStr]; ok {
				dp.logger.Error("DecryptionReporting Report: the observation has multiple decryption shares for the same ciphertext id", commontypes.LogFields{
					"ciphertextID": ciphertextId.String(),
					"observer":     ob.Observer,
				})
				continue
			}
			ciphertextIDs[ciphertextIdRawStr] = true

			ciphertext, ok := ciphertexts[ciphertextIdRawStr]
			if !ok {
				dp.logger.Error("DecryptionReporting Report: there is not ciphertext in the query with matching id", commontypes.LogFields{
					"ciphertextID": ciphertextId.String(),
					"observer":     ob.Observer,
				})
				continue
			}

			validDecryptionShare, err := dp.getValidDecryptionShare(ob.Observer,
				ciphertext, decryptionShareWithID.DecryptionShare)
			if err != nil {
				dp.logger.Error("DecryptionReporting Report: invalid decryption share", commontypes.LogFields{
					"error":        err,
					"ciphertextID": ciphertextId.String(),
					"observer":     ob.Observer,
				})
				continue
			}

			if len(validDecryptionShares[ciphertextIdRawStr]) < int(dp.specificConfig.Config.K) {
				validDecryptionShares[ciphertextIdRawStr] = append(validDecryptionShares[ciphertextIdRawStr], validDecryptionShare)
			} else {
				dp.logger.Trace("DecryptionReporting Report: we have already k valid decryption shares", commontypes.LogFields{
					"ciphertextID": ciphertextId.String(),
					"observer":     ob.Observer,
				})
			}
		}
	}

	reportProto := Report{}
	for _, request := range queryProto.DecryptionRequests {
		ciphertextId := CiphertextId(request.CiphertextId)
		ciphertextIdRawStr := string(ciphertextId)
		decrShares, ok := validDecryptionShares[ciphertextIdRawStr]
		if !ok {
			// Request not included in any observation in the current round.
			dp.logger.Debug("DecryptionReporting Report: ciphertextID was not included in any observation in the current round", commontypes.LogFields{
				"ciphertextID": ciphertextId.String(),
			})
			continue
		}
		ciphertext, ok := ciphertexts[ciphertextIdRawStr]
		if !ok {
			dp.logger.Error("DecryptionReporting Report: there is not ciphertext in the query with matching id, skipping aggregation of decryption shares", commontypes.LogFields{
				"ciphertextID": ciphertextId.String(),
			})
			continue
		}

		if len(decrShares) < int(dp.specificConfig.Config.K) {
			dp.logger.Debug("DecryptionReporting Report: not enough valid decryption shares after processing all observations, skipping aggregation of decryption shares", commontypes.LogFields{
				"ciphertextID": ciphertextId.String(),
			})
			continue
		}

		plaintext, err := tdh2easy.Aggregate(ciphertext, decrShares, dp.genericConfig.N)
		if err != nil {
			dp.decryptionQueue.SetResult(ciphertextId, nil, ErrAggregation)
			dp.logger.Error("DecryptionReporting Report: cannot aggregate decryption shares", commontypes.LogFields{
				"error":        err,
				"ciphertextID": ciphertextId.String(),
			})
			continue
		}

		dp.logger.Debug("DecryptionReporting Report: plaintext aggregated successfully", commontypes.LogFields{
			"epoch":        ts.Epoch,
			"round":        ts.Round,
			"ciphertextID": ciphertextId.String(),
		})
		reportProto.ProcessedDecryptedRequests = append(reportProto.ProcessedDecryptedRequests, &ProcessedDecryptionRequest{
			CiphertextId: ciphertextId,
			Plaintext:    plaintext,
		})
	}

	dp.logger.Debug("DecryptionReporting Report: end", commontypes.LogFields{
		"epoch":                      ts.Epoch,
		"round":                      ts.Round,
		"aggregatedDecryptionShares": len(reportProto.ProcessedDecryptedRequests),
		"reporting":                  len(reportProto.ProcessedDecryptedRequests) > 0,
	})

	if len(reportProto.ProcessedDecryptedRequests) == 0 {
		return false, nil, nil
	}

	reportBytes, err := proto.Marshal(&reportProto)
	if err != nil {
		return false, nil, fmt.Errorf("cannot marshal report: %w", err)
	}
	return true, reportBytes, nil
}

func (dp *decryptionPlugin) getValidDecryptionShare(observer commontypes.OracleID,
	ciphertext *tdh2easy.Ciphertext, decryptionShareBytes []byte) (*tdh2easy.DecryptionShare, error) {
	decryptionShare := &tdh2easy.DecryptionShare{}
	if err := decryptionShare.Unmarshal(decryptionShareBytes); err != nil {
		return nil, fmt.Errorf("cannot unmarshal decryption share: %w", err)
	}

	expectedKeyShareIndex, ok := dp.oracleToKeyShare[observer]
	if !ok {
		return nil, fmt.Errorf("invalid observer ID")
	}

	if expectedKeyShareIndex != decryptionShare.Index() {
		return nil, fmt.Errorf("invalid decryption share index: expected %d and got %d", expectedKeyShareIndex, decryptionShare.Index())
	}

	if err := tdh2easy.VerifyShare(ciphertext, dp.publicKey, decryptionShare); err != nil {
		return nil, fmt.Errorf("decryption share verification failed: %w", err)
	}
	return decryptionShare, nil
}

// ShouldAcceptFinalizedReport updates the decryption queue.
// Returns always false as the report will not be transmitted on-chain.
func (dp *decryptionPlugin) ShouldAcceptFinalizedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	dp.logger.Debug("DecryptionReporting ShouldAcceptFinalizedReport: start", commontypes.LogFields{
		"epoch": ts.Epoch,
		"round": ts.Round,
	})

	reportProto := &Report{}
	if err := proto.Unmarshal(report, reportProto); err != nil {
		return false, fmt.Errorf("cannot unmarshal report: %w", err)
	}

	for _, item := range reportProto.ProcessedDecryptedRequests {
		dp.decryptionQueue.SetResult(item.CiphertextId, item.Plaintext, nil)
	}

	dp.logger.Debug("DecryptionReporting ShouldAcceptFinalizedReport: end", commontypes.LogFields{
		"epoch":     ts.Epoch,
		"round":     ts.Round,
		"accepting": false,
	})

	return false, nil
}

// ShouldTransmitAcceptedReport is a no-op
func (dp *decryptionPlugin) ShouldTransmitAcceptedReport(ctx context.Context, ts types.ReportTimestamp, report types.Report) (bool, error) {
	return false, nil
}

// Close complies with ReportingPlugin
func (dp *decryptionPlugin) Close() error {
	dp.logger.Debug("DecryptionReporting Close", nil)
	return nil
}
