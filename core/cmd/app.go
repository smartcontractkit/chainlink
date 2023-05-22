package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/build"
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func removeHidden(cmds ...cli.Command) []cli.Command {
	var ret []cli.Command
	for _, cmd := range cmds {
		if cmd.Hidden {
			continue
		}
		ret = append(ret, cmd)
	}
	return ret
}

// NewApp returns the command-line parser/function-router for the given client
func NewApp(client *Client) *cli.App {
	app := cli.NewApp()
	app.Usage = "CLI for Chainlink"
	app.Version = fmt.Sprintf("%v@%v", static.Version, static.Sha)
	// TOML
	var opts chainlink.GeneralConfigOpts

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "json output as opposed to table",
		},
		cli.StringFlag{
			Name:  "admin-credentials-file",
			Usage: fmt.Sprintf("optional, applies only in client mode when making remote API calls. If provided, `FILE` containing admin credentials will be used for logging in, allowing to avoid an additional login step. If `FILE` is missing, it will be ignored. Defaults to %s", filepath.Join("<RootDir>", "apicredentials")),
		},
		cli.StringFlag{
			Name:  "remote-node-url",
			Usage: "optional, applies only in client mode when making remote API calls. If provided, `URL` will be used as the remote Chainlink API endpoint",
			Value: "http://localhost:6688",
		},
		cli.BoolFlag{
			Name:  "insecure-skip-verify",
			Usage: "optional, applies only in client mode when making remote API calls. If turned on, SSL certificate verification will be disabled. This is mostly useful for people who want to use Chainlink with a self-signed TLS certificate",
		},
		cli.StringSliceFlag{
			Name:  "config, c",
			Usage: "TOML configuration file(s) via flag, or raw TOML via env var. If used, legacy env vars must not be set. Multiple files can be used (-c configA.toml -c configB.toml), and they are applied in order with duplicated fields overriding any earlier values. If the 'CL_CONFIG' env var is specified, it is always processed last with the effect of being the final override. [$CL_CONFIG]",
			// Note: we cannot use the EnvVar field since it will combine with the flags.
			Hidden: true,
		},
		cli.StringFlag{
			Name:   "secrets, s",
			Usage:  "TOML configuration file for secrets. Must be set if and only if config is set.",
			Hidden: true,
		},
	}
	app.Before = func(c *cli.Context) error {
		client.configFiles = c.StringSlice("config")
		client.configFilesIsSet = c.IsSet("config")
		client.secretsFile = c.String("secrets")
		client.secretsFileIsSet = c.IsSet("secrets")

		// Default to using a stdout logger only.
		// This is overidden for server commands which may start a rotating
		// logger instead.
		lggr, closeFn := logger.NewLogger()

		cfg, err := opts.New()
		if err != nil {
			return err
		}

		client.Logger = lggr
		client.CloseLogger = closeFn
		client.Config = cfg

		if c.Bool("json") {
			client.Renderer = RendererJSON{Writer: os.Stdout}
		}

		cookieJar, err := NewUserCache("cookies", func() logger.Logger { return client.Logger })
		if err != nil {
			return fmt.Errorf("error initialize chainlink cookie cache: %w", err)
		}

		urlStr := c.String("remote-node-url")
		remoteNodeURL, err := url.Parse(urlStr)
		if err != nil {
			return errors.Wrapf(err, "%s is not a valid URL", urlStr)
		}

		insecureSkipVerify := c.Bool("insecure-skip-verify")
		clientOpts := ClientOpts{RemoteNodeURL: *remoteNodeURL, InsecureSkipVerify: insecureSkipVerify}
		cookieAuth := NewSessionCookieAuthenticator(clientOpts, DiskCookieStore{Config: cookieJar}, client.Logger)
		sessionRequestBuilder := NewFileSessionRequestBuilder(client.Logger)

		credentialsFile := c.String("admin-credentials-file")
		sr, err := sessionRequestBuilder.Build(credentialsFile)
		if err != nil && !errors.Is(errors.Cause(err), ErrNoCredentialFile) && !os.IsNotExist(err) {
			return errors.Wrapf(err, "failed to load API credentials from file %s", credentialsFile)
		}

		client.HTTP = NewAuthenticatedHTTPClient(client.Logger, clientOpts, cookieAuth, sr)
		client.CookieAuthenticator = cookieAuth
		client.FileSessionRequestBuilder = sessionRequestBuilder
		return nil

	}
	app.After = func(c *cli.Context) error {
		if client.CloseLogger != nil {
			return client.CloseLogger()
		}
		return nil
	}
	app.Commands = removeHidden([]cli.Command{
		{
			Name:        "admin",
			Usage:       "Commands for remotely taking admin related actions",
			Subcommands: initAdminSubCmds(client),
		},
		{
			Name:        "attempts",
			Aliases:     []string{"txas"},
			Usage:       "Commands for managing Ethereum Transaction Attempts",
			Subcommands: initAttemptsSubCmds(client),
		},
		{
			Name:        "blocks",
			Aliases:     []string{},
			Usage:       "Commands for managing blocks",
			Subcommands: initBlocksSubCmds(client),
		},
		{
			Name:        "bridges",
			Usage:       "Commands for Bridges communicating with External Adapters",
			Subcommands: initBrideSubCmds(client),
		},
		{
			Name:        "config",
			Usage:       "Commands for the node's configuration",
			Subcommands: initRemoteConfigSubCmds(client),
		},
		{
			Name:        "jobs",
			Usage:       "Commands for managing Jobs",
			Subcommands: initJobsSubCmds(client),
		},
		{
			Name:  "keys",
			Usage: "Commands for managing various types of keys used by the Chainlink node",
			Subcommands: []cli.Command{
				// TODO unify init vs keysCommand
				// out of scope for initial refactor because it breaks usage messages.
				initEthKeysSubCmd(client),
				initP2PKeysSubCmd(client),
				initCSAKeysSubCmd(client),
				initOCRKeysSubCmd(client),
				initOCR2KeysSubCmd(client),

				keysCommand("Cosmos", NewCosmosKeysClient(client)),
				keysCommand("Solana", NewSolanaKeysClient(client)),
				keysCommand("StarkNet", NewStarkNetKeysClient(client)),
				keysCommand("DKGSign", NewDKGSignKeysClient(client)),
				keysCommand("DKGEncrypt", NewDKGEncryptKeysClient(client)),

				initVRFKeysSubCmd(client),
			},
		},
		{
			Name:        "node",
			Aliases:     []string{"local"},
			Usage:       "Commands for admin actions that must be run locally",
			Description: "Commands can only be run from on the same machine as the Chainlink node.",
			Subcommands: initLocalSubCmds(client, build.IsProd()),
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "config, c",
					Usage: "TOML configuration file(s) via flag, or raw TOML via env var. If used, legacy env vars must not be set. Multiple files can be used (-c configA.toml -c configB.toml), and they are applied in order with duplicated fields overriding any earlier values. If the 'CL_CONFIG' env var is specified, it is always processed last with the effect of being the final override. [$CL_CONFIG]",
				},
				cli.StringFlag{
					Name:  "secrets, s",
					Usage: "TOML configuration file for secrets. Must be set if and only if config is set.",
				},
			},
			Before: func(c *cli.Context) error {
				errNoDuplicateFlags := fmt.Errorf("multiple commands with --config or --secrets flags. only one command may specify these flags. when secrets are used, they must be specific together in the same command")
				if c.IsSet("config") {
					if client.configFilesIsSet || client.secretsFileIsSet {
						return errNoDuplicateFlags
					} else {
						client.configFiles = c.StringSlice("config")
					}
				}

				if c.IsSet("secrets") {
					if client.configFilesIsSet || client.secretsFileIsSet {
						return errNoDuplicateFlags
					} else {
						client.secretsFile = c.String("secrets")
					}
				}

				// flags here, or ENV VAR only
				cfg, err := initServerConfig(&opts, client.configFiles, client.secretsFile)
				if err != nil {
					return err
				}
				client.Config = cfg

				logFileMaxSizeMB := client.Config.LogFileMaxSize() / utils.MB
				if logFileMaxSizeMB > 0 {
					err = utils.EnsureDirAndMaxPerms(client.Config.LogFileDir(), os.FileMode(0700))
					if err != nil {
						return err
					}
				}

				// Swap out the logger, replacing the old one.
				err = client.CloseLogger()
				if err != nil {
					return err
				}

				lggrCfg := logger.Config{
					LogLevel:       client.Config.LogLevel(),
					Dir:            client.Config.LogFileDir(),
					JsonConsole:    client.Config.JSONConsole(),
					UnixTS:         client.Config.LogUnixTimestamps(),
					FileMaxSizeMB:  int(logFileMaxSizeMB),
					FileMaxAgeDays: int(client.Config.LogFileMaxAge()),
					FileMaxBackups: int(client.Config.LogFileMaxBackups()),
				}
				l, closeFn := lggrCfg.New()

				client.Logger = l
				client.CloseLogger = closeFn

				return nil
			},
		},
		{
			Name:        "initiators",
			Usage:       "Commands for managing External Initiators",
			Subcommands: initInitiatorsSubCmds(client),
		},
		{
			Name:  "txs",
			Usage: "Commands for handling transactions",
			Subcommands: []cli.Command{
				initEVMTxSubCmd(client),
				initCosmosTxSubCmd(client),
				initSolanaTxSubCmd(client),
			},
		},
		{
			Name:  "chains",
			Usage: "Commands for handling chain configuration",
			Subcommands: cli.Commands{
				chainCommand("EVM", EVMChainClient(client), cli.Int64Flag{Name: "id", Usage: "chain ID"}),
				chainCommand("Cosmos", CosmosChainClient(client), cli.StringFlag{Name: "id", Usage: "chain ID"}),
				chainCommand("Solana", SolanaChainClient(client),
					cli.StringFlag{Name: "id", Usage: "chain ID, options: [mainnet, testnet, devnet, localnet]"}),
				chainCommand("StarkNet", StarkNetChainClient(client), cli.StringFlag{Name: "id", Usage: "chain ID"}),
			},
		},
		{
			Name:  "nodes",
			Usage: "Commands for handling node configuration",
			Subcommands: cli.Commands{
				initEVMNodeSubCmd(client),
				initCosmosNodeSubCmd(client),
				initSolanaNodeSubCmd(client),
				initStarkNetNodeSubCmd(client),
			},
		},
		{
			Name:        "forwarders",
			Usage:       "Commands for managing forwarder addresses.",
			Subcommands: initFowardersSubCmds(client),
		},
	}...)
	return app
}

var whitespace = regexp.MustCompile(`\s+`)

// format returns result of replacing all whitespace in s with a single space
func format(s string) string {
	return string(whitespace.ReplaceAll([]byte(s), []byte(" ")))
}

func initServerConfig(opts *chainlink.GeneralConfigOpts, configFiles []string, secretsFile string) (chainlink.GeneralConfig, error) {
	configs := []string{}
	for _, fileName := range configFiles {
		b, err := os.ReadFile(fileName)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read config file: %s", fileName)
		}
		configs = append(configs, string(b))
	}

	if configTOML := v2.EnvConfig.Get(); configTOML != "" {
		configs = append(configs, configTOML)
	}

	opts.ConfigStrings = configs

	secrets := ""
	if secretsFile != "" {
		b, err := os.ReadFile(secretsFile)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read secrets file: %s", secretsFile)
		}

		secrets = string(b)
	}

	opts.SecretsString = secrets

	return opts.New()
}
