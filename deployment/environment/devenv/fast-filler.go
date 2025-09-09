package devenv

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	tcLog "github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
)

var (
	E2eFastFillerImage     = "E2E_FAST_FILLER_IMAGE"
	E2eFastFillerVersion   = "E2E_FAST_FILLER_VERSION"
	DefaultFastFillerImage = "ccip-fast-filler:latest"
)

type ListenerConfig struct {
	RPCURL               string `json:"rpcUrl"`
	TokenPoolAddress     string `json:"tokenPoolAddress"`
	ChainSelector        string `json:"chainSelector"`
	DestinationTokenPool string `json:"destinationTokenPool"`
}

type SignerProvider struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	PrivateKey  string `json:"privateKey,omitempty"`
	EnvVariable string `json:"envVariable,omitempty"`
}

type FillerConfig struct {
	RPCURL           string `json:"rpcUrl"`
	TokenPoolAddress string `json:"tokenPoolAddress"`
	SignerProvider   string `json:"signerProvider"`
	ChainSelector    string `json:"chainSelector"`
	SourceTokenPool  string `json:"sourceTokenPool"`
}

type CCIPFastFillerConfig struct {
	SignerProviders []SignerProvider `json:"signerProviders"`
	Listeners       []ListenerConfig `json:"listeners"`
	Fillers         []FillerConfig   `json:"fillers"`
}

type CCIPFastFiller struct {
	image     string
	config    CCIPFastFillerConfig
	container testcontainers.Container
	logger    zerolog.Logger
	networks  []string
}

func NewCCIPFastFiller(config CCIPFastFillerConfig, l zerolog.Logger, networks []string, image string) *CCIPFastFiller {
	return &CCIPFastFiller{
		image:    image,
		config:   config,
		logger:   l,
		networks: networks,
	}
}

func (f *CCIPFastFiller) Start(ctx context.Context, t *testing.T) error {
	if f.container != nil {
		return nil
	}

	configContent, err := json.Marshal(f.config)
	if err != nil {
		return err
	}
	l := tcLog.Default()
	if t != nil {
		l = logging.CustomT{
			T: t,
			L: f.logger,
		}
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Networks:   f.networks,
			Image:      f.image,
			WaitingFor: wait.ForLog("Relayer started"),
			Files: []testcontainers.ContainerFile{
				{
					Reader:            bytes.NewReader(configContent),
					ContainerFilePath: "/app/config.json",
					FileMode:          0644,
				},
			},
		},
		Logger:  l,
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return err
	}

	f.container = container
	return nil
}

func (f *CCIPFastFiller) Stop(ctx context.Context) error {
	if f.container == nil {
		return nil
	}

	if err := f.container.Terminate(ctx); err != nil {
		return err
	}
	f.container = nil
	return nil
}
