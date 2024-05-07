package loadfunctions

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	mrand "math/rand"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/seth"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"

	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/types"
)

type FunctionsTest struct {
	SethClient                seth.Client
	LinkToken                 contracts.LinkToken
	Coordinator               contracts.FunctionsCoordinator
	Router                    contracts.FunctionsRouter
	LoadTestClient            contracts.FunctionsLoadTestClient
	EthereumPrivateKey        *ecdsa.PrivateKey
	EthereumPublicKey         *ecdsa.PublicKey
	ThresholdPublicKey        *tdh2easy.PublicKey
	DONPublicKey              []byte
	ThresholdPublicKeyBytes   []byte
	ThresholdEncryptedSecrets string
}

type S4SecretsCfg struct {
	GatewayURL            string
	PrivateKey            string
	RecieverAddr          string
	MessageID             string
	Method                string
	DonID                 string
	S4SetSlotID           uint
	S4SetVersion          uint64
	S4SetExpirationPeriod int64
	S4SetPayload          string
}

func SetupLocalLoadTestEnv(globalConfig ctf_config.GlobalTestConfig, functionsConfig types.FunctionsTestConfig) (*FunctionsTest, error) {
	selectedNetwork := networks.MustGetSelectedNetworkConfig(globalConfig.GetNetworkConfig())[0]
	seth, err := actions_seth.GetChainClient(globalConfig, selectedNetwork)
	if err != nil {
		return nil, err
	}

	cfg := functionsConfig.GetFunctionsConfig()

	lt, err := contracts.DeployLinkTokenContract(log.Logger, seth)
	if err != nil {
		return nil, err
	}
	coord, err := contracts.LoadFunctionsCoordinator(seth, *cfg.Common.Coordinator)
	if err != nil {
		return nil, err
	}
	router, err := contracts.LoadFunctionsRouter(log.Logger, seth, *cfg.Common.Router)
	if err != nil {
		return nil, err
	}
	var loadTestClient contracts.FunctionsLoadTestClient
	if cfg.Common.LoadTestClient != nil && *cfg.Common.LoadTestClient != "" {
		loadTestClient, err = contracts.LoadFunctionsLoadTestClient(seth, *cfg.Common.LoadTestClient)
	} else {
		loadTestClient, err = contracts.DeployFunctionsLoadTestClient(seth, *cfg.Common.Router)
	}
	if err != nil {
		return nil, err
	}
	if cfg.Common.SubscriptionID == nil {
		log.Info().Msg("Creating new subscription")
		subID, err := router.CreateSubscriptionWithConsumer(loadTestClient.Address())
		if err != nil {
			return nil, fmt.Errorf("failed to create a new subscription: %w", err)
		}
		encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint64"}]`, subID)
		if err != nil {
			return nil, fmt.Errorf("failed to encode subscription ID for funding: %w", err)
		}
		_, err = lt.TransferAndCall(router.Address(), big.NewInt(0).Mul(cfg.Common.SubFunds, big.NewInt(1e18)), encodedSubId)
		if err != nil {
			return nil, fmt.Errorf("failed to transferAndCall router, LINK funding: %w", err)
		}
		cfg.Common.SubscriptionID = &subID
	}
	pKey, pubKey, err := parseEthereumPrivateKey(selectedNetwork.PrivateKeys[0])
	if err != nil {
		return nil, fmt.Errorf("failed to load Ethereum private key: %w", err)
	}
	tpk, err := coord.GetThresholdPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get Threshold public key: %w", err)
	}
	log.Info().Hex("ThresholdPublicKeyBytesHex", tpk).Msg("Loaded coordinator keys")
	donPubKey, err := coord.GetDONPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get DON public key: %w", err)
	}
	log.Info().Hex("DONPublicKeyHex", donPubKey).Msg("Loaded DON key")
	tdh2pk, err := ParseTDH2Key(tpk)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tdh2 public key: %w", err)
	}
	var encryptedSecrets string
	if cfg.Common.Secrets != nil && *cfg.Common.Secrets != "" {
		encryptedSecrets, err = EncryptS4Secrets(pKey, tdh2pk, donPubKey, *cfg.Common.Secrets)
		if err != nil {
			return nil, fmt.Errorf("failed to generate tdh2 secrets: %w", err)
		}
		slotID, slotVersion, err := UploadS4Secrets(resty.New(), &S4SecretsCfg{
			GatewayURL:            *cfg.Common.GatewayURL,
			PrivateKey:            selectedNetwork.PrivateKeys[0],
			MessageID:             strconv.Itoa(mrand.Intn(100000-1) + 1),
			Method:                "secrets_set",
			DonID:                 *cfg.Common.DONID,
			S4SetSlotID:           uint(mrand.Intn(5)),
			S4SetVersion:          uint64(time.Now().UnixNano()),
			S4SetExpirationPeriod: 60 * 60 * 1000,
			S4SetPayload:          encryptedSecrets,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to upload secrets to S4: %w", err)
		}
		cfg.Common.SecretsSlotID = &slotID
		cfg.Common.SecretsVersionID = &slotVersion
		log.Info().
			Uint8("SlotID", slotID).
			Uint64("SlotVersion", slotVersion).
			Msg("Set new secret")
	}
	return &FunctionsTest{
		SethClient:                *seth,
		LinkToken:                 lt,
		Coordinator:               coord,
		Router:                    router,
		LoadTestClient:            loadTestClient,
		EthereumPrivateKey:        pKey,
		EthereumPublicKey:         pubKey,
		ThresholdPublicKey:        tdh2pk,
		ThresholdPublicKeyBytes:   tpk,
		ThresholdEncryptedSecrets: encryptedSecrets,
		DONPublicKey:              donPubKey,
	}, nil
}

func parseEthereumPrivateKey(pk string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	pKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert Ethereum key from hex: %w", err)
	}

	publicKey := pKey.Public()
	pubKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, fmt.Errorf("failed to get public key from Ethereum private key: %w", err)
	}
	log.Info().Str("Address", crypto.PubkeyToAddress(*pubKey).Hex()).Msg("Parsed private key for address")
	return pKey, pubKey, nil
}
