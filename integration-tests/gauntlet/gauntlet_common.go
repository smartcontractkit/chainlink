package gauntlet

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/go-gauntlet/pkg/gauntlet"
	"os"
	"strings"
)

type GauntletRunner struct {
	G      *gauntlet.Gauntlet
	config *gauntlet.DefaultConfig
	ctx    context.Context
}

// Setup TO-DO - set logic for remote load of binary
// Setup needs GAUNTLET_LOCAL_BINARY defined to load gauntlet binary from local
func Setup(nodeUrl string, privateKey string) (*GauntletRunner, error) {
	var g *gauntlet.Gauntlet
	var err error
	cfg := gauntlet.NewDefaultConfig(nodeUrl, privateKey, os.Stdout)
	gauntletBinary, pathSet := os.LookupEnv("GAUNTLET_LOCAL_BINARY")
	if pathSet {
		g, err = gauntlet.NewGauntletFromLocal(gauntletBinary, cfg)
	} else {
		g, err = gauntlet.NewGauntletFromRelease("", "", "", cfg)
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

func (g *GauntletRunner) ExecuteCommand(args []string) (gauntlet.Report, error) {
	log.Info().Str("Gauntlet", strings.Join(args, " ")).Msg("Executing Gauntlet command")
	output, err := g.G.ExecWriteCommand(g.ctx, args)
	if err != nil {
		return gauntlet.Report{}, err
	}
	return output, nil
}
