package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
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
	lvls = append(lvls, strconv.FormatBool(cc.App.GetConfig().LogSQL()))

	logSvcs := logger.GetLogServices()
	logORM := logger.NewORM(cc.App.GetSqlxDB(), cc.App.GetLogger())
	for _, svcName := range logSvcs {
		lvl, _ := logORM.GetServiceLogLevel(svcName)

		svcs = append(svcs, svcName)
		lvls = append(lvls, lvl)
	}

	response := &presenters.ServiceLogConfigResource{
		JAID: presenters.JAID{
			ID: "log",
		},
		ServiceName:     svcs,
		LogLevel:        lvls,
		DefaultLogLevel: cc.App.GetConfig().DefaultLogLevel().String(),
	}

	jsonAPIResponse(c, response, "log")
}

// Patch sets a log level and enables sql logging for the logger
func (cc *LogController) Patch(c *gin.Context) {
	ctx := c.Request.Context()
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
		if err := cc.App.SetLogLevel(ll); err != nil {
			jsonAPIError(c, http.StatusInternalServerError, err)
			return
		}
	}
	svcs = append(svcs, "Global")
	lvls = append(lvls, cc.App.GetConfig().LogLevel().String())

	if request.SqlEnabled != nil {
		cc.App.GetConfig().SetLogSQL(*request.SqlEnabled)
	}

	svcs = append(svcs, "IsSqlEnabled")
	lvls = append(lvls, strconv.FormatBool(cc.App.GetConfig().LogSQL()))

	if len(request.ServiceLogLevel) > 0 {
		logORM := logger.NewORM(cc.App.GetSqlxDB(), cc.App.GetLogger())
		for _, svcLogLvl := range request.ServiceLogLevel {
			svcName := svcLogLvl[0]
			svcLvl := svcLogLvl[1]

			var lvl zapcore.Level
			err := lvl.UnmarshalText([]byte(svcLvl))
			if err != nil {
				jsonAPIError(c, http.StatusBadRequest, err)
				return
			}

			if err := cc.App.SetServiceLogLevel(ctx, svcName, lvl); err != nil {
				jsonAPIError(c, http.StatusInternalServerError, err)
				return
			}

			ll, _ := logORM.GetServiceLogLevel(svcName)

			svcs = append(svcs, svcName)
			lvls = append(lvls, ll)
		}
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
