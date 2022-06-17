package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap/zapcore"
)

// Static audit log event type constants
const (
	AUTH_LOGIN_FAILED_EMAIL     = "AUTH_LOGIN_FAILED_EMAIL"
	AUTH_LOGIN_FAILED_PASSWORD  = "AUTH_LOGIN_FAILED_PASSWORD"
	AUTH_LOGIN_FAILED_2FA       = "AUTH_LOGIN_FAILED_2FA"
	AUTH_LOGIN_SUCCESS_WITH_2FA = "AUTH_LOGIN_SUCCESS_WITH_2FA"
	AUTH_LOGIN_SUCCESS_NO_2FA   = "AUTH_LOGIN_SUCCESS_NO_2FA"
	AUTH_2FA_ENROLLED           = "AUTH_2FA_ENROLLED"
	AUTH_SESSION_DELETED        = "SESSION_DELETED"

	PASSWORD_RESET_ATTEMPT_FAILED_MISMATCH = "PASSWORD_RESET_ATTEMPT_FAILED_MISMATCH"
	PASSWORD_RESET_SUCCESS                 = "PASSWORD_RESET_SUCCESS"

	API_TOKEN_CREATE_ATTEMPT_PASSWORD_MISMATCH = "API_TOKEN_CREATE_ATTEMPT_PASSWORD_MISMATCH"
	API_TOKEN_CREATED                          = "API_TOKEN_CREATED"
	API_TOKEN_DELETE_ATTEMPT_PASSWORD_MISMATCH = "API_TOKEN_DELETE_ATTEMPT_PASSWORD_MISMATCH"
	API_TOKEN_DELETED                          = "API_TOKEN_DELETED"

	CSA_KEY_CREATED  = "CSA_KEY_CREATED"
	CSA_KEY_IMPORTED = "CSA_KEY_IMPORTED"
	CSA_KEY_EXPORTED = "CSA_KEY_EXPORTED"
	CSA_KEY_DELETED  = "CSA_KEY_DELETED"

	FEEDS_MAN_CREATED = "FEEDS_MAN_CREATED"
	FEEDS_MAN_UPDATED = "FEEDS_MAN_UPDATED"

	FEEDS_MAN_CHAIN_CONFIG_CREATED = "FEEDS_MAN_CHAIN_CONFIG_CREATED"
	FEEDS_MAN_CHAIN_CONFIG_UPDATED = "FEEDS_MAN_CHAIN_CONFIG_UPDATED"
	FEEDS_MAN_CHAIN_CONFIG_DELETED = "FEEDS_MAN_CHAIN_CONFIG_DELETED"

	OCR_KEY_BUNDLE_CREATED  = "OCR_KEY_BUNDLE_CREATED"
	OCR_KEY_BUNDLE_IMPORTED = "OCR_KEY_BUNDLE_IMPORTED"
	OCR_KEY_BUNDLE_EXPORTED = "OCR_KEY_BUNDLE_EXPORTED"
	OCR_KEY_BUNDLE_DELETED  = "OCR_KEY_BUNDLE_DELETED"

	OCR2_KEY_BUNDLE_CREATED  = "OCR2_KEY_BUNDLE_CREATED"
	OCR2_KEY_BUNDLE_IMPORTED = "OCR2_KEY_BUNDLE_IMPORTED"
	OCR2_KEY_BUNDLE_EXPORTED = "OCR2_KEY_BUNDLE_EXPORTED"
	OCR2_KEY_BUNDLE_DELETED  = "OCR2_KEY_BUNDLE_DELETED"

	ETH_KEY_CREATED  = "ETH_KEY_CREATED"
	ETH_KEY_UPDATED  = "ETH_KEY_UPDATED"
	ETH_KEY_IMPORTED = "ETH_KEY_IMPORTED"
	ETH_KEY_EXPORTED = "ETH_KEY_EXPORTED"
	ETH_KEY_DELETED  = "ETH_KEY_DELETED"

	P2P_KEY_CREATED  = "P2P_KEY_CREATED"
	P2P_KEY_IMPORTED = "P2P_KEY_IMPORTED"
	P2P_KEY_EXPORTED = "P2P_KEY_EXPORTED"
	P2P_KEY_DELETED  = "P2P_KEY_DELETED"

	VRF_KEY_CREATED  = "VRF_KEY_CREATED"
	VRF_KEY_IMPORTED = "VRF_KEY_IMPORTED"
	VRF_KEY_EXPORTED = "VRF_KEY_EXPORTED"
	VRF_KEY_DELETED  = "VRF_KEY_DELETED"

	TERRA_KEY_CREATED  = "TERRA_KEY_CREATED"
	TERRA_KEY_IMPORTED = "TERRA_KEY_IMPORTED"
	TERRA_KEY_EXPORTED = "TERRA_KEY_EXPORTED"
	TERRA_KEY_DELETED  = "TERRA_KEY_DELETED"

	SOLANA_KEY_CREATED  = "SOLANA_KEY_CREATED"
	SOLANA_KEY_IMPORTED = "SOLANA_KEY_IMPORTED"
	SOLANA_KEY_EXPORTED = "SOLANA_KEY_EXPORTED"
	SOLANA_KEY_DELETED  = "SOLANA_KEY_DELETED"

	ETH_TRANSACTION_CREATED    = "ETH_TRANSACTION_CREATED"
	TERRA_TRANSACTION_CREATED  = "TERRA_TRANSACTION_CREATED"
	SOLANA_TRANSACTION_CREATED = "SOLANA_TRANSACTION_CREATED"

	JOB_CREATED = "JOB_CREATED"
	JOB_DELETED = "JOB_DELETED"

	CHAIN_ADDED        = "CHAIN_ADDED"
	CHAIN_SPEC_UPDATED = "CHAIN_SPEC_UPDATED"
	CHAIN_DELETED      = "CHAIN_DELETED"

	CHAIN_RPC_NODE_ADDED   = "CHAIN_RPC_NODE_ADDED"
	CHAIN_RPC_NODE_DELETED = "CHAIN_RPC_NODE_DELETED"

	BRIDGE_CREATED = "BRIDGE_CREATED"
	BRIDGE_UPDATED = "BRIDGE_UPDATED"
	BRIDGE_DELETED = "BRIDGE_DELETED"

	FORWARDER_CREATED = "FORWARDER_CREATED"
	FORWARDER_DELETED = "FORWARDER_DELETED"

	EXTERNAL_INITIATOR_CREATED = "EXTERNAL_INITIATOR_CREATED"
	EXTERNAL_INITIATOR_DELETED = "EXTERNAL_INITIATOR_DELETED"

	JOB_PROPOSAL_SPEC_APPROVED = "JOB_PROPOSAL_SPEC_APPROVED"
	JOB_PROPOSAL_SPEC_UPDATED  = "JOB_PROPOSAL_SPEC_UPDATED"
	JOB_PROPOSAL_SPEC_CANCELED = "JOB_PROPOSAL_SPEC_CANCELED"
	JOB_PROPOSAL_SPEC_REJECTED = "JOB_PROPOSAL_SPEC_REJECTED"

	CONFIG_UPDATED              = "CONFIG_UPDATED"
	CONFIG_SQL_LOGGING_ENABLED  = "CONFIG_SQL_LOGGING_ENABLED"
	CONFIG_SQL_LOGGING_DISABLED = "CONFIG_SQL_LOGGING_DISABLED"
	GLOBAL_LOG_LEVEL_SET        = "GLOBAL_LOG_LEVEL_SET"

	JOB_ERROR_DISMISSED = "JOB_ERROR_DISMISSED"
	JOB_RUN_SET         = "JOB_RUN_SET"

	ENV_NONCRITICAL_ENV_DUMPED = "ENV_NONCRITICAL_ENV_DUMPED"

	UNAUTHED_RUN_RESUMED = "UNAUTHED_RUN_RESUMED"
)

type splunkLogger struct {
	logger          Logger
	splunkToken     string
	splunkURL       string
	environmentName string
	hostname        string
	localIP         string
}

func newSplunkLogger(logger Logger, splunkToken string, splunkURL string, hostname string, environment string) Logger {
	// Initialize and return Splunk logger struct with required state for HEC calls
	return &splunkLogger{
		logger:          logger.Helper(1),
		splunkToken:     splunkToken,
		splunkURL:       splunkURL,
		environmentName: environment,
		hostname:        hostname,
		localIP:         getLocalIP(),
	}
}

func (l *splunkLogger) Audit(eventID string, data map[string]interface{}) {
	// goroutine to async POST to splunk HTTP Event Collector (HEC)
	go l.postLogToSplunk(eventID, data)
	l.logger.Audit(eventID, data)
}

// getLocalIP returns the first non- loopback local IP of the host
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func (l *splunkLogger) postLogToSplunk(eventID string, data map[string]interface{}) {
	// Splunk JSON data
	splunkLog := map[string]interface{}{
		"eventID":  eventID,
		"hostname": l.hostname,
		"localIP":  l.localIP,
		"env":      l.environmentName,
	}
	if len(data) != 0 {
		splunkLog["data"] = data
	}

	// Wrap serialized audit log map into JSON object `event` for API call
	serializedArgs, _ := json.Marshal(splunkLog)
	splunkLog = map[string]interface{}{"event": string(serializedArgs)}
	serializedSplunkLog, _ := json.Marshal(splunkLog)

	// Send up to HEC log collector
	httpClient := &http.Client{Timeout: time.Second * 60}
	req, _ := http.NewRequest("POST", l.splunkURL, bytes.NewReader(serializedSplunkLog))
	req.Header.Add("Authorization", "Splunk "+l.splunkToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		l.logger.Errorw("Failed to send audit log to Splunk", "err", err, "splunkLog", splunkLog)
	}
	if resp.StatusCode != 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			l.logger.Errorw("Error reading errored Splunk webhook response body", "err", err, "splunkLog", splunkLog)
		}
		l.logger.Errorw("Error sending log to Splunk", "statusCode", resp.StatusCode, "bodyString", string(bodyBytes))
	}
}

func (l *splunkLogger) With(args ...interface{}) Logger {
	return &splunkLogger{
		logger:          l.logger.With(args...),
		splunkToken:     l.splunkToken,
		splunkURL:       l.splunkURL,
		environmentName: l.environmentName,
		hostname:        l.hostname,
		localIP:         getLocalIP(),
	}
}

func (l *splunkLogger) Named(name string) Logger {
	return &splunkLogger{
		logger:          l.logger.Named(name),
		splunkToken:     l.splunkToken,
		splunkURL:       l.splunkURL,
		environmentName: l.environmentName,
		hostname:        l.hostname,
		localIP:         getLocalIP(),
	}
}

func (l *splunkLogger) SetLogLevel(level zapcore.Level) {
	l.logger.SetLogLevel(level)
}

func (l *splunkLogger) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *splunkLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *splunkLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *splunkLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *splunkLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *splunkLogger) Critical(args ...interface{}) {
	l.logger.Critical(args...)
}

func (l *splunkLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *splunkLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *splunkLogger) Tracef(format string, values ...interface{}) {
	l.logger.Tracef(format, values...)
}

func (l *splunkLogger) Debugf(format string, values ...interface{}) {
	l.logger.Debugf(format, values...)
}

func (l *splunkLogger) Infof(format string, values ...interface{}) {
	l.logger.Infof(format, values...)
}

func (l *splunkLogger) Warnf(format string, values ...interface{}) {
	l.logger.Warnf(format, values...)
}

func (l *splunkLogger) Errorf(format string, values ...interface{}) {
	l.logger.Errorf(format, values...)
}

func (l *splunkLogger) Criticalf(format string, values ...interface{}) {
	l.logger.Criticalf(format, values...)
}

func (l *splunkLogger) Panicf(format string, values ...interface{}) {
	l.logger.Panicf(format, values...)
}

func (l *splunkLogger) Fatalf(format string, values ...interface{}) {
	l.logger.Fatalf(format, values...)
}

func (l *splunkLogger) Tracew(msg string, keysAndValues ...interface{}) {
	l.logger.Tracew(msg, keysAndValues...)
}

func (l *splunkLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *splunkLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *splunkLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *splunkLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *splunkLogger) Criticalw(msg string, keysAndValues ...interface{}) {
	l.logger.Criticalw(msg, keysAndValues...)
}

func (l *splunkLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.logger.Panicw(msg, keysAndValues...)
}

func (l *splunkLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

func (l *splunkLogger) ErrorIf(err error, msg string) {
	if err != nil {
		l.logger.Errorw(msg, "err", err)
	}
}

func (l *splunkLogger) ErrorIfClosing(c io.Closer, name string) {
	if err := c.Close(); err != nil {
		l.logger.Errorw(fmt.Sprintf("Error closing %s", name), "err", err)
	}
}

func (l *splunkLogger) Sync() error {
	return l.logger.Sync()
}

func (l *splunkLogger) Helper(add int) Logger {
	return &splunkLogger{
		logger:          l.logger.Helper(add),
		splunkToken:     l.splunkToken,
		splunkURL:       l.splunkURL,
		environmentName: l.environmentName,
		hostname:        l.hostname,
		localIP:         getLocalIP(),
	}
}

func (l *splunkLogger) Recover(panicErr interface{}) {
	l.logger.Recover(panicErr)
}
