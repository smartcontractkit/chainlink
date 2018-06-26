// +build !windows

package services_test

import (
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/tevino/abool"
)

func TestServices_ApplicationSignalShutdown(t *testing.T) {
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

func TestRunManager_Start_ResumeSleepingRuns(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()
	rm := services.NewRunManager(store)

	j := models.NewJob()
	i := models.Initiator{Type: models.InitiatorWeb}
	j.Initiators = []models.Initiator{i}
	json := fmt.Sprintf(`{"until":"%v"}`, utils.ISO8601UTC(time.Now().Add(time.Second)))
	j.Tasks = []models.TaskSpec{cltest.NewTask("sleep", json)}
	assert.NoError(t, store.Save(&j))

	jr := j.NewRun(i)
	jr.Status = models.RunStatusPendingSleep
	assert.NoError(t, store.Save(&jr))

	assert.NoError(t, rm.ResumeSleepingRuns())
	input, open := store.RunQueue.Pop()
	assert.Equal(t, jr.ID, input.JobRunID)
	assert.True(t, open)
}
