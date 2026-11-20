package web

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

// ReloadController applies the Dynamic Config file the sidecar rendered.
type ReloadController struct {
	App chainlink.Application
}

// ReloadResponse is what the sidecar's reload hook checks: the hash of the
// config the node actually loaded, and the revision it came from.
type ReloadResponse struct {
	ConfigHash string `json:"configHash"`
	Revision   uint64 `json:"revision"`
}

// Reload re-reads the Dynamic Config file and reconciles the node's jobs
// against it. It is reachable only from loopback: the sidecar shares the Pod's
// network namespace, and nothing outside the Pod may drive the node's jobs.
// Example:
// "POST <application>/reload"
func (rc *ReloadController) Reload(c *gin.Context) {
	if !isLoopback(c.Request.RemoteAddr) {
		jsonAPIError(c, http.StatusForbidden, errors.New("reload is only allowed from localhost"))

		return
	}

	result, err := rc.App.ReloadDynamicConfig(c.Request.Context())
	if err != nil {
		rc.App.GetLogger().Errorw("dynamic config reload failed", "err", err)
		jsonAPIError(c, http.StatusInternalServerError, err)

		return
	}
	hash, err := configHash(result.Raw)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)

		return
	}
	c.JSON(http.StatusOK, ReloadResponse{ConfigHash: hash, Revision: result.Result.Revision})
}

// configHash is the hash of the config the node loaded, computed the way the
// sidecar's verifier does: sha256 over the canonical JSON. The sidecar only
// records the reload as applied when the two match.
func configHash(raw []byte) (string, error) {
	canonical, err := canonicalize(raw)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)

	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// canonicalize re-encodes a JSON object with sorted keys and no insignificant
// whitespace, matching the sidecar verifier's canonical form.
func canonicalize(raw []byte) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return nil, errors.New("config has trailing data")
	}
	if _, ok := v.(map[string]any); !ok {
		return nil, errors.New("config must be a JSON object")
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

func isLoopback(remote string) bool {
	host, _, err := net.SplitHostPort(remote)
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)

	return ip != nil && ip.IsLoopback()
}
