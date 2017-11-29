package models_test

import (
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSave(t *testing.T) {
	cltest.SetUpDB()
	defer cltest.TearDownDB()

	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}

	models.Save(&j1)

	var j2 models.Job
	models.Find("ID", j1.ID, &j2)

	assert.Equal(t, j1.Schedule, j2.Schedule)
}
