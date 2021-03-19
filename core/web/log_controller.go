package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"go.uber.org/zap/zapcore"
)

// LogController manages the logger config
type LogController struct {
	App chainlink.Application
}

type LoglevelPatchRequest struct {
	LogLevel string `json:"logLevel"`
	LogSql   string `json:"logSql"`
}

func getLogLevelFromStr(logLevel string) (zapcore.Level, error) {
	switch strings.ToLower(logLevel) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "dpanic":
		return zapcore.DPanicLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("could not parse %s as log level (debug, info, warn, error)", logLevel)
	}
}

// SetDebug sets the debug log mode for the logger
func (cc *LogController) SetDebug(c *gin.Context) {
	request := &LoglevelPatchRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	if request.LogLevel == "" && request.LogSql == "" {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("please set either logLevel or logSql as params in order to set the log level"))
		return
	}

	if request.LogLevel != "" {
		ll, err := getLogLevelFromStr(request.LogLevel)
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		cc.App.GetStore().Config.Set("LOG_LEVEL", ll.String())
		err = cc.App.GetStore().SetConfigStrValue("LogLevel", ll.String())
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}

	if request.LogSql != "" {
		logSql, err := strconv.ParseBool(request.LogSql)
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		cc.App.GetStore().Config.Set("LOG_SQL", request.LogSql)
		err = cc.App.GetStore().SetConfigStrValue("LogSQLStatements", request.LogSql)
		if err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		cc.App.GetStore().SetLogging(logSql)
	}

	// Set default logger with new configurations
	logger.Default = cc.App.GetStore().Config.CreateProductionLogger()

	response := &presenters.LogResource{
		JAID: presenters.JAID{
			ID: "log",
		},
		LogLevel: cc.App.GetStore().Config.LogLevel().String(),
		LogSql:   cc.App.GetStore().Config.LogSQLStatements(),
	}

	jsonAPIResponse(c, response, "log")
}
