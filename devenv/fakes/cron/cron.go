package cron

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "fake-cron"}).Logger()

var result = 200

func RegisterRoutes() error {
	return fake.Func("POST", "/cron_response", func(ctx *gin.Context) {
		L.Info().Int("Result", result).Msg("Returning feed value result")
		ctx.JSON(200, gin.H{
			"data": map[string]any{
				"result": result,
			},
		})
	})
}
