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

const AUDIT_LOGS_CAPACITY = 2048

type AuditLogger interface {
	//Audit(ctx context.Context, eventID EventID, data map[string]interface{})
	Audit(ctx context.Context, eventID EventID, data map[string]interface{})
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
	data    map[string]interface{}
}

// NewAuditLogger returns a buffer push system that ingests audit log events and
// asynchronously pushes them up to an HTTP log service.
// Parses and validates the AUDIT_LOGS_* environment values and returns an enabled
// AuditLogger instance. If the environment variables are not set, the logger
// is disabled and short circuits execution via enabled flag.
func NewAuditLogger(logger logger.Logger) (AuditLogger, error) {
	// Start parsing environment variables for audit logger
	auditLogsURL := os.Getenv("AUDIT_LOGS_FORWARDER_URL")
	if auditLogsURL == "" {
		// Unset, return a disabled audit logger
		logger.Info("No AUDIT_LOGS_FORWARDER_URL environment set, audit log events will not be captured")

		return &AuditLoggerService{}, nil
	}

	env := "production"
	if os.Getenv("CHAINLINK_DEV") == "true" {
		env = "develop"
	}
	hostname, err := os.Hostname()
	if err != nil {
		return &AuditLoggerService{}, errors.Errorf("Audit Log initialization error - unable to get hostname: %s", err)
	}

	// Split and prepare optional service client headers from env variable
	headers := []serviceHeader{}
	headersEncoded := os.Getenv("AUDIT_LOGS_FORWARDER_HEADERS")
	if headersEncoded != "" {
		headerLines := strings.Split(headersEncoded, "\\")
		for _, header := range headerLines {
			keyValue := strings.Split(header, "||")
			if len(keyValue) != 2 {
				return &AuditLoggerService{}, errors.Errorf("Invalid AUDIT_LOGS_FORWARDER_HEADERS value, single pair split on || required, got: %s", keyValue)
			}
			headers = append(headers, serviceHeader{
				header: keyValue[0],
				value:  keyValue[1],
			})
		}
	}

	loggingChannel := make(chan WrappedAuditLog, AUDIT_LOGS_CAPACITY)

	// Finally, create new auditLogger with parameters
	auditLogger := AuditLoggerService{
		logger:          logger.Helper(1),
		enabled:         true,
		serviceURL:      auditLogsURL,
		serviceHeaders:  headers,
		jsonWrapperKey:  os.Getenv("AUDIT_LOGS_FORWARDER_JSON_WRAPPER_KEY"),
		environmentName: env,
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
func (l *AuditLoggerService) Audit(ctx context.Context, eventID EventID, data map[string]interface{}) {
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
		if len(l.loggingChannel) == AUDIT_LOGS_CAPACITY {
			l.logger.Errorw("Audit log buffer is full. Dropping log with eventID: %s", eventID)
		} else {
			l.logger.Errorw("Could not send log to audit subsystem even though queue has %d space", AUDIT_LOGS_CAPACITY-len(l.loggingChannel))
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
func (l *AuditLoggerService) postLogToLogService(eventID EventID, data map[string]interface{}) {
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
