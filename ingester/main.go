package main

import (
	"os"
	"os/signal"
	"time"

	"ingester/service"

	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
)

type runner func(*service.Config) (*service.Application, error)

func init() {
	time.LoadLocation("UTC")
}

func generateCmd() *cobra.Command {
	newcmd := &cobra.Command{
		Use:  "ingester",
		Args: cobra.MaximumNArgs(0),
		Long: "Manages ingestion tasks for the ethereum blockchain",
		Run:  func(_ *cobra.Command, _ []string) { run(service.NewApplication) },
	}
	return newcmd
}

func run(r runner) {
	a, err := r(service.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to create application: %+v", err)
	}

	a.Start(func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
	})
	a.Stop()
}

func main() {
	if err := generateCmd().Execute(); err != nil {
		log.Warn(err)
	}
}
