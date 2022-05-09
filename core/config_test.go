package main_test

import (
	_ "embed"
	"strings"
	"testing"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	solanadb "github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	terradb "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type Config struct {
	Global configGlobal // separate field so that the type can be imported w/o other chain types

	EVM map[string]evmChain

	Solana map[string]solanaChain

	Terra map[string]terraChain
}

type (
	configGlobal struct {
		FeatureFeedsManager       bool
		FeatureOffchainReporting2 bool
		LogLevel                  zapcore.Level
	}

	evmNode struct {
		HTTPURL string
		WSURL   string
	}
	evmChain struct {
		evmtypes.ChainCfg
		Nodes map[string]evmNode
	}

	solanaNode struct {
		URL string
	}
	solanaChain struct {
		solanadb.ChainCfg
		Nodes map[string]solanaNode
	}

	terraNode struct {
		TendermintURL string
	}
	terraChain struct {
		terradb.ChainCfg
		Nodes map[string]terraNode
	}
)

//go:embed config-example.toml
var s string

func TestConfig(t *testing.T) {
	var got Config
	d := toml.NewDecoder(strings.NewReader(s)).Strict(true)
	err := d.Decode(&got)
	if err != nil {
		t.Fatal(err)
	}

	hour := models.MustMakeDuration(time.Hour)
	minute := models.MustMakeDuration(time.Minute)
	exp := Config{
		Global: configGlobal{
			FeatureFeedsManager:       true,
			FeatureOffchainReporting2: true,
			LogLevel:                  zapcore.WarnLevel,
		},
		EVM: map[string]evmChain{
			"1": {
				ChainCfg: evmtypes.ChainCfg{
					OCRObservationTimeout: &hour,
				},
				Nodes: map[string]evmNode{
					"primary-foo": {
						WSURL:   "wss://example.com/ws",
						HTTPURL: "https://example.com",
					},
					"sendonly-mirror": {
						HTTPURL: "https://broadcast-mirror.internet",
					},
				},
			},
		},
		Terra: map[string]terraChain{
			"Columbus-5": {
				ChainCfg: terradb.ChainCfg{
					BlocksUntilTxTimeout: null.IntFrom(10),
					FCDURL:               null.StringFrom("http://fcd.url.com"),
				},
				Nodes: map[string]terraNode{
					"primary": {
						TendermintURL: "http://terra.url",
					},
				},
			},
		},
		Solana: map[string]solanaChain{
			"mainnet": {
				ChainCfg: solanadb.ChainCfg{
					BalancePollPeriod: &minute,
					SkipPreflight:     null.BoolFrom(true),
				},
				Nodes: map[string]solanaNode{
					"node-name": {
						URL: "http://foo.bar",
					},
				},
			},
		},
	}
	require.Equal(t, exp, got)
}
