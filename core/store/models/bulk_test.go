package models_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
)

func TestBulkDeleteRunTask_ValidateBulkDeleteRunRequest(t *testing.T) {
	err := models.ValidateBulkDeleteRunRequest(&models.BulkDeleteRunRequest{})
	assert.NoError(t, err)

	err = models.ValidateBulkDeleteRunRequest(&models.BulkDeleteRunRequest{
		Status: []models.RunStatus{models.RunStatusCompleted},
	})
	assert.NoError(t, err)

	err = models.ValidateBulkDeleteRunRequest(&models.BulkDeleteRunRequest{
		Status: []models.RunStatus{""},
	})
	assert.Error(t, err)

	err = models.ValidateBulkDeleteRunRequest(&models.BulkDeleteRunRequest{
		Status: []models.RunStatus{models.RunStatusInProgress},
	})
	assert.Error(t, err)
}
