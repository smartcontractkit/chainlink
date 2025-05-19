package src

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustReadNodesList(t *testing.T) {
	t.Run("valid nodes list", func(t *testing.T) {
		content := "localhost:50100 chainlink.core.1:50100 user1 pass1\nlocalhost:50101 chainlink.core.2:50101 user2 pass2"
		filePath := writeTempFile(t, content)
		defer os.Remove(filePath)

		nodes := mustReadNodesList(filePath)
		assert.Len(t, nodes, 2)

		assert.Equal(t, "user1", nodes[0].APILogin)
		assert.Equal(t, "user2", nodes[1].APILogin)

		assert.Equal(t, "pass1", nodes[0].APIPassword)
		assert.Equal(t, "pass2", nodes[1].APIPassword)

		assert.Equal(t, "http://localhost:50100", nodes[0].RemoteURL.String())
		assert.Equal(t, "http://localhost:50101", nodes[1].RemoteURL.String())

		assert.Equal(t, "chainlink.core.1", nodes[0].ServiceName)
		assert.Equal(t, "chainlink.core.2", nodes[1].ServiceName)

		assert.Equal(t, "http://chainlink.core.1:50100", nodes[0].URL.String())
		assert.Equal(t, "http://chainlink.core.2:50101", nodes[1].URL.String())
	})
}

func writeTempFile(t *testing.T, content string) string {
	file, err := os.CreateTemp("", "nodeslist")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	return file.Name()
}
