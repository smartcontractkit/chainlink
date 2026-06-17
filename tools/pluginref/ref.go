package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type pluginsFile struct {
	Plugins map[string][]pluginEntry `yaml:"plugins"`
}

type pluginEntry struct {
	ModuleURI string `yaml:"moduleURI"`
	GitRef    string `yaml:"gitRef"`
}

var (
	pseudoWithSHARe = regexp.MustCompile(
		`^v\d+\.\d+\.\d+(?:-\d{14}|-(?:0|[0-9A-Za-z-]+\.0)\.\d{14})-([0-9a-f]{7,40})$`,
	)
	plainTagRe     = regexp.MustCompile(`^v\d+\.\d+\.\d+([.-].*)?$`)
	prefixedTagRe  = regexp.MustCompile(`^(.+?)/+(v\d+\.\d+\.\d+(?:[.-].*)?)$`)
	shaOnlyRe      = regexp.MustCompile(`^[0-9a-f]{7,40}$`)
	pseudoMiddleRe = regexp.MustCompile(`^(?:\d{14}|(?:0|[0-9A-Za-z-]+\.0)\.\d{14})$`)
)

type moduleVersion struct {
	Raw       string
	SHA       string
	Tag       string
	TagPrefix string
}

func normalizeVersion(raw string) moduleVersion {
	mv := moduleVersion{Raw: raw}
	low := strings.ToLower(strings.TrimSpace(raw))

	if m := pseudoWithSHARe.FindStringSubmatch(low); m != nil {
		mv.SHA = m[1]
		return mv
	}

	if strings.HasPrefix(low, "v") && strings.Count(low, "-") == 2 {
		parts := strings.Split(low, "-")
		middle := parts[1]
		last := strings.TrimPrefix(parts[2], "g")
		if pseudoMiddleRe.MatchString(middle) && shaOnlyRe.MatchString(last) {
			mv.SHA = last
			return mv
		}
	}

	if shaOnlyRe.MatchString(low) {
		mv.SHA = low
		return mv
	}

	if m := prefixedTagRe.FindStringSubmatch(low); len(m) == 3 && plainTagRe.MatchString(m[2]) {
		orig := strings.TrimSpace(raw)
		if pos := strings.LastIndex(orig, "/"); pos >= 0 && pos+1 < len(orig) {
			return moduleVersion{
				Raw:       raw,
				Tag:       orig[pos+1:],
				TagPrefix: strings.TrimSuffix(orig[:pos], "/"),
			}
		}
	}

	if plainTagRe.MatchString(low) {
		mv.Tag = raw
		return mv
	}

	return mv
}

// gitRefForPlugin returns the gitRef for a named plugin in a plugins YAML file.
func gitRefForPlugin(path, name string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read plugin file %s: %w", path, err)
	}

	var pf pluginsFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return "", fmt.Errorf("parse plugin file %s: %w", path, err)
	}

	entries, ok := pf.Plugins[name]
	if !ok || len(entries) == 0 {
		return "", fmt.Errorf("plugin %q not found in %s", name, path)
	}

	for _, entry := range entries {
		if entry.GitRef != "" {
			return entry.GitRef, nil
		}
	}

	return "", fmt.Errorf("plugin %q has no gitRef in %s", name, path)
}

// commitRefForGitRef returns a git object name suitable for rev-parse.
func commitRefForGitRef(gitRef string) (string, error) {
	gitRef = strings.TrimSpace(gitRef)
	if gitRef == "" {
		return "", errors.New("gitRef is empty")
	}

	mv := normalizeVersion(gitRef)
	if mv.SHA != "" {
		return mv.SHA, nil
	}
	if mv.Tag != "" {
		if mv.TagPrefix != "" {
			return mv.Raw, nil
		}
		return mv.Tag, nil
	}

	// Unrecognized formats (e.g. branch names) pass through for git rev-parse.
	return mv.Raw, nil
}
