package docs

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/utils"
)

const (
	fieldDefault = "# Default"
	fieldExample = "# Example"

	tokenAdvanced = "**ADVANCED**"
	tokenExtended = "**EXTENDED**"
)

var (
	//go:embed secrets.toml
	secretsTOML string
	//go:embed core.toml
	coreTOML string
	//go:embed chains-evm.toml
	chainsEVMTOML string
	//go:embed chains-solana.toml
	chainsSolanaTOML string
	//go:embed chains-starknet.toml
	chainsStarknetTOML string

	//go:embed example-config.toml
	exampleConfig string
	//go:embed example-secrets.toml
	exampleSecrets string

	docsTOML = coreTOML + chainsEVMTOML + chainsSolanaTOML + chainsStarknetTOML
)

// GenerateConfig returns MarkDown documentation generated from core.toml & chains-*.toml.
func GenerateConfig() (string, error) {
	return generateDocs(docsTOML, `[//]: # (Documentation generated from docs/*.toml - DO NOT EDIT.)

This document describes the TOML format for configuration.

See also [SECRETS.md](SECRETS.md)
`, exampleConfig)
}

// GenerateSecrets returns MarkDown documentation generated from secrets.toml.
func GenerateSecrets() (string, error) {
	return generateDocs(secretsTOML, `[//]: # (Documentation generated from docs/secrets.toml - DO NOT EDIT.)

This document describes the TOML format for secrets.

Each secret has an alternative corresponding environment variable.

See also [CONFIG.md](CONFIG.md)
`, exampleSecrets)
}

// generateDocs returns MarkDown documentation generated from the TOML string.
func generateDocs(toml, header, example string) (string, error) {
	items, err := parseTOMLDocs(toml)
	var sb strings.Builder

	sb.WriteString(header)
	sb.WriteString(`
## Example

`)
	sb.WriteString("```toml\n")
	sb.WriteString(example)
	sb.WriteString("```\n\n")

	for _, item := range items {
		sb.WriteString(item.String())
		sb.WriteString("\n\n")
	}

	return sb.String(), err
}

func advancedWarning(msg string) string {
	return fmt.Sprintf(":warning: **_ADVANCED_**: _%s_\n", msg)
}

// lines holds a set of contiguous lines
type lines []string

func (d lines) String() string {
	return strings.Join(d, "\n")
}

type table struct {
	name  string
	codes lines
	adv   bool
	desc  lines
	ext   bool
}

func newTable(line string, desc lines) *table {
	t := &table{
		name:  strings.Trim(line, "[]"),
		codes: []string{line},
		desc:  desc,
	}
	if len(desc) > 0 {
		if strings.HasPrefix(strings.TrimSpace(desc[0]), tokenAdvanced) {
			t.adv = true
			t.desc = t.desc[1:]
		} else if strings.HasPrefix(strings.TrimSpace(desc[len(desc)-1]), tokenExtended) {
			t.ext = true
			t.desc = t.desc[:len(desc)-1]
		}
	}
	return t
}

func (t table) advanced() string {
	if t.adv {
		return advancedWarning("Do not change these settings unless you know what you are doing.")
	}
	return ""
}

func (t table) code() string {
	if !t.ext {
		return fmt.Sprint("```toml\n", t.codes, "\n```\n")
	}
	return ""
}

func (t table) extended() string {
	if t.ext {
		if t.name != "EVM" {
			log.Fatalf("%s: no extended description available", t.name)
		}
		s, err := evmChainDefaults()
		if err != nil {
			log.Fatalf("%s: failed to generate evm chain defaults: %v", t.name, err)
		}
		return s
	}
	return ""
}

// String prints a table as an H2, followed by a code block and description.
func (t *table) String() string {
	return fmt.Sprint("## ", t.name, "\n",
		t.advanced(),
		t.code(),
		t.desc,
		t.extended())
}

type keyval struct {
	name string
	code string
	adv  bool
	desc lines
}

func newKeyval(line string, desc lines) keyval {
	line = strings.TrimSpace(line)
	kv := keyval{
		name: line[:strings.Index(line, " ")],
		code: line,
		desc: desc,
	}
	if len(desc) > 0 && strings.HasPrefix(strings.TrimSpace(desc[0]), tokenAdvanced) {
		kv.adv = true
		kv.desc = kv.desc[1:]
	}
	return kv
}

func (k keyval) advanced() string {
	if k.adv {
		return advancedWarning("Do not change this setting unless you know what you are doing.")
	}
	return ""
}

// String prints a keyval as an H3, followed by a code block and description.
func (k keyval) String() string {
	name := k.name
	if i := strings.LastIndex(name, "."); i > -1 {
		name = name[i+1:]
	}
	return fmt.Sprint("### ", name, "\n",
		k.advanced(),
		"```toml\n",
		k.code,
		"\n```\n",
		k.desc)
}

func parseTOMLDocs(s string) (items []fmt.Stringer, err error) {
	defer func() { _, err = utils.MultiErrorList(err) }()
	globalTable := table{name: "Global"}
	currentTable := &globalTable
	items = append(items, currentTable)
	var desc lines
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(line, "#") {
			// comment
			desc = append(desc, strings.TrimSpace(line[1:]))
		} else if strings.TrimSpace(line) == "" {
			// empty
			if len(desc) > 0 {
				items = append(items, desc)
				desc = nil
			}
		} else if strings.HasPrefix(line, "[") {
			currentTable = newTable(line, desc)
			items = append(items, currentTable)
			desc = nil
		} else {
			kv := newKeyval(line, desc)
			shortName := kv.name
			if currentTable != &globalTable {
				// update to full name
				kv.name = currentTable.name + "." + kv.name
			}
			if len(kv.desc) == 0 {
				err = multierr.Append(err, fmt.Errorf("%s: missing description", kv.name))
			} else if !strings.HasPrefix(kv.desc[0], shortName) {
				err = multierr.Append(err, fmt.Errorf("%s: description does not begin with %q", kv.name, shortName))
			}
			if !strings.HasSuffix(line, fieldDefault) && !strings.HasSuffix(line, fieldExample) {
				err = multierr.Append(err, fmt.Errorf(`%s: is not one of %v`, kv.name, []string{fieldDefault, fieldExample}))
			}

			items = append(items, kv)
			currentTable.codes = append(currentTable.codes, kv.code)
			desc = nil
		}
	}
	if len(globalTable.codes) == 0 {
		//drop it
		items = items[1:]
	}
	if len(desc) > 0 {
		items = append(items, desc)
	}
	return
}
