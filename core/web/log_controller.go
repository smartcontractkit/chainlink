package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"go.uber.org/zap/zapcore"
)

// LogController manages the logger config
type LogController struct {
	App chainlink.Application
}

type LogPatchRequest struct {
	Level           string      `json:"level"`
	SqlEnabled      *bool       `json:"sqlEnabled"`
	ServiceLogLevel [][2]string `json:"serviceLogLevel"`
}

// Get retrieves the current log config settings
func (cc *LogController) Get(c *gin.Context) {
	var svcs, lvls []string
	svcs = append(svcs, "Global")
	lvls = append(lvls, cc.App.GetConfig().LogLevel().String())

	svcs = append(svcs, "IsSqlEnabled")
	lvls = append(lvls, strconv.FormatBool(cc.App.GetConfig().LogSQLStatements()))

	logSvcs := logger.GetLogServices()
	logORM := logger.NewORM(cc.App.GetDB())
	for _, svcName := range logSvcs {
		lvl, _ := logORM.GetServiceLogLevel(svcName)

		svcs = append(svcs, svcName)
		lvls = append(lvls, lvl)
	}

	response := &presenters.ServiceLogConfigResource{
		JAID: presenters.JAID{
			ID: "log",
		},
		ServiceName: svcs,
		LogLevel:    lvls,
	}

	jsonAPIResponse(c, response, "log")
}

// Patch sets a log level and enables sql logging for the logger
func (cc *LogController) Patch(c *gin.Context) {
	request := &LogPatchRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Build log config response
	var svcs, lvls []string

	// Validate request params
	if request.Level == "" && request.SqlEnabled == nil && len(request.ServiceLogLevel) == 0 {
		jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("please check request params, no params configured"))
		return
	}

	if request.Level != "" {
		var ll zapcore.Level
		err := ll.UnmarshalText([]byte(request.Level))
		if err != nil {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		if err = cc.App.GetConfig().SetLogLevel(c.Request.Context(), ll.String()); err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}
	svcs = append(svcs, "Global")
	lvls = append(lvls, cc.App.GetConfig().LogLevel().String())

	if request.SqlEnabled != nil {
		if err := cc.App.GetConfig().SetLogSQLStatements(c.Request.Context(), *request.SqlEnabled); err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
		postgres.SetLogAllQueries(cc.App.GetDB(), *request.SqlEnabled)
	}
	svcs = append(svcs, "IsSqlEnabled")
	lvls = append(lvls, strconv.FormatBool(cc.App.GetConfig().LogSQLStatements()))

	if len(request.ServiceLogLevel) > 0 {
		logORM := logger.NewORM(cc.App.GetDB())
		for _, svcLogLvl := range request.ServiceLogLevel {
			svcName := svcLogLvl[0]
			svcLvl := svcLogLvl[1]

			if err := cc.App.SetServiceLogger(c.Request.Context(), svcName, svcLvl); err != nil {
				jsonAPIError(c, http.StatusInternalServerError, err)
				return
			}

			ll, _ := logORM.GetServiceLogLevel(svcName)

			svcs = append(svcs, svcName)
			lvls = append(lvls, ll)
		}
	}

	// Set default logger with new configurations
	logger.SetLogger(cc.App.GetConfig().CreateProductionLogger())

	response := &presenters.ServiceLogConfigResource{
		JAID: presenters.JAID{
			ID: "log",
		},
		ServiceName: svcs,
		LogLevel:    lvls,
	}

	jsonAPIResponse(c, response, "log")
}
