//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/pkg/errors"
	"os"

	"github.com/google/wire"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/sessions"
)

func newProductionHttpClient(
	cfg config.GeneralConfig,
	lggr logger.Logger,
	cookieAuth cmd.CookieAuthenticator,
	sessionRequestBuilder cmd.FileSessionRequestBuilder,
) cmd.HTTPClient {
	sr := sessions.SessionRequest{}
	if credentialsFile := cfg.AdminCredentialsFile(); credentialsFile != "" {
		var err error
		sr, err = sessionRequestBuilder.Build(credentialsFile)
		if err != nil && errors.Cause(err) != cmd.ErrNoCredentialFile && !os.IsNotExist(err) {
			lggr.Fatalw("Error loading API credentials", "error", err, "credentialsFile", credentialsFile)
		}
	}
	return cmd.NewAuthenticatedHTTPClient(cfg, cookieAuth, sr)
}

func InitializeProductionClient() *cmd.Client {
	wire.Build(
		config.NewGeneralConfig,
		logger.NewLogger,
		wire.Struct(new(cmd.TerminalKeyStoreAuthenticator), "*"),
		wire.Struct(new(cmd.DiskCookieStore), "*"),
		wire.Struct(new(cmd.ChainlinkAppFactory)),
		wire.Struct(new(cmd.ChainlinkRunner)),
		wire.Bind(new(logger.Config), new(config.GeneralConfig)),
		wire.Bind(new(cmd.DiskCookieConfig), new(config.GeneralConfig)),
		wire.Bind(new(cmd.SessionCookieAuthenticatorConfig), new(config.GeneralConfig)),
		wire.Bind(new(cmd.CookieStore), new(cmd.DiskCookieStore)),
		wire.Bind(new(cmd.AppFactory), new(cmd.ChainlinkAppFactory)),
		wire.Bind(new(cmd.Runner), new(cmd.ChainlinkRunner)),
		newProductionHttpClient,
		cmd.ProviderSet,
	)
	return nil
}
