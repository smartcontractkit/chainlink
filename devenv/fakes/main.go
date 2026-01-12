package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "ocr2"}).Logger()

const (
	DefaultJuelsPerLinkRatio = "15"
)

// some initial value, otherwise OCR2 jobs won't start
var result = "200"

// a very simple mock that allow us to control EA answers in tests
func main() {
	_, err := fake.NewFakeDataProvider(&fake.Input{Port: fake.DefaultFakeServicePort})
	if err != nil {
		panic(err)
	}
	err = fake.Func("POST", "/juelsPerFeeCoinSource", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"data": map[string]any{
				"result": DefaultJuelsPerLinkRatio,
			},
		})
	})
	if err != nil {
		panic(err)
	}

	err = fake.Func("POST", "/trigger_deviation", func(ctx *gin.Context) {
		result = ctx.Query("result")
		L.Info().Str("Result", result).Msg("Changing returned result")
		ctx.JSON(200, gin.H{
			"result": "ok",
		})
	})

	if err != nil {
		panic(err)
	}

	err = fake.Func("POST", "/ea", func(ctx *gin.Context) {
		L.Info().Str("Result", result).Msg("Returning feed value result")
		ctx.JSON(200, gin.H{
			"data": map[string]any{
				"result": result,
			},
		})
	})
	if err != nil {
		panic(err)
	}
	select {}
}
