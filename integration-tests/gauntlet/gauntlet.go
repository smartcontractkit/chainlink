package gauntlet

import (
	"encoding/json"
	"fmt"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/gauntlet"
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
				Type  string `json:"type"`
				Nonce string `json:"nonce"`
				Hash  string `json:"hash"`
			} `json:"tx"`
		} `json:"tx"`
		Contract string `json:"contract"`
	} `json:"responses"`
}

// New Creates a default gauntlet config
func New(workingDir string) (*Gauntlet, error) {
	config, err := gauntlet.NewGauntlet("gauntlet", "")
	config.SetWorkingDir(workingDir)
	if err != nil {
		return nil, err
	}

	return &Gauntlet{
		dir:       workingDir,
		Config:    config,
		Contracts: &Contracts{},
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
func (g *Gauntlet) SetupNetwork(addr string, account string, privateKey string) error {
	err := os.Setenv("NODE_URL", addr)
	if err != nil {
		return err
	}
	err = os.Setenv("ACCOUNT", account)
	if err != nil {
		return err
	}
	err = os.Setenv("PRIVATE_KEY", privateKey)
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
