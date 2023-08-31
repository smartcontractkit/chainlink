package loadfunctions

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"github.com/smartcontractkit/wasp"
	"math/rand"
	"os"
	"strconv"
	"time"
)

/* SingleFunctionCallGun is a gun that constantly requests randomness for one feed  */

type GatewaySecretsSetGun struct {
	Cfg                *PerformanceConfig
	EthereumPrivateKey *ecdsa.PrivateKey
	ThresholdPublicKey *tdh2easy.PublicKey
	DONPublicKey       []byte
}

func NewGatewaySecretsSetGun(cfg *PerformanceConfig, pKey *ecdsa.PrivateKey, tdh2PubKey *tdh2easy.PublicKey, donPubKey []byte) *GatewaySecretsSetGun {
	return &GatewaySecretsSetGun{
		Cfg:                cfg,
		EthereumPrivateKey: pKey,
		ThresholdPublicKey: tdh2PubKey,
		DONPublicKey:       donPubKey,
	}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *GatewaySecretsSetGun) Call(l *wasp.Generator) *wasp.CallResult {
	secrets, err := EncryptS4Secrets(m.EthereumPrivateKey, m.ThresholdPublicKey, m.DONPublicKey, "{\"ltsecret\": \"1\"}")
	if err != nil {
		return &wasp.CallResult{Error: err.Error(), Failed: true}
	}
	if err := UploadS4Secrets(&S4SecretsCfg{
		GatewayURL:            fmt.Sprintf("%s/user", m.Cfg.Common.GatewayURL),
		PrivateKey:            os.Getenv("MUMBAI_KEYS"),
		MessageID:             strconv.Itoa(rand.Intn(100000-1) + 1),
		Method:                "secrets_set",
		DonID:                 m.Cfg.Common.DONID,
		S4SetSlotID:           0,
		S4SetVersion:          uint64(time.Now().UnixNano()),
		S4SetExpirationPeriod: 60 * 60 * 1000,
		S4SetPayload:          secrets,
	}); err != nil {
		return &wasp.CallResult{Error: err.Error(), Failed: true}
	}
	return &wasp.CallResult{}
}
