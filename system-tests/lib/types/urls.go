package types

// all structs are copies of identical structs in ${CRIB_REPO}/dependencies/donut/scripts/urls/main.go
// in the future we should move these types to a dedicated module that would be imported both by CRIB and this module
type JdURLs struct {
	GRPCExternalURL string `json:"grpc_host_url"`
	GRCPInternalURL string `json:"grpc_internal_url"`
	WSExternalURL   string `json:"ws_host_url"`
	WSInternalURL   string `json:"ws_internal_url"`
}

type DonURL struct {
	ExternalURL    string `json:"host_url"`
	InternalURL    string `json:"internal_url"`
	P2PInternalURL string `json:"p2p_internal_url"`
	InternalIP     string `json:"internal_ip"`
}

type DonURLs struct {
	BootstrapNodes []DonURL `json:"bootstrap_nodes"`
	WorkerNodes    []DonURL `json:"worker_nodes"`
}

type DonAPICredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChainURLs struct {
	HTTPExternalURL string `json:"http_host_url"`
	WSExternalURL   string `json:"ws_host_url"`
	HTTPInternalURL string `json:"http_internal_url"`
	WSInternalURL   string `json:"ws_internal_url"`
}
