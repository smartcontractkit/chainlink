package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func writeManifest(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseLine(t *testing.T) {
	t.Parallel()
	tests := []struct {
		line            string
		wantName        string
		wantVersion     string
		wantOK          bool
	}{
		{"", "", "", false},
		{"# comment", "", "", false},
		{"golang        1.26.4", "golang", "1.26.4", true},
		{"github.com/jmank88/gomods  0.1.7", "github.com/jmank88/gomods", "0.1.7", true},
	}
	for _, tt := range tests {
		name, version, ok := parseLine(tt.line)
		if ok != tt.wantOK || name != tt.wantName || version != tt.wantVersion {
			t.Errorf("parseLine(%q) = (%q, %q, %v), want (%q, %q, %v)",
				tt.line, name, version, ok, tt.wantName, tt.wantVersion, tt.wantOK)
		}
	}
}

func TestLookup(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	tv := writeManifest(t, dir, ".tool-versions", `golang 1.26.4
mockery 2.53.0
protoc 29.3
`)
	gt := writeManifest(t, dir, "go-tools.txt", `# comment
github.com/jmank88/gomods 0.1.7
`)

	store, err := New(tv, gt)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		key     string
		want    string
		wantErr bool
	}{
		{"mockery", "2.53.0", false},
		{"protoc", "29.3", false},
		{"github.com/jmank88/gomods", "0.1.7", false},
		{"missing", "", true},
		{"github.com/missing/tool", "", true},
	}
	for _, tt := range tests {
		got, err := store.Lookup(tt.key)
		if (err != nil) != tt.wantErr {
			t.Errorf("Lookup(%q) error = %v, wantErr %v", tt.key, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("Lookup(%q) = %q, want %q", tt.key, got, tt.want)
		}
	}
}

func TestList(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	tv := writeManifest(t, dir, ".tool-versions", "mockery 2.53.0\n")
	gt := writeManifest(t, dir, "go-tools.txt", "github.com/jmank88/gomods 0.1.7\n")

	store, err := New(tv, gt)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("len(List()) = %d, want 2", len(entries))
	}
	if entries[0].Name != "mockery" || entries[1].Name != "github.com/jmank88/gomods" {
		t.Fatalf("List() = %+v", entries)
	}
}

func TestNewFileNotFound(t *testing.T) {
	t.Parallel()
	_, err := New("/nonexistent/.tool-versions", "/nonexistent/go-tools.txt")
	if err == nil {
		t.Fatal("expected error for missing files")
	}
}
