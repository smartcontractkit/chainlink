package cmd_test

import (
	"io/ioutil"
	"testing"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestRendererTableRenderJobs(t *testing.T) {
	r := cmd.RendererTable{ioutil.Discard}
	job := cltest.NewJob()
	jobs := []models.Job{*job}
	assert.Nil(t, r.Render(&jobs))
}
