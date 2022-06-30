package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"go.uber.org/zap/zapcore"
)

type AuditLogger struct {
	logger logger.Logger
	AuditLoggerConfig
}

type AuditLoggerConfig struct {
	serviceURL      string
	serviceHeaders  []serviceHeader
	jsonWrapperKey  string
	environmentName string
	hostname        string
	localIP         string
}

// Configurable headers to include in POST to log service
type serviceHeader struct {
	header string
	value  string
}

// NewAuditLoggerConfig parses and validatesthe passed AUDIT_LOGS_* environment values and populates fields
func NewAuditLoggerConfig(serviceURL string, headersEncoded, jsonWrapperKey, hostname, environment string) (AuditLoggerConfig, error) {
	newConfig := AuditLoggerConfig{}

	// Split and prepare optional service client headers from env variable
	headers := []serviceHeader{}
	if headersEncoded != "" {
		headerLines := strings.Split(headersEncoded, "\\")
		for _, header := range headerLines {
			keyValue := strings.Split(header, "||")
			if len(keyValue) != 2 {
				return AuditLoggerConfig{}, errors.Errorf("Invalid AUDIT_LOGS_FORWARDER_HEADERS value, single pair split on || required, got: %s", keyValue)
			}
			headers = append(headers, serviceHeader{
				header: keyValue[0],
				value:  keyValue[1],
			})
		}
	}

	newConfig.serviceURL = serviceURL
	newConfig.serviceHeaders = headers
	newConfig.jsonWrapperKey = jsonWrapperKey
	newConfig.environmentName = environment
	newConfig.hostname = hostname
	newConfig.localIP = getLocalIP()

	return newConfig, nil
}

// The .Audit function in the logger interface is exclusively used by this AuditLogger implementation.
// .Auditf function implementations should continue the pattern. Audit logs here must be emitted
// regardless of log level, hence the separate 'Audit' log level
// Audit log events are posted up to with an HTTP forwarder
func newAuditLogger(auditLoggerCfg AuditLoggerConfig, logger Logger) AuditLogger {
	sLogger := AuditLogger{
		logger: logger.Helper(1),
	}
	sLogger.AuditLoggerConfig = auditLoggerCfg
	return &sLogger
}

func (l *auditLogger) Audit(eventID audit.EventID, data map[string]interface{}) {
	l.postLogToLogService(eventID, data)
	l.logger.Audit(eventID, data)
}

// getLocalIP returns the first non-loopback local IP of the host
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// filter and return address types for first non loopback address
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func (l *auditLogger) postLogToLogService(eventID audit.EventID, data map[string]interface{}) {
	// Audit log JSON data
	logItem := map[string]interface{}{
		"eventID":  eventID,
		"hostname": l.hostname,
		"localIP":  l.localIP,
		"env":      l.environmentName,
		"data":     data,
	}

	// Optionally wrap audit log data into JSON object to help dynamically structure for an HTTP log service call
	if l.jsonWrapperKey != "" {
		logItem = map[string]interface{}{l.jsonWrapperKey: logItem}
	}

	serializedLog, err := json.Marshal(logItem)
	if err != nil {
		l.logger.Errorw("Unable to serialize wrapped audit log item to JSON", "err", err, "logItem", logItem)
		return
	}

	// Send up to HEC log collector
	httpClient := &http.Client{Timeout: time.Second * 10}
	req, _ := http.NewRequest("POST", l.serviceURL, bytes.NewReader(serializedLog))
	for _, header := range l.serviceHeaders {
		req.Header.Add(header.header, header.value)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		l.logger.Errorw("Failed to send audit log to HTTP log service", "err", err, "logItem", logItem)
		return
	}
	if resp.StatusCode != 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			l.logger.Errorw("Error reading errored HTTP log service webhook response body", "err", err, "logItem", logItem)
			return
		}
		l.logger.Errorw("Error sending log to HTTP log service", "statusCode", resp.StatusCode, "bodyString", string(bodyBytes))
		return
	}
}

func (l *auditLogger) With(args ...interface{}) Logger {
	return &auditLogger{
		logger:            l.logger.With(args...),
		AuditLoggerConfig: l.AuditLoggerConfig,
	}
}

func (l *auditLogger) Named(name string) Logger {
	return &auditLogger{
		logger:            l.logger.Named(name),
		AuditLoggerConfig: l.AuditLoggerConfig,
	}
}

func (l *auditLogger) SetLogLevel(level zapcore.Level) {
	l.logger.SetLogLevel(level)
}

func (l *auditLogger) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *auditLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *auditLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *auditLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *auditLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *auditLogger) Critical(args ...interface{}) {
	l.logger.Critical(args...)
}

func (l *auditLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *auditLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *auditLogger) Tracef(format string, values ...interface{}) {
	l.logger.Tracef(format, values...)
}

func (l *auditLogger) Debugf(format string, values ...interface{}) {
	l.logger.Debugf(format, values...)
}

func (l *auditLogger) Infof(format string, values ...interface{}) {
	l.logger.Infof(format, values...)
}

func (l *auditLogger) Warnf(format string, values ...interface{}) {
	l.logger.Warnf(format, values...)
}

func (l *auditLogger) Errorf(format string, values ...interface{}) {
	l.logger.Errorf(format, values...)
}

func (l *auditLogger) Criticalf(format string, values ...interface{}) {
	l.logger.Criticalf(format, values...)
}

func (l *auditLogger) Panicf(format string, values ...interface{}) {
	l.logger.Panicf(format, values...)
}

func (l *auditLogger) Fatalf(format string, values ...interface{}) {
	l.logger.Fatalf(format, values...)
}

func (l *auditLogger) Tracew(msg string, keysAndValues ...interface{}) {
	l.logger.Tracew(msg, keysAndValues...)
}

func (l *auditLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *auditLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *auditLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *auditLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *auditLogger) Criticalw(msg string, keysAndValues ...interface{}) {
	l.logger.Criticalw(msg, keysAndValues...)
}

func (l *auditLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.logger.Panicw(msg, keysAndValues...)
}

func (l *auditLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

func (l *auditLogger) ErrorIf(err error, msg string) {
	if err != nil {
		l.logger.Errorw(msg, "err", err)
	}
}

func (l *auditLogger) ErrorIfClosing(c io.Closer, name string) {
	if err := c.Close(); err != nil {
		l.logger.Errorw(fmt.Sprintf("Error closing %s", name), "err", err)
	}
}

func (l *auditLogger) Sync() error {
	return l.logger.Sync()
}

func (l *auditLogger) Helper(add int) Logger {
	return &auditLogger{
		logger:            l.logger.Helper(add),
		AuditLoggerConfig: l.AuditLoggerConfig,
	}
}

func (l *auditLogger) Recover(panicErr interface{}) {
	l.logger.Recover(panicErr)
}
