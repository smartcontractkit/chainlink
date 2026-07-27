package ui

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/nodeconfig"
)

type Server struct {
	desiredPath    string
	statePath      string
	chartValuesDir string
	env            string
	namespace      string
	kubeconfig     string
	log            zerolog.Logger
	mux            *http.ServeMux
}

func NewServer(desiredPath, statePath, chartValuesDir, env, namespace, kubeconfig string, log zerolog.Logger) *Server {
	s := &Server{
		desiredPath:    desiredPath,
		statePath:      statePath,
		chartValuesDir: chartValuesDir,
		env:            env,
		namespace:      namespace,
		kubeconfig:     kubeconfig,
		log:            log,
		mux:            http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/api/nodes", s.handleNodes)
	s.mux.HandleFunc("/api/chains/discovered", s.handleChainsDiscovered)
	s.mux.HandleFunc("/api/desired", s.handleDesired)
	s.mux.HandleFunc("/api/state", s.handleState)
	s.mux.HandleFunc("/api/capabilities", s.handleCapabilities)
	s.mux.HandleFunc("/api/preview-toml", s.handlePreviewTOML)
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/jd/check", s.handleJDCheck)
	s.mux.HandleFunc("/", s.handleStatic)
}

func (s *Server) ListenAndServe(addr string) error {
	s.log.Info().Msgf("Web UI available at http://%s", addr)
	srv := &http.Server{
		Addr:              addr,
		Handler:           s.withCORS(s.mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return srv.ListenAndServe()
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type NodeResponse struct {
	Name         string `json:"name"`
	NodeType     string `json:"nodeType"`
	DonName      string `json:"donName,omitempty"`
	ChartDonName string `json:"chartDonName,omitempty"`
}

type NodesResponse struct {
	Nodes     []NodeResponse `json:"nodes"`
	Namespace string         `json:"namespace"`
	Bootstrap string         `json:"bootstrap"`
	Gateways  []string       `json:"gateways"`
	ChartDir  string         `json:"chartDir"`
}

// ChainResponse mirrors domain.Chain for the UI's Chains section.
type ChainResponse struct {
	ChainID  uint64 `json:"chainId"`
	WSURL    string `json:"wsUrl"`
	HTTPURL  string `json:"httpUrl"`
	Registry bool   `json:"registry"`
}

type DesiredResponse struct {
	Infra struct {
		Type        string `json:"type"`
		ChartValues string `json:"chartValues"`
		Namespace   string `json:"namespace"`
	} `json:"infra"`
	JD struct {
		GRPC        string `json:"grpc"`
		Domain      string `json:"domain"`
		Environment string `json:"environment"`
		UseTLS      bool   `json:"useTLS"`
	} `json:"jd"`
	Chains            []ChainResponse                    `json:"chains"`
	DONs              []DONResponse                      `json:"dons"`
	GatewayNodes      []GatewayNodeResponse              `json:"gatewayNodes"`
	CapabilityConfigs map[string]domain.CapabilityConfig `json:"capabilityConfigs"`
}

type GatewayNodeResponse struct {
	Node string `json:"node"`
	DON  string `json:"don"`
}

type DONResponse struct {
	Name                   string                             `json:"name"`
	DONTypes               []string                           `json:"donTypes"`
	Capabilities           []string                           `json:"capabilities"`
	Nodes                  []string                           `json:"nodes"`
	BootstrapNode          string                             `json:"bootstrapNode"`
	ExposesRemoteCaps      bool                               `json:"exposesRemoteCapabilities"`
	RegistryBasedAllowlist []string                           `json:"registryBasedLaunchAllowlist"`
	CapabilityConfigs      map[string]domain.CapabilityConfig `json:"capabilityConfigs,omitempty"`
}

type StateResponse struct {
	Phase     string              `json:"phase"`
	Addresses []domain.AddressRef `json:"addresses"`
	DONIDs    map[string]uint64   `json:"donIds"`
	HasState  bool                `json:"hasState"`
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleNodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cv, err := domain.LoadChartValues(s.chartValuesDir, s.env)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "failed to load chart values: "+err.Error())
		return
	}

	var nodes []NodeResponse
	for _, n := range cv.Nodes {
		resp := NodeResponse{
			Name:         n.Name,
			NodeType:     string(n.NodeType),
			ChartDonName: n.DONName,
			DonName:      n.DONName,
		}
		nodes = append(nodes, resp)
	}

	var gateways []string
	for _, n := range cv.Nodes {
		if n.NodeType == domain.RoleGateway {
			gateways = append(gateways, n.Name)
		}
	}

	s.writeJSON(w, http.StatusOK, NodesResponse{
		Nodes:     nodes,
		Namespace: s.namespace,
		Bootstrap: cv.FindBootstrap(),
		Gateways:  gateways,
		ChartDir:  s.chartValuesDir,
	})
}

// handleChainsDiscovered returns EVM chains parsed from the chart nodes'
// existing typed config, for prepopulating the UI's Chains section. This is
// prepopulation only — the authoritative chain list is whatever is saved
// into desired.toml via /api/desired.
func (s *Server) handleChainsDiscovered(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cv, err := domain.LoadChartValues(s.chartValuesDir, s.env)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "failed to load chart values: "+err.Error())
		return
	}

	chains, err := nodeconfig.DiscoverChains(cv)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "failed to discover chains: "+err.Error())
		return
	}

	resp := make([]ChainResponse, 0, len(chains))
	for _, ch := range chains {
		resp = append(resp, ChainResponse{ChainID: ch.ChainID, WSURL: ch.WSURL, HTTPURL: ch.HTTPURL})
	}

	s.writeJSON(w, http.StatusOK, map[string]any{"chains": resp})
}

func (s *Server) handleDesired(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ds, err := domain.LoadDesiredState(s.desiredPath)
		if err != nil {
			if os.IsNotExist(errors.Cause(err)) {
				s.writeJSON(w, http.StatusOK, DesiredResponse{})
				return
			}
			s.writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		resp := DesiredResponse{
			CapabilityConfigs: ds.CapabilityConfigs,
		}
		resp.Infra.Type = ds.Infra.Type
		resp.Infra.ChartValues = ds.Infra.ChartValues
		resp.Infra.Namespace = ds.Infra.Namespace
		resp.JD.GRPC = ds.JD.GRPC
		resp.JD.Domain = ds.JD.Domain
		resp.JD.Environment = ds.JD.Environment
		resp.JD.UseTLS = ds.JD.UseTLS

		for _, ch := range ds.Chains {
			resp.Chains = append(resp.Chains, ChainResponse{
				ChainID:  ch.ChainID,
				WSURL:    ch.WSURL,
				HTTPURL:  ch.HTTPURL,
				Registry: ch.Registry,
			})
		}

		cv, cvErr := domain.LoadChartValues(s.chartValuesDir, s.env)

		for _, don := range ds.DONs {
			donResp := DONResponse{
				Name:                   don.Name,
				DONTypes:               don.DONTypes,
				Capabilities:           don.Capabilities,
				BootstrapNode:          don.BootstrapNode,
				ExposesRemoteCaps:      don.ExposesRemoteCaps,
				RegistryBasedAllowlist: don.RegistryBasedAllowlist,
				CapabilityConfigs:      don.CapabilityConfigs,
			}
			if cvErr == nil && cv != nil {
				donResp.Nodes = cv.NodeNamesForDONName(don.Name)
			} else {
				donResp.Nodes = []string{}
			}
			resp.DONs = append(resp.DONs, donResp)
		}

		for _, gwn := range ds.GatewayNodes {
			resp.GatewayNodes = append(resp.GatewayNodes, GatewayNodeResponse{Node: gwn.Node, DON: gwn.DON})
		}

		s.writeJSON(w, http.StatusOK, resp)

	case http.MethodPost:
		var req DesiredResponse
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}

		ds := responseToDesiredState(req)
		if err := ds.Validate(); err != nil {
			s.writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := storeDesiredState(s.desiredPath, ds); err != nil {
			s.writeError(w, http.StatusInternalServerError, "failed to save: "+err.Error())
			return
		}

		s.writeJSON(w, http.StatusOK, map[string]string{"status": "saved"})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	st, err := domain.LoadState(s.statePath)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if st == nil {
		s.writeJSON(w, http.StatusOK, StateResponse{HasState: false})
		return
	}

	s.writeJSON(w, http.StatusOK, StateResponse{
		Phase:     string(st.Phase),
		Addresses: st.Addresses,
		DONIDs:    st.DONIDs,
		HasState:  true,
	})
}

func (s *Server) handleCapabilities(w http.ResponseWriter, _ *http.Request) {
	caps := []map[string]any{
		{"name": "cron", "label": "Cron", "description": "Schedule-based workflow triggers", "chainScoped": false},
		{"name": "evm", "label": "EVM", "description": "EVM chain interaction (logs, calls, writes)", "chainScoped": true},
		{"name": "http-action", "label": "HTTP Action", "description": "Gateway-routed HTTP requests", "chainScoped": false},
		{"name": "http-trigger", "label": "HTTP Trigger", "description": "Gateway-routed HTTP-triggered workflows", "chainScoped": false},
		{"name": "vault", "label": "Vault", "description": "Gateway-routed secret management", "chainScoped": false},
		{"name": "consensus", "label": "Consensus", "description": "OCR3 consensus for offchain reporting", "chainScoped": false},
		{"name": "don-time", "label": "DON Time", "description": "DON-wide timestamp capability", "chainScoped": false},
		{"name": "solana", "label": "Solana", "description": "Solana chain interaction", "chainScoped": true},
	}

	defaults := LoadCapabilityDefaults()

	s.writeJSON(w, http.StatusOK, map[string]any{
		"capabilities": caps,
		"defaults":     defaults,
	})
}

func (s *Server) handlePreviewTOML(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DesiredResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	ds := responseToDesiredState(req)
	data, err := toml.Marshal(ds)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "failed to generate TOML: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write(data); err != nil {
		s.log.Warn().Err(err).Msg("failed to write preview TOML response")
	}
}

// JDCheckRequest deliberately has no AccessToken field — the JD bearer token
// is a live credential and must come only from the server's environment
// (infra.JDAccessToken), never from the browser.
type JDCheckRequest struct {
	GRPC      string   `json:"grpc"`
	UseTLS    bool     `json:"useTLS"`
	Namespace string   `json:"namespace"`
	NodeNames []string `json:"nodeNames"`
}

type JDCheckResponseFull struct {
	infra.JDCheckResponse
	K8sErrors    []string `json:"k8sErrors,omitempty"`
	CSAKeysFound int      `json:"csaKeysFound,omitempty"`
}

func (s *Server) handleJDCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req JDCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if req.GRPC == "" {
		s.writeError(w, http.StatusBadRequest, "grpc endpoint is required")
		return
	}

	// Fail fast with a clear, actionable message instead of letting the dial attempt
	// below fail with a confusing gRPC/auth error — JD always requires this token.
	if infra.JDAccessToken() == "" {
		s.log.Warn().Msg("JD check requested but " + infra.JDAccessTokenEnv + " is not set on the server")
		s.writeJSON(w, http.StatusOK, JDCheckResponseFull{
			JDCheckResponse: infra.JDCheckResponse{
				Connected: false,
				Error:     infra.JDAccessTokenEnv + " is not set on the server — export it and restart `serve`",
			},
		})
		return
	}

	namespace := req.Namespace
	if namespace == "" {
		namespace = s.namespace
	}

	s.log.Info().
		Str("grpc", req.GRPC).
		Str("namespace", namespace).
		Str("kubeconfig", s.kubeconfig).
		Int("nodeNames", len(req.NodeNames)).
		Bool("hasToken", infra.JDAccessToken() != "").
		Msg("JD check starting")

	nodeCSAKeys := make(map[string]string)
	k8sErrors := []string{}
	var mu sync.Mutex

	if namespace != "" && len(req.NodeNames) > 0 {
		s.log.Info().Str("namespace", namespace).Str("kubeconfig", s.kubeconfig).Msg("Creating K8s client")
		k8sClient, err := infra.NewK8sClient(s.kubeconfig, namespace, s.log)
		if err != nil {
			k8sErrors = append(k8sErrors, fmt.Sprintf("K8s client creation failed: %v", err))
			s.log.Error().Err(err).Str("kubeconfig", s.kubeconfig).Str("namespace", namespace).Msg("failed to create K8s client")
		} else {
			s.log.Info().Msg("K8s client created successfully")
			cv, cvErr := domain.LoadChartValues(s.chartValuesDir, s.env)

			g, gctx := errgroup.WithContext(r.Context())
			g.SetLimit(min(len(req.NodeNames), 8))
			for _, nodeName := range req.NodeNames {
				g.Go(func() error {
					nodeNamespace := namespace
					if cvErr == nil && cv != nil {
						nodeNamespace = cv.GetNodeNamespace(nodeName)
					}
					s.log.Info().Str("node", nodeName).Str("namespace", nodeNamespace).Msg("Fetching CSA key via node API")
					csaKey, err := k8sClient.GetNodeCSAKey(gctx, nodeName, nodeNamespace)
					if err != nil {
						mu.Lock()
						k8sErrors = append(k8sErrors, fmt.Sprintf("node %s: %v", nodeName, err))
						mu.Unlock()
						s.log.Error().Err(err).Str("node", nodeName).Msg("failed to get CSA key")
						return nil
					}
					if csaKey != "" {
						mu.Lock()
						nodeCSAKeys[nodeName] = csaKey
						mu.Unlock()
						s.log.Info().Str("node", nodeName).Str("csaKey", csaKey[:min(20, len(csaKey))]+"...").Msg("Got CSA key")
					}
					return nil
				})
			}
			_ = g.Wait()
		}
	} else {
		s.log.Warn().Str("namespace", namespace).Int("nodeNames", len(req.NodeNames)).Msg("Skipping K8s: namespace or nodeNames empty")
		k8sErrors = append(k8sErrors, fmt.Sprintf("namespace=%q or nodeNames empty — cannot read K8s secrets", namespace))
	}

	s.log.Info().Int("csaKeysFound", len(nodeCSAKeys)).Int("k8sErrors", len(k8sErrors)).Msg("K8s phase done, checking JD")

	jdResult := infra.CheckJD(r.Context(), req.GRPC, infra.JDAccessToken(), req.UseTLS, nodeCSAKeys, s.log)

	resp := JDCheckResponseFull{
		JDCheckResponse: jdResult,
		K8sErrors:       k8sErrors,
		CSAKeysFound:    len(nodeCSAKeys),
	}

	s.writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	cleanPath := strings.TrimPrefix(path, "/")
	content, ok := webFiles[cleanPath]
	if !ok {
		http.NotFound(w, r)
		return
	}

	setContentType(w, cleanPath)
	if _, err := w.Write([]byte(content)); err != nil {
		s.log.Warn().Err(err).Str("path", cleanPath).Msg("failed to write static file response")
	}
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		s.log.Warn().Err(err).Msg("failed to encode JSON response")
	}
}

func (s *Server) writeError(w http.ResponseWriter, status int, msg string) {
	s.writeJSON(w, status, map[string]string{"error": msg})
}

func setContentType(w http.ResponseWriter, path string) {
	switch filepath.Ext(path) {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
}

func normalizeCapabilityConfigNumbers(configs map[string]domain.CapabilityConfig) {
	for _, cc := range configs {
		normalizeJSONNumberMap(cc.Values)
	}
}

func normalizeJSONNumberMap(m map[string]any) {
	for k, v := range m {
		m[k] = normalizeJSONNumberValue(v)
	}
}

func normalizeJSONNumberValue(v any) any {
	switch val := v.(type) {
	case float64:
		if !math.IsInf(val, 0) && !math.IsNaN(val) && val == math.Trunc(val) {
			return int64(val)
		}
		return val
	case map[string]any:
		normalizeJSONNumberMap(val)
		return val
	case []any:
		for i, e := range val {
			val[i] = normalizeJSONNumberValue(e)
		}
		return val
	default:
		return v
	}
}

func responseToDesiredState(req DesiredResponse) *domain.DesiredState {
	normalizeCapabilityConfigNumbers(req.CapabilityConfigs)

	ds := &domain.DesiredState{
		Infra: domain.Infra{
			Type:        req.Infra.Type,
			ChartValues: req.Infra.ChartValues,
			Namespace:   req.Infra.Namespace,
		},
		JD: domain.JDConfig{
			GRPC:        req.JD.GRPC,
			Domain:      req.JD.Domain,
			Environment: req.JD.Environment,
			UseTLS:      req.JD.UseTLS,
		},
		CapabilityConfigs: req.CapabilityConfigs,
	}

	if ds.Infra.Type == "" {
		ds.Infra.Type = "griddle"
	}

	for _, ch := range req.Chains {
		ds.Chains = append(ds.Chains, domain.Chain{
			ChainID:  ch.ChainID,
			WSURL:    ch.WSURL,
			HTTPURL:  ch.HTTPURL,
			Registry: ch.Registry,
		})
	}

	for _, d := range req.DONs {
		normalizeCapabilityConfigNumbers(d.CapabilityConfigs)
		ds.DONs = append(ds.DONs, domain.DON{
			Name:                   d.Name,
			DONTypes:               d.DONTypes,
			Capabilities:           d.Capabilities,
			BootstrapNode:          d.BootstrapNode,
			ExposesRemoteCaps:      d.ExposesRemoteCaps,
			RegistryBasedAllowlist: d.RegistryBasedAllowlist,
			CapabilityConfigs:      d.CapabilityConfigs,
		})
	}

	for _, gwn := range req.GatewayNodes {
		ds.GatewayNodes = append(ds.GatewayNodes, domain.GatewayNodeAssignment{Node: gwn.Node, DON: gwn.DON})
	}

	return ds
}

func storeDesiredState(path string, ds *domain.DesiredState) error {
	data, err := toml.Marshal(ds)
	if err != nil {
		return errors.Wrap(err, "failed to marshal desired state")
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.Wrapf(err, "failed to create dir %s", dir)
		}
	}

	return os.WriteFile(path, data, 0600)
}
