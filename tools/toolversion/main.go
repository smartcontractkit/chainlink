package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/drift"
	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/manifest"
	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/paths"
	"github.com/smartcontractkit/chainlink/v2/tools/toolversion/internal/resolve"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "toolversion: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "toolversion",
		Short: "Read dev-tool versions from .tool-versions and tools/go-tools.txt",
	}

	root.AddCommand(
		cmdGet(),
		cmdRef(),
		cmdTarget(),
		cmdGoInstall(),
		cmdList(),
		cmdModules(),
		cmdCheck(),
		cmdMakeVars(),
	)
	return root
}

func loadResolver() (*resolve.Resolver, paths.Config, error) {
	cfg, err := paths.FromEnv()
	if err != nil {
		return nil, cfg, err
	}
	store, err := manifest.New(cfg.ToolVersionsFile, cfg.GoToolsFile)
	if err != nil {
		return nil, cfg, err
	}
	return resolve.New(store), cfg, nil
}

func cmdGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Print the raw version for a tool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, _, err := loadResolver()
			if err != nil {
				return err
			}
			v, err := r.Get(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), v)
			return nil
		},
	}
}

func cmdRef() *cobra.Command {
	return &cobra.Command{
		Use:   "ref <key>",
		Short: "Print a consumer-ready version reference (v-prefix for semver)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, _, err := loadResolver()
			if err != nil {
				return err
			}
			v, err := r.Ref(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), v)
			return nil
		},
	}
}

func cmdTarget() *cobra.Command {
	return &cobra.Command{
		Use:   "target <key>",
		Short: "Print the go install argument: module@version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, _, err := loadResolver()
			if err != nil {
				return err
			}
			t, err := r.Target(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), t)
			return nil
		},
	}
}

func cmdGoInstall() *cobra.Command {
	return &cobra.Command{
		Use:   "go-install <key>",
		Short: "Run go install for the tool at the pinned version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, _, err := loadResolver()
			if err != nil {
				return err
			}
			target, err := r.Target(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "go install %s\n", target)
			c := exec.Command("go", "install", target)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			return c.Run()
		},
	}
}

func cmdList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Print all name/version pairs from both manifests",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, _, err := loadResolver()
			if err != nil {
				return err
			}
			entries, err := r.List()
			if err != nil {
				return err
			}
			for _, e := range entries {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", e.Name, e.Version)
			}
			return nil
		},
	}
}

func cmdModules() *cobra.Command {
	return &cobra.Command{
		Use:   "modules",
		Short: "Print all managed go module import paths",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, _, err := loadResolver()
			if err != nil {
				return err
			}
			for _, m := range r.ManagedModules() {
				fmt.Fprintln(cmd.OutOrStdout(), m)
			}
			return nil
		},
	}
}

func cmdCheck() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Fail if version pins drift from the manifests",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, cfg, err := loadResolver()
			if err != nil {
				return err
			}
			if err := drift.NewChecker(cfg, r).Check(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "check-tool-versions: ok")
			return nil
		},
	}
}

func cmdMakeVars() *cobra.Command {
	return &cobra.Command{
		Use:   "make-vars",
		Short: "Print Makefile variable assignments for common tool versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, _, err := loadResolver()
			if err != nil {
				return err
			}
			vars := []struct {
				key  string
				name string
			}{
				{"golangci-lint", "GOLANGCI_LINT_VERSION"},
				{"protoc", "PROTOC_VERSION"},
			}
			var b strings.Builder
			for _, v := range vars {
				var val string
				if v.key == "protoc" {
					val, err = r.Get(v.key)
				} else {
					val, err = r.Ref(v.key)
				}
				if err != nil {
					return err
				}
				fmt.Fprintf(&b, "%s=%s\n", v.name, val)
			}
			fmt.Fprint(cmd.OutOrStdout(), b.String())
			return nil
		},
	}
}
