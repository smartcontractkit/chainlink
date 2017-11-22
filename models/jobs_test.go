package models_test

import (
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/orm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSave(t *testing.T) {
	j1 := models.NewJob()
	j1.Schedule = "1 * * * *"
	orm.Init()
	defer orm.Close()

	db := orm.GetDB()
	db.Save(&j1)

	var j2 models.Job
	db.One("ID", j1.ID, &j2)

	assert.Equal(t, j1.Schedule, j2.Schedule)
}
