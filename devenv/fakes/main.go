package main

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink/devenv/fakes/automation"
	"github.com/smartcontractkit/chainlink/devenv/fakes/cron"
	"github.com/smartcontractkit/chainlink/devenv/fakes/directrequest"
	"github.com/smartcontractkit/chainlink/devenv/fakes/ocr2"
)

// a very simple mock that allow us to control EA answers in tests
func main() {
	_, err := fake.NewFakeDataProvider(&fake.Input{Port: fake.DefaultFakeServicePort})
	if err != nil {
		panic(err)
	}
	if err := automation.RegisterRoutes(); err != nil {
		panic(err)
	}
	if err := ocr2.RegisterRoutes(); err != nil {
		panic(err)
	}
	if err := cron.RegisterRoutes(); err != nil {
		panic(err)
	}
	if err := directrequest.RegisterRoutes(); err != nil {
		panic(err)
	}
	select {}
}
