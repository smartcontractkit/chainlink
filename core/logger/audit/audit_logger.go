package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

const bufferCapacity = 2048
const webRequestTimeout = 10

type Data = map[string]any

type AuditLogger interface {
	services.Service

	Audit(eventID EventID, data Data)
}

type HTTPAuditLoggerInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type AuditLoggerService struct {
	logger          logger.Logger            // The standard logger configured in the node
	enabled         bool                     // Whether the audit logger is enabled or not
	forwardToUrl    commonconfig.URL         // Location we are going to send logs to
	headers         []models.ServiceHeader   // Headers to be sent along with logs for identification/authentication
	jsonWrapperKey  string                   // Wrap audit data as a map under this key if present
	environmentName string                   // Decorate the environment this is coming from
	hostname        string                   // The self-reported hostname of the machine
	localIP         string                   // A non-loopback IP address as reported by the machine
	loggingClient   HTTPAuditLoggerInterface // Abstract type for sending logs onward

	loggingChannel chan wrappedAuditLog
	chStop         services.StopChan
	chDone         chan struct{}
}

type wrappedAuditLog struct {
	eventID EventID
	data    Data
}

var NoopLogger AuditLogger = &AuditLoggerService{}

// NewAuditLogger returns a buffer push system that ingests audit log events and
// asynchronously pushes them up to an HTTP log service.
// Parses and validates the AUDIT_LOGS_* environment values and returns an enabled
// AuditLogger instance. If the environment variables are not set, the logger
// is disabled and short circuits execution via enabled flag.
func NewAuditLogger(logger logger.Logger, config config.AuditLogger) (AuditLogger, error) {
	// If the unverified config is nil, then we assume this came from the
	// configuration system and return a nil logger.
	if config == nil || !config.Enabled() {
		return &AuditLoggerService{}, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("initialization error - unable to get hostname: %w", err)
	}

	forwardToUrl, err := config.ForwardToUrl()
	if err != nil {
		return &AuditLoggerService{}, nil
	}

	headers, err := config.Headers()
	if err != nil {
		return &AuditLoggerService{}, nil
	}

	loggingChannel := make(chan wrappedAuditLog, bufferCapacity)

	// Create new AuditLoggerService
	auditLogger := AuditLoggerService{
		logger:          logger.Helper(1),
		enabled:         true,
		forwardToUrl:    forwardToUrl,
		headers:         headers,
		jsonWrapperKey:  config.JsonWrapperKey(),
		environmentName: config.Environment(),
		hostname:        hostname,
		localIP:         getLocalIP(),
		loggingClient:   &http.Client{Timeout: time.Second * webRequestTimeout},

		loggingChannel: loggingChannel,
		chStop:         make(chan struct{}),
		chDone:         make(chan struct{}),
	}

	return &auditLogger, nil
}

func (l *AuditLoggerService) SetLoggingClient(newClient HTTPAuditLoggerInterface) {
	l.loggingClient = newClient
}

// Entrypoint for new audit logs. This buffers all logs that come in they will
// sent out by the goroutine that was started when the AuditLoggerService was
// created. If this service was not enabled, this immeidately returns.
//
// This function never blocks.
func (l *AuditLoggerService) Audit(eventID EventID, data Data) {
	if !l.enabled {
		return
	}

	wrappedLog := wrappedAuditLog{
		eventID: eventID,
		data:    data,
	}

	select {
	case l.loggingChannel <- wrappedLog:
	default:
		l.logger.Errorf("buffer is full. Dropping log with eventID: %s", eventID)
	}
}

// Start the audit logger and begin processing logs on the channel
func (l *AuditLoggerService) Start(context.Context) error {
	if !l.enabled {
		return errors.New("The audit logger is not enabled")
	}

	go l.runLoop()
	return nil
}

// Stops the logger and will close the channel.
func (l *AuditLoggerService) Close() error {
	if !l.enabled {
		return errors.New("The audit logger is not enabled")
	}

	l.logger.Warnf("Disabled the audit logger service")
	close(l.chStop)
	<-l.chDone

	return nil
}

func (l *AuditLoggerService) Name() string {
	return l.logger.Name()
}

func (l *AuditLoggerService) HealthReport() map[string]error {
	var err error
	if !l.enabled {
		err = errors.New("the audit logger is not enabled")
	} else if len(l.loggingChannel) == bufferCapacity {
		err = errors.New("buffer is full")
	}
	return map[string]error{l.Name(): err}
}

func (l *AuditLoggerService) Ready() error {
	if !l.enabled {
		return errors.New("the audit logger is not enabled")
	}

	return nil
}

// Entrypoint for our log handling goroutine. This waits on the channel and sends out
// logs as they come in.
//
// This function calls postLogToLogService which blocks.
func (l *AuditLoggerService) runLoop() {
	defer close(l.chDone)

	for {
		select {
		case <-l.chStop:
			l.logger.Warn("The audit logger is shutting down")
			return
		case event := <-l.loggingChannel:
			l.postLogToLogService(event.eventID, event.data)
		}
	}
}

// Takes an EventID and associated data and sends it to the configured logging
// endpoint. This function blocks on the send by timesout after a period of
// several seconds. This helps us prevent getting stuck on a single log
// due to transient network errors.
//
// This function blocks when called.
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
		l.logger.Errorw("unable to serialize wrapped audit log item to JSON", "err", err, "logItem", logItem)
		return
	}
	ctx, cancel := l.chStop.NewCtx()
	defer cancel()

	// Send to remote service
	req, err := http.NewRequestWithContext(ctx, "POST", (*url.URL)(&l.forwardToUrl).String(), bytes.NewReader(serializedLog))
	if err != nil {
		l.logger.Error("failed to create request to remote logging service!")
	}
	for _, header := range l.headers {
		req.Header.Add(header.Header, header.Value)
	}
	resp, err := l.loggingClient.Do(req)
	if err != nil {
		l.logger.Errorw("failed to send audit log to HTTP log service", "err", err, "logItem", logItem)
		return
	}
	if resp.StatusCode != 200 {
		if resp.Body == nil {
			l.logger.Errorw("no body to read. Possibly an error occurred sending", "logItem", logItem)
			return
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			l.logger.Errorw("error reading errored HTTP log service webhook response body", "err", err, "logItem", logItem)
			return
		}
		l.logger.Errorw("error sending log to HTTP log service", "statusCode", resp.StatusCode, "bodyString", string(bodyBytes))
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
