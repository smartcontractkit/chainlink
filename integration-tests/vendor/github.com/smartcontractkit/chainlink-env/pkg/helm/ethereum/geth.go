package ethereum

import (
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/client"
	"github.com/smartcontractkit/chainlink-env/config"
	"github.com/smartcontractkit/chainlink-env/environment"
)

type Props struct {
	NetworkName string   `envconfig:"network_name"`
	Simulated   bool     `envconfig:"network_simulated"`
	HttpURLs    []string `envconfig:"http_url"`
	WsURLs      []string `envconfig:"ws_url"`
	Values      map[string]interface{}
}

type HelmProps struct {
	Name    string
	Path    string
	Version string
	Values  *map[string]interface{}
}

type Chart struct {
	HelmProps *HelmProps
	Props     *Props
}

func (m Chart) IsDeploymentNeeded() bool {
	return m.Props.Simulated
}

func (m Chart) GetProps() interface{} {
	return m.Props
}

func (m Chart) GetName() string {
	return m.HelmProps.Name
}

func (m Chart) GetPath() string {
	return m.HelmProps.Path
}

func (m Chart) GetVersion() string {
	return m.HelmProps.Version
}

func (m Chart) GetValues() *map[string]interface{} {
	return m.HelmProps.Values
}

func (m Chart) ExportData(e *environment.Environment) error {
	if m.Props.Simulated {
		gethLocalHttp, err := e.Fwd.FindPort("geth:0", "geth-network", "http-rpc").As(client.LocalConnection, client.HTTP)
		if err != nil {
			return err
		}
		gethInternalHttp, err := e.Fwd.FindPort("geth:0", "geth-network", "http-rpc").As(client.RemoteConnection, client.HTTP)
		if err != nil {
			return err
		}
		gethLocalWs, err := e.Fwd.FindPort("geth:0", "geth-network", "ws-rpc").As(client.LocalConnection, client.WS)
		if err != nil {
			return err
		}
		gethInternalWs, err := e.Fwd.FindPort("geth:0", "geth-network", "ws-rpc").As(client.RemoteConnection, client.WS)
		if err != nil {
			return err
		}
		if e.Cfg.InsideK8s {
			e.URLs[m.Props.NetworkName] = []string{gethInternalWs}
		} else {
			e.URLs[m.Props.NetworkName] = []string{gethLocalWs}
		}
		e.URLs[m.Props.NetworkName+"_http"] = []string{gethLocalHttp}

		// For cases like starknet we need the internalHttp address to set up the L1<>L2 interaction
		e.URLs[m.Props.NetworkName+"_internal"] = []string{gethInternalWs}
		e.URLs[m.Props.NetworkName+"_internal_http"] = []string{gethInternalHttp}

		log.Info().Str("Name", "Geth").Str("URLs", gethLocalWs).Msg("Geth network")
	} else {
		e.URLs[m.Props.NetworkName] = m.Props.WsURLs
		log.Info().Str("Name", m.Props.NetworkName).Strs("URLs", m.Props.WsURLs).Msg("Ethereum network")
	}
	return nil
}

func defaultProps() *Props {
	return &Props{
		NetworkName: "Simulated Geth",
		Simulated:   true,
		Values: map[string]interface{}{
			"replicas": "1",
			"geth": map[string]interface{}{
				"image": map[string]interface{}{
					"image":   "ethereum/client-go",
					"version": "v1.10.25",
				},
			},
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    "1000m",
					"memory": "768Mi",
				},
				"limits": map[string]interface{}{
					"cpu":    "1000m",
					"memory": "768Mi",
				},
			},
		},
	}
}

func New(props *Props) environment.ConnectedChart {
	return NewVersioned("", props)
}

// NewVersioned enables choosing a specific helm chart version
func NewVersioned(helmVersion string, props *Props) environment.ConnectedChart {
	targetProps := defaultProps()
	if props == nil {
		props = targetProps
	}
	config.MustMerge(targetProps, props)
	config.MustMerge(&targetProps.Values, props.Values)
	targetProps.Simulated = props.Simulated // Mergo has issues with boolean merging for simulated networks
	if targetProps.Simulated {
		return Chart{
			HelmProps: &HelmProps{
				Name:   "geth",
				Path:   "chainlink-qa/geth",
				Values: &targetProps.Values,
			},
			Props: targetProps,
		}
	}
	return Chart{
		Props: targetProps,
		HelmProps: &HelmProps{
			Version: helmVersion,
		},
	}
}
