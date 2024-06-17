package internal

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type Manifest struct {
	Entries []ManifestEntry
	m       map[string]ManifestEntry
}

func (m Manifest) Latest() (ManifestEntry, error) {
	if len(m.Entries) == 0 {
		return ManifestEntry{}, errors.New("no entries in manifest")
	}
	return m.Entries[0], nil
}

func (m Manifest) After(e ManifestEntry) ([]ManifestEntry, error) {
	indexed, exists := m.m[e.id()]
	if !exists {
		return nil, fmt.Errorf("entry not found in manifest: %v key %s", e, e.id())
	}
	var entries []ManifestEntry
	// reverse order index
	for i := len(m.Entries) - 1; i > indexed.index; i-- {
		entries = append(entries, m.Entries[i])
	}
	return entries, nil

}

func (m Manifest) Before(e ManifestEntry) ([]ManifestEntry, error) {
	var entries []ManifestEntry
	for _, entry := range m.Entries {
		if entry.Version < e.Version {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

type ManifestEntry struct {
	Type          string // core, plugin
	PluginKind    string // relayer, app
	PluginVariant string // evm, optimism, arbitrum, functions, ccip
	Version       int

	index int    // 0 ==> most recent
	path  string // migration path

}

func (m ManifestEntry) id() string {
	if m.Type == "core" {
		return fmt.Sprintf("%s_%d", "core", m.Version)
	}
	return fmt.Sprintf("%s_%s_%s_%d", "plugin", m.PluginKind, m.PluginVariant, m.Version)
}

func validateMigrationEntry(m ManifestEntry) error {
	if m.Version == 0 {
		return fmt.Errorf("missing version")
	}
	if m.Type != "core" && m.Type != "plugin" {
		return fmt.Errorf("unknown migration type: %s", m.Type)
	}
	if m.Type == "core" {
		if m.PluginKind != "" || m.PluginVariant != "" {
			return fmt.Errorf("core migration: expected empty plugin configruation but got plugin kind '%s', variant '%s'", m.PluginKind, m.PluginVariant)
		}
	}
	if m.Type == "plugin" {
		if m.PluginKind != "relayer" && m.PluginKind != "app" {
			return fmt.Errorf("unknown plugin kind: %s", m.PluginKind)
		}
		if m.PluginVariant == "" {
			return fmt.Errorf("missing plugin variant")
		}
	}
	return nil
}

func LoadManifest(txt string) (Manifest, error) {
	lines := strings.Split(txt, "\n")
	var m Manifest
	m.m = make(map[string]ManifestEntry, len(lines))
	for i, l := range lines {
		if l == "" {
			continue
		}
		e, err := parseEntry(l)
		if err != nil {
			return Manifest{}, fmt.Errorf("failed to parse line %s: %w", l, err)
		}
		e.index = i
		m.Entries = append(m.Entries, e)
		m.m[e.id()] = e
	}
	return m, nil
}

var (
	errInvalidManifestEntryName   = fmt.Errorf("invalid migration name")
	errInvalidPluginManifestEntry = fmt.Errorf("invalid plugin migration path")
)

func parseEntry(path string) (e ManifestEntry, err error) {
	p := strings.TrimPrefix(path, "core/store/migrate/")
	e, err = parseCoreEntry(p)
	if err != nil {
		var err2 error
		e, err2 = parsePluginEntry(p)
		if err2 != nil {
			return e, errors.Join(fmt.Errorf("failed to parse path '%s' into entry", path), err, err2)
		}
	}
	return e, validateMigrationEntry(e)
}

func parseCoreEntry(path string) (ManifestEntry, error) {
	version, err := extractVersion(filepath.Base(path))
	if err != nil {
		return ManifestEntry{}, fmt.Errorf("failed to extract version for %s: %w", path, err)
	}
	parts := strings.Split(path, "/")
	path = strings.TrimPrefix(path, "core/store/migrate/")
	if len(parts) != 2 {
		return ManifestEntry{}, fmt.Errorf("invalid core migration path: %s", path)
	}
	return ManifestEntry{
		path:    path,
		Type:    "core",
		Version: version,
	}, nil
}

func parsePluginEntry(path string) (ManifestEntry, error) {
	version, err := extractVersion(filepath.Base(path))
	if err != nil {
		return ManifestEntry{}, fmt.Errorf("failed to extract version for %s: %w", path, err)
	}
	path = strings.TrimPrefix(path, "core/store/migrate/")

	// plugins/<kind>/<variant>/<version>_<name>.sql
	parts := strings.Split(path, "/")
	if len(parts) != 4 {
		return ManifestEntry{}, fmt.Errorf("invalid plugin migration path: %s", path)
	}
	return ManifestEntry{
		path:          path,
		Type:          "plugin",
		PluginKind:    parts[1],
		PluginVariant: parts[2],
		Version:       version,
	}, nil
}

func extractVersion(migrationName string) (int, error) {
	if migrationName == "" {
		return 0, fmt.Errorf("%w: empty migration name", errInvalidManifestEntryName)
	}
	parts := strings.Split(migrationName, "_")
	if len(parts) < 2 {
		return 0, fmt.Errorf("%w: %s", errInvalidManifestEntryName, migrationName)
	}
	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("%w: could not parse version: %s", errInvalidManifestEntryName, migrationName)
	}
	return version, nil
}
