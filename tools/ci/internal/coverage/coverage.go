// Package coverage computes Go coverage summaries and PR patch coverage.
package coverage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// ReportOptions configures coverage report generation.
type ReportOptions struct {
	Profiles     []string
	Base         string
	Head         string
	Repo         string
	Category     string
	JSONPath     string
	MarkdownPath string
	MinGlobal    float64
	MinPatch     float64
}

// ReportResult contains a generated report and its Markdown rendering.
type ReportResult struct {
	Report   Report
	Markdown string
}

// LineRange is an inclusive range of added lines.
type LineRange struct {
	Start int
	End   int
}

// FilePatch contains patch coverage for one changed file.
type FilePatch struct {
	Path        string
	TestFile    bool
	Added       []LineRange
	Statements  int
	Covered     int
	Missed      int
	MissedSpans []string
}

// Stats contains statement coverage totals.
type Stats struct {
	Statements int
	Covered    int
	Missed     int
	Percent    float64
}

// PatchStats contains coverage totals for changed files.
type PatchStats struct {
	ChangedFiles int
	ChangedLines int
	Statements   int
	Covered      int
	Missed       int
	Percent      float64
	Files        []FilePatch
}

// Report contains global and optional patch coverage.
type Report struct {
	Category         string
	Profile          []string
	Base             string
	Head             string
	Global           Stats
	GlobalSourceOnly Stats
	Patch            *PatchStats
	PatchSourceOnly  *PatchStats
}

type block struct {
	File      string
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
	NumStmt   int
	Count     int64
}

var modulePrefixes = []string{
	"github.com/smartcontractkit/chainlink/v2/",
	"github.com/smartcontractkit/chainlink/",
}

// GenerateReport loads profiles, computes coverage, writes requested outputs,
// and checks the configured thresholds.
func GenerateReport(ctx context.Context, opts ReportOptions) (ReportResult, error) {
	if len(opts.Profiles) == 0 {
		return ReportResult{}, errors.New("--profile is required")
	}

	blocks, err := loadProfiles(opts.Profiles)
	if err != nil {
		return ReportResult{}, err
	}
	rep := Report{
		Category: opts.Category,
		Profile:  opts.Profiles,
		Base:     opts.Base,
		Head:     opts.Head,
	}
	rep.Global = summarize(blocks, func(string) bool { return true })
	rep.GlobalSourceOnly = summarize(blocks, func(path string) bool { return !excluded(path) })

	if opts.Base != "" || opts.Head != "" {
		if opts.Base == "" || opts.Head == "" {
			return ReportResult{}, errors.New("--base and --head must be provided together")
		}
		diff, err := gitDiff(ctx, opts.Repo, opts.Base, opts.Head)
		if err != nil {
			return ReportResult{}, err
		}
		rep.Patch = computePatch(blocks, diff, func(string) bool { return true })
		rep.PatchSourceOnly = computePatch(blocks, diff, func(path string) bool {
			return !excluded(path) && !strings.HasSuffix(path, "_test.go")
		})
	}

	result := ReportResult{Report: rep, Markdown: RenderMarkdown(rep)}
	if opts.MarkdownPath != "" {
		if err := os.WriteFile(opts.MarkdownPath, []byte(result.Markdown), 0o600); err != nil {
			return result, err
		}
	}
	if opts.JSONPath != "" {
		data, err := json.MarshalIndent(rep, "", "  ")
		if err != nil {
			return result, err
		}
		if err := os.WriteFile(opts.JSONPath, append(data, '\n'), 0o600); err != nil {
			return result, err
		}
	}

	var failures []string
	if opts.MinGlobal >= 0 && rep.Global.Percent < opts.MinGlobal {
		failures = append(failures, fmt.Sprintf("global coverage %.2f%% is below %.2f%%", rep.Global.Percent, opts.MinGlobal))
	}
	if opts.MinPatch >= 0 {
		if rep.Patch == nil {
			failures = append(failures, "patch coverage requested but --base/--head not provided")
		} else if rep.Patch.Percent < opts.MinPatch {
			failures = append(failures, fmt.Sprintf("patch coverage %.2f%% is below %.2f%%", rep.Patch.Percent, opts.MinPatch))
		}
	}
	if len(failures) > 0 {
		return result, fmt.Errorf("%s", strings.Join(failures, "; "))
	}
	return result, nil
}

// Merge merges binary Go coverage data directories and writes a text profile.
func Merge(ctx context.Context, covdataDirs []string, out string) error {
	if len(covdataDirs) == 0 {
		return errors.New("--covdata is required")
	}
	for _, dir := range covdataDirs {
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			return fmt.Errorf("covdata directory %q is not accessible", dir)
		}
	}
	tmp, err := os.MkdirTemp("", "coverage-merge-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	merge := exec.CommandContext(ctx, "go", "tool", "covdata", "merge", "-i="+strings.Join(covdataDirs, ","), "-o="+tmp) //nolint:gosec // paths are validated above
	if err := merge.Run(); err != nil {
		return fmt.Errorf("covdata merge: %w", err)
	}
	textfmt := exec.CommandContext(ctx, "go", "tool", "covdata", "textfmt", "-i="+tmp, "-o="+out) //nolint:gosec // output path is caller-provided
	if err := textfmt.Run(); err != nil {
		return fmt.Errorf("covdata textfmt: %w", err)
	}
	return nil
}

// MergeProfiles merges text profiles, retaining the greatest count for each block.
func MergeProfiles(_ context.Context, profiles []string, out string) error {
	if len(profiles) == 0 {
		return errors.New("--profile is required")
	}
	blocks, err := loadProfiles(profiles)
	if err != nil {
		return err
	}
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := fmt.Fprintln(f, "mode: atomic"); err != nil {
		return err
	}
	keys := make([]string, 0, len(blocks))
	for key := range blocks {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		a, b := blocks[keys[i]], blocks[keys[j]]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.StartLine != b.StartLine {
			return a.StartLine < b.StartLine
		}
		if a.StartCol != b.StartCol {
			return a.StartCol < b.StartCol
		}
		if a.EndLine != b.EndLine {
			return a.EndLine < b.EndLine
		}
		return a.EndCol < b.EndCol
	})
	for _, key := range keys {
		b := blocks[key]
		if _, err := fmt.Fprintf(f, "%s:%d.%d,%d.%d %d %d\n", b.File, b.StartLine, b.StartCol, b.EndLine, b.EndCol, b.NumStmt, b.Count); err != nil {
			return err
		}
	}
	return nil
}

func loadProfiles(paths []string) (map[string]block, error) {
	merged := make(map[string]block)
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 64*1024), 1024*1024)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "mode:") {
				continue
			}
			b, err := parseBlock(line)
			if err != nil {
				f.Close()
				return nil, fmt.Errorf("%s: %w", path, err)
			}
			key := blockKey(b)
			if old, ok := merged[key]; !ok || b.Count > old.Count {
				merged[key] = b
			}
		}
		if err := sc.Err(); err != nil {
			f.Close()
			return nil, err
		}
		f.Close()
	}
	return merged, nil
}

func parseBlock(line string) (block, error) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return block{}, fmt.Errorf("malformed coverage line %q", line)
	}
	loc := fields[0]
	filePart, span, ok := strings.CutLast(loc, ":")
	if !ok {
		return block{}, fmt.Errorf("malformed coverage location %q", loc)
	}
	file := normalizePath(filePart)
	startPart, endPart, ok := strings.Cut(span, ",")
	if !ok {
		return block{}, fmt.Errorf("malformed coverage span %q", span)
	}
	start, err := parsePos(startPart)
	if err != nil {
		return block{}, err
	}
	end, err := parsePos(endPart)
	if err != nil {
		return block{}, err
	}
	numStmt, err := strconv.Atoi(fields[1])
	if err != nil {
		return block{}, err
	}
	count, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return block{}, err
	}
	return block{
		File:      file,
		StartLine: start[0],
		StartCol:  start[1],
		EndLine:   end[0],
		EndCol:    end[1],
		NumStmt:   numStmt,
		Count:     count,
	}, nil
}

func parsePos(s string) ([2]int, error) {
	var pos [2]int
	linePart, colPart, ok := strings.Cut(s, ".")
	if !ok {
		return pos, fmt.Errorf("malformed position %q", s)
	}
	line, err := strconv.Atoi(linePart)
	if err != nil {
		return pos, err
	}
	col, err := strconv.Atoi(colPart)
	if err != nil {
		return pos, err
	}
	pos[0], pos[1] = line, col
	return pos, nil
}

func blockKey(b block) string {
	return fmt.Sprintf("%s:%d.%d,%d.%d", b.File, b.StartLine, b.StartCol, b.EndLine, b.EndCol)
}

func normalizePath(path string) string {
	for _, prefix := range modulePrefixes {
		if stripped, ok := strings.CutPrefix(path, prefix); ok {
			path = stripped
			break
		}
	}
	path = strings.TrimPrefix(path, "./")
	return path
}

func excluded(path string) bool {
	base := filepath.Base(path)
	for part := range strings.SplitSeq(path, "/") {
		switch part {
		case "mocks", "testutils", "testhelpers", "cltest", "testdata", "fixtures":
			return true
		}
	}
	if strings.HasPrefix(base, "mock_") || strings.HasSuffix(base, "_mock.go") {
		return true
	}
	for _, suffix := range []string{".pb.go", "_gen.go", "_codecgen.go"} {
		if strings.HasSuffix(base, suffix) {
			return true
		}
	}
	return false
}

func summarize(blocks map[string]block, include func(string) bool) Stats {
	var s Stats
	for _, b := range blocks {
		if !include(b.File) {
			continue
		}
		s.Statements += b.NumStmt
		if b.Count > 0 {
			s.Covered += b.NumStmt
		}
	}
	s.Missed = s.Statements - s.Covered
	s.Percent = percent(s.Covered, s.Statements)
	return s
}

func computePatch(blocks map[string]block, diff map[string][]LineRange, include func(string) bool) *PatchStats {
	ps := &PatchStats{}
	paths := make([]string, 0, len(diff))
	for path := range diff {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	ps.ChangedFiles = len(paths)
	for _, path := range paths {
		ranges := diff[path]
		for _, r := range ranges {
			ps.ChangedLines += r.End - r.Start + 1
		}
		fp := FilePatch{
			Path:     path,
			TestFile: strings.HasSuffix(path, "_test.go"),
			Added:    ranges,
		}
		if include(path) && !fp.TestFile {
			seen := make(map[string]bool)
			missed := make(map[string]bool)
			for key, b := range blocks {
				if b.File != path || seen[key] || !intersects(b, ranges) {
					continue
				}
				seen[key] = true
				fp.Statements += b.NumStmt
				if b.Count > 0 {
					fp.Covered += b.NumStmt
				} else {
					fp.Missed += b.NumStmt
					missed[fmt.Sprintf("%d-%d", b.StartLine, b.EndLine)] = true
				}
			}
			for span := range missed {
				fp.MissedSpans = append(fp.MissedSpans, span)
			}
			sort.Strings(fp.MissedSpans)
			ps.Statements += fp.Statements
			ps.Covered += fp.Covered
			ps.Missed += fp.Missed
		}
		ps.Files = append(ps.Files, fp)
	}
	ps.Percent = percent(ps.Covered, ps.Statements)
	return ps
}

func intersects(b block, ranges []LineRange) bool {
	for _, r := range ranges {
		if b.EndLine >= r.Start && b.StartLine <= r.End {
			return true
		}
	}
	return false
}

func percent(covered, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(int(10000*float64(covered)/float64(total)+0.5)) / 100
}

func gitDiff(ctx context.Context, repo, base, head string) (map[string][]LineRange, error) {
	cmd := exec.CommandContext(ctx, "git", "-c", "core.quotePath=false", "diff", "-U0", "--no-color", "--no-ext-diff", base+".."+head, "--", "*.go") //nolint:gosec // revisions come from the caller
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git diff: %w", err)
	}
	return parseDiff(string(out)), nil
}

func parseDiff(diff string) map[string][]LineRange {
	files := make(map[string][]LineRange)
	var current string
	sc := bufio.NewScanner(strings.NewReader(diff))
	sc.Buffer(make([]byte, 64*1024), 8*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "+++ "):
			path := strings.TrimPrefix(line, "+++ ")
			path = strings.TrimPrefix(path, "b/")
			if path == "/dev/null" {
				current = ""
			} else {
				current = path
			}
		case strings.HasPrefix(line, "@@ ") && current != "":
			if r, ok := parseHunk(line); ok {
				files[current] = append(files[current], r)
			}
		}
	}
	return files
}

func parseHunk(line string) (LineRange, bool) {
	plus := strings.Index(line, " +")
	if plus < 0 {
		return LineRange{}, false
	}
	rest := line[plus+1:]
	end := strings.Index(rest, " @@")
	if end < 0 {
		return LineRange{}, false
	}
	rest = rest[:end]
	startStr := rest
	count := 1
	if comma := strings.Index(rest, ","); comma >= 0 {
		startStr = rest[:comma]
		n, err := strconv.Atoi(rest[comma+1:])
		if err != nil {
			return LineRange{}, false
		}
		count = n
	}
	start, err := strconv.Atoi(startStr)
	if err != nil || count <= 0 || start <= 0 {
		return LineRange{}, false
	}
	return LineRange{Start: start, End: start + count - 1}, true
}

// RenderMarkdown renders a report in the legacy coverage Markdown format.
func RenderMarkdown(rep Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## Coverage — %s\n\n", rep.Category)
	if rep.Base != "" || rep.Head != "" {
		fmt.Fprintf(&b, "`%s..%s`\n\n", rep.Base, rep.Head)
	}
	b.WriteString("### Global\n\n")
	b.WriteString("| Scope | Statements | Covered | Missed | Coverage |\n")
	b.WriteString("|---|---:|---:|---:|---:|\n")
	writeStatsRow(&b, "raw", rep.Global)
	writeStatsRow(&b, "source-only", rep.GlobalSourceOnly)
	if rep.Patch != nil {
		b.WriteString("\n### Patch\n\n")
		b.WriteString("| Scope | Changed files | Changed lines | Statements | Covered | Missed | Coverage |\n")
		b.WriteString("|---|---:|---:|---:|---:|---:|---:|\n")
		writePatchRow(&b, "raw", rep.Patch)
		if rep.PatchSourceOnly != nil {
			writePatchRow(&b, "source-only", rep.PatchSourceOnly)
		}
		if len(rep.Patch.Files) > 0 {
			b.WriteString("\n### Changed files\n\n")
			b.WriteString("| File | Statements | Covered | Missed | Coverage | Missed spans |\n")
			b.WriteString("|---|---:|---:|---:|---:|---|\n")
			for _, f := range rep.Patch.Files {
				if f.TestFile || f.Statements == 0 {
					continue
				}
				fmt.Fprintf(&b, "| `%s` | %d | %d | %d | %.2f%% | %s |\n",
					f.Path, f.Statements, f.Covered, f.Missed, percent(f.Covered, f.Statements), strings.Join(f.MissedSpans, ", "))
			}
		}
	}
	return b.String()
}

func writeStatsRow(b *strings.Builder, name string, s Stats) {
	fmt.Fprintf(b, "| %s | %d | %d | %d | **%.2f%%** |\n", name, s.Statements, s.Covered, s.Missed, s.Percent)
}

func writePatchRow(b *strings.Builder, name string, p *PatchStats) {
	fmt.Fprintf(b, "| %s | %d | %d | %d | %d | %d | **%.2f%%** |\n",
		name, p.ChangedFiles, p.ChangedLines, p.Statements, p.Covered, p.Missed, p.Percent)
}
