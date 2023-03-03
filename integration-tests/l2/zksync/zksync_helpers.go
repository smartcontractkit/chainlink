package zksync

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/gauntlet"
	"math/big"
	"strings"
)

type ZKSyncClient struct {
	GRunner              *gauntlet.GauntletRunner
	LinkAddr             string
	OCRAddr              string
	AccessControllerAddr string
	L2RPC                string
	NKeys                []client.NodeKeysBundle
	Transmitters         []string
	Payees               []string
	Signers              []string
	PeerIds              []string
	OcrConfigPubKeys     []string
}

func Setup(L2RPC string, privateKey string) (*ZKSyncClient, error) {
	g, err := gauntlet.Setup(L2RPC, privateKey)
	if err != nil {
		return nil, err
	}

	return &ZKSyncClient{
		GRunner:  g,
		LinkAddr: "",
		OCRAddr:  "",
		L2RPC:    L2RPC,
	}, nil

}

func (z *ZKSyncClient) DeployLinkToken() error {
	output, err := z.GRunner.ExecuteCommand([]string{"token:deploy"})
	if err != nil {
		return err
	}
	z.LinkAddr = output.Responses[0].Tx.Address
	return nil
}

func (z *ZKSyncClient) DeployAccessController() error {
	output, err := z.GRunner.ExecuteCommand([]string{"access_controller:deploy"})
	if err != nil {
		return err
	}
	z.AccessControllerAddr = output.Responses[0].Tx.Address
	return nil
}

func (z *ZKSyncClient) DeployOCR(
	maxGasPrice string,
	reasonableGasPrice string,
	microLinkPerEth string,
	linkGweiPerObservation string,
	linkGweiPerTransmission string,
	minAnswer string,
	maxAnswer string,
	decimals string,
	description string,
) error {
	output, err := z.GRunner.ExecuteCommand([]string{"ocr:deploy",
		fmt.Sprintf("--maxGasPrice=%s", maxGasPrice),
		fmt.Sprintf("--maxGasPrice=%s", maxGasPrice),
		fmt.Sprintf("--reasonableGasPrice=%s", reasonableGasPrice),
		fmt.Sprintf("--microLinkPerEth=%s", microLinkPerEth),
		fmt.Sprintf("--linkGweiPerObservation=%s", linkGweiPerObservation),
		fmt.Sprintf("--linkGweiPerTransmission=%s", linkGweiPerTransmission),
		fmt.Sprintf("--minAnswer=%s", minAnswer),
		fmt.Sprintf("--maxAnswer=%s", maxAnswer),
		fmt.Sprintf("--billingController=%s", z.AccessControllerAddr),
		fmt.Sprintf("--requesterController=%s", z.AccessControllerAddr),
		fmt.Sprintf("--decimals=%s", decimals),
		fmt.Sprintf("--description=%s", description),
		fmt.Sprintf("--link=%s", z.LinkAddr),
	},
	)
	if err != nil {
		return err
	}
	z.OCRAddr = output.Responses[0].Tx.Address
	return nil
}

func (z *ZKSyncClient) AddAccess(ocrAddress string) error {
	_, err := z.GRunner.ExecuteCommand([]string{"access_controller:add_access", fmt.Sprintf("--address=%s", ocrAddress), z.AccessControllerAddr})
	if err != nil {
		return err
	}
	return nil
}

func (z *ZKSyncClient) SetPayees(ocrAddress string, payees []string, transmitters []string) error {
	_, err := z.GRunner.ExecuteCommand([]string{"ocr:set_payees", ocrAddress, fmt.Sprintf("--transmitters=%s", strings.Join(transmitters, ",")), fmt.Sprintf("--payees=%s", strings.Join(payees, ","))})
	if err != nil {
		return err
	}
	return nil
}
func (z *ZKSyncClient) CreateKeys(chainlinkNodes []*client.Chainlink) error {
	var err error

	z.NKeys, _, err = client.CreateNodeKeysBundle(chainlinkNodes, "evm", "280")
	if err != nil {
		return err
	}
	for _, key := range z.NKeys {
		z.OcrConfigPubKeys = append(z.OcrConfigPubKeys, strings.Replace(key.OCRKey.Data.Attributes.ConfigPublicKey, "ocrcfg_", "", 1))
		z.PeerIds = append(z.PeerIds, key.PeerID)
		z.Transmitters = append(z.Transmitters, key.TXKey.Data.ID)
		z.Signers = append(z.Signers, strings.Replace(key.OCRKey.Data.Attributes.OnChainSigningAddress, "ocrsad_", "", 1))
		z.Payees = append(z.Payees, key.EthAddress)
	}

	return nil
}

func (z *ZKSyncClient) SetConfig(
	ocrAddress string,
	threshold string,
	badEpochTimeout string,
	resendInterval string,
	roundInterval string,
	observationGracePeriod string,
	maxContractValueAge string,
	relativeDeviationThresholdPPB string,
	transmissionStageTimeout string,
	maxRoundCount string,
	transmissionStages string,
) error {
	_, err := z.GRunner.ExecuteCommand([]string{"ocr:set_config",
		ocrAddress,
		fmt.Sprintf("--signers=%s", strings.Join(z.Signers, ",")),
		fmt.Sprintf("--transmitters=%s", strings.Join(z.Transmitters, ",")),
		fmt.Sprintf("--threshold=%s", threshold),
		fmt.Sprintf("--badEpochTimeout=%s", badEpochTimeout),
		fmt.Sprintf("--resendInterval=%s", resendInterval),
		fmt.Sprintf("--roundInterval=%s", roundInterval),
		fmt.Sprintf("--observationGracePeriod=%s", observationGracePeriod),
		fmt.Sprintf("--maxContractValueAge=%s", maxContractValueAge),
		fmt.Sprintf("--relativeDeviationThresholdPPB=%s", relativeDeviationThresholdPPB),
		fmt.Sprintf("--transmissionStageTimeout=%s", transmissionStageTimeout),
		fmt.Sprintf("--maxRoundCount=%s", maxRoundCount),
		fmt.Sprintf("--transmissionStages=%s", transmissionStages),
		fmt.Sprintf("--ocrConfigPublicKeys=%s", strings.Join(z.OcrConfigPubKeys, ",")),
		fmt.Sprintf("--operatorsPeerIds=%v", strings.Join(z.PeerIds, ",")),
		"--secret=Test123"})
	if err != nil {
		return err
	}
	return nil
}

func (z *ZKSyncClient) FundNodes(chainlinkClient blockchain.EVMClient) error {
	for _, key := range z.NKeys {
		log.Info().Str("ZKSync=", fmt.Sprintf("Funding %s", key.TXKey.Data.ID)).Msg("Executing ZKSync command")
		amount := big.NewFloat(100000000000000000)
		err := chainlinkClient.Fund(key.TXKey.Data.ID, amount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *ZKSyncClient) DeployContracts(
	chainlinkClient blockchain.EVMClient,
	maxGasPrice string,
	reasonableGasPrice string,
	microLinkPerEth string,
	linkGweiPerObservation string,
	linkGweiPerTransmission string,
	minAnswer string,
	maxAnswer string,
	decimals string,
	description string,
	threshold string,
	badEpochTimeout string,
	resendInterval string,
	roundInterval string,
	observationGracePeriod string,
	maxContractValueAge string,
	relativeDeviationThresholdPPB string,
	transmissionStageTimeout string,
	maxRoundCount string,
	transmissionStages string,
) error {
	err := z.DeployLinkToken()
	if err != nil {
		return err
	}

	err = z.DeployAccessController()
	if err != nil {
		return err
	}

	err = z.DeployOCR(maxGasPrice, reasonableGasPrice, microLinkPerEth, linkGweiPerObservation, linkGweiPerTransmission, minAnswer, maxAnswer, decimals, description)
	if err != nil {
		return err
	}

	err = z.AddAccess(z.OCRAddr)
	if err != nil {
		return err
	}

	err = z.SetPayees(z.OCRAddr, z.Payees, z.Transmitters)
	if err != nil {
		return err
	}

	err = z.SetConfig(z.OCRAddr, threshold, badEpochTimeout, resendInterval, roundInterval, observationGracePeriod, maxContractValueAge, relativeDeviationThresholdPPB, transmissionStageTimeout, maxRoundCount, transmissionStages)
	if err != nil {
		return err
	}

	err = z.FundNodes(chainlinkClient)
	if err != nil {
		return err
	}

	return nil
}
