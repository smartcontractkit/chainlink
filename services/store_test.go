package services_test

import (
	"io/ioutil"
	"testing"

	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestCreateKey(t *testing.T) {
	t.Parallel()
	store := cltest.Store()
	defer store.Close()

	password := "p@ssword"
	store.CreateKey(password)

	files, _ := ioutil.ReadDir(store.Config.KeysDir())
	assert.Equal(t, 1, len(files))
}
