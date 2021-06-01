package orm

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/dialects"

	"gorm.io/gorm"

	"github.com/multiformats/go-multiaddr"

	ocr "github.com/smartcontractkit/libocr/offchainreporting"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	ethCore "github.com/ethereum/go-ethereum/core"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gorilla/securecookie"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

// this permission grants read / write accccess to file owners only
const readWritePerms = os.FileMode(0600)

var (
	ErrUnset   = errors.New("env var unset")
	ErrInvalid = errors.New("env var invalid")

	configFileNotFoundError = reflect.TypeOf(viper.ConfigFileNotFoundError{})

	// keyed by ChainID
	ChainSpecificDefaults map[int64]ChainSpecificDefaultSet
	// If the chain is unknown, fallback to the general defaults
	GeneralDefaults ChainSpecificDefaultSet
)

type (
	// Config holds parameters used by the application which can be overridden by
	// setting environment variables.
	//
	// If you add an entry here which does not contain sensitive information, you
	// should also update presenters.ConfigWhitelist and cmd_test.TestClient_RunNodeShowsEnv.
	Config struct {
		viper            *viper.Viper
		SecretGenerator  SecretGenerator
		runtimeStore     *ORM
		randomP2PPort    uint16
		randomP2PPortMtx *sync.RWMutex
		Dialect          dialects.DialectName
		AdvisoryLockID   int64
	}

	// ChainSpecificDefaultSet us a list of defaults specific to a particular chain ID
	ChainSpecificDefaultSet struct {
		EnableLegacyJobPipeline          bool
		EthGasLimitDefault               uint64
		EthGasBumpThreshold              uint64
		EthGasBumpWei                    big.Int
		EthGasPriceDefault               big.Int
		EthMaxGasPriceWei                big.Int
		EthMinGasPriceWei                big.Int
		EthFinalityDepth                 uint
		EthHeadTrackerHistoryDepth       uint
		EthHeadTrackerSamplingInterval   time.Duration
		EthBalanceMonitorBlockDelay      uint16
		EthTxResendAfterThreshold        time.Duration
		GasUpdaterBatchSize              *uint32
		GasUpdaterBlockDelay             uint16
		GasUpdaterBlockHistorySize       uint16
		GasUpdaterEnabled                bool
		MinIncomingConfirmations         uint32
		MinRequiredOutgoingConfirmations uint64
		OptimismGasFees                  bool
		MinimumContractPayment           *assets.Link
	}
)

func init() {
	ChainSpecificDefaults = make(map[int64]ChainSpecificDefaultSet)
	// FIXME: Workaround `websocket: read limit exceeded` until https://app.clubhouse.io/chainlinklabs/story/6717/geth-websockets-can-sometimes-go-bad-under-heavy-load-proposal-for-eth-node-balancer
	var defaultGasUpdaterBatchSize uint32 = 4

	mainnet := ChainSpecificDefaultSet{
		EnableLegacyJobPipeline:          true,
		EthGasBumpThreshold:              3,
		EthGasBumpWei:                    *big.NewInt(5000000000),    // 5 Gwei
		EthGasPriceDefault:               *big.NewInt(20000000000),   // 20 Gwei
		EthMaxGasPriceWei:                *big.NewInt(5000000000000), // 5000 Gwei
		EthMinGasPriceWei:                *big.NewInt(1000000000),    // 1 Gwei
		EthFinalityDepth:                 50,
		EthHeadTrackerHistoryDepth:       100,
		EthHeadTrackerSamplingInterval:   1 * time.Second,
		EthBalanceMonitorBlockDelay:      1,
		EthTxResendAfterThreshold:        1 * time.Minute,
		GasUpdaterBlockDelay:             1,
		GasUpdaterBlockHistorySize:       24,
		GasUpdaterBatchSize:              &defaultGasUpdaterBatchSize,
		GasUpdaterEnabled:                true,
		MinIncomingConfirmations:         3,
		MinRequiredOutgoingConfirmations: 12,
	}

	// NOTE: There are probably other variables we can tweak for Kovan and other
	// test chains, but requires more in-depth research on their consensus
	// mechanisms. For now, mainnet defaults ought to be safe enough for testnet.
	kovan := mainnet

	// xDai currently uses AuRa (like Parity) consensus so finality rules will be similar to parity
	// See: https://www.poa.network/for-users/whitepaper/poadao-v1/proof-of-authority
	// NOTE: xDai is planning to move to Honeybadger BFT which might have different finality guarantees
	// https://www.xdaichain.com/for-validators/consensus/honeybadger-bft-consensus
	// For worst case re-org depth on AuRa, assume 2n+2 (see: https://github.com/poanetwork/wiki/wiki/Aura-Consensus-Protocol-Audit)
	// With xDai's current maximum of 19 validators then 40 blocks is the maximum possible re-org)
	// The mainnet default of 50 blocks is ok here
	xDai := mainnet
	xDai.EthGasBumpThreshold = 8                       // mainnet * 2.8 ish (5s vs 13s block time)
	xDai.EthGasPriceDefault = *big.NewInt(1000000000)  // 1 Gwei
	xDai.EthMinGasPriceWei = *big.NewInt(1000000000)   // 1 Gwei is the minimum accepted by the validators (unless whitelisted)
	xDai.EthMaxGasPriceWei = *big.NewInt(500000000000) // 500 Gwei

	// BSC uses Clique consensus with ~3s block times
	// Clique offers finality within (N/2)+1 blocks where N is number of signers
	// There are 21 BSC validators so theoretically finality should occur after 21/2+1 = 11 blocks
	bscMainnet := ChainSpecificDefaultSet{
		EnableLegacyJobPipeline:          true,
		EthGasBumpThreshold:              12,                        // mainnet * 4 (3s vs 13s block time)
		EthGasBumpWei:                    *big.NewInt(5000000000),   // 5 Gwei
		EthGasPriceDefault:               *big.NewInt(5000000000),   // 5 Gwei
		EthMaxGasPriceWei:                *big.NewInt(500000000000), // 500 Gwei
		EthMinGasPriceWei:                *big.NewInt(1000000000),   // 1 Gwei
		EthFinalityDepth:                 50,                        // Keeping this >> 11 because it's not expensive and gives us a safety margin
		EthHeadTrackerHistoryDepth:       100,
		EthHeadTrackerSamplingInterval:   1 * time.Second,
		EthBalanceMonitorBlockDelay:      2,
		EthTxResendAfterThreshold:        1 * time.Minute,
		GasUpdaterBlockDelay:             2,
		GasUpdaterBlockHistorySize:       24,
		GasUpdaterBatchSize:              &defaultGasUpdaterBatchSize,
		GasUpdaterEnabled:                true,
		MinIncomingConfirmations:         3,
		MinRequiredOutgoingConfirmations: 12,
	}

	hecoMainnet := bscMainnet

	// Matic has a 1s block time and looser finality guarantees than Ethereum.
	polygonMatic := ChainSpecificDefaultSet{
		EnableLegacyJobPipeline:          true,
		EthGasBumpThreshold:              39,                        // mainnet * 13
		EthGasBumpWei:                    *big.NewInt(5000000000),   // 5 Gwei
		EthGasPriceDefault:               *big.NewInt(1000000000),   // 1 Gwei
		EthMaxGasPriceWei:                *big.NewInt(500000000000), // 500 Gwei
		EthMinGasPriceWei:                *big.NewInt(1000000000),   // 1 Gwei
		EthFinalityDepth:                 200,                       // A sprint is 64 blocks long and doesn't guarantee finality. To be safe, we take three sprints (192 blocks) plus a safety margin
		EthHeadTrackerHistoryDepth:       250,                       // EthFinalityDepth + safety margin
		EthHeadTrackerSamplingInterval:   1 * time.Second,
		EthBalanceMonitorBlockDelay:      13,              // equivalent of 1 eth block seems reasonable
		EthTxResendAfterThreshold:        5 * time.Minute, // 5 minutes is roughly 300 blocks on Matic. Since re-orgs occur often and can be deep, we want to avoid overloading the node with a ton of re-sent unconfirmed transactions.
		GasUpdaterBlockDelay:             32,              // Delay needs to be large on matic since re-orgs are so frequent at the top level
		GasUpdaterBlockHistorySize:       128,
		GasUpdaterBatchSize:              &defaultGasUpdaterBatchSize,
		GasUpdaterEnabled:                true,
		MinIncomingConfirmations:         39, // mainnet * 13 (1s vs 13s block time)
		MinRequiredOutgoingConfirmations: 39, // mainnet * 13
	}

	// Optimism is an L2 chain. Pending proper L2 support, for now we rely on Optimism's sequencer
	optimism := ChainSpecificDefaultSet{
		EnableLegacyJobPipeline:          false,
		EthGasBumpThreshold:              0, // Never bump gas on optimism
		EthFinalityDepth:                 1, // Sequencer offers absolute finality as long as no re-org longer than 20 blocks occurs on main chain, this event would require special handling (new txm)
		EthHeadTrackerHistoryDepth:       10,
		EthHeadTrackerSamplingInterval:   1 * time.Second,
		EthBalanceMonitorBlockDelay:      0,
		EthTxResendAfterThreshold:        15 * time.Second,
		GasUpdaterBlockHistorySize:       0, // Force an error if someone set GAS_UPDATER_ENABLED=true by accident; we never want to run the gas updater on optimism
		GasUpdaterEnabled:                false,
		MinIncomingConfirmations:         1,
		MinRequiredOutgoingConfirmations: 0,
		OptimismGasFees:                  true,
	}

	// Fantom
	fantomTestnet := ChainSpecificDefaultSet{
		// Fantom specific
		EthGasLimitDefault:               500000,
		EthGasPriceDefault:               *big.NewInt(15000000000),
		EthMaxGasPriceWei:                *big.NewInt(100000000000),
		MinIncomingConfirmations:         3,
		MinRequiredOutgoingConfirmations: 2,
		MinimumContractPayment:           assets.NewLink(3500000000000000),

		// Same as mainnet
		EnableLegacyJobPipeline:        true,
		EthGasBumpThreshold:            3,
		EthGasBumpWei:                  *big.NewInt(5000000000), // 5 Gwei
		EthMinGasPriceWei:              *big.NewInt(1000000000), // 1 Gwei
		EthFinalityDepth:               50,
		EthHeadTrackerHistoryDepth:     100,
		EthHeadTrackerSamplingInterval: 1 * time.Second,
		EthBalanceMonitorBlockDelay:    1,
		EthTxResendAfterThreshold:      1 * time.Minute,
		GasUpdaterBlockDelay:           1,
		GasUpdaterBlockHistorySize:     24,
		GasUpdaterBatchSize:            &defaultGasUpdaterBatchSize,
		GasUpdaterEnabled:              true,
		HeadTimeBudget:                 13 * time.Second,
	}
	fantomMainnet := fantomTestnet

	GeneralDefaults = mainnet
	ChainSpecificDefaults[1] = mainnet
	ChainSpecificDefaults[42] = kovan

	ChainSpecificDefaults[10] = optimism
	ChainSpecificDefaults[69] = optimism
	ChainSpecificDefaults[420] = optimism

	ChainSpecificDefaults[56] = bscMainnet
	ChainSpecificDefaults[128] = hecoMainnet
	ChainSpecificDefaults[250] = fantomMainnet
	ChainSpecificDefaults[4002] = fantomTestnet
	ChainSpecificDefaults[80001] = polygonMatic

	ChainSpecificDefaults[100] = xDai
}

func chainSpecificConfig(c Config) ChainSpecificDefaultSet {
	chainID := c.ChainID().Int64()
	if cset, exists := ChainSpecificDefaults[chainID]; exists {
		return cset
	}
	return GeneralDefaults
}

// NewConfig returns the config with the environment variables set to their
// respective fields, or their defaults if environment variables are not set.
func NewConfig() *Config {
	v := viper.New()
	return newConfigWithViper(v)
}

func newConfigWithViper(v *viper.Viper) *Config {
	schemaT := reflect.TypeOf(ConfigSchema{})
	for index := 0; index < schemaT.NumField(); index++ {
		item := schemaT.FieldByIndex([]int{index})
		name := item.Tag.Get("env")
		def, exists := item.Tag.Lookup("default")
		if exists {
			v.SetDefault(name, def)
		}
		_ = v.BindEnv(name, name)
	}

	config := &Config{
		viper:            v,
		SecretGenerator:  filePersistedSecretGenerator{},
		randomP2PPortMtx: new(sync.RWMutex),
	}

	if err := utils.EnsureDirAndMaxPerms(config.RootDir(), os.FileMode(0700)); err != nil {
		logger.Fatalf(`Error creating root directory "%s": %+v`, config.RootDir(), err)
	}

	v.SetConfigName("chainlink")
	v.AddConfigPath(config.RootDir())
	err := v.ReadInConfig()
	if err != nil && reflect.TypeOf(err) != configFileNotFoundError {
		logger.Warnf("Unable to load config file: %v\n", err)
	}

	return config
}

// Validate performs basic sanity checks on config and returns error if any
// misconfiguration would be fatal to the application
func (c *Config) Validate() error {
	ethGasBumpPercent := c.EthGasBumpPercent()
	if uint64(ethGasBumpPercent) < ethCore.DefaultTxPoolConfig.PriceBump {
		return errors.Errorf(
			"ETH_GAS_BUMP_PERCENT of %v may not be less than Geth's default of %v",
			c.EthGasBumpPercent(),
			ethCore.DefaultTxPoolConfig.PriceBump,
		)
	}

	if uint32(c.EthGasBumpTxDepth()) > c.EthMaxInFlightTransactions() {
		return errors.New("ETH_GAS_BUMP_TX_DEPTH must be less than or equal to ETH_MAX_IN_FLIGHT_TRANSACTIONS")
	}
	if c.EthMinGasPriceWei().Cmp(c.EthGasPriceDefault()) > 0 {
		return errors.New("ETH_MIN_GAS_PRICE_WEI must be less than or equal to ETH_GAS_PRICE_DEFAULT")
	}
	if c.EthMaxGasPriceWei().Cmp(c.EthGasPriceDefault()) < 0 {
		return errors.New("ETH_MAX_GAS_PRICE_WEI must be greater than or equal to ETH_GAS_PRICE_DEFAULT")
	}

	if c.EthHeadTrackerHistoryDepth() < c.EthFinalityDepth() {
		return errors.New("ETH_HEAD_TRACKER_HISTORY_DEPTH must be equal to or greater than ETH_FINALITY_DEPTH")
	}

	if c.GasUpdaterEnabled() && c.GasUpdaterBlockHistorySize() <= 0 {
		return errors.New("GAS_UPDATER_BLOCK_HISTORY_SIZE must be greater than or equal to 1 if gas updater is enabled")
	}

	if c.P2PAnnouncePort() != 0 && c.P2PAnnounceIP() == nil {
		return errors.Errorf("P2P_ANNOUNCE_PORT was given as %v but P2P_ANNOUNCE_IP was unset. You must also set P2P_ANNOUNCE_IP if P2P_ANNOUNCE_PORT is set", c.P2PAnnouncePort())
	}

	if c.EthFinalityDepth() < 1 {
		return errors.New("ETH_FINALITY_DEPTH must be greater than or equal to 1")
	}

	if c.MinIncomingConfirmations() < 1 {
		return errors.New("MIN_INCOMING_CONFIRMATIONS must be greater than or equal to 1")
	}

	var override time.Duration
	lc := ocrtypes.LocalConfig{
		BlockchainTimeout:                      c.OCRBlockchainTimeout(override),
		ContractConfigConfirmations:            c.OCRContractConfirmations(0),
		ContractConfigTrackerPollInterval:      c.OCRContractPollInterval(override),
		ContractConfigTrackerSubscribeInterval: c.OCRContractSubscribeInterval(override),
		ContractTransmitterTransmitTimeout:     c.OCRContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                        c.OCRDatabaseTimeout(),
		DataSourceTimeout:                      c.OCRObservationTimeout(override),
		DataSourceGracePeriod:                  c.OCRObservationGracePeriod(),
	}
	if err := ocr.SanityCheckLocalConfig(lc); err != nil {
		return err
	}
	if _, err := c.P2PPeerID(nil); errors.Cause(err) == ErrInvalid {
		return err
	}
	if _, err := c.OCRKeyBundleID(nil); errors.Cause(err) == ErrInvalid {
		return err
	}
	if _, err := c.OCRTransmitterAddress(nil); errors.Cause(err) == ErrInvalid {
		return err
	}
	if peers, err := c.P2PBootstrapPeers(nil); err == nil {
		for i := range peers {
			if _, err := multiaddr.NewMultiaddr(peers[i]); err != nil {
				return errors.Errorf("p2p bootstrap peer %d is invalid: err %v", i, err)
			}
		}
	}
	if me := c.OCRMonitoringEndpoint(""); me != "" {
		if _, err := url.Parse(me); err != nil {
			return errors.Wrapf(err, "invalid monitoring url: %s", me)
		}
	}
	return nil
}

// SetRuntimeStore tells the configuration system to use a store for retrieving
// configuration variables that can be configured at runtime.
func (c *Config) SetRuntimeStore(orm *ORM) {
	c.runtimeStore = orm
}

// Set a specific configuration variable
func (c Config) Set(name string, value interface{}) {
	schemaT := reflect.TypeOf(ConfigSchema{})
	for index := 0; index < schemaT.NumField(); index++ {
		item := schemaT.FieldByIndex([]int{index})
		envName := item.Tag.Get("env")
		if envName == name {
			c.viper.Set(name, value)
			return
		}
	}
	logger.Panicf("No configuration parameter for %s", name)
}

const defaultPostgresAdvisoryLockID int64 = 1027321974924625846

func (c Config) GetAdvisoryLockIDConfiguredOrDefault() int64 {
	if c.AdvisoryLockID == 0 {
		return defaultPostgresAdvisoryLockID
	}
	return c.AdvisoryLockID
}

func (c Config) GetDatabaseDialectConfiguredOrDefault() dialects.DialectName {
	if c.Dialect == "" {
		return dialects.Postgres
	}
	return c.Dialect
}

// AllowOrigins returns the CORS hosts used by the frontend.
func (c Config) AllowOrigins() string {
	return c.viper.GetString(EnvVarName("AllowOrigins"))
}

// AdminCredentialsFile points to text file containing admnn credentials for logging in
func (c Config) AdminCredentialsFile() string {
	fieldName := "AdminCredentialsFile"
	file := c.viper.GetString(EnvVarName(fieldName))
	defaultValue, _ := defaultValue(fieldName)
	if file == defaultValue {
		return filepath.Join(c.RootDir(), "apicredentials")
	}
	return file
}

// AuthenticatedRateLimit defines the threshold to which requests authenticated requests get limited
func (c Config) AuthenticatedRateLimit() int64 {
	return c.viper.GetInt64(EnvVarName("AuthenticatedRateLimit"))
}

// AuthenticatedRateLimitPeriod defines the period to which authenticated requests get limited
func (c Config) AuthenticatedRateLimitPeriod() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("AuthenticatedRateLimitPeriod", parseDuration).(time.Duration))
}

// BalanceMonitorEnabled enables the balance monitor
func (c Config) BalanceMonitorEnabled() bool {
	return c.viper.GetBool(EnvVarName("BalanceMonitorEnabled"))
}

// BlockBackfillDepth specifies the number of blocks before the current HEAD that the
// log broadcaster will try to re-consume logs from
func (c Config) BlockBackfillDepth() uint64 {
	return c.getWithFallback("BlockBackfillDepth", parseUint64).(uint64)
}

// BridgeResponseURL represents the URL for bridges to send a response to.
func (c Config) BridgeResponseURL() *url.URL {
	return c.getWithFallback("BridgeResponseURL", parseURL).(*url.URL)
}

// ChainID represents the chain ID to use for transactions.
func (c Config) ChainID() *big.Int {
	return c.getWithFallback("ChainID", parseBigInt).(*big.Int)
}

// ClientNodeURL is the URL of the Ethereum node this Chainlink node should connect to.
func (c Config) ClientNodeURL() string {
	return c.viper.GetString(EnvVarName("ClientNodeURL"))
}

// FeatureCronV2 enables the Cron v2 feature.
func (c Config) FeatureCronV2() bool {
	return c.getWithFallback("FeatureCronV2", parseBool).(bool)
}

func (c Config) DatabaseListenerMinReconnectInterval() time.Duration {
	return c.getWithFallback("DatabaseListenerMinReconnectInterval", parseDuration).(time.Duration)
}

func (c Config) DatabaseListenerMaxReconnectDuration() time.Duration {
	return c.getWithFallback("DatabaseListenerMaxReconnectDuration", parseDuration).(time.Duration)
}

func (c Config) DatabaseMaximumTxDuration() time.Duration {
	return c.getWithFallback("DatabaseMaximumTxDuration", parseDuration).(time.Duration)
}

// DatabaseBackupMode sets the database backup mode
func (c Config) DatabaseBackupMode() DatabaseBackupMode {
	return c.getWithFallback("DatabaseBackupMode", parseDatabaseBackupMode).(DatabaseBackupMode)
}

// DatabaseBackupFrequency turns on the periodic database backup if set to a positive value
// DatabaseBackupMode must be then set to a value other than "none"
func (c Config) DatabaseBackupFrequency() time.Duration {
	return c.getWithFallback("DatabaseBackupFrequency", parseDuration).(time.Duration)
}

// DatabaseBackupURL configures the URL for the database to backup, if it's to be different from the main on
func (c Config) DatabaseBackupURL() *url.URL {
	s := c.viper.GetString(EnvVarName("DatabaseBackupURL"))
	if s == "" {
		return nil
	}
	uri, err := url.Parse(s)
	if err != nil {
		logger.Error("invalid database backup url %s", s)
		return nil
	}
	return uri
}

// DatabaseBackupDir configures the directory for saving the backup file, if it's to be different from default one located in the RootDir
func (c Config) DatabaseBackupDir() string {
	return c.viper.GetString(EnvVarName("DatabaseBackupDir"))
}

// DatabaseTimeout represents how long to tolerate non response from the DB.
func (c Config) DatabaseTimeout() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("DatabaseTimeout", parseDuration).(time.Duration))
}

// GlobalLockRetryInterval represents how long to wait before trying again to get the global advisory lock.
func (c Config) GlobalLockRetryInterval() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("GlobalLockRetryInterval", parseDuration).(time.Duration))
}

// DatabaseURL configures the URL for chainlink to connect to. This must be
// a properly formatted URL, with a valid scheme (postgres://)
func (c Config) DatabaseURL() url.URL {
	s := c.viper.GetString(EnvVarName("DatabaseURL"))
	uri, err := url.Parse(s)
	if err != nil {
		logger.Error("invalid database url %s", s)
		return url.URL{}
	}
	if uri.String() == "" {
		return *uri
	}
	static.SetConsumerName(uri, "Default")
	return *uri
}

// MigrateDatabase determines whether the database will be automatically
// migrated on application startup if set to true
func (c Config) MigrateDatabase() bool {
	return c.viper.GetBool(EnvVarName("MigrateDatabase"))
}

// DefaultMaxHTTPAttempts defines the limit for HTTP requests.
func (c Config) DefaultMaxHTTPAttempts() uint {
	return uint(c.getWithFallback("DefaultMaxHTTPAttempts", parseUint64).(uint64))
}

// DefaultHTTPLimit defines the size limit for HTTP requests and responses
func (c Config) DefaultHTTPLimit() int64 {
	return c.viper.GetInt64(EnvVarName("DefaultHTTPLimit"))
}

// DefaultHTTPTimeout defines the default timeout for http requests
func (c Config) DefaultHTTPTimeout() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("DefaultHTTPTimeout", parseDuration).(time.Duration))
}

// DefaultHTTPAllowUnrestrictedNetworkAccess controls whether http requests are unrestricted by default
// It is recommended that this be left disabled
func (c Config) DefaultHTTPAllowUnrestrictedNetworkAccess() bool {
	return c.viper.GetBool(EnvVarName("DefaultHTTPAllowUnrestrictedNetworkAccess"))
}

// Dev configures "development" mode for chainlink.
func (c Config) Dev() bool {
	return c.viper.GetBool(EnvVarName("Dev"))
}

// EnableExperimentalAdapters enables support for experimental adapters
func (c Config) EnableExperimentalAdapters() bool {
	return c.viper.GetBool(EnvVarName("EnableExperimentalAdapters"))
}

// EnableLegacyJobPipeline enables the v1 job pipeline (JSON job specs)
func (c Config) EnableLegacyJobPipeline() bool {
	if c.viper.IsSet(EnvVarName("EnableLegacyJobPipeline")) {
		return c.viper.GetBool(EnvVarName("EnableLegacyJobPipeline"))
	}
	return chainSpecificConfig(c).EnableLegacyJobPipeline
}

// FeatureExternalInitiators enables the External Initiator feature.
func (c Config) FeatureExternalInitiators() bool {
	return c.viper.GetBool(EnvVarName("FeatureExternalInitiators"))
}

// FeatureFluxMonitor enables the Flux Monitor job type.
func (c Config) FeatureFluxMonitor() bool {
	return c.viper.GetBool(EnvVarName("FeatureFluxMonitor"))
}

// FeatureFluxMonitorV2 enables the Flux Monitor v2 job type.
func (c Config) FeatureFluxMonitorV2() bool {
	return c.getWithFallback("FeatureFluxMonitorV2", parseBool).(bool)
}

// FeatureOffchainReporting enables the Flux Monitor job type.
func (c Config) FeatureOffchainReporting() bool {
	return c.viper.GetBool(EnvVarName("FeatureOffchainReporting"))
}

// FeatureWebhookV2 enables the Webhook v2 job type
func (c Config) FeatureWebhookV2() bool {
	return c.getWithFallback("FeatureWebhookV2", parseBool).(bool)
}

// MaximumServiceDuration is the maximum time that a service agreement can run
// from after the time it is created. Default 1 year = 365 * 24h = 8760h
func (c Config) MaximumServiceDuration() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("MaximumServiceDuration", parseDuration).(time.Duration))
}

// MinimumServiceDuration is the shortest duration from now that a service is
// allowed to run.
func (c Config) MinimumServiceDuration() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("MinimumServiceDuration", parseDuration).(time.Duration))
}

// EthBalanceMonitorBlockDelay is the number of blocks that the balance monitor
// trails behind head. This is required e.g. for Infura because they will often
// announce a new head, then route a request to a different node which does not
// have this head yet.
func (c Config) EthBalanceMonitorBlockDelay() uint16 {
	if c.viper.IsSet(EnvVarName("EthBalanceMonitorBlockDelay")) {
		return uint16(c.viper.GetUint32(EnvVarName("EthBalanceMonitorBlockDelay")))
	}
	return chainSpecificConfig(c).EthBalanceMonitorBlockDelay
}

// EthRPCDefaultBatchSize controls the number of receipts fetched in each
// request in the EthConfirmer
func (c Config) EthRPCDefaultBatchSize() uint32 {
	return c.viper.GetUint32(EnvVarName("EthRPCDefaultBatchSize"))
}

// EthGasBumpThreshold is the number of blocks to wait before bumping gas again on unconfirmed transactions
// Set to 0 to disable gas bumping
func (c Config) EthGasBumpThreshold() uint64 {
	if c.viper.IsSet(EnvVarName("EthGasBumpThreshold")) {
		return c.viper.GetUint64(EnvVarName("EthGasBumpThreshold"))
	}
	return chainSpecificConfig(c).EthGasBumpThreshold
}

// EthGasBumpTxDepth is the number of transactions to gas bump starting from oldest.
// Set to 0 for no limit (i.e. bump all)
func (c Config) EthGasBumpTxDepth() uint16 {
	return c.getWithFallback("EthGasBumpTxDepth", parseUint16).(uint16)
}

// EthGasBumpPercent is the minimum percentage by which gas is bumped on each transaction attempt
// Change with care since values below geth's default will fail with "underpriced replacement transaction"
func (c Config) EthGasBumpPercent() uint16 {
	return c.getWithFallback("EthGasBumpPercent", parseUint16).(uint16)
}

// EthGasBumpWei is the minimum fixed amount of wei by which gas is bumped on each transaction attempt
func (c Config) EthGasBumpWei() *big.Int {
	str := c.viper.GetString(EnvVarName("EthGasBumpWei"))
	if str != "" {
		n, err := parseBigInt(str)
		if err != nil {
			logger.Errorw(
				"Invalid value provided for EthGasBumpWei, falling back to default.",
				"value", str,
				"error", err)
		} else {
			return n.(*big.Int)
		}
	}
	n := chainSpecificConfig(c).EthGasBumpWei
	return &n
}

// EthMaxInFlightTransactions controls how many transactions are allowed to be
// "in-flight" i.e. broadcast but unconfirmed at any one time
// 0 value disables the limit
func (c Config) EthMaxInFlightTransactions() uint32 {
	return c.getWithFallback("EthMaxInFlightTransactions", parseUint32).(uint32)
}

// EthMaxGasPriceWei is the maximum amount in Wei that a transaction will be
// bumped to before abandoning it and marking it as errored.
func (c Config) EthMaxGasPriceWei() *big.Int {
	str := c.viper.GetString(EnvVarName("EthMaxGasPriceWei"))
	if str != "" {
		n, err := parseBigInt(str)
		if err != nil {
			logger.Errorw(
				"Invalid value provided for EthMaxGasPriceWei, falling back to default.",
				"value", str,
				"error", err)
		} else {
			return n.(*big.Int)
		}
	}
	n := chainSpecificConfig(c).EthMaxGasPriceWei
	return &n
}

// EthMaxQueuedTransactions is the maximum number of unbroadcast
// transactions per key that are allowed to be enqueued before jobs will start
// failing and rejecting send of any further transactions.
// 0 value disables
func (c Config) EthMaxQueuedTransactions() uint64 {
	return c.getWithFallback("EthMaxQueuedTransactions", parseUint64).(uint64)
}

// EthMinGasPriceWei is the minimum amount in Wei that a transaction may be priced.
// Chainlink will never send a transaction priced below this amount.
func (c Config) EthMinGasPriceWei() *big.Int {
	str := c.viper.GetString(EnvVarName("EthMinGasPriceWei"))
	if str != "" {
		n, err := parseBigInt(str)
		if err != nil {
			logger.Errorw(
				"Invalid value provided for EthMinGasPriceWei, falling back to default.",
				"value", str,
				"error", err)
		} else {
			return n.(*big.Int)
		}
	}
	n := chainSpecificConfig(c).EthMinGasPriceWei
	return &n
}

// EthNonceAutoSync enables/disables running the NonceSyncer on application start
func (c Config) EthNonceAutoSync() bool {
	return c.getWithFallback("EthNonceAutoSync", parseBool).(bool)
}

// EthGasLimitDefault sets the default gas limit for outgoing transactions.
func (c Config) EthGasLimitDefault() uint64 {
	ethGasLimitDefault := chainSpecificConfig(c).EthGasLimitDefault
	if c.viper.IsSet(EnvVarName("EthGasLimitDefault")) || ethGasLimitDefault == 0 {
		return c.getWithFallback("EthGasLimitDefault", parseUint64).(uint64)
	}
	return ethGasLimitDefault
}

// EthBaseTransactionGasLimit represents how much gas is required to send a eth->eth transfer
func (c Config) EthBaseTransactionGasLimit() uint64 {
	return uint64(21000)
}

// EthGasPriceDefault is the starting gas price for every transaction
func (c Config) EthGasPriceDefault() *big.Int {
	if c.runtimeStore != nil {
		var value big.Int
		if err := c.runtimeStore.GetConfigValue("EthGasPriceDefault", &value); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnw("Error while trying to fetch EthGasPriceDefault.", "error", err)
		} else if err == nil {
			return &value
		}
	}
	str := c.viper.GetString(EnvVarName("EthGasPriceDefault"))
	if str != "" {
		n, err := parseBigInt(str)
		if err != nil {
			logger.Errorw(
				"Invalid value provided for EthGasPriceDefault, falling back to default.",
				"value", str,
				"error", err)
		} else {
			return n.(*big.Int)
		}
	}
	n := chainSpecificConfig(c).EthGasPriceDefault
	return &n
}

// EthGasLimitMultiplier is a factory by which a transaction's GasLimit is
// multiplied before transmission. So if the value is 1.1, and the GasLimit for
// a transaction is 10, 10% will be added before transmission.
//
// This factor is always applied, so includes Optimism L2 transactions which
// uses a default gas limit of 1 and is also applied to EthGasLimitDefault.
func (c Config) EthGasLimitMultiplier() float32 {
	return (float32)(c.getWithFallback("EthGasLimitMultiplier", parseF32).(float64))
}

// SetEthGasPriceDefault saves a runtime value for the default gas price for transactions
func (c Config) SetEthGasPriceDefault(value *big.Int) error {
	min := c.EthMinGasPriceWei()
	max := c.EthMaxGasPriceWei()
	if value.Cmp(min) < 0 {
		return errors.Errorf("cannot set default gas price to %s, it is below the minimum allowed value of %s", value.String(), min.String())
	}
	if value.Cmp(max) > 0 {
		return errors.Errorf("cannot set default gas price to %s, it is above the maximum allowed value of %s", value.String(), max.String())
	}
	if c.runtimeStore == nil {
		return errors.New("No runtime store installed")
	}
	return c.runtimeStore.SetConfigValue("EthGasPriceDefault", value)
}

// EthFinalityDepth is the number of blocks after which an ethereum transaction is considered "final"
// BlocksConsideredFinal determines how deeply we look back to ensure that transactions are confirmed onto the longest chain
// There is not a large performance penalty to setting this relatively high (on the order of hundreds)
// It is practically limited by the number of heads we store in the database and should be less than this with a comfortable margin.
// If a transaction is mined in a block more than this many blocks ago, and is reorged out, we will NOT retransmit this transaction and undefined behaviour can occur including gaps in the nonce sequence that require manual intervention to fix.
// Therefore this number represents a number of blocks we consider large enough that no re-org this deep will ever feasibly happen.
//
// Special cases:
// ETH_FINALITY_DEPTH=0 would imply that transactions can be final even before they were mined into a block. This is not supported.
// ETH_FINALITY_DEPTH=1 implies that transactions are final after we see them in one block.
//
// Examples:
//
// Transaction sending:
// A transaction is sent at block height 42
//
// ETH_FINALITY_DEPTH is set to 5
// A re-org occurs at height 44 starting at block 41, transaction is marked for rebroadcast
// A re-org occurs at height 46 starting at block 41, transaction is marked for rebroadcast
// A re-org occurs at height 47 starting at block 41, transaction is NOT marked for rebroadcast
func (c Config) EthFinalityDepth() uint {
	if c.viper.IsSet(EnvVarName("EthFinalityDepth")) {
		return uint(c.viper.GetUint64(EnvVarName("EthFinalityDepth")))
	}
	return chainSpecificConfig(c).EthFinalityDepth
}

// EthHeadTrackerHistoryDepth tracks the top N block numbers to keep in the `heads` database table.
// Note that this can easily result in MORE than N records since in the case of re-orgs we keep multiple heads for a particular block height.
// This number should be at least as large as `EthFinalityDepth`.
// There may be a small performance penalty to setting this to something very large (10,000+)
func (c Config) EthHeadTrackerHistoryDepth() uint {
	if c.viper.IsSet(EnvVarName("EthHeadTrackerHistoryDepth")) {
		return uint(c.viper.GetUint64(EnvVarName("EthHeadTrackerHistoryDepth")))
	}
	return chainSpecificConfig(c).EthHeadTrackerHistoryDepth
}

// EthHeadTrackerMaxBufferSize is the maximum number of heads that may be
// buffered in front of the head tracker before older heads start to be
// dropped. You may think of it as something like the maximum permittable "lag"
// for the head tracker before we start dropping heads to keep up.
func (c Config) EthHeadTrackerMaxBufferSize() uint {
	return uint(c.getWithFallback("EthHeadTrackerMaxBufferSize", parseUint64).(uint64))
}

// EthHeadTrackerSamplingInterval is the interval between sampled head callbacks
// to services that are only interested in the latest head every some time
func (c Config) EthHeadTrackerSamplingInterval() time.Duration {
	str := c.viper.GetString(EnvVarName("EthHeadTrackerSamplingInterval"))
	if str != "" {
		n, err := parseDuration(str)
		if err != nil {
			logger.Errorw(
				"Invalid value provided for EthHeadTrackerSamplingInterval, falling back to default.",
				"value", str,
				"error", err)
		} else {
			return n.(time.Duration)
		}
	}
	return chainSpecificConfig(c).EthHeadTrackerSamplingInterval
}

// EthTxResendAfterThreshold controls how long the ethResender will wait before
// re-sending the latest eth_tx_attempt. This is designed a as a fallback to
// protect against the eth nodes dropping txes (it has been anecdotally
// observed to happen), networking issues or txes being ejected from the
// mempool.
// See eth_resender.go for more details
func (c Config) EthTxResendAfterThreshold() time.Duration {
	str := c.viper.GetString(EnvVarName("EthTxResendAfterThreshold"))
	if str != "" {
		n, err := parseDuration(str)
		if err != nil {
			logger.Errorw(
				"Invalid value provided for EthTxResendAfterThreshold, falling back to default.",
				"value", str,
				"error", err)
		} else {
			return n.(time.Duration)
		}
	}
	return chainSpecificConfig(c).EthTxResendAfterThreshold
}

// EthTxReaperInterval controls how often the eth tx reaper should run
func (c Config) EthTxReaperInterval() time.Duration {
	return c.getWithFallback("EthTxReaperInterval", parseDuration).(time.Duration)
}

// EthTxReaperThreshold represents how long any confirmed/fatally_errored eth_txes will hang around in the database.
// If the eth_tx is confirmed but still below ETH_FINALITY_DEPTH it will not be deleted even if it was created at a time older than this value.
// EXAMPLE
// With:
// EthTxReaperThreshold=1h
// EthFinalityDepth=50
//
// Current head is 142, any eth_tx confirmed in block 91 or below will be reaped as long as its created_at was more than EthTxReaperThreshold ago
// Set to 0 to disable eth_tx reaping
func (c Config) EthTxReaperThreshold() time.Duration {
	return c.getWithFallback("EthTxReaperThreshold", parseDuration).(time.Duration)
}

// EthLogBackfillBatchSize sets the batch size for calling FilterLogs when we backfill missing logs
func (c Config) EthLogBackfillBatchSize() uint32 {
	return c.getWithFallback("EthLogBackfillBatchSize", parseUint32).(uint32)
}

// EthereumURL represents the URL of the Ethereum node to connect Chainlink to.
func (c Config) EthereumURL() string {
	return c.viper.GetString(EnvVarName("EthereumURL"))
}

// EthereumHTTPURL is an optional but recommended url that points to the HTTP port of the primary node
func (c Config) EthereumHTTPURL() (uri *url.URL) {
	urlStr := c.viper.GetString(EnvVarName("EthereumHTTPURL"))
	if urlStr == "" {
		return nil
	}
	var err error
	uri, err = url.Parse(urlStr)
	if err != nil || !(uri.Scheme == "http" || uri.Scheme == "https") {
		logger.Fatalf("Invalid Ethereum HTTP URL: %s, got error: %s", urlStr, err)
	}
	return
}

// EthereumSecondaryURLs is an optional backup RPC URL
// Must be http(s) format
// If specified, transactions will also be broadcast to this ethereum node
func (c Config) EthereumSecondaryURLs() []url.URL {
	oldConfig := c.viper.GetString(EnvVarName("EthereumSecondaryURL"))
	newConfig := c.viper.GetString(EnvVarName("EthereumSecondaryURLs"))

	config := ""
	if newConfig != "" {
		config = newConfig
	} else if oldConfig != "" {
		config = oldConfig
	}

	urlStrings := regexp.MustCompile(`\s*[;,]\s*`).Split(config, -1)
	urls := []url.URL{}
	for _, urlString := range urlStrings {
		if urlString == "" {
			continue
		}
		url, err := url.Parse(urlString)
		if err != nil {
			logger.Fatalf("Invalid Secondary Ethereum URL: %s, got error: %v", urlString, err)
		}
		urls = append(urls, *url)
	}

	return urls
}

// EthereumDisabled shows whether Ethereum interactions are supported.
func (c Config) EthereumDisabled() bool {
	return c.viper.GetBool(EnvVarName("EthereumDisabled"))
}

// FlagsContractAddress represents the Flags contract address
func (c Config) FlagsContractAddress() string {
	return c.viper.GetString(EnvVarName("FlagsContractAddress"))
}

// GasUpdaterBatchSize sets the maximum number of blocks to fetch in one batch in the gas updater
// If the env var GAS_UPDATER_BATCH_SIZE is unset, it defaults to ETH_RPC_DEFAULT_BATCH_SIZE
func (c Config) GasUpdaterBatchSize() uint32 {
	if c.viper.IsSet(EnvVarName("GasUpdaterBatchSize")) {
		return c.viper.GetUint32(EnvVarName("GasUpdaterBatchSize"))
	}
	defaultGasUpdaterBatchSize := chainSpecificConfig(c).GasUpdaterBatchSize
	if defaultGasUpdaterBatchSize != nil {
		return *defaultGasUpdaterBatchSize
	}
	return c.EthRPCDefaultBatchSize()
}

// GasUpdaterBlockDelay is the number of blocks that the gas updater trails behind head.
// E.g. if this is set to 3, and we receive block 10, gas updater will
// fetch block 7.
// CAUTION: You might be tempted to set this to 0 to use the latest possible
// block, but it is possible to receive a head BEFORE that block is actually
// available from the connected node via RPC. In this case you will get false
// "zero" blocks that are missing transactions.
func (c Config) GasUpdaterBlockDelay() uint16 {
	if c.viper.IsSet(EnvVarName("GasUpdaterBlockDelay")) {
		return uint16(c.viper.GetUint32(EnvVarName("GasUpdaterBlockDelay")))
	}
	return chainSpecificConfig(c).GasUpdaterBlockDelay
}

// GasUpdaterBlockHistorySize is the number of past blocks to keep in memory to
// use as a basis for calculating a percentile gas price
func (c Config) GasUpdaterBlockHistorySize() uint16 {
	if c.viper.IsSet(EnvVarName("GasUpdaterBlockHistorySize")) {
		return uint16(c.viper.GetUint32(EnvVarName("GasUpdaterBlockHistorySize")))
	}
	return chainSpecificConfig(c).GasUpdaterBlockHistorySize
}

// GasUpdaterTransactionPercentile is the percentile gas price to choose. E.g.
// if the past transaction history contains four transactions with gas prices:
// [100, 200, 300, 400], picking 25 for this number will give a value of 200
func (c Config) GasUpdaterTransactionPercentile() uint16 {
	return c.getWithFallback("GasUpdaterTransactionPercentile", parseUint16).(uint16)
}

// GasUpdaterEnabled turns on the automatic gas updater if set to true
// It is enabled by default on most chains
func (c Config) GasUpdaterEnabled() bool {
	if c.viper.IsSet(EnvVarName("GasUpdaterEnabled")) {
		return c.viper.GetBool(EnvVarName("GasUpdaterEnabled"))
	}
	return chainSpecificConfig(c).GasUpdaterEnabled
}

// InsecureFastScrypt causes all key stores to encrypt using "fast" scrypt params instead
// This is insecure and only useful for local testing. DO NOT SET THIS IN PRODUCTION
func (c Config) InsecureFastScrypt() bool {
	return c.viper.GetBool(EnvVarName("InsecureFastScrypt"))
}

func (c Config) TriggerFallbackDBPollInterval() time.Duration {
	return c.getWithFallback("TriggerFallbackDBPollInterval", parseDuration).(time.Duration)
}

// JobPipelineMaxRunDuration is the maximum time that a job run may take
func (c Config) JobPipelineMaxRunDuration() time.Duration {
	return c.getWithFallback("JobPipelineMaxRunDuration", parseDuration).(time.Duration)
}

func (c Config) JobPipelineResultWriteQueueDepth() uint64 {
	return c.getWithFallback("JobPipelineResultWriteQueueDepth", parseUint64).(uint64)
}

func (c Config) JobPipelineReaperInterval() time.Duration {
	return c.getWithFallback("JobPipelineReaperInterval", parseDuration).(time.Duration)
}

func (c Config) JobPipelineReaperThreshold() time.Duration {
	return c.getWithFallback("JobPipelineReaperThreshold", parseDuration).(time.Duration)
}

func (c Config) KeeperRegistrySyncInterval() time.Duration {
	return c.getWithFallback("KeeperRegistrySyncInterval", parseDuration).(time.Duration)
}

func (c Config) KeeperMinimumRequiredConfirmations() uint64 {
	return c.viper.GetUint64(EnvVarName("KeeperMinimumRequiredConfirmations"))
}

func (c Config) KeeperMaximumGracePeriod() int64 {
	return c.viper.GetInt64(EnvVarName("KeeperMaximumGracePeriod"))
}

// JSONConsole when set to true causes logging to be made in JSON format
// If set to false, logs in console format
func (c Config) JSONConsole() bool {
	return c.viper.GetBool(EnvVarName("JSONConsole"))
}

// LinkContractAddress represents the address
func (c Config) LinkContractAddress() string {
	return c.viper.GetString(EnvVarName("LinkContractAddress"))
}

// ExplorerURL returns the websocket URL for this node to push stats to, or nil.
func (c Config) ExplorerURL() *url.URL {
	rval := c.getWithFallback("ExplorerURL", parseURL)
	switch t := rval.(type) {
	case nil:
		return nil
	case *url.URL:
		return t
	default:
		logger.Panicf("invariant: ExplorerURL returned as type %T", rval)
		return nil
	}
}

// ExplorerAccessKey returns the access key for authenticating with explorer
func (c Config) ExplorerAccessKey() string {
	return c.viper.GetString(EnvVarName("ExplorerAccessKey"))
}

// ExplorerSecret returns the secret for authenticating with explorer
func (c Config) ExplorerSecret() string {
	return c.viper.GetString(EnvVarName("ExplorerSecret"))
}

// FIXME: Add comments to all of these
func (c Config) OCRBootstrapCheckInterval() time.Duration {
	return c.getWithFallback("OCRBootstrapCheckInterval", parseDuration).(time.Duration)
}

func (c Config) OCRContractTransmitterTransmitTimeout() time.Duration {
	return c.getWithFallback("OCRContractTransmitterTransmitTimeout", parseDuration).(time.Duration)
}

func (c Config) getDurationWithOverride(override time.Duration, field string) time.Duration {
	if override != time.Duration(0) {
		return override
	}
	return c.getWithFallback(field, parseDuration).(time.Duration)
}

func (c Config) OCRObservationTimeout(override time.Duration) time.Duration {
	return c.getDurationWithOverride(override, "OCRObservationTimeout")
}

func (c Config) OCRObservationGracePeriod() time.Duration {
	return c.getWithFallback("OCRObservationGracePeriod", parseDuration).(time.Duration)
}

func (c Config) OCRBlockchainTimeout(override time.Duration) time.Duration {
	return c.getDurationWithOverride(override, "OCRBlockchainTimeout")
}

func (c Config) OCRContractSubscribeInterval(override time.Duration) time.Duration {
	return c.getDurationWithOverride(override, "OCRContractSubscribeInterval")
}

func (c Config) OCRContractPollInterval(override time.Duration) time.Duration {
	return c.getDurationWithOverride(override, "OCRContractPollInterval")
}

func (c Config) OCRContractConfirmations(override uint16) uint16 {
	if override != uint16(0) {
		return override
	}
	return c.getWithFallback("OCRContractConfirmations", parseUint16).(uint16)
}

func (c Config) OCRDatabaseTimeout() time.Duration {
	return c.getWithFallback("OCRDatabaseTimeout", parseDuration).(time.Duration)
}

func (c Config) OCRDHTLookupInterval() int {
	return int(c.getWithFallback("OCRDHTLookupInterval", parseUint16).(uint16))
}

func (c Config) OCRIncomingMessageBufferSize() int {
	return int(c.getWithFallback("OCRIncomingMessageBufferSize", parseUint16).(uint16))
}

func (c Config) OCRNewStreamTimeout() time.Duration {
	return c.getWithFallback("OCRNewStreamTimeout", parseDuration).(time.Duration)
}

func (c Config) OCROutgoingMessageBufferSize() int {
	return int(c.getWithFallback("OCRIncomingMessageBufferSize", parseUint16).(uint16))
}

// OCRTraceLogging determines whether OCR logs at TRACE level are enabled. The
// option to turn them off is given because they can be very verbose
func (c Config) OCRTraceLogging() bool {
	return c.viper.GetBool(EnvVarName("OCRTraceLogging"))
}

func (c Config) OCRMonitoringEndpoint(override string) string {
	if override != "" {
		return override
	}
	return c.viper.GetString(EnvVarName("OCRMonitoringEndpoint"))
}

func (c Config) OCRTransmitterAddress(override *models.EIP55Address) (models.EIP55Address, error) {
	if override != nil {
		return *override, nil
	}
	taStr := c.viper.GetString(EnvVarName("OCRTransmitterAddress"))
	if taStr != "" {
		ta, err := models.NewEIP55Address(taStr)
		if err != nil {
			return "", errors.Wrapf(ErrInvalid, "OCR_TRANSMITTER_ADDRESS is invalid EIP55 %v", err)
		}
		return ta, nil
	}
	return "", errors.Wrap(ErrUnset, "OCR_TRANSMITTER_ADDRESS")
}

func (c Config) OCRKeyBundleID(override *models.Sha256Hash) (models.Sha256Hash, error) {
	if override != nil {
		return *override, nil
	}
	kbStr := c.viper.GetString(EnvVarName("OCRKeyBundleID"))
	if kbStr != "" {
		kb, err := models.Sha256HashFromHex(kbStr)
		if err != nil {
			return models.Sha256Hash{}, errors.Wrapf(ErrInvalid, "OCR_KEY_BUNDLE_ID is an invalid sha256 hash hex string %v", err)
		}
		return kb, nil
	}
	return models.Sha256Hash{}, errors.Wrap(ErrUnset, "OCR_KEY_BUNDLE_ID")
}

func (c Config) ORMMaxOpenConns() int {
	return int(c.getWithFallback("ORMMaxOpenConns", parseUint16).(uint16))
}

func (c Config) ORMMaxIdleConns() int {
	return int(c.getWithFallback("ORMMaxIdleConns", parseUint16).(uint16))
}

// OperatorContractAddress represents the address where the Operator.sol
// contract is deployed, this is used for filtering RunLog requests
func (c Config) OperatorContractAddress() common.Address {
	if c.viper.GetString(EnvVarName("OperatorContractAddress")) == "" {
		return common.Address{}
	}
	address, ok := c.getWithFallback("OperatorContractAddress", parseAddress).(*common.Address)
	if !ok {
		return common.Address{}
	}
	return *address
}

// OptimismGasFees enables asking the network for gas price before submitting
// transactions, enabling compatibility with Optimism's L2 chain
func (c Config) OptimismGasFees() bool {
	if c.viper.IsSet(EnvVarName("OptimismGasFees")) {
		return c.viper.GetBool(EnvVarName("OptimismGasFees"))
	}
	return chainSpecificConfig(c).OptimismGasFees
}

// LogLevel represents the maximum level of log messages to output.
func (c Config) LogLevel() LogLevel {
	if c.runtimeStore != nil {
		var value LogLevel
		if err := c.runtimeStore.GetConfigValue("LogLevel", &value); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnw("Error while trying to fetch LogLevel.", "error", err)
		} else if err == nil {
			return value
		}
	}
	return c.getWithFallback("LogLevel", parseLogLevel).(LogLevel)
}

// SetLogLevel saves a runtime value for the default logger level
func (c Config) SetLogLevel(ctx context.Context, value string) error {
	if c.runtimeStore == nil {
		return errors.New("No runtime store installed")
	}
	var ll LogLevel
	err := ll.Set(value)
	if err != nil {
		return err
	}
	return c.runtimeStore.SetConfigStrValue(ctx, "LogLevel", ll.String())
}

// LogToDisk configures disk preservation of logs.
func (c Config) LogToDisk() bool {
	return c.viper.GetBool(EnvVarName("LogToDisk"))
}

// LogSQLStatements tells chainlink to log all SQL statements made using the default logger
func (c Config) LogSQLStatements() bool {
	if c.runtimeStore != nil {
		logSqlStatements, err := c.runtimeStore.GetConfigBoolValue("LogSQLStatements")
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnw("Error while trying to fetch LogSQLStatements.", "error", err)
		} else if err == nil {
			return *logSqlStatements
		}
	}
	return c.viper.GetBool(EnvVarName("LogSQLStatements"))
}

// SetLogSQLStatements saves a runtime value for enabling/disabling logging all SQL statements on the default logger
func (c Config) SetLogSQLStatements(ctx context.Context, sqlEnabled bool) error {
	if c.runtimeStore == nil {
		return errors.New("No runtime store installed")
	}

	return c.runtimeStore.SetConfigStrValue(ctx, "LogSQLStatements", strconv.FormatBool(sqlEnabled))
}

// LogSQLMigrations tells chainlink to log all SQL migrations made using the default logger
func (c Config) LogSQLMigrations() bool {
	return c.viper.GetBool(EnvVarName("LogSQLMigrations"))
}

// MinIncomingConfirmations represents the minimum number of block
// confirmations that need to be recorded since a job run started before a task
// can proceed.
// MIN_INCOMING_CONFIRMATIONS=1 would kick off a job after seeing the transaction in a block
// MIN_INCOMING_CONFIRMATIONS=0 would kick off a job even before the transaction is mined, which is not supported
func (c Config) MinIncomingConfirmations() uint32 {
	if c.viper.IsSet(EnvVarName("MinIncomingConfirmations")) {
		return c.viper.GetUint32(EnvVarName("MinIncomingConfirmations"))
	}
	return chainSpecificConfig(c).MinIncomingConfirmations
}

// MinRequiredOutgoingConfirmations represents the default minimum number of block
// confirmations that need to be recorded on an outgoing ethtx task before the run can move onto the next task.
// This can be overridden on a per-task basis by setting the `MinRequiredOutgoingConfirmations` parameter.
// MIN_OUTGOING_CONFIRMATIONS=1 considers a transaction as "done" once it has been mined into one block
// MIN_OUTGOING_CONFIRMATIONS=0 would consider a transaction as "done" even before it has been mined
func (c Config) MinRequiredOutgoingConfirmations() uint64 {
	if c.viper.IsSet(EnvVarName("MinRequiredOutgoingConfirmations")) {
		return c.viper.GetUint64(EnvVarName("MinRequiredOutgoingConfirmations"))
	}
	return chainSpecificConfig(c).MinRequiredOutgoingConfirmations
}

// MinimumContractPayment represents the minimum amount of LINK that must be
// supplied for a contract to be considered.
func (c Config) MinimumContractPayment() *assets.Link {
	minimumContractPayment := chainSpecificConfig(c).MinimumContractPayment
	if c.viper.IsSet(EnvVarName("MinimumContractPayment")) || minimumContractPayment == nil {
		return c.getWithFallback("MinimumContractPayment", parseLink).(*assets.Link)
	}
	return minimumContractPayment
}

// MinimumRequestExpiration is the minimum allowed request expiration for a Service Agreement.
func (c Config) MinimumRequestExpiration() uint64 {
	return c.getWithFallback("MinimumRequestExpiration", parseUint64).(uint64)
}

// P2PListenIP is the ip that libp2p willl bind to and listen on
func (c Config) P2PListenIP() net.IP {
	return c.getWithFallback("P2PListenIP", parseIP).(net.IP)
}

// P2PListenPort is the port that libp2p will bind to and listen on
func (c *Config) P2PListenPort() uint16 {
	if c.viper.IsSet(EnvVarName("P2PListenPort")) {
		return uint16(c.viper.GetUint32(EnvVarName("P2PListenPort")))
	}
	// Fast path in case it was already set
	c.randomP2PPortMtx.RLock()
	if c.randomP2PPort > 0 {
		c.randomP2PPortMtx.RUnlock()
		return c.randomP2PPort
	}
	c.randomP2PPortMtx.RUnlock()
	// Path for initial set
	c.randomP2PPortMtx.Lock()
	defer c.randomP2PPortMtx.Unlock()
	if c.randomP2PPort > 0 {
		return c.randomP2PPort
	}
	r, err := rand.Int(rand.Reader, big.NewInt(65535-1023))
	if err != nil {
		logger.Fatalw("unexpected error generating random port", "err", err)
	}
	randPort := uint16(r.Int64() + 1024)
	logger.Warnw(fmt.Sprintf("P2P_LISTEN_PORT was not set, listening on random port %d. A new random port will be generated on every boot, for stability it is recommended to set P2P_LISTEN_PORT to a fixed value in your environment", randPort), "p2pPort", randPort)
	c.randomP2PPort = randPort
	return c.randomP2PPort
}

// P2PListenPortRaw returns the raw string value of P2P_LISTEN_PORT
func (c Config) P2PListenPortRaw() string {
	return c.viper.GetString(EnvVarName("P2PListenPort"))
}

// P2PAnnounceIP is an optional override. If specified it will force the p2p
// layer to announce this IP as the externally reachable one to the DHT
// If this is set, P2PAnnouncePort MUST also be set.
func (c Config) P2PAnnounceIP() net.IP {
	str := c.viper.GetString(EnvVarName("P2PAnnounceIP"))
	return net.ParseIP(str)
}

// P2PAnnouncePort is an optional override. If specified it will force the p2p
// layer to announce this port as the externally reachable one to the DHT.
// If this is set, P2PAnnounceIP MUST also be set.
func (c Config) P2PAnnouncePort() uint16 {
	return uint16(c.viper.GetUint32(EnvVarName("P2PAnnouncePort")))
}

// P2PDHTAnnouncementCounterUserPrefix can be used to restore the node's
// ability to announce its IP/port on the P2P network after a database
// rollback. Make sure to only increase this value, and *never* decrease it.
// Don't use this variable unless you really know what you're doing, since you
// could semi-permanently exclude your node from the P2P network by
// misconfiguring it.
func (c Config) P2PDHTAnnouncementCounterUserPrefix() uint32 {
	return c.viper.GetUint32(EnvVarName("P2PDHTAnnouncementCounterUserPrefix"))
}

func (c Config) P2PPeerstoreWriteInterval() time.Duration {
	return c.getWithFallback("P2PPeerstoreWriteInterval", parseDuration).(time.Duration)
}

// P2PPeerID is the default peer ID that will be used, if not overridden
func (c Config) P2PPeerID(override *models.PeerID) (models.PeerID, error) {
	if override != nil {
		return *override, nil
	}
	pidStr := c.viper.GetString(EnvVarName("P2PPeerID"))
	if pidStr != "" {
		var pid models.PeerID
		err := pid.UnmarshalText([]byte(pidStr))
		if err != nil {
			return "", errors.Wrapf(ErrInvalid, "P2P_PEER_ID is invalid %v", err)
		}
		return pid, nil
	}
	return "", errors.Wrap(ErrUnset, "P2P_PEER_ID")
}

func (c Config) P2PPeerIDIsSet() bool {
	return c.viper.GetString(EnvVarName("P2PPeerID")) != ""
}

// P2PPeerIDRaw returns the string value of whatever P2P_PEER_ID was set to with no parsing
func (c Config) P2PPeerIDRaw() string {
	return c.viper.GetString(EnvVarName("P2PPeerID"))
}

func (c Config) P2PBootstrapPeers(override []string) ([]string, error) {
	if override != nil {
		return override, nil
	}
	if c.viper.IsSet(EnvVarName("P2PBootstrapPeers")) {
		bps := c.viper.GetStringSlice(EnvVarName("P2PBootstrapPeers"))
		if bps != nil {
			return bps, nil
		}
		return nil, errors.Wrap(ErrUnset, "P2P_BOOTSTRAP_PEERS")
	}
	return []string{}, nil
}

// Port represents the port Chainlink should listen on for client requests.
func (c Config) Port() uint16 {
	return c.getWithFallback("Port", parseUint16).(uint16)
}

func (c Config) HTTPServerWriteTimeout() time.Duration {
	return c.getWithFallback("HTTPServerWriteTimeout", parseDuration).(time.Duration)
}

// ReaperExpiration represents
func (c Config) ReaperExpiration() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("ReaperExpiration", parseDuration).(time.Duration))
}

func (c Config) ReplayFromBlock() int64 {
	return c.viper.GetInt64(EnvVarName("ReplayFromBlock"))
}

// RootDir represents the location on the file system where Chainlink should
// keep its files.
func (c Config) RootDir() string {
	return c.getWithFallback("RootDir", parseHomeDir).(string)
}

// SecureCookies allows toggling of the secure cookies HTTP flag
func (c Config) SecureCookies() bool {
	return c.viper.GetBool(EnvVarName("SecureCookies"))
}

// SessionTimeout is the maximum duration that a user session can persist without any activity.
func (c Config) SessionTimeout() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("SessionTimeout", parseDuration).(time.Duration))
}

// StatsPusherLogging toggles very verbose logging of raw messages for the StatsPusher (also telemetry)
func (c Config) StatsPusherLogging() bool {
	return c.getWithFallback("StatsPusherLogging", parseBool).(bool)
}

// TLSCertPath represents the file system location of the TLS certificate
// Chainlink should use for HTTPS.
func (c Config) TLSCertPath() string {
	return c.viper.GetString(EnvVarName("TLSCertPath"))
}

// TLSHost represents the hostname to use for TLS clients. This should match
// the TLS certificate.
func (c Config) TLSHost() string {
	return c.viper.GetString(EnvVarName("TLSHost"))
}

// TLSKeyPath represents the file system location of the TLS key Chainlink
// should use for HTTPS.
func (c Config) TLSKeyPath() string {
	return c.viper.GetString(EnvVarName("TLSKeyPath"))
}

// TLSPort represents the port Chainlink should listen on for encrypted client requests.
func (c Config) TLSPort() uint16 {
	return c.getWithFallback("TLSPort", parseUint16).(uint16)
}

// TLSRedirect forces TLS redirect for unencrypted connections
func (c Config) TLSRedirect() bool {
	return c.viper.GetBool(EnvVarName("TLSRedirect"))
}

// UnAuthenticatedRateLimit defines the threshold to which requests unauthenticated requests get limited
func (c Config) UnAuthenticatedRateLimit() int64 {
	return c.viper.GetInt64(EnvVarName("UnAuthenticatedRateLimit"))
}

// UnAuthenticatedRateLimitPeriod defines the period to which unauthenticated requests get limited
func (c Config) UnAuthenticatedRateLimitPeriod() models.Duration {
	return models.MustMakeDuration(c.getWithFallback("UnAuthenticatedRateLimitPeriod", parseDuration).(time.Duration))
}

func (c Config) tlsDir() string {
	return filepath.Join(c.RootDir(), "tls")
}

// KeyFile returns the path where the server key is kept
func (c Config) KeyFile() string {
	if c.TLSKeyPath() == "" {
		return filepath.Join(c.tlsDir(), "server.key")
	}
	return c.TLSKeyPath()
}

// CertFile returns the path where the server certificate is kept
func (c Config) CertFile() string {
	if c.TLSCertPath() == "" {
		return filepath.Join(c.tlsDir(), "server.crt")
	}
	return c.TLSCertPath()
}

// CreateProductionLogger returns a custom logger for the config's root
// directory and LogLevel, with pretty printing for stdout. If LOG_TO_DISK is
// false, the logger will only log to stdout.
func (c Config) CreateProductionLogger() *logger.Logger {
	return logger.CreateProductionLogger(c.RootDir(), c.JSONConsole(), c.LogLevel().Level, c.LogToDisk())
}

// SessionSecret returns a sequence of bytes to be used as a private key for
// session signing or encryption.
func (c Config) SessionSecret() ([]byte, error) {
	return c.SecretGenerator.Generate(c)
}

// SessionOptions returns the sesssions.Options struct used to configure
// the session store.
func (c Config) SessionOptions() sessions.Options {
	return sessions.Options{
		Secure:   c.SecureCookies(),
		HttpOnly: true,
		MaxAge:   86400 * 30,
	}
}

func (c Config) getWithFallback(name string, parser func(string) (interface{}, error)) interface{} {
	str := c.viper.GetString(EnvVarName(name))
	defaultValue, hasDefault := defaultValue(name)
	if str != "" {
		v, err := parser(str)
		if err == nil {
			return v
		}
		logger.Errorw(
			fmt.Sprintf("Invalid value provided for %s, falling back to default.", name),
			"value", str,
			"default", defaultValue,
			"error", err)
	}

	if !hasDefault {
		return zeroValue(name)
	}

	v, err := parser(defaultValue)
	if err != nil {
		log.Fatalf(fmt.Sprintf(`Invalid default for %s: "%s"`, name, defaultValue))
	}
	return v
}

// SecretGenerator is the interface for objects that generate a secret
// used to sign or encrypt.
type SecretGenerator interface {
	Generate(Config) ([]byte, error)
}

type filePersistedSecretGenerator struct{}

func (f filePersistedSecretGenerator) Generate(c Config) ([]byte, error) {
	sessionPath := filepath.Join(c.RootDir(), "secret")
	if utils.FileExists(sessionPath) {
		data, err := ioutil.ReadFile(sessionPath)
		if err != nil {
			return data, err
		}
		return base64.StdEncoding.DecodeString(string(data))
	}
	key := securecookie.GenerateRandomKey(32)
	str := base64.StdEncoding.EncodeToString(key)
	err := utils.WriteFileWithMaxPerms(sessionPath, []byte(str), readWritePerms)
	return key, err
}

func parseAddress(str string) (interface{}, error) {
	if str == "" {
		return nil, nil
	} else if common.IsHexAddress(str) {
		val := common.HexToAddress(str)
		return &val, nil
	} else if i, ok := new(big.Int).SetString(str, 10); ok {
		val := common.BigToAddress(i)
		return &val, nil
	}
	return nil, fmt.Errorf("unable to parse '%s' into EIP55-compliant address", str)
}

func parseLink(str string) (interface{}, error) {
	i, ok := new(assets.Link).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("unable to parse '%v' into *assets.Link(base 10)", str)
	}
	return i, nil
}

func parseLogLevel(str string) (interface{}, error) {
	var lvl LogLevel
	err := lvl.Set(str)
	return lvl, err
}

func parseUint16(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 16)
	return uint16(v), err
}

func parseUint32(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	return uint32(v), err
}

func parseUint64(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	return v, err
}

func parseF32(s string) (interface{}, error) {
	v, err := strconv.ParseFloat(s, 32)
	return v, err
}

func parseURL(s string) (interface{}, error) {
	return url.Parse(s)
}

func parseIP(s string) (interface{}, error) {
	return net.ParseIP(s), nil
}

func parseDuration(s string) (interface{}, error) {
	return time.ParseDuration(s)
}

func parseBool(s string) (interface{}, error) {
	return strconv.ParseBool(s)
}

func parseBigInt(str string) (interface{}, error) {
	i, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("unable to parse %v into *big.Int(base 10)", str)
	}
	return i, nil
}

func parseHomeDir(str string) (interface{}, error) {
	exp, err := homedir.Expand(str)
	if err != nil {
		return nil, err
	}
	return filepath.ToSlash(exp), nil
}

// LogLevel determines the verbosity of the events to be logged.
type LogLevel struct {
	zapcore.Level
}

type DatabaseBackupMode string

var (
	DatabaseBackupModeNone DatabaseBackupMode = "none"
	DatabaseBackupModeLite DatabaseBackupMode = "lite"
	DatabaseBackupModeFull DatabaseBackupMode = "full"
)

func parseDatabaseBackupMode(s string) (interface{}, error) {
	switch DatabaseBackupMode(s) {
	case DatabaseBackupModeNone, DatabaseBackupModeLite, DatabaseBackupModeFull:
		return DatabaseBackupMode(s), nil
	default:
		return "", fmt.Errorf("unable to parse %v into DatabaseBackupMode. Must be one of values: \"%s\", \"%s\", \"%s\"", s, DatabaseBackupModeNone, DatabaseBackupModeLite, DatabaseBackupModeFull)
	}
}
