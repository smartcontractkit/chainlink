package models

import (
	"github.com/smartcontractkit/chainlink-go/orm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSave(t *testing.T) {
	j1 := Job{Schedule: "0 1 2 * *"}
	orm.Init()
	defer orm.Close()

	db := orm.GetDB()
	db.Create(&j1)

	j2 := Job{}
	db.First(&j2, j1.ID)

	assert.Equal(t, j1.Schedule, j2.Schedule)
}
