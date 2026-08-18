package modules_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

func BenchmarkFindAffectedModules(b *testing.B) {
	tmpDir := b.TempDir()

	singleFile := []string{"core/services/telemetry/ingress.go"}
	multiModFiles := []string{
		"core/logger/logger.go",
		"core/services/cron/cron.go",
		"deployment/environment.go",
		"tools/ci-testshard/main.go",
		"system-tests/lib/suite.go",
	}
	largeFileSet := make([]string, 0, 100)
	for i := range 50 {
		largeFileSet = append(largeFileSet, filepath.Join("core/services", string(rune('a'+i%26)), "service.go"))
		largeFileSet = append(largeFileSet, filepath.Join("deployment", string(rune('a'+i%26)), "env.go"))
	}

	b.Run("SingleFile", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, err := modules.FindAffectedModules(tmpDir, singleFile)
			require.NoError(b, err)
		}
	})

	b.Run("MultiModule", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, err := modules.FindAffectedModules(tmpDir, multiModFiles)
			require.NoError(b, err)
		}
	})

	b.Run("LargeFileSet100", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, err := modules.FindAffectedModules(tmpDir, largeFileSet)
			require.NoError(b, err)
		}
	})
}

func BenchmarkFindTestPackages(b *testing.B) {
	tmpDir := b.TempDir()

	singleFile := []string{"core/logger/logger_test.go"}
	mixedFiles := []string{
		"core/logger/logger.go",
		"core/services/cron/cron_test.go",
		"tools/ci-testshard/main.go",
		"deployment/environment.go",
		"system-tests/lib/suite.go",
		"integration-tests/smoke/vrf.go",
	}

	b.Run("SingleFile", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, err := modules.FindTestPackages(tmpDir, singleFile)
			require.NoError(b, err)
		}
	})

	b.Run("MixedWithExclusions", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, err := modules.FindTestPackages(tmpDir, mixedFiles)
			require.NoError(b, err)
		}
	})
}
