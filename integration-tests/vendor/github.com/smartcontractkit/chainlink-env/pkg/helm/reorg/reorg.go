package reorg

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/client"
	"github.com/smartcontractkit/chainlink-env/config"
	"github.com/smartcontractkit/chainlink-env/environment"
)

const (
	URLsKey            = "geth"
	TXNodesAppLabel    = "geth-ethereum-geth"
	MinerNodesAppLabel = "geth-ethereum-miner-node"
)

type Props struct {
	NetworkName string `envconfig:"network_name"`
	NetworkType string `envconfig:"network_type"`
	Values      map[string]interface{}
}

type Chart struct {
	Name    string
	Path    string
	Version string
	Props   *Props
	Values  *map[string]interface{}
}

func (m Chart) IsDeploymentNeeded() bool {
	return true
}

func (m Chart) GetName() string {
	return m.Name
}

func (m Chart) GetProps() interface{} {
	return m.Props
}

func (m Chart) GetPath() string {
	return m.Path
}

func (m Chart) GetVersion() string {
	return m.Version
}

func (m Chart) GetValues() *map[string]interface{} {
	return m.Values
}

func (m Chart) ExportData(e *environment.Environment) error {
	urls := make([]string, 0)
	minerPods, err := e.Client.ListPods(e.Cfg.Namespace, fmt.Sprintf("app=%s-ethereum-miner-node", m.Props.NetworkName))
	if err != nil {
		return err
	}
	txPods, err := e.Client.ListPods(e.Cfg.Namespace, fmt.Sprintf("app=%s-ethereum-geth", m.Props.NetworkName))
	if err != nil {
		return err
	}
	if len(txPods.Items) > 0 {
		for i := range txPods.Items {
			podName := fmt.Sprintf("%s-ethereum-geth:%d", m.Props.NetworkName, i)
			txNodeLocalWS, err := e.Fwd.FindPort(podName, "geth", "ws-rpc").As(client.LocalConnection, client.WS)
			if err != nil {
				return err
			}
			txNodeInternalWs, err := e.Fwd.FindPort(podName, "geth", "ws-rpc").As(client.RemoteConnection, client.WS)
			if err != nil {
				return err
			}
			if e.Cfg.InsideK8s {
				urls = append(urls, txNodeInternalWs)
				log.Info().Str("URL", txNodeInternalWs).Msgf("Geth network (TX Node) - %d", i)
			} else {
				urls = append(urls, txNodeLocalWS)
				log.Info().Str("URL", txNodeLocalWS).Msgf("Geth network (TX Node) - %d", i)
			}
		}
	}

	if len(minerPods.Items) > 0 {
		for i := range minerPods.Items {
			podName := fmt.Sprintf("%s-ethereum-miner-node:%d", m.Props.NetworkName, i)
			minerNodeLocalWS, err := e.Fwd.FindPort(podName, "geth-miner", "ws-rpc-miner").As(client.LocalConnection, client.WS)
			if err != nil {
				return err
			}
			minerNodeInternalWs, err := e.Fwd.FindPort(podName, "geth-miner", "ws-rpc-miner").As(client.RemoteConnection, client.WS)
			if err != nil {
				return err
			}
			if e.Cfg.InsideK8s {
				urls = append(urls, minerNodeInternalWs)
				log.Info().Str("URL", minerNodeInternalWs).Msgf("Geth network (Miner Node) - %d", i)
			} else {
				urls = append(urls, minerNodeLocalWS)
				log.Info().Str("URL", minerNodeLocalWS).Msgf("Geth network (Miner Node) - %d", i)
			}
		}
	}

	e.URLs[m.Props.NetworkName] = urls
	return nil
}

func defaultProps() *Props {
	return &Props{
		NetworkName: "geth",
		NetworkType: "geth-reorg",
		Values: map[string]interface{}{
			"imagePullPolicy": "IfNotPresent",
			"bootnode": map[string]interface{}{
				"replicas": "2",
				"image": map[string]interface{}{
					"repository": "ethereum/client-go",
					"tag":        "alltools-v1.10.25",
				},
			},
			"bootnodeRegistrar": map[string]interface{}{
				"replicas": "1",
				"image": map[string]interface{}{
					"repository": "jpoon/bootnode-registrar",
					"tag":        "v1.0.0",
				},
			},
			"geth": map[string]interface{}{
				"image": map[string]interface{}{
					"repository": "ethereum/client-go",
					"tag":        "v1.10.25",
				},
				"tx": map[string]interface{}{
					"replicas": "1",
					"service": map[string]interface{}{
						"type": "ClusterIP",
					},
				},
				"miner": map[string]interface{}{
					"replicas": "2",
					"account": map[string]interface{}{
						"secret": "",
					},
				},
				"genesis": map[string]interface{}{
					"networkId": "1337",
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
	config.MustMerge(targetProps, props)
	config.MustMerge(&targetProps.Values, props.Values)
	return Chart{
		Name:    targetProps.NetworkName,
		Path:    "chainlink-qa/ethereum",
		Values:  &targetProps.Values,
		Props:   targetProps,
		Version: helmVersion,
	}
}
