package models_test

import (
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSave(t *testing.T) {
	store := cltest.Store()
	defer store.Close()

	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}

	store.Save(&j1)

	var j2 models.Job
	store.One("ID", j1.ID, &j2)

	assert.Equal(t, j1.Schedule, j2.Schedule)
}
