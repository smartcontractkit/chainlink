package legacygasstation

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation/types"
)

type StatusUpdater struct {
	endpointURL string
	client      *http.Client
	lggr        logger.Logger
}

func NewStatusUpdater(endpointURL, mtlsCert, mtlsKey string, lggr logger.Logger) (*StatusUpdater, error) {
	if mtlsCert == "" || mtlsKey == "" {
		lggr.Info("Instantiating status updater without mTLS")
		return &StatusUpdater{
			endpointURL: endpointURL,
			client:      &http.Client{},
			lggr:        lggr,
		}, nil
	}

	cert, err := tls.X509KeyPair([]byte(mtlsCert), []byte(mtlsKey))
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return &StatusUpdater{
		endpointURL: endpointURL,
		client:      client,
		lggr:        lggr,
	}, nil
}

func (s *StatusUpdater) Update(tx types.LegacyGaslessTx) error {
	req, err := tx.ToSendTransactionStatusRequest()
	if err != nil {
		return errors.Wrap(err, "to send transaction status request")
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "json marshal request body")
	}

	resp, err := s.client.Post(s.endpointURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return errors.Wrap(err, "post failed")
	}

	s.lggr.Infof("posted request: %+v. raw response: %+v", req, resp)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		s.lggr.Infow("POST succeeded", "status", resp.Status)
		return nil
	}

	// best-effort attempt to parse error message
	var jsonResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		s.lggr.Warnw("decode json response", "err", err)
		return errors.Errorf("error while calling send transaction status API status: %s", resp.Status)
	}

	respErr, ok := jsonResp["error"]
	if !ok {
		s.lggr.Infow("no error found in response")
		return errors.Errorf("error while calling send transaction status API status: %s", resp.Status)
	}

	respErrString, ok := respErr.(string)
	if !ok {
		s.lggr.Warn("failed to parse error in response")
		return errors.Errorf("error while calling send transaction status API status: %s", resp.Status)
	}

	return errors.Errorf("error while calling send transaction status API status: %s, error: %s", resp.Status, respErrString)
}
