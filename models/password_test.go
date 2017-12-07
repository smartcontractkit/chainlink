package models_test

import (
	"encoding/base64"
	"testing"

	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/stretchr/testify/assert"
)

const phrase = "p@ssword"

func TestNewPassword(t *testing.T) {
	t.Parallel()
	store := cltest.Store()
	defer store.Close()
	var passwords []models.Password

	password := models.NewPassword(phrase)
	assert.NotEqual(t, password.Hash, phrase)
	store.AddPassword(password)

	store.ORM.All(&passwords)
	assert.Equal(t, 1, len(passwords))
}

func TestPasswordCheck(t *testing.T) {
	t.Parallel()
	hash, _ := base64.StdEncoding.DecodeString("/qlGjiuorAcE4tFURJzq4vG+SEaYb0KExvGpfMZ0jF0=")
	salt, _ := base64.StdEncoding.DecodeString("+FTNZKzpDKl5BYTCp0CD/LYiHtzEQCph3s/UHKUZHyQ=")
	password := models.Password{Hash: hash, Salt: salt}

	assert.True(t, password.Check(phrase))
	assert.False(t, password.Check("NotThePassword"))
}
