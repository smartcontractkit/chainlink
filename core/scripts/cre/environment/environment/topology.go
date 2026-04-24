package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/topologyviz"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

const environmentRootDir = "core/scripts/cre/environment"

type topologyProbe struct {
	Blockchains []any `toml:"blockchains"`
	NodeSets    []any `toml:"nodesets"`
	JD          any   `toml:"jd"`
	Infra       any   `toml:"infra"`
}

type discoveredTopology struct {
	ConfigPath string
	ConfigAbs  string
	Summary    *topologyviz.TopologySummary
}

func TopologyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "topology",
		Short: "Topology discovery and visualization commands",
		Long:  "List, inspect, and document CRE DON topology TOML files.",
	}

	cmd.AddCommand(topologyListCmd())
	cmd.AddCommand(topologyShowCmd())
	cmd.AddCommand(topologyGenerateCmd())
	return cmd
}

func topologyListCmd() *cobra.Command {
	return &cobra.Command{
		Use:              "list",
		Short:            "List available topology configs",
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			discovered, err := discoverTopologies()
			if err != nil {
				return err
			}
			if len(discovered) == 0 {
				fmt.Println("No topology TOML files found in configs/")
				return nil
			}
			fmt.Println()
			fmt.Println("Available topologies:")
			fmt.Print(renderTopologiesTable(discovered))
			fmt.Println()
			return nil
		},
	}
}

func topologyShowCmd() *cobra.Command {
	var (
		configPath string
		outputDir  string
	)

	cmd := &cobra.Command{
		Use:              "show",
		Short:            "Show and generate topology visualization for one config",
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPathAbs, err := resolveConfigPath(configPath)
			if err != nil {
				return err
			}
			_, summary, err := loadAndSummarizeConfig(cfgPathAbs)
			if err != nil {
				return err
			}

			fmt.Print(topologyviz.RenderASCII(summary))

			outputDirAbs, err := resolveOutputPath(outputDir)
			if err != nil {
				return err
			}
			artifacts, err := topologyviz.WriteArtifacts(summary, outputDirAbs)
			if err != nil {
				return err
			}

			fmt.Printf("\nTopology artifacts:\n- %s\n- %s\n", artifacts.ASCIIPath, artifacts.MarkdownPath)
			fmt.Println()
			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "configs/workflow-gateway-don.toml", "Path to topology TOML file")
	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", "state", "Output directory for generated artifacts")
	return cmd
}

func topologyGenerateCmd() *cobra.Command {
	var (
		outputDir string
		indexPath string
		checkOnly bool
	)

	cmd := &cobra.Command{
		Use:              "generate",
		Short:            "Generate topology docs for all topology configs",
		PersistentPreRun: globalPreRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			discovered, err := discoverTopologies()
			if err != nil {
				return err
			}
			if len(discovered) == 0 {
				return errors.New("no topology configs discovered")
			}

			outputDirAbs, err := resolveOutputPath(outputDir)
			if err != nil {
				return err
			}
			indexPathAbs, err := resolveOutputPath(indexPath)
			if err != nil {
				return err
			}
			if mkErr := os.MkdirAll(outputDirAbs, 0o755); mkErr != nil {
				return errors.Wrap(mkErr, "failed to create topology docs directory")
			}
			if mkErr := os.MkdirAll(filepath.Dir(indexPathAbs), 0o755); mkErr != nil {
				return errors.Wrap(mkErr, "failed to create topology index directory")
			}

			var outOfDate []string
			for _, item := range discovered {
				base := strings.TrimSuffix(filepath.Base(item.ConfigPath), filepath.Ext(item.ConfigPath))
				mdTarget := filepath.Join(outputDirAbs, base+".md")
				mdContent := topologyviz.RenderMarkdown(item.Summary)
				outdatedMD, writeErr := writeOrCheck(mdTarget, mdContent, checkOnly)
				if writeErr != nil {
					return writeErr
				}
				if outdatedMD {
					outOfDate = append(outOfDate, mdTarget)
				}

				jsonTarget := filepath.Join(outputDirAbs, base+".json")
				if !checkOnly {
					_ = os.Remove(jsonTarget)
				}
			}

			indexContent := renderTopologyIndex(discovered, outputDirAbs, indexPathAbs)
			outdatedIndex, writeErr := writeOrCheck(indexPathAbs, indexContent, checkOnly)
			if writeErr != nil {
				return writeErr
			}
			if outdatedIndex {
				outOfDate = append(outOfDate, indexPathAbs)
			}

			if checkOnly {
				if len(outOfDate) > 0 {
					sort.Strings(outOfDate)
					return fmt.Errorf("topology docs are out of date; regenerate with `go run . topology generate`:\n- %s", strings.Join(outOfDate, "\n- "))
				}
				fmt.Println("Topology docs are up to date")
				return nil
			}

			fmt.Printf("Generated topology docs for %d topology configs\n", len(discovered))
			fmt.Printf("Index: %s\n", indexPathAbs)
			fmt.Printf("Per-topology docs directory: %s\n", outputDirAbs)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", "docs/topologies", "Directory for generated per-topology docs")
	cmd.Flags().StringVarP(&indexPath, "index-path", "i", "docs/TOPOLOGIES.md", "Path to generated topology index markdown")
	cmd.Flags().BoolVar(&checkOnly, "check", false, "Check if generated topology docs are up to date")
	return cmd
}

func generateTopologyArtifactsForLoadedConfig(cfg *envconfig.Config) (*topologyviz.TopologySummary, *topologyviz.Artifacts, error) {
	stateDirAbs, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, envconfig.StateDirname))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to resolve topology artifact output directory")
	}
	summary, err := topologyviz.BuildSummary(cfg, os.Getenv("CTF_CONFIGS"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to build topology summary")
	}
	artifacts, err := topologyviz.WriteArtifacts(summary, stateDirAbs)
	if err != nil {
		return nil, nil, err
	}
	return summary, artifacts, nil
}

func discoverTopologies() ([]discoveredTopology, error) {
	envDirAbs, err := environmentDirAbs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve environment directory")
	}
	configsDirAbs, err := environmentPathAbs("configs")
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve configs directory")
	}
	defaultsAbs, err := environmentPathAbs(defaultCapabilitiesConfigFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve capability defaults path")
	}

	var candidates []string
	walkErr := filepath.WalkDir(configsDirAbs, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".toml") {
			return nil
		}
		if path == defaultsAbs {
			return nil
		}
		ok, err := isTopologyConfig(path)
		if err != nil || !ok {
			return nil
		}
		candidates = append(candidates, path)
		return nil
	})
	if walkErr != nil {
		return nil, errors.Wrap(walkErr, "failed to scan topology configs")
	}

	sort.Strings(candidates)
	discovered := make([]discoveredTopology, 0, len(candidates))
	for _, candidate := range candidates {
		_, summary, err := loadAndSummarizeConfig(candidate)
		if err != nil {
			continue
		}
		relPath, relErr := filepath.Rel(envDirAbs, candidate)
		if relErr != nil {
			relPath = candidate
		}
		discovered = append(discovered, discoveredTopology{
			ConfigPath: filepath.ToSlash(relPath),
			ConfigAbs:  candidate,
			Summary:    summary,
		})
	}
	return discovered, nil
}

func isTopologyConfig(path string) (bool, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	var probe topologyProbe
	if err := toml.Unmarshal(raw, &probe); err != nil {
		return false, err
	}
	return len(probe.NodeSets) > 0 &&
		len(probe.Blockchains) > 0 &&
		probe.JD != nil &&
		probe.Infra != nil, nil
}

func loadAndSummarizeConfig(configPath string) (*envconfig.Config, *topologyviz.TopologySummary, error) {
	cfg := &envconfig.Config{}
	cfgArg := defaultCapabilitiesConfigFile + "," + configPath
	if err := cfg.Load(cfgArg); err != nil {
		return nil, nil, errors.Wrapf(err, "failed to load topology config: %s", configPath)
	}

	envDirAbs, err := environmentDirAbs()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to resolve environment directory")
	}
	configRef := configPath
	if rel, relErr := filepath.Rel(envDirAbs, configPath); relErr == nil {
		configRef = filepath.ToSlash(rel)
	}
	summary, err := topologyviz.BuildSummary(cfg, configRef)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to build topology summary")
	}
	return cfg, summary, nil
}

func resolveConfigPath(configPath string) (string, error) {
	if configPath == "" {
		return "", errors.New("config path must not be empty")
	}
	if filepath.IsAbs(configPath) {
		return configPath, nil
	}
	envDirAbs, err := environmentDirAbs()
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve environment directory")
	}
	return filepath.Join(envDirAbs, configPath), nil
}

func resolveOutputPath(pathArg string) (string, error) {
	if pathArg == "" {
		return "", errors.New("output path must not be empty")
	}
	if filepath.IsAbs(pathArg) {
		return pathArg, nil
	}
	envDirAbs, err := environmentDirAbs()
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve environment directory")
	}
	return filepath.Join(envDirAbs, pathArg), nil
}

func writeOrCheck(targetPath, expected string, checkOnly bool) (bool, error) {
	currentBytes, err := os.ReadFile(targetPath)
	if err != nil && !os.IsNotExist(err) {
		return false, errors.Wrapf(err, "failed to read file: %s", targetPath)
	}
	outdated := string(currentBytes) != expected
	if checkOnly {
		return outdated, nil
	}
	if !outdated {
		return false, nil
	}
	if err := os.WriteFile(targetPath, []byte(expected), 0o644); err != nil { //nolint:gosec // we want the documentation to be readable by everyone
		return false, errors.Wrapf(err, "failed to write file: %s", targetPath)
	}
	return false, nil
}

func renderTopologyIndex(items []discoveredTopology, outputDirAbs, indexPathAbs string) string {
	var b strings.Builder
	b.WriteString("# CRE Topologies\n\n")
	b.WriteString("This file is generated by `go run . topology generate`. Do not edit manually.\n\n")
	b.WriteString("| Config | Class | DONs |\n")
	b.WriteString("|---|---|---:|\n")
	indexDir := filepath.Dir(indexPathAbs)
	for _, item := range items {
		base := strings.TrimSuffix(filepath.Base(item.ConfigPath), filepath.Ext(item.ConfigPath))
		docPath := filepath.Join(outputDirAbs, base+".md")
		relDocPath, relErr := filepath.Rel(indexDir, docPath)
		if relErr != nil {
			relDocPath = docPath
		}
		relDocPath = filepath.ToSlash(relDocPath)
		b.WriteString(fmt.Sprintf("| `%s` | `%s` | %d ([details](%s)) |\n",
			item.ConfigPath,
			item.Summary.Topology,
			len(item.Summary.DONs),
			relDocPath,
		))
	}
	b.WriteString("\nTip: run `go run . topology list` for quick terminal guidance.\n")
	return b.String()
}

func renderTopologiesTable(items []discoveredTopology) string {
	headers := []string{"Topology", "Class", "DONs"}
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	rows := make([][]string, 0, len(items))
	for _, item := range items {
		row := []string{
			item.ConfigPath,
			item.Summary.Topology,
			strconv.Itoa(len(item.Summary.DONs)),
		}
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
		rows = append(rows, row)
	}

	var b strings.Builder
	b.WriteString(renderASCIIBorder(widths))
	b.WriteString(renderASCIIRow(headers, widths))
	b.WriteString(renderASCIIBorder(widths))
	for _, row := range rows {
		b.WriteString(renderASCIIRow(row, widths))
	}
	b.WriteString(renderASCIIBorder(widths))
	return b.String()
}

func renderASCIIBorder(widths []int) string {
	var b strings.Builder
	b.WriteString("+")
	for _, w := range widths {
		b.WriteString(strings.Repeat("-", w+2))
		b.WriteString("+")
	}
	b.WriteString("\n")
	return b.String()
}

func renderASCIIRow(values []string, widths []int) string {
	var b strings.Builder
	b.WriteString("|")
	for i, v := range values {
		b.WriteString(" ")
		b.WriteString(v)
		if len(v) < widths[i] {
			b.WriteString(strings.Repeat(" ", widths[i]-len(v)))
		}
		b.WriteString(" |")
	}
	b.WriteString("\n")
	return b.String()
}

func environmentDirAbs() (string, error) {
	return environmentPathAbs()
}

func environmentPathAbs(parts ...string) (string, error) {
	pathParts := append([]string{relativePathToRepoRoot, environmentRootDir}, parts...)
	return filepath.Abs(filepath.Join(pathParts...))
}
