// +build !windows

package services_test

import (
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/tevino/abool"
)

func TestServices_ApplicationSignalShutdown(t *testing.T) {
	app, _ := cltest.NewApplication()

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

func TestApplication_ResumeSleptRuns(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(0))

	j := models.NewJob()
	i := models.Initiator{Type: models.InitiatorWeb}
	j.Initiators = []models.Initiator{i}
	json := fmt.Sprintf(`{"until":"%v"}`, utils.ISO8601UTC(time.Now().Add(time.Second)))
	j.Tasks = []models.TaskSpec{cltest.NewTask("sleep", json)}
	assert.NoError(t, app.Store.Save(&j))

	jr := j.NewRun(i)
	jr.Status = models.RunStatusPendingSleep
	assert.NoError(t, app.Store.Save(&jr))

	assert.NoError(t, app.ResumeSleptRuns())
	input := <-app.Store.RunChannel
	assert.Equal(t, jr.ID, input.JobRunID)
}
