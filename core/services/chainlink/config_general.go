package chainlink

import (
	_ "embed"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"

	ocrnetworking "github.com/smartcontractkit/libocr/networking"

	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/chains/solana"
	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/parse"
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// generalConfig is a wrapper to adapt Config to the config.GeneralConfig interface.
type generalConfig struct {
	inputTOML     string // user input, normalized via de/re-serialization
	effectiveTOML string // with default values included
	secretsTOML   string // with env overdies includes, redacted

	c       *Config // all fields non-nil (unless the legacy method signature return a pointer)
	secrets *Secrets

	logLevelDefault zapcore.Level

	appIDOnce sync.Once

	logMu sync.RWMutex // for the mutable fields Log.Level & Log.SQL

	passwordMu sync.RWMutex // passwords are set after initialization
}

// GeneralConfigOpts holds configuration options for creating a coreconfig.GeneralConfig via New().
//
// See ParseTOML to initilialize Config and Secrets from TOML.
type GeneralConfigOpts struct {
	ConfigStrings []string
	SecretsString string

	Config
	Secrets

	// OverrideFn is a *test-only* hook to override effective values.
	OverrideFn func(*Config, *Secrets)

	SkipEnv bool
}

// parseConfig sets Config from the given TOML string, overriding any existing duplicate Config fields.
func (o *GeneralConfigOpts) parseConfig(config string) error {
	var c Config
	if err2 := v2.DecodeTOML(strings.NewReader(config), &c); err2 != nil {
		return fmt.Errorf("failed to decode config TOML: %w", err2)
	}

	// Overrides duplicate fields
	if err4 := o.Config.SetFrom(&c); err4 != nil {
		return fmt.Errorf("invalid configuration: %w", err4)
	}
	return nil
}

// parseSecrets sets Secrets from the given TOML string.
func (o *GeneralConfigOpts) parseSecrets() (err error) {
	if err2 := v2.DecodeTOML(strings.NewReader(o.SecretsString), &o.Secrets); err2 != nil {
		return fmt.Errorf("failed to decode secrets TOML: %w", err2)
	}
	return nil
}

// New returns a coreconfig.GeneralConfig for the given options.
func (o GeneralConfigOpts) New() (GeneralConfig, error) {
	for _, c := range o.ConfigStrings {
		err := o.parseConfig(c)
		if err != nil {
			return nil, err
		}
	}

	if o.SecretsString != "" {
		err := o.parseSecrets()
		if err != nil {
			return nil, err
		}
	}

	input, err := o.Config.TOMLString()
	if err != nil {
		return nil, err
	}

	o.Config.setDefaults()
	if !o.SkipEnv {
		err = o.Secrets.setEnv()
		if err != nil {
			return nil, err
		}
	}

	if fn := o.OverrideFn; fn != nil {
		fn(&o.Config, &o.Secrets)
	}

	effective, err := o.Config.TOMLString()
	if err != nil {
		return nil, err
	}

	secrets, err := o.Secrets.TOMLString()
	if err != nil {
		return nil, err
	}

	cfg := &generalConfig{
		inputTOML:     input,
		effectiveTOML: effective,
		secretsTOML:   secrets,
		c:             &o.Config,
		secrets:       &o.Secrets,
	}
	if lvl := o.Config.Log.Level; lvl != nil {
		cfg.logLevelDefault = zapcore.Level(*lvl)
	}

	return cfg, nil
}

func (g *generalConfig) EVMConfigs() evmcfg.EVMConfigs {
	return g.c.EVM
}

func (g *generalConfig) CosmosConfigs() cosmos.CosmosConfigs {
	return g.c.Cosmos
}

func (g *generalConfig) SolanaConfigs() solana.SolanaConfigs {
	return g.c.Solana
}

func (g *generalConfig) StarknetConfigs() starknet.StarknetConfigs {
	return g.c.Starknet
}

func (g *generalConfig) Validate() error {
	return g.validate(g.secrets.Validate)
}

func (g *generalConfig) validate(secretsValidationFn func() error) error {
	err := multierr.Combine(
		validateEnv(),
		g.c.Validate(),
		secretsValidationFn(),
	)

	_, errList := utils.MultiErrorList(err)
	return errList
}

func (g *generalConfig) ValidateDB() error {
	return g.validate(g.secrets.ValidateDB)
}

//go:embed legacy.env
var emptyStringsEnv string

// validateEnv returns an error if any legacy environment variables are set, unless a v2 equivalent exists with the same value.
func validateEnv() (err error) {
	defer func() {
		if err != nil {
			_, err = utils.MultiErrorList(err)
			err = fmt.Errorf("invalid environment: %w", err)
		}
	}()
	for _, kv := range strings.Split(emptyStringsEnv, "\n") {
		if strings.TrimSpace(kv) == "" {
			continue
		}
		i := strings.Index(kv, "=")
		if i == -1 {
			return errors.Errorf("malformed .env file line: %s", kv)
		}
		k := kv[:i]
		_, ok := os.LookupEnv(k)
		if ok {
			err = multierr.Append(err, fmt.Errorf("environment variable %s must not be set: %v", k, v2.ErrUnsupported))
		}
	}
	return
}

func (g *generalConfig) LogConfiguration(log coreconfig.LogfFn) {
	log("# Secrets:\n%s\n", g.secretsTOML)
	log("# Input Configuration:\n%s\n", g.inputTOML)
	log("# Effective Configuration, with defaults applied:\n%s\n", g.effectiveTOML)
}

// ConfigTOML implements chainlink.ConfigV2
func (g *generalConfig) ConfigTOML() (user, effective string) {
	return g.inputTOML, g.effectiveTOML
}

func (g *generalConfig) Feature() coreconfig.Feature {
	return &featureConfig{c: g.c.Feature}
}

func (g *generalConfig) FeatureFeedsManager() bool {
	return *g.c.Feature.FeedsManager
}

func (g *generalConfig) OCR() config.OCR {
	return &ocrConfig{c: g.c.OCR}
}

func (g *generalConfig) OCR2Enabled() bool {
	return *g.c.OCR2.Enabled
}

func (g *generalConfig) FeatureLogPoller() bool {
	return *g.c.Feature.LogPoller
}

func (g *generalConfig) FeatureUICSAKeys() bool {
	return *g.c.Feature.UICSAKeys
}

func (g *generalConfig) AutoPprof() config.AutoPprof {
	return &autoPprofConfig{c: g.c.AutoPprof, rootDir: g.RootDir}
}

func (g *generalConfig) EVMEnabled() bool {
	for _, c := range g.c.EVM {
		if c.IsEnabled() {
			return true
		}
	}
	return false
}

func (g *generalConfig) EVMRPCEnabled() bool {
	for _, c := range g.c.EVM {
		if c.IsEnabled() {
			if len(c.Nodes) > 0 {
				return true
			}
		}
	}
	return false
}

func (g *generalConfig) DefaultChainID() *big.Int {
	for _, c := range g.c.EVM {
		if c.IsEnabled() {
			return (*big.Int)(c.ChainID)
		}
	}
	return nil
}

func (g *generalConfig) SolanaEnabled() bool {
	for _, c := range g.c.Solana {
		if c.IsEnabled() {
			return true
		}
	}
	return false
}

func (g *generalConfig) CosmosEnabled() bool {
	for _, c := range g.c.Cosmos {
		if c.IsEnabled() {
			return true
		}
	}
	return false
}

func (g *generalConfig) StarkNetEnabled() bool {
	for _, c := range g.c.Starknet {
		if c.IsEnabled() {
			return true
		}
	}
	return false
}

func (g *generalConfig) WebServer() config.WebServer {
	return &webServerConfig{c: g.c.WebServer, rootDir: g.RootDir}
}

func (g *generalConfig) AutoPprofBlockProfileRate() int {
	return int(*g.c.AutoPprof.BlockProfileRate)
}

func (g *generalConfig) AutoPprofCPUProfileRate() int {
	return int(*g.c.AutoPprof.CPUProfileRate)
}

func (g *generalConfig) AutoPprofGatherDuration() models.Duration {
	return models.MustMakeDuration(g.c.AutoPprof.GatherDuration.Duration())
}

func (g *generalConfig) AutoPprofGatherTraceDuration() models.Duration {
	return models.MustMakeDuration(g.c.AutoPprof.GatherTraceDuration.Duration())
}

func (g *generalConfig) AutoPprofGoroutineThreshold() int {
	return int(*g.c.AutoPprof.GoroutineThreshold)
}

func (g *generalConfig) AutoPprofMaxProfileSize() utils.FileSize {
	return *g.c.AutoPprof.MaxProfileSize
}

func (g *generalConfig) AutoPprofMemProfileRate() int {
	return int(*g.c.AutoPprof.MemProfileRate)
}

func (g *generalConfig) AutoPprofMemThreshold() utils.FileSize {
	return *g.c.AutoPprof.MemThreshold
}

func (g *generalConfig) AutoPprofMutexProfileFraction() int {
	return int(*g.c.AutoPprof.MutexProfileFraction)
}

func (g *generalConfig) AutoPprofPollInterval() models.Duration {
	return *g.c.AutoPprof.PollInterval
}

func (g *generalConfig) AutoPprofProfileRoot() string {
	s := *g.c.AutoPprof.ProfileRoot
	if s == "" {
		s = filepath.Join(g.RootDir(), "pprof")
	}
	return s
}

func (g *generalConfig) Database() coreconfig.Database {
	return &databaseConfig{c: g.c.Database, s: g.secrets.Secrets.Database, logSQL: g.logSQL}
}

func (g *generalConfig) ShutdownGracePeriod() time.Duration {
	return g.c.ShutdownGracePeriod.Duration()
}

func (g *generalConfig) Explorer() config.Explorer {
	return &explorerConfig{s: g.secrets.Explorer, explorerURL: g.c.ExplorerURL}
}

func (g *generalConfig) ExplorerURL() *url.URL {
	u := (*url.URL)(g.c.ExplorerURL)
	if *u == zeroURL {
		u = nil
	}
	return u
}

func (g *generalConfig) FluxMonitor() config.FluxMonitor {
	return &fluxMonitorConfig{c: g.c.FluxMonitor}
}

func (g *generalConfig) InsecureFastScrypt() bool {
	return *g.c.InsecureFastScrypt
}

func (g *generalConfig) JobPipelineReaperInterval() time.Duration {
	return g.c.JobPipeline.ReaperInterval.Duration()
}

func (g *generalConfig) JobPipelineResultWriteQueueDepth() uint64 {
	return uint64(*g.c.JobPipeline.ResultWriteQueueDepth)
}

func (g *generalConfig) JobPipeline() coreconfig.JobPipeline {
	return &jobPipelineConfig{c: g.c.JobPipeline}
}

func (g *generalConfig) Keeper() config.Keeper {
	return &keeperConfig{c: g.c.Keeper}
}

func (g *generalConfig) Log() config.Log {
	return &logConfig{c: g.c.Log, rootDir: g.RootDir, level: g.logLevel, defaultLevel: g.logLevelDefault}
}

func (g *generalConfig) OCR2() config.OCR2 {
	return &ocr2Config{c: g.c.OCR2}
}

func (g *generalConfig) P2P() config.P2P {
	return &p2p{c: g.c.P2P}
}

func (g *generalConfig) P2PNetworkingStack() (n ocrnetworking.NetworkingStack) {
	return g.c.P2P.NetworkStack()
}

func (g *generalConfig) P2PNetworkingStackRaw() string {
	return g.c.P2P.NetworkStack().String()
}

func (g *generalConfig) P2PPeerID() p2pkey.PeerID {
	return *g.c.P2P.PeerID
}

func (g *generalConfig) P2PPeerIDRaw() string {
	return g.c.P2P.PeerID.String()
}

func (g *generalConfig) P2PIncomingMessageBufferSize() int {
	return int(*g.c.P2P.IncomingMessageBufferSize)
}

func (g *generalConfig) P2POutgoingMessageBufferSize() int {
	return int(*g.c.P2P.OutgoingMessageBufferSize)
}

func (g *generalConfig) Pyroscope() config.Pyroscope {
	return &pyroscopeConfig{c: g.c.Pyroscope, s: g.secrets.Pyroscope}
}

func (g *generalConfig) RootDir() string {
	d := *g.c.RootDir
	h, err := parse.HomeDir(d)
	if err != nil {
		panic(err) // never happens since we validate that the RootDir is expandable in config.Core.ValidateConfig().
	}
	return h
}

func (g *generalConfig) TelemetryIngress() coreconfig.TelemetryIngress {
	return &telemetryIngressConfig{
		c: g.c.TelemetryIngress,
	}
}

func (g *generalConfig) AuditLogger() coreconfig.AuditLogger {
	return auditLoggerConfig{c: g.c.AuditLogger}
}

func (g *generalConfig) Insecure() config.Insecure {
	return &insecureConfig{c: g.c.Insecure}
}

func (g *generalConfig) Sentry() coreconfig.Sentry {
	return sentryConfig{g.c.Sentry}
}

func (g *generalConfig) Password() coreconfig.Password {
	return &passwordConfig{keystore: g.keystorePassword, vrf: g.vrfPassword}
}

func (g *generalConfig) Prometheus() coreconfig.Prometheus {
	return &prometheusConfig{s: g.secrets.Prometheus}
}

func (g *generalConfig) Mercury() coreconfig.Mercury {
	return &mercuryConfig{s: g.secrets.Mercury}
}

func (g *generalConfig) Threshold() coreconfig.Threshold {
	return &thresholdConfig{s: g.secrets.Threshold}
}

var (
	zeroURL        = url.URL{}
	zeroSha256Hash = models.Sha256Hash{}
)
