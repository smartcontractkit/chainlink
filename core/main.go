package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/recovery"
	"github.com/smartcontractkit/chainlink/core/sessions"
)

var (
	simpleCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "simple_counter",
		Help: "Some help message",
	})

	simpleCounterVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "simple_counter_vec",
		Help: "Some help message",
	}, []string{"label1", "label2"})
)

func init() {
	// prometheus.DefaultRegisterer.MustRegister(simpleCounterVec)
	simpleCounterVec.WithLabelValues("foo", "bar")
	simpleCounterVec.WithLabelValues("qux", "baz")
	simpleCounterVec.WithLabelValues("1", "2")

	metricsFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		panic(err)
	}
	for _, f := range metricsFamilies {
		fmt.Println(f)
	}
	panic("DONE")
}

func main() {
	recovery.ReportPanics(func() {
		Run(NewProductionClient(), os.Args...)
	})
}

// Run runs the CLI, providing further command instructions by default.
func Run(client *cmd.Client, args ...string) {
	app := cmd.NewApp(client)
	client.Logger.ErrorIf(app.Run(args), "Error running app")
	if err := client.CloseLogger(); err != nil {
		log.Fatal(err)
	}
}

// NewProductionClient configures an instance of the CLI to be used
// in production.
func NewProductionClient() *cmd.Client {
	lggr, closeLggr := logger.NewLogger()
	cfg := config.NewGeneralConfig(lggr)

	prompter := cmd.NewTerminalPrompter()
	cookieAuth := cmd.NewSessionCookieAuthenticator(cfg, cmd.DiskCookieStore{Config: cfg}, lggr)
	sr := sessions.SessionRequest{}
	sessionRequestBuilder := cmd.NewFileSessionRequestBuilder(lggr)
	if credentialsFile := cfg.AdminCredentialsFile(); credentialsFile != "" {
		var err error
		sr, err = sessionRequestBuilder.Build(credentialsFile)
		if err != nil && !errors.Is(errors.Cause(err), cmd.ErrNoCredentialFile) && !os.IsNotExist(err) {
			lggr.Fatalw("Error loading API credentials", "error", err, "credentialsFile", credentialsFile)
		}
	}
	return &cmd.Client{
		Renderer:                       cmd.RendererTable{Writer: os.Stdout},
		Config:                         cfg,
		Logger:                         lggr,
		CloseLogger:                    closeLggr,
		AppFactory:                     cmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          cmd.TerminalKeyStoreAuthenticator{Prompter: prompter},
		FallbackAPIInitializer:         cmd.NewPromptingAPIInitializer(prompter),
		Runner:                         cmd.ChainlinkRunner{},
		HTTP:                           cmd.NewAuthenticatedHTTPClient(cfg, cookieAuth, sr),
		CookieAuthenticator:            cookieAuth,
		FileSessionRequestBuilder:      sessionRequestBuilder,
		PromptingSessionRequestBuilder: cmd.NewPromptingSessionRequestBuilder(prompter),
		ChangePasswordPrompter:         cmd.NewChangePasswordPrompter(),
		PasswordPrompter:               cmd.NewPasswordPrompter(),
	}
}
