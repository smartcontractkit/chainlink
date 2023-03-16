package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/config"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/static"
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

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func isDevMode() bool {
	var clDev string
	v1, v2 := os.Getenv("CHAINLINK_DEV"), os.Getenv("CL_DEV")
	if v1 != "" && v2 != "" {
		if v1 != v2 {
			panic("you may only set one of CHAINLINK_DEV and CL_DEV environment variables, not both")
		}
	} else if v1 == "" {
		clDev = v2
	} else if v2 == "" {
		clDev = v1
	}
	return strings.ToLower(clDev) == "true"
}

// NewApp returns the command-line parser/function-router for the given client
func NewApp(client *Client) *cli.App {
	devMode := isDevMode()
	app := cli.NewApp()
	app.Usage = "CLI for Chainlink"
	app.Version = fmt.Sprintf("%v@%v", static.Version, static.Sha)
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
		},
		cli.StringFlag{
			Name:  "secrets, s",
			Usage: "TOML configuration file for secrets. Must be set if and only if config is set.",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.IsSet("config") || v2.EnvConfig.Get() != "" {
			// TOML
			var opts chainlink.GeneralConfigOpts

			fileNames := c.StringSlice("config")
			if err := loadOpts(&opts, fileNames...); err != nil {
				return err
			}

			secretsTOML := ""
			if c.IsSet("secrets") {
				secretsFileName := c.String("secrets")
				b, err := os.ReadFile(secretsFileName)
				if err != nil {
					return errors.Wrapf(err, "failed to read secrets file: %s", secretsFileName)
				}
				secretsTOML = string(b)
			}
			if err := opts.ParseSecrets(secretsTOML); err != nil {
				return err
			}

			if cfg, lggr, closeLggr, err := opts.NewAndLogger(); err != nil {
				return err
			} else {
				client.Config = cfg
				client.Logger = lggr
				client.CloseLogger = closeLggr
			}
		} else {
			// Legacy ENV
			if c.IsSet("secrets") {
				panic("secrets file must not be used without a core config file")
			}
			client.Logger, client.CloseLogger = logger.NewLogger()
			client.Config = config.NewGeneralConfig(client.Logger)
		}
		logDeprecatedClientEnvWarnings(client.Logger)
		if c.Bool("json") {
			client.Renderer = RendererJSON{Writer: os.Stdout}
		}
		urlStr := c.String("remote-node-url")
		if envUrlStr := os.Getenv("CLIENT_NODE_URL"); envUrlStr != "" {
			urlStr = envUrlStr
		}
		remoteNodeURL, err := url.Parse(urlStr)
		if err != nil {
			return errors.Wrapf(err, "%s is not a valid URL", urlStr)
		}
		insecureSkipVerify := c.Bool("insecure-skip-verify")
		if envInsecureSkipVerify := os.Getenv("INSECURE_SKIP_VERIFY"); envInsecureSkipVerify == "true" {
			insecureSkipVerify = true
		}
		clientOpts := ClientOpts{RemoteNodeURL: *remoteNodeURL, InsecureSkipVerify: insecureSkipVerify}
		cookieAuth := NewSessionCookieAuthenticator(clientOpts, DiskCookieStore{Config: client.Config}, client.Logger)
		sessionRequestBuilder := NewFileSessionRequestBuilder(client.Logger)

		credentialsFile := c.String("admin-credentials-file")
		if envCredentialsFile := os.Getenv("ADMIN_CREDENTIALS_FILE"); envCredentialsFile != "" {
			credentialsFile = envCredentialsFile
		}
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
			Subcommands: initLocalSubCmds(client, devMode),
		},
		{
			Name:        "initiators",
			Usage:       "Commands for managing External Initiators",
			Hidden:      !devMode,
			Subcommands: initInitiatorsSubCmds(client, devMode),
		},
		{
			Name:  "txs",
			Usage: "Commands for handling transactions",
			Subcommands: []cli.Command{
				initEVMTxSubCmd(client),
				initSolanaTxSubCmd(client),
			},
		},
		{
			Name:  "chains",
			Usage: "Commands for handling chain configuration",
			Subcommands: cli.Commands{
				chainCommand("EVM", EVMChainClient(client), cli.Int64Flag{Name: "id", Usage: "chain ID"}),
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

func logDeprecatedClientEnvWarnings(lggr logger.Logger) {
	if s := os.Getenv("INSECURE_SKIP_VERIFY"); s != "" {
		lggr.Error("INSECURE_SKIP_VERIFY env var has been deprecated and will be removed in a future release. Use flag instead: --insecure-skip-verify")
	}
	if s := os.Getenv("CLIENT_NODE_URL"); s != "" {
		lggr.Errorf("CLIENT_NODE_URL env var has been deprecated and will be removed in a future release. Use flag instead: --remote-node-url=%s", s)
	}
	if s := os.Getenv("ADMIN_CREDENTIALS_FILE"); s != "" {
		lggr.Errorf("ADMIN_CREDENTIALS_FILE env var has been deprecated and will be removed in a future release. Use flag instead: --admin-credentials-file=%s", s)
	}
}

// loadOpts applies file configs and then overlays env config
func loadOpts(opts *chainlink.GeneralConfigOpts, fileNames ...string) error {
	for _, fileName := range fileNames {
		b, err := os.ReadFile(fileName)
		if err != nil {
			return errors.Wrapf(err, "failed to read config file: %s", fileName)
		}
		if err := opts.ParseConfig(string(b)); err != nil {
			return errors.Wrapf(err, "failed to parse file: %s", fileName)
		}
	}
	if configTOML := v2.EnvConfig.Get(); configTOML != "" {
		if err := opts.ParseConfig(configTOML); err != nil {
			return errors.Wrapf(err, "failed to parse env var %q", v2.EnvConfig)
		}
	}
	return nil
}
