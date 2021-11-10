//go:build smoke

package integration_test

import (
	"os"
	"testing"

	"github.com/smartcontractkit/integrations-framework/config"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestIntegration(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	RegisterFailHandler(Fail)
	conf, err := config.NewConfig("../")
	if err != nil {
		Fail("failed to load config")
	}
	log.Logger = log.Logger.Level(zerolog.Level(conf.Logging.Level))
	junitReporter := reporters.NewJUnitReporter("../logs/tests-integration.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Integration suite", []Reporter{junitReporter})
}
