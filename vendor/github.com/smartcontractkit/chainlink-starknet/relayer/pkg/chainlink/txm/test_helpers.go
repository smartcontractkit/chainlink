package txm

import (
	"bytes"
	"net/http"
	"os/exec"
	"testing"
	"time"

	starknetutils "github.com/NethermindEth/starknet.go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

var (
	// seed = 0 keys for starknet-devnet
	PrivateKeys0Seed = []string{
		"0xe3e70682c2094cac629f6fbed82c07cd",
		"0xf728b4fa42485e3a0a5d2f346baa9455",
		"0xeb1167b367a9c3787c65c1e582e2e662",
		"0xf7c1bd874da5e709d4713d60c8a70639",
		"0xe443df789558867f5ba91faf7a024204",
		"0x23a7711a8133287637ebdcd9e87a1613",
		"0x1846d424c17c627923c6612f48268673",
		"0xfcbd04c340212ef7cca5a5a19e4d6e3c",
		"0xb4862b21fb97d43588561712e8e5216a",
		"0x259f4329e6f4590b9a164106cf6a659e",
	}
)

// SetupLocalStarknetNode sets up a local starknet node via cli, and returns the url
func SetupLocalStarknetNode(t *testing.T) string {
	ctx := tests.Context(t)
	port := utils.MustRandomPort(t)
	url := "http://127.0.0.1:" + port
	cmd := exec.Command("starknet-devnet",
		"--seed", "0", // use same seed for testing
		"--port", port,
	)
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	require.NoError(t, cmd.Start())
	t.Cleanup(func() {
		assert.NoError(t, cmd.Process.Kill())
		if err2 := cmd.Wait(); assert.Error(t, err2) {
			if !assert.Contains(t, err2.Error(), "signal: killed", cmd.ProcessState.String()) {
				t.Log("starknet-devnet stderr:", stdErr.String())
			}
		}
		t.Log("starknet-devnet server closed")
	})

	// Wait for api server to boot
	var ready bool
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/is_alive", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != 200 {
			t.Logf("API server not ready yet (attempt %d)\n", i+1)
			continue
		}
		ready = true
		t.Logf("API server ready at %s\n", url)
		break
	}
	require.True(t, ready)
	return url
}

func TestKeys(t *testing.T, count int) (rawkeys [][]byte) {
	require.True(t, len(PrivateKeys0Seed) >= count, "requested more keys than available")
	for i, k := range PrivateKeys0Seed {
		// max number of keys to generate
		if i >= count {
			break
		}
		f, _ := starknetutils.HexToFelt(k)
		keyBytes := f.Bytes()
		rawkeys = append(rawkeys, keyBytes[:])
	}
	return rawkeys
}
