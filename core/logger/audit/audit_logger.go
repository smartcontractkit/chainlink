package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/pkg/errors"
)

const bufferCapacity = 2048

type Data = map[string]any

type AuditLogger interface {
	Audit(ctx context.Context, eventID EventID, data Data)
}

type AuditLoggerConfig struct {
	forwardToUrl   string
	environment    string
	jsonWrapperKey string
	headers        []serviceHeader
}

func NewAuditLoggerConfig(forwardToUrl string, isDev bool, jsonWrapperKey string, encodedHeaders string) (AuditLoggerConfig, error) {
	environment := "production"
	if isDev {
		environment = "develop"
	}

	if forwardToUrl == "" {
		return AuditLoggerConfig{}, errors.Errorf("No forwardToURL provided")
	}

	// Split and prepare optional service client headers from env variable
	headers := []serviceHeader{}
	if encodedHeaders != "" {
		headerLines := strings.Split(encodedHeaders, "\\")
		for _, header := range headerLines {
			keyValue := strings.Split(header, "||")
			if len(keyValue) != 2 {
				return AuditLoggerConfig{}, errors.Errorf("Invalid headers provided for the audit logger. Value, single pair split on || required, got: %s", keyValue)
			}
			headers = append(headers, serviceHeader{
				header: keyValue[0],
				value:  keyValue[1],
			})
		}
	}

	return AuditLoggerConfig{
		forwardToUrl:   forwardToUrl,
		environment:    environment,
		jsonWrapperKey: jsonWrapperKey,
		headers:        headers,
	}, nil
}

type AuditLoggerService struct {
	logger          logger.Logger
	enabled         bool
	serviceURL      string
	serviceHeaders  []serviceHeader
	jsonWrapperKey  string
	environmentName string
	hostname        string
	localIP         string
	loggingChannel  chan WrappedAuditLog
}

// Configurable headers to include in POST to log service
type serviceHeader struct {
	header string
	value  string
}

type WrappedAuditLog struct {
	eventID EventID
	data    Data
}

// NewAuditLogger returns a buffer push system that ingests audit log events and
// asynchronously pushes them up to an HTTP log service.
// Parses and validates the AUDIT_LOGS_* environment values and returns an enabled
// AuditLogger instance. If the environment variables are not set, the logger
// is disabled and short circuits execution via enabled flag.
func NewAuditLogger(logger logger.Logger, config *AuditLoggerConfig) (AuditLogger, error) {
	if config == nil {
		logger.Info("Audit logger configuration is nil. Cannot start audit logger subsystem and audit events will not be captured.")
		return &AuditLoggerService{}, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		return &AuditLoggerService{}, errors.Errorf("Audit Log initialization error - unable to get hostname: %s", err)
	}

	loggingChannel := make(chan WrappedAuditLog, bufferCapacity)

	// Finally, create new auditLogger with parameters
	auditLogger := AuditLoggerService{
		logger:          logger.Helper(1),
		enabled:         true,
		serviceURL:      config.forwardToUrl,
		serviceHeaders:  config.headers,
		jsonWrapperKey:  config.jsonWrapperKey,
		environmentName: config.environment,
		hostname:        hostname,
		localIP:         getLocalIP(),
		loggingChannel:  loggingChannel,
	}

	// Start our go routine that will receive logs and send them out to the
	// configured service.
	go auditLogger.auditLoggerRoutine()

	return &auditLogger, nil
}

// / Entrypoint for new audit logs. This buffers all logs that come in they will
// / sent out by the goroutine that was started when the AuditLoggerService was
// / created. If this service was not enabled, this immeidately returns.
// /
// / This function never blocks.
func (l *AuditLoggerService) Audit(ctx context.Context, eventID EventID, data Data) {
	if !l.enabled {
		return
	}

	wrappedLog := WrappedAuditLog{
		eventID: eventID,
		data:    data,
	}

	select {
	case l.loggingChannel <- wrappedLog:
	default:
		if len(l.loggingChannel) == bufferCapacity {
			l.logger.Errorw("Audit log buffer is full. Dropping log with eventID: %s", eventID)
		} else if l.loggingChannel == nil {
			l.logger.Errorw("Could not send log to audit subsystem because it has gone away!")
		} else {
			l.logger.Errorw("An unknown error has occured in the audit logging subsystem and the audit log was dropped")
		}
	}
}

// / Entrypoint for our log handling goroutine. This waits on the channel and sends out
// / logs as they come in.
// /
// / This function calls postLogToLogService which blocks.
func (l *AuditLoggerService) auditLoggerRoutine() {
	for event := range l.loggingChannel {
		l.postLogToLogService(event.eventID, event.data)
	}

	l.logger.Errorw("Audit logger is shut down. Will not send requested audit log")
}

// / Takes an EventID and associated data and sends it to the configured logging
// / endpoint. This function blocks on the send by timesout after a period of
// / several seconds. This helps us prevent getting stuck on a single log
// / due to transient network errors.
// /
// / This function blocks when called.
func (l *AuditLoggerService) postLogToLogService(eventID EventID, data Data) {
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
