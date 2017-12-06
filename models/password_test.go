package models_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/stretchr/testify/assert"
)

const phrase = "p@ssword"

func TestNewPassword(t *testing.T) {
	store := cltest.Store()
	defer store.Close()
	var passwords []models.Password

	password := models.NewPassword(phrase)
	assert.NotEqual(t, password.Hash, phrase)
	store.AddPassword(password)

	store.ORM.All(&passwords)
	assert.Equal(t, 1, len(passwords))
}
