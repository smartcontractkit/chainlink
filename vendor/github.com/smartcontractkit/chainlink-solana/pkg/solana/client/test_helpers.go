package client

import (
	"bytes"
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SetupLocalSolNode sets up a local solana node via solana cli, and returns the url
func SetupLocalSolNode(t *testing.T) string {
	port := utils.MustRandomPort(t)
	faucetPort := utils.MustRandomPort(t)
	url := "http://127.0.0.1:" + port
	cmd := exec.Command("solana-test-validator",
		"--reset",
		"--rpc-port", port,
		"--faucet-port", faucetPort,
	)
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut
	require.NoError(t, cmd.Start())
	t.Cleanup(func() {
		assert.NoError(t, cmd.Process.Kill())
		if err2 := cmd.Wait(); assert.Error(t, err2) {
			if !assert.Contains(t, err2.Error(), "signal: killed", cmd.ProcessState.String()) {
				t.Log("solana-test-validator stderr:", stdErr.String())
			}
		}
	})

	// Wait for api server to boot
	var ready bool
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		client := rpc.New(url)
		out, err := client.GetHealth(context.Background())
		if err != nil || out != rpc.HealthOk {
			t.Logf("API server not ready yet (attempt %d)\n", i+1)
			continue
		}
		ready = true
		break
	}
	if !ready {
		t.Logf("Cmd output: %s\nCmd error: %s\n", stdOut.String(), stdErr.String())
	}
	require.True(t, ready)
	return url
}

func FundTestAccounts(t *testing.T, keys []solana.PublicKey, url string) {
	for i := range keys {
		account := keys[i].String()
		_, err := exec.Command("solana", "airdrop", "100",
			account,
			"--url", url,
		).Output()
		require.NoError(t, err)
	}
}
