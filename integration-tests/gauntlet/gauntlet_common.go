package gauntlet

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	gogauntlet "github.com/smartcontractkit/chainlink/integration-tests/gauntlet/go-gauntlet"
	"os"
	"strings"
)

type GauntletRunner struct {
	G      *gogauntlet.Gauntlet
	config *gogauntlet.DefaultConfig
	ctx    context.Context
}

// Setup TO-DO - set logic for remote load of binary
// Setup needs GAUNTLET_LOCAL_BINARY defined to load gauntlet binary from local
func Setup(nodeUrl string, privateKey string) (*GauntletRunner, error) {
	var g *gogauntlet.Gauntlet
	var err error
	cfg := gogauntlet.NewDefaultConfig(nodeUrl, privateKey, os.Stdout)
	gauntletBinary, pathSet := os.LookupEnv("GAUNTLET_LOCAL_BINARY")
	if pathSet {
		g, err = gogauntlet.NewGauntletFromLocal(gauntletBinary, cfg)
	} else {
		return nil, errors.New("Please set GAUNTLET_LOCAL_BINARY")
	}

	if err != nil {
		return nil, err
	}
	return &GauntletRunner{
		G:      g,
		config: cfg,
		ctx:    context.Background(),
	}, nil
}

func (g *GauntletRunner) ExecuteCommand(args []string) (gogauntlet.Report, error) {
	log.Info().Str("Gauntlet", strings.Join(args, " ")).Msg("Executing Gauntlet command")
	output, err := g.G.ExecWriteCommand(g.ctx, args)
	if err != nil {
		return gogauntlet.Report{}, err
	}
	return output, nil
}
