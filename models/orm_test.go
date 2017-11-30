package models_test

import (
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhereNotFound(t *testing.T) {
	cltest.SetUpDB()
	defer cltest.TearDownDB()

	j1 := models.NewJob()
	jobs := []models.Job{j1}

	err := models.Where("ID", "bogus", &jobs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestAllIndexedNotFound(t *testing.T) {
	cltest.SetUpDB()
	defer cltest.TearDownDB()

	j1 := models.NewJob()
	jobs := []models.Job{j1}

	err := models.AllIndexed("Cron", &jobs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}
