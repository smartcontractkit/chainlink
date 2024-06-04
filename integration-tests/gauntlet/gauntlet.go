package gauntlet

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/gauntlet"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/gauntlet/configs"
)

type Gauntlet struct {
	dir       string
	Config    *gauntlet.Gauntlet
	Contracts *Contracts
	options   *gauntlet.ExecCommandOptions
}

type Contracts struct {
	LinkContract            *LinkContract
	OCRContract             *OCRContract
	AccessControllerAddress string
}

type LinkContract struct {
	Address  string
	Contract contracts.LinkToken
}

type OCRContract struct {
	Address  string
	Contract contracts.OffchainAggregator
}

// Response Default response output for gauntlet commands
type Response struct {
	Responses []struct {
		Tx struct {
			Hash    string `json:"hash"`
			Address string `json:"address"`
			Status  string `json:"status"`

			Tx struct {
				Type  int    `json:"type"`
				Nonce int    `json:"nonce"`
				Hash  string `json:"hash"`
			} `json:"tx"`
		} `json:"tx"`
		Contract string `json:"contract"`
	} `json:"responses"`
}

// New Creates a default gauntlet config
func New(binaryName string, workingDir string) (*Gauntlet, error) {
	config, err := gauntlet.NewGauntlet("yarn", binaryName)
	config.SetWorkingDir(workingDir)
	if err != nil {
		return nil, err
	}

	return &Gauntlet{
		dir:    workingDir,
		Config: config,
		Contracts: &Contracts{
			LinkContract: &LinkContract{
				Address:  "",
				Contract: nil,
			},
			OCRContract: &OCRContract{
				Address:  "",
				Contract: nil,
			},
			AccessControllerAddress: "",
		},
		options: &gauntlet.ExecCommandOptions{
			ErrHandling:       []string{},
			CheckErrorsInRead: true,
		},
	}, nil
}

// FetchGauntletJsonOutput Parse gauntlet json response that is generated after yarn gauntlet command execution
func (g *Gauntlet) FetchGauntletJsonOutput() (*Response, error) {
	var payload = &Response{}
	gauntletOutput, err := os.ReadFile(g.dir + "report.json")
	if err != nil {
		return payload, err
	}
	err = json.Unmarshal(gauntletOutput, &payload)
	if err != nil {
		return payload, err
	}
	return payload, nil
}

// SetupNetwork Sets up a new network and sets the NODE_URL for the RPC
func (g *Gauntlet) SetupNetwork(addr string, privateKey string) error {
	g.Config.AddNetworkConfigVar("NODE_URL", addr)
	g.Config.AddNetworkConfigVar("PRIVATE_KEY", privateKey)

	err := g.Config.WriteNetworkConfigMap(g.dir + "networks/")
	if err != nil {
		return err
	}

	return nil
}

func (g *Gauntlet) DeployLinkToken() error {
	_, err := g.Config.ExecCommand([]string{"token:deploy"}, *g.options)
	if err != nil {
		return err
	}
	res, err := g.FetchGauntletJsonOutput()
	if err != nil {
		return err
	}

	g.Contracts.LinkContract.Address = res.Responses[0].Tx.Address
	return nil
}

func (g *Gauntlet) DeployAccessController() error {
	_, err := g.Config.ExecCommand([]string{"access_controller:deploy"}, *g.options)
	if err != nil {
		return err
	}
	res, err := g.FetchGauntletJsonOutput()
	if err != nil {
		return err
	}

	g.Contracts.AccessControllerAddress = res.Responses[0].Tx.Address
	return nil
}

func (g *Gauntlet) DeployOCR(ocrContractValues string) error {
	_, err := g.Config.ExecCommand([]string{"ocr:deploy", fmt.Sprintf("--input=%s", ocrContractValues)}, *g.options)
	if err != nil {
		return err
	}
	res, err := g.FetchGauntletJsonOutput()
	if err != nil {
		return err
	}

	g.Contracts.OCRContract.Address = res.Responses[0].Tx.Address
	return nil
}

func (g *Gauntlet) AddAccess(ocrAddress string) error {
	_, err := g.Config.ExecCommand([]string{"access_controller:add_access", fmt.Sprintf("--address=%s", ocrAddress), g.Contracts.AccessControllerAddress}, *g.options)
	if err != nil {
		return err
	}
	_, err = g.FetchGauntletJsonOutput()
	if err != nil {
		return err
	}

	return nil
}

func (g *Gauntlet) SetPayees(ocrAddress string, payees []string, transmitters []string) error {
	_, err := g.Config.ExecCommand([]string{
		"ocr:set_payees",
		ocrAddress,
		fmt.Sprintf("--transmitters=%s", strings.Join(transmitters, ",")),
		fmt.Sprintf("--payees=%s", strings.Join(payees, ","))}, *g.options)
	if err != nil {
		return err
	}
	_, err = g.FetchGauntletJsonOutput()
	if err != nil {
		return err
	}

	return nil
}

func (g *Gauntlet) SetConfig(ocrAddress string, ocrConfigValues string) error {
	_, err := g.Config.ExecCommand([]string{"ocr:set_config", ocrAddress, fmt.Sprintf("--input=%s", ocrConfigValues)}, *g.options)
	if err != nil {
		return err
	}
	_, err = g.FetchGauntletJsonOutput()
	if err != nil {
		return err
	}

	return nil
}

func (g *Gauntlet) DeployContracts(ocrConfig *configs.OCRConfig, transmitters []string, signers []string, peerIDs []string, payees []string) error {
	err := g.DeployLinkToken()
	if err != nil {
		return err
	}

	err = g.DeployAccessController()
	if err != nil {
		return err
	}

	ocrConfig.Contract.Link = g.Contracts.LinkContract.Address
	ocrConfig.Contract.BillingAccessController = g.Contracts.AccessControllerAddress
	ocrConfig.Contract.RequesterAccessController = g.Contracts.AccessControllerAddress
	ocrConfig.Config.Transmitters = transmitters
	ocrConfig.Config.Signers = signers
	ocrConfig.Config.OperatorsPeerIds = strings.Join(peerIDs, ",")

	ocrJsonContract, err := ocrConfig.MarshalContract()
	if err != nil {
		return err
	}

	err = g.DeployOCR(ocrJsonContract)
	if err != nil {
		return err
	}

	err = g.AddAccess(g.Contracts.OCRContract.Address)
	if err != nil {
		return err
	}

	err = g.SetPayees(g.Contracts.OCRContract.Address, payees, transmitters)
	if err != nil {
		return err
	}

	ocrJsonConfig, err := ocrConfig.MarshalConfig()
	if err != nil {
		return err
	}
	err = g.SetConfig(g.Contracts.OCRContract.Address, ocrJsonConfig)
	if err != nil {
		return err
	}

	return nil
}
