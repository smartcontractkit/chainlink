package config

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/services/chainlink/cfgtest"
)

//go:embed docs.toml
var docsTOML string

// GenerateDocs returns MarkDown documentation generated from docs.toml.
func GenerateDocs() (string, error) {
	return generateDocs(docsTOML)
}

// generateDocs returns MarkDown documentation generated from the TOML string.
func generateDocs(toml string) (string, error) {
	items, err := parseTOMLDocs(toml)
	var sb strings.Builder

	// Header.
	sb.WriteString(`[//]: # (Documentation generated from docs.toml - DO NOT EDIT.)

## Table of contents

`)
	// Link to each table group.
	for _, item := range items {
		switch t := item.(type) {
		case *table:
			indent := strings.Repeat("\t", strings.Count(t.name, "."))
			name := t.name
			if i := strings.LastIndex(name, "."); i > -1 {
				name = name[i+1:]
			}
			link := strings.ReplaceAll(t.name, ".", "-")
			fmt.Fprintf(&sb, "%s- [%s](#%s)\n", indent, name, link)
		}
	}
	fmt.Fprintln(&sb)

	for _, item := range items {
		fmt.Fprintln(&sb, item)
		fmt.Fprintln(&sb)
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
		s, err := extended(t.name)
		if err != nil {
			log.Fatalf("%s: failed to sprint extended table description: %v", t.name, err)
		}
		return s
	}
	return ""
}

// String prints a table as an H2, followed by a code block and description.
func (t *table) String() string {
	link := strings.ReplaceAll(t.name, ".", "-")
	return fmt.Sprint("## ", t.name, "<a id='", link, "'></a>\n",
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
	link := strings.ReplaceAll(k.name, ".", "-")
	return fmt.Sprint("### ", name, "<a id='", link, "'></a>\n",
		k.advanced(),
		"```toml\n",
		k.code,
		"\n```\n",
		k.desc)
}

func parseTOMLDocs(s string) (items []fmt.Stringer, err error) {
	defer func() { err = cfgtest.MultiErrorList(err) }()
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
			currentTable = &globalTable
			if len(desc) > 0 {
				items = append(items, desc)
				desc = nil
			}
		} else if strings.HasPrefix(line, "[") {
			currentTable = &table{
				name:  strings.Trim(line, "[]"),
				codes: []string{line},
				desc:  desc,
			}
			if len(desc) > 0 && strings.HasPrefix(strings.TrimSpace(desc[0]), "**ADVANCED**") {
				currentTable.adv = true
				currentTable.desc = currentTable.desc[1:]
			} else if len(desc) > 0 && strings.HasPrefix(strings.TrimSpace(desc[len(desc)-1]), "**EXTENDED**") {
				currentTable.ext = true
				currentTable.desc = currentTable.desc[:len(desc)-1]
			}
			items = append(items, currentTable)
			desc = nil
		} else {
			line = strings.TrimSpace(line)
			kv := keyval{
				name: line[:strings.Index(line, " ")],
				code: line,
				desc: desc,
			}
			if len(desc) > 0 && strings.HasPrefix(strings.TrimSpace(desc[0]), "**ADVANCED**") {
				kv.adv = true
				kv.desc = kv.desc[1:]
			}
			shortName := kv.name
			if currentTable != &globalTable {
				kv.name = currentTable.name + "." + kv.name
			}
			if len(kv.desc) == 0 {
				err = multierr.Append(err, fmt.Errorf("%s: missing description", kv.name))
			} else if !strings.HasPrefix(kv.desc[0], shortName) {
				err = multierr.Append(err, fmt.Errorf("%s: description does not begin with %q", kv.name, shortName))
			}
			if !strings.HasSuffix(line, "# Default") && !strings.HasSuffix(line, "# Example") {
				err = multierr.Append(err, fmt.Errorf(`%s: is neither a "# Default" or "# Example"`, kv.name))
			}

			items = append(items, kv)
			currentTable.codes = append(currentTable.codes, kv.code)
			desc = nil
		}
	}
	if len(desc) > 0 {
		items = append(items, desc)
	}
	return
}
