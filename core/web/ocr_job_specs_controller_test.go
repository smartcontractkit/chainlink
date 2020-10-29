package web_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOCRJobSpecsController_Create_ValidationFailure(t *testing.T) {
	_, client, cleanup := setupOCRJobSpecsControllerTests(t)
	defer cleanup()

	fixtureBytes := cltest.MustReadFile(t, "testdata/oracle-spec-invalid-key.toml")

	resp, cleanup := client.Post("/v2/ocr/specs", bytes.NewReader(fixtureBytes))
	defer cleanup()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "{\"errors\":[{\"detail\":\"unrecognised key: isBootstrapNode\"}]}", string(b))
}

func TestOCRJobSpecsController_Create_HappyPath(t *testing.T) {
	app, client, cleanup := setupOCRJobSpecsControllerTests(t)
	defer cleanup()

	fixtureBytes := cltest.MustReadFile(t, "testdata/oracle-spec.toml")

	resp, cleanup := client.Post("/v2/ocr/specs", bytes.NewReader(fixtureBytes))
	defer cleanup()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	job := models.JobSpecV2{}
	require.NoError(t, app.Store.DB.Preload("OffchainreportingOracleSpec").First(&job).Error)

	b, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("{\"jobID\":%v}", job.ID), string(b))

	// Sanity check to make sure it inserted correctly
	require.Equal(t, models.EIP55Address("0x613a38AC1659769640aaE063C651F48E0250454C"), job.OffchainreportingOracleSpec.ContractAddress)
}

func setupOCRJobSpecsControllerTests(t *testing.T) (*cltest.TestApplication, cltest.HTTPClientCleaner, func()) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	return app, client, cleanup
}
