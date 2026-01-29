package products

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	EnvVarTestConfigs = "CTF_CONFIGS"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "product_config"}).Logger()

func Load[T any]() (*T, error) {
	var config T
	paths := strings.Split(os.Getenv(EnvVarTestConfigs), ",")
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read product config file path %s: %w", path, err)
		}
		L.Trace().Str("ProductConfig", string(data)).Send()

		decoder := toml.NewDecoder(strings.NewReader(string(data)))

		if err := decoder.Decode(&config); err != nil {
			return nil, fmt.Errorf("failed to decode TOML config, strict mode: %w", err)
		}
	}
	return &config, nil
}

// Store writes config to a file, adds -cache.toml suffix if it's an initial configuration.
func Store[T any](path string, cfg *T) error {
	baseConfigPath, err := BaseConfigPath(EnvVarTestConfigs)
	if err != nil {
		return err
	}
	newCacheName := strings.ReplaceAll(baseConfigPath, ".toml", "")
	var outCacheName string
	if strings.Contains(newCacheName, "cache") {
		L.Info().Str("Cache", baseConfigPath).Msg("Cache file already exists, overriding")
		outCacheName = baseConfigPath
	} else {
		outCacheName = strings.ReplaceAll(baseConfigPath, ".toml", "") + "-out.toml"
	}
	L.Info().Str("OutputFile", outCacheName).Msg("Storing configuration output")
	d, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(path, outCacheName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(d); err != nil {
		return err
	}
	return nil
}

// LoadOutput loads config output file from path.
func LoadOutput[T any](path string) (*T, error) {
	_ = os.Setenv(EnvVarTestConfigs, path)
	return Load[T]()
}

// BaseConfigPath returns base config path, ex. env.toml,overrides.toml -> env.toml.
func BaseConfigPath(envVar string) (string, error) {
	configs := os.Getenv(envVar)
	if configs == "" {
		return "", fmt.Errorf("no %s env var is provided, you should provide at least one test config in TOML", envVar)
	}
	L.Debug().Str("Configs", configs).Msg("Getting base config path")
	return strings.Split(configs, ",")[0], nil
}
