package loadfunctions

import (
	"crypto/ecdsa"
	"math/big"
	mrand "math/rand"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/utils"
)

type FunctionsTest struct {
	EVMClient                 blockchain.EVMClient
	ContractDeployer          contracts.ContractDeployer
	ContractLoader            contracts.ContractLoader
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

func SetupLocalLoadTestEnv(cfg *PerformanceConfig) (*FunctionsTest, error) {
	bc, err := blockchain.NewEVMClientFromNetwork(networks.SelectedNetwork, log.Logger)
	if err != nil {
		return nil, err
	}
	cd, err := contracts.NewContractDeployer(bc, log.Logger)
	if err != nil {
		return nil, err
	}

	cl, err := contracts.NewContractLoader(bc, log.Logger)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	lt, err := cl.LoadLINKToken(cfg.Common.LINKTokenAddr)
	if err != nil {
		return nil, err
	}
	coord, err := cl.LoadFunctionsCoordinator(cfg.Common.Coordinator)
	if err != nil {
		return nil, err
	}
	router, err := cl.LoadFunctionsRouter(cfg.Common.Router)
	if err != nil {
		return nil, err
	}
	var loadTestClient contracts.FunctionsLoadTestClient
	if cfg.Common.LoadTestClient != "" {
		loadTestClient, err = cl.LoadFunctionsLoadTestClient(cfg.Common.LoadTestClient)
	} else {
		loadTestClient, err = cd.DeployFunctionsLoadTestClient(cfg.Common.Router)
	}
	if err != nil {
		return nil, err
	}
	if cfg.Common.SubscriptionID == 0 {
		log.Info().Msg("Creating new subscription")
		subID, err := router.CreateSubscriptionWithConsumer(loadTestClient.Address())
		if err != nil {
			return nil, errors.Wrap(err, "failed to create a new subscription")
		}
		encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint64"}]`, subID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to encode subscription ID for funding")
		}
		_, err = lt.TransferAndCall(router.Address(), big.NewInt(0).Mul(cfg.Common.Funding.SubFunds, big.NewInt(1e18)), encodedSubId)
		if err != nil {
			return nil, errors.Wrap(err, "failed to transferAndCall router, LINK funding")
		}
		cfg.Common.SubscriptionID = subID
	}
	pKey, pubKey, err := parseEthereumPrivateKey(os.Getenv("MUMBAI_KEYS"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to load Ethereum private key")
	}
	tpk, err := coord.GetThresholdPublicKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Threshold public key")
	}
	log.Info().Hex("ThresholdPublicKeyBytesHex", tpk).Msg("Loaded coordinator keys")
	donPubKey, err := coord.GetDONPublicKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get DON public key")
	}
	log.Info().Hex("DONPublicKeyHex", donPubKey).Msg("Loaded DON key")
	tdh2pk, err := ParseTDH2Key(tpk)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal tdh2 public key")
	}
	var encryptedSecrets string
	if cfg.Common.Secrets != "" {
		encryptedSecrets, err = EncryptS4Secrets(pKey, tdh2pk, donPubKey, cfg.Common.Secrets)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate tdh2 secrets")
		}
		slotID, slotVersion, err := UploadS4Secrets(resty.New(), &S4SecretsCfg{
			GatewayURL:            cfg.Common.GatewayURL,
			PrivateKey:            cfg.MumbaiPrivateKey,
			MessageID:             strconv.Itoa(mrand.Intn(100000-1) + 1),
			Method:                "secrets_set",
			DonID:                 cfg.Common.DONID,
			S4SetSlotID:           uint(mrand.Intn(5)),
			S4SetVersion:          uint64(time.Now().UnixNano()),
			S4SetExpirationPeriod: 60 * 60 * 1000,
			S4SetPayload:          encryptedSecrets,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to upload secrets to S4")
		}
		cfg.Common.SecretsSlotID = slotID
		cfg.Common.SecretsVersionID = slotVersion
		log.Info().
			Uint8("SlotID", slotID).
			Uint64("SlotVersion", slotVersion).
			Msg("Set new secret")
	}
	return &FunctionsTest{
		EVMClient:                 bc,
		ContractDeployer:          cd,
		ContractLoader:            cl,
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
		return nil, nil, errors.Wrap(err, "failed to convert Ethereum key from hex")
	}

	publicKey := pKey.Public()
	pubKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, errors.Wrap(err, "failed to get public key from Ethereum private key")
	}
	log.Info().Str("Address", crypto.PubkeyToAddress(*pubKey).Hex()).Msg("Parsed private key for address")
	return pKey, pubKey, nil
}
