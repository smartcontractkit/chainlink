package loadfunctions

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"io"
	"math/big"
	"net/http"
	"os"
	"time"
)

const (
	SecretJSON           = "{\"num\": \"1\"}"
	JSPayloadWithSecrets = "return Functions.encodeUint256(BigInt(secrets.num))"
	DefaultJSPayload     = "const response = await Functions.makeHttpRequest({ url: 'http://dummyjson.com/products/1' }); return Functions.encodeUint256(response.data.id)"
)

type FunctionsTest struct {
	LinkToken                 contracts.LinkToken
	Coordinator               contracts.FunctionsCoordinator
	Router                    contracts.FunctionsRouter
	LoadTestClient            contracts.FunctionsLoadTestClient
	EthereumPrivateKey        *ecdsa.PrivateKey
	EthereumPublicKey         *ecdsa.PublicKey
	ThresholdPublicKey        *tdh2easy.PublicKey
	ThresholdPublicKeyBytes   []byte
	ThresholdEncryptedSecrets string
	DONPublicKey              []byte
}

type S4SecretsCfg struct {
	GatewayURL            string
	PrivateKey            string
	MessageID             string
	Method                string
	DonID                 string
	S4SetSlotID           uint
	S4SetVersion          uint64
	S4SetExpirationPeriod int64
	S4SetPayload          string
}

func UploadS4Secrets(s4Cfg *S4SecretsCfg) error {
	key, err := crypto.HexToECDSA(s4Cfg.PrivateKey)
	if err != nil {
		return err
	}
	address := crypto.PubkeyToAddress(key.PublicKey)
	var payloadJSON []byte
	if s4Cfg.Method == functions.MethodSecretsSet {
		envelope := s4.Envelope{
			Address:    address.Bytes(),
			SlotID:     s4Cfg.S4SetSlotID,
			Version:    s4Cfg.S4SetVersion,
			Payload:    []byte(s4Cfg.S4SetPayload),
			Expiration: time.Now().UnixMilli() + s4Cfg.S4SetExpirationPeriod,
		}
		signature, err := envelope.Sign(key)
		if err != nil {
			return err
		}

		s4SetPayload := functions.SecretsSetRequest{
			SlotID:     envelope.SlotID,
			Version:    envelope.Version,
			Expiration: envelope.Expiration,
			Payload:    []byte(s4Cfg.S4SetPayload),
			Signature:  signature,
		}

		payloadJSON, err = json.Marshal(s4SetPayload)
		if err != nil {
			return err
		}
	}

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: s4Cfg.MessageID,
			Method:    s4Cfg.Method,
			DonId:     s4Cfg.DonID,
			Payload:   json.RawMessage(payloadJSON),
		},
	}

	err = msg.Sign(key)
	if err != nil {
		return err
	}
	codec := api.JsonRPCCodec{}
	rawMsg, err := codec.EncodeRequest(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(context.Background(), "POST", s4Cfg.GatewayURL, bytes.NewBuffer(rawMsg))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Info().Str("Response", string(body)).Msg("S4 Gateway response")
	return nil
}

func SetupLocalLoadTestEnv(cfg *PerformanceConfig) (*test_env.CLClusterTestEnv, *FunctionsTest, error) {
	env, err := test_env.NewCLTestEnvBuilder().
		Build()
	if err != nil {
		return env, nil, err
	}
	lt, err := env.ContractLoader.LoadLINKToken(cfg.Common.LINKTokenAddr)
	if err != nil {
		return env, nil, err
	}
	coord, err := env.ContractLoader.LoadFunctionsCoordinator(cfg.Common.Coordinator)
	if err != nil {
		return env, nil, err
	}
	router, err := env.ContractLoader.LoadFunctionsRouter(cfg.Common.Router)
	if err != nil {
		return env, nil, err
	}
	loadTestClient, err := env.ContractLoader.LoadFunctionsLoadTestClient(cfg.Common.LoadTestClient)
	if err != nil {
		return env, nil, err
	}
	if cfg.Common.SubscriptionID == 0 {
		log.Info().Msg("Creating new subscription")
		subID, err := router.CreateSubscriptionWithConsumer(loadTestClient.Address())
		if err != nil {
			return env, nil, errors.Wrap(err, "failed to create a new subscription")
		}
		encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint64"}]`, subID)
		if err != nil {
			return env, nil, errors.Wrap(err, "failed to encode subscription ID for funding")
		}
		_, err = lt.TransferAndCall(router.Address(), big.NewInt(0).Mul(cfg.Common.Funding.SubFunds, big.NewInt(1e18)), encodedSubId)
		if err != nil {
			return env, nil, errors.Wrap(err, "failed to transferAndCall router, LINK funding")
		}
		cfg.Common.SubscriptionID = subID
	}
	pKey, pubKey, err := parseEthereumPrivateKey(os.Getenv("MUMBAI_KEYS"))
	if err != nil {
		return env, nil, errors.Wrap(err, "failed to load Ethereum private key")
	}
	tpk, err := coord.GetThresholdPublicKey()
	if err != nil {
		return env, nil, errors.Wrap(err, "failed to get threshold public key")
	}
	log.Info().Hex("ThresholdPublicKeyBytesHex", tpk).Msg("Loaded coordinator keys")
	tdh2pk, err := parseTDH2Key(tpk)
	if err != nil {
		return env, nil, errors.Wrap(err, "failed to unmarshal tdh2 public key")
	}
	secrets, err := generateSecrets(tdh2pk, "{\"num\": \"1\"}")
	if err != nil {
		return env, nil, errors.Wrap(err, "failed to generate tdh2 secrets")
	}
	if err := UploadS4Secrets(&S4SecretsCfg{
		GatewayURL:            "https://gateway-staging1.main.stage.cldev.sh/user",
		PrivateKey:            os.Getenv("MUMBAI_KEYS"),
		MessageID:             "1",
		Method:                "secrets_set",
		DonID:                 "functions_staging_mumbai",
		S4SetSlotID:           4,
		S4SetVersion:          0,
		S4SetExpirationPeriod: 60 * 60 * 1000,
		S4SetPayload:          secrets,
	}); err != nil {
		return env, nil, errors.Wrap(err, "failed to upload secrets to S4")
	}
	donPk, err := coord.GetDONPublicKey()
	if err != nil {
		return env, nil, errors.Wrap(err, "failed to get DON public key")
	}
	log.Info().Hex("DONPublicKeyHex", donPk).Msg("Loaded coordinator keys")
	return env, &FunctionsTest{
		LinkToken:                 lt,
		Coordinator:               coord,
		Router:                    router,
		LoadTestClient:            loadTestClient,
		EthereumPrivateKey:        pKey,
		EthereumPublicKey:         pubKey,
		ThresholdPublicKey:        tdh2pk,
		ThresholdPublicKeyBytes:   tpk,
		ThresholdEncryptedSecrets: secrets,
		DONPublicKey:              donPk,
	}, nil
}

func parseTDH2Key(data []byte) (*tdh2easy.PublicKey, error) {
	pk := &tdh2easy.PublicKey{}
	if err := pk.Unmarshal(data); err != nil {
		return nil, err
	}
	return pk, nil
}

func generateSecrets(pk *tdh2easy.PublicKey, msg string) (string, error) {
	ctxt, err := tdh2easy.Encrypt(pk, []byte(msg))
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt tdh2 encoded secrets")
	}
	mctxt, err := ctxt.Marshal()
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal tdh2 encoded secrets")
	}
	log.Info().RawJSON("SecretsJSON", mctxt).Msg("Encrypted secrets")
	secretsHex := hex.EncodeToString(mctxt)
	log.Info().Str("SecretsHex", secretsHex).Msg("Encoded tdh2 secrets")
	return secretsHex, nil
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
