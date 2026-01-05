package devenv

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type Cfg struct {
	OCR2        *OCR2               `toml:"ocr2"`
	Blockchains []*blockchain.Input `toml:"blockchains" validate:"required"`
	FakeServer  *fake.Input         `toml:"fake_server" validate:"required"`
	NodeSets    []*ns.Input         `toml:"nodesets"    validate:"required"`
	JD          *jd.Input           `toml:"jd"`
}

func NewEnvironment() (*Cfg, error) {
	if err := framework.DefaultNetwork(nil); err != nil {
		return nil, err
	}
	in, err := Load[Cfg]()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	_, err = blockchain.NewBlockchainNetwork(in.Blockchains[0])
	if err != nil {
		return nil, fmt.Errorf("failed to create blockchain network 1337: %w", err)
	}
	_, err = fake.NewDockerFakeDataProvider(in.FakeServer)
	if err != nil {
		return nil, fmt.Errorf("failed to create fake data provider: %w", err)
	}
	if err := OCR2ProductConfiguration(in, ConfigureNodesNetwork); err != nil {
		return nil, fmt.Errorf("failed to setup default CLDF orchestration: %w", err)
	}
	_, err = ns.NewSharedDBNodeSet(in.NodeSets[0], nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new shared db node set: %w", err)
	}
	if err := OCR2ProductConfiguration(in, ConfigureProductContractsJobs); err != nil {
		return nil, fmt.Errorf("failed to setup default CLDF orchestration: %w", err)
	}
	return in, Store[Cfg](in)
}

// verifyEnvironment internal function describing how to verify your environment is working.
func verifyEnvironment(in *Cfg, c *ethclient.Client, addr string) error {
	if !in.OCR2.Verify {
		return nil
	}
	Plog.Info().Msg("Verifying environment")
	ocr2i, err := ocr2aggregator.NewOCR2Aggregator(common.HexToAddress(addr), c)
	if err != nil {
		return fmt.Errorf("could not create ocr2 aggregator: %w", err)
	}
	timeout := time.After(in.OCR2.VerificationTimeoutSec * time.Second) //noling:staticcheck
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timed out waiting for non-zero answer from OCR2 aggregator")
		case <-ticker.C:
			result, err := ocr2i.LatestAnswer(&bind.CallOpts{})
			if err != nil {
				return fmt.Errorf("error getting latest answer: %w", err)
			}
			if result.Int64() != 0 {
				Plog.Info().Int64("Answer", result.Int64()).Msg("OCR2 is working, latest answer for price")
				return nil
			}
		}
	}
}
