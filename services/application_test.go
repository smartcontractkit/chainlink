// +build !windows

package services_test

import (
	"syscall"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services/mock_services"
	"github.com/tevino/abool"
)

func TestChainlinkApplication_SignalShutdown(t *testing.T) {
	config, cleanup := cltest.NewConfig()
	defer cleanup()
	app, _ := cltest.NewApplicationWithConfig(config)

	completed := abool.New()
	app.Exiter = func(code int) {
		completed.Set()
	}

	app.Start()
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		return completed.IsSet()
	}).Should(gomega.BeTrue())
}

func TestChainlinkApplication_AddJob(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	ctrl := gomock.NewController(t)
	jobSubscriberMock := mock_services.NewMockJobSubscriber(ctrl)
	app.ChainlinkApplication.JobSubscriber = jobSubscriberMock
	jobSubscriberMock.EXPECT().AddJob(gomock.Any(), nil) // nil to represent "latest" block
	app.AddJob(cltest.NewJob())
}
