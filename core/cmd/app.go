package cmd

import (
	"cmp"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/build"
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
func NewApp(s *Shell) *cli.App {
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
		cli.StringSliceFlag{
			Name:   "secrets, s",
			Usage:  "TOML configuration file for secrets. Must be set if and only if config is set. Multiple files can be used (-s secretsA.toml -s secretsB.toml), and they are applied in order. No overrides are allowed.",
			Hidden: true,
		},
	}
	app.Before = func(c *cli.Context) error {
		s.configFiles = c.StringSlice("config")
		s.configFilesIsSet = c.IsSet("config")
		s.secretsFiles = c.StringSlice("secrets")
		s.secretsFileIsSet = c.IsSet("secrets")

		// Default to using a stdout logger only.
		// This is overidden for server commands which may start a rotating
		// logger instead.
		lggr, closeFn := logger.NewLogger()

		cfg, err := opts.New()
		if err != nil {
			return err
		}

		s.Logger = lggr
		s.Registerer = prometheus.DefaultRegisterer // use the global DefaultRegisterer, should be safe since we only ever run one instance of the app per shell
		s.CloseLogger = closeFn
		s.Config = cfg

		if c.Bool("json") {
			s.Renderer = RendererJSON{Writer: os.Stdout}
		}

		cookieJar, err := NewUserCache("cookies", func() logger.Logger { return s.Logger })
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
		cookieAuth := NewSessionCookieAuthenticator(clientOpts, DiskCookieStore{Config: cookieJar}, s.Logger)
		sessionRequestBuilder := NewFileSessionRequestBuilder(s.Logger)

		credentialsFile := c.String("admin-credentials-file")
		sr, err := sessionRequestBuilder.Build(credentialsFile)
		if err != nil && !errors.Is(errors.Cause(err), ErrNoCredentialFile) && !os.IsNotExist(err) {
			return errors.Wrapf(err, "failed to load API credentials from file %s", credentialsFile)
		}

		s.HTTP = NewAuthenticatedHTTPClient(s.Logger, clientOpts, cookieAuth, sr)
		s.CookieAuthenticator = cookieAuth
		s.FileSessionRequestBuilder = sessionRequestBuilder

		// Allow for initServerConfig to be called if the flag is provided.
		if c.Bool("applyInitServerConfig") {
			cfg, err = initServerConfig(&opts, s.configFiles, s.secretsFiles)
			if err != nil {
				return err
			}
			s.Config = cfg
		}

		return nil
	}
	app.After = func(c *cli.Context) error {
		if s.CloseLogger != nil {
			return s.CloseLogger()
		}
		return nil
	}
	app.Commands = removeHidden([]cli.Command{
		{
			Name:        "admin",
			Usage:       "Commands for remotely taking admin related actions",
			Subcommands: initAdminSubCmds(s),
		},
		{
			Name:        "attempts",
			Aliases:     []string{"txas"},
			Usage:       "Commands for managing Ethereum Transaction Attempts",
			Subcommands: initAttemptsSubCmds(s),
		},
		{
			Name:        "blocks",
			Aliases:     []string{},
			Usage:       "Commands for managing blocks",
			Subcommands: initBlocksSubCmds(s),
		},
		{
			Name:        "bridges",
			Usage:       "Commands for Bridges communicating with External Adapters",
			Subcommands: initBrideSubCmds(s),
		},
		{
			Name:        "config",
			Usage:       "Commands for the node's configuration",
			Subcommands: initRemoteConfigSubCmds(s),
		},
		{
			Name:   "health",
			Usage:  "Prints a health report",
			Action: s.Health,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "failing, f",
					Usage: "filter for failing services",
				},
				cli.BoolFlag{
					Name:  "json, j",
					Usage: "json output",
				},
			},
		},
		{
			Name:        "jobs",
			Usage:       "Commands for managing Jobs",
			Subcommands: initJobsSubCmds(s),
		},
		{
			Name:  "keys",
			Usage: "Commands for managing various types of keys used by the Chainlink node",
			Subcommands: []cli.Command{
				// TODO unify init vs keysCommand
				// out of scope for initial refactor because it breaks usage messages.
				initEthKeysSubCmd(s),
				initP2PKeysSubCmd(s),
				initCSAKeysSubCmd(s),
				initOCRKeysSubCmd(s),
				initOCR2KeysSubCmd(s),

				keysCommand("Cosmos", NewCosmosKeysClient(s)),
				keysCommand("Solana", NewSolanaKeysClient(s)),
				keysCommand("StarkNet", NewStarkNetKeysClient(s)),
				keysCommand("Aptos", NewAptosKeysClient(s)),
				keysCommand("Tron", NewTronKeysClient(s)),

				initVRFKeysSubCmd(s),
			},
		},
		{
			Name:        "node",
			Aliases:     []string{"local"},
			Usage:       "Commands for admin actions that must be run locally",
			Description: "Commands can only be run from on the same machine as the Chainlink node.",
			Subcommands: initLocalSubCmds(s, build.IsProd()),
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "config, c",
					Usage: "TOML configuration file(s) via flag, or raw TOML via env var. If used, legacy env vars must not be set. Multiple files can be used (-c configA.toml -c configB.toml), and they are applied in order with duplicated fields overriding any earlier values. If the 'CL_CONFIG' env var is specified, it is always processed last with the effect of being the final override. [$CL_CONFIG]",
				},
				cli.StringSliceFlag{
					Name:  "secrets, s",
					Usage: "TOML configuration file for secrets. Must be set if and only if config is set. Multiple files can be used (-s secretsA.toml -s secretsB.toml), and fields from the files will be merged. No overrides are allowed.",
				},
			},
			Before: func(c *cli.Context) error {
				errNoDuplicateFlags := fmt.Errorf("multiple commands with --config or --secrets flags. only one command may specify these flags. when secrets are used, they must be specific together in the same command")
				if c.IsSet("config") {
					if s.configFilesIsSet || s.secretsFileIsSet {
						return errNoDuplicateFlags
					}
					s.configFiles = c.StringSlice("config")
				}

				if c.IsSet("secrets") {
					if s.configFilesIsSet || s.secretsFileIsSet {
						return errNoDuplicateFlags
					}
					s.secretsFiles = c.StringSlice("secrets")
				}

				// flags here, or ENV VAR only
				cfg, err := initServerConfig(&opts, s.configFiles, s.secretsFiles)
				if err != nil {
					return err
				}
				s.Config = cfg

				logFileMaxSizeMB := s.Config.Log().File().MaxSize() / utils.MB
				if logFileMaxSizeMB > 0 {
					err = utils.EnsureDirAndMaxPerms(s.Config.Log().File().Dir(), os.FileMode(0700))
					if err != nil {
						return err
					}
				}

				// Swap out the logger, replacing the old one.
				err = s.CloseLogger()
				if err != nil {
					return err
				}

				lggrCfg := logger.Config{
					LogLevel:       s.Config.Log().Level(),
					Dir:            s.Config.Log().File().Dir(),
					JsonConsole:    s.Config.Log().JSONConsole(),
					UnixTS:         s.Config.Log().UnixTimestamps(),
					FileMaxSizeMB:  int(logFileMaxSizeMB),
					FileMaxAgeDays: int(s.Config.Log().File().MaxAgeDays()),
					FileMaxBackups: int(s.Config.Log().File().MaxBackups()),
					SentryEnabled:  s.Config.Sentry().DSN() != "",
				}
				l, closeFn := lggrCfg.New()

				s.Logger = l
				s.CloseLogger = closeFn

				return nil
			},
		},
		{
			Name:        "initiators",
			Usage:       "Commands for managing External Initiators",
			Subcommands: initInitiatorsSubCmds(s),
		},
		{
			Name:  "txs",
			Usage: "Commands for handling transactions",
			Subcommands: []cli.Command{
				initEVMTxSubCmd(s),
				initCosmosTxSubCmd(s),
				initSolanaTxSubCmd(s),
			},
		},
		{
			Name:        "chains",
			Usage:       "Commands for handling chain configuration",
			Subcommands: initChainSubCmds(s),
		},
		{
			Name:        "nodes",
			Usage:       "Commands for handling node configuration",
			Subcommands: initNodeSubCmds(s),
		},
		{
			Name:        "forwarders",
			Usage:       "Commands for managing forwarder addresses.",
			Subcommands: initFowardersSubCmds(s),
		},
		{
			Name:  "help-all",
			Usage: "Shows a list of all commands and sub-commands",
			Action: func(c *cli.Context) error {
				printCommands("", c.App.Commands)
				return nil
			},
		},
	}...)
	return app
}

var whitespace = regexp.MustCompile(`\s+`)

// format returns result of replacing all whitespace in s with a single space
func format(s string) string {
	return string(whitespace.ReplaceAll([]byte(s), []byte(" ")))
}

func initServerConfig(opts *chainlink.GeneralConfigOpts, configFiles []string, secretsFiles []string) (chainlink.GeneralConfig, error) {
	err := opts.Setup(configFiles, secretsFiles)
	if err != nil {
		return nil, err
	}
	return opts.New()
}

func printCommands(parent string, cs cli.Commands) {
	slices.SortFunc(cs, func(a, b cli.Command) int {
		return cmp.Compare(a.Name, b.Name)
	})
	for i := range cs {
		c := cs[i]
		name := c.Name
		if parent != "" {
			name = parent + " " + name
		}
		fmt.Printf("%s # %s\n", name, c.Usage)
		printCommands(name, c.Subcommands)
	}
}
