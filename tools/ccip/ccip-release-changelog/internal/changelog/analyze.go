package changelog

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Report is the full changelog analysis between two core-repo refs.
type Report struct {
	Old, New  DepSnapshot
	Repos     []RepoReport
	Flags     []string // top-level audit flags
	Generated time.Time
}

// RepoReport is the per-repository analysis.
type RepoReport struct {
	Config       RepoConfig
	Old, New     repoPin
	Status       string // ahead / behind / diverged / identical / "" when unknown
	TotalInRange int    // commits in range before path filtering
	Commits      []CommitEntry
	KeywordHits  []CommitEntry
	Notes        []string // low-severity observations (rendered in repo section)
	Err          string   // set when the commit changelog could not be produced
	Truncated    bool     // compare API commit list was capped
}

// keywordPattern flags commit titles interesting to a release audit.
var keywordPattern = regexp.MustCompile(`(?i)\b(breaking|revert|hotfix|security|config)\b|fix!`)

// Analyze computes the full report between two core-repo refs.
// gh may be nil if only local processing is desired (external repos will
// record errors instead of commit logs).
func Analyze(ctx context.Context, g gitRunner, gh *ghClient, oldRef, newRef string) (*Report, error) {
	oldSnap, err := LoadSnapshot(ctx, g, oldRef)
	if err != nil {
		return nil, fmt.Errorf("loading old snapshot: %w", err)
	}
	newSnap, err := LoadSnapshot(ctx, g, newRef)
	if err != nil {
		return nil, fmt.Errorf("loading new snapshot: %w", err)
	}

	rep := &Report{Old: oldSnap, New: newSnap, Generated: time.Now().UTC()}

	for _, cfg := range TrackedRepos {
		rr := analyzeRepo(ctx, g, gh, cfg, oldSnap, newSnap)
		rep.Repos = append(rep.Repos, rr)
		rep.Flags = append(rep.Flags, repoFlags(rr, oldSnap, newSnap)...)
	}
	return rep, nil
}

func analyzeRepo(ctx context.Context, g gitRunner, gh *ghClient, cfg RepoConfig, oldSnap, newSnap DepSnapshot) RepoReport {
	rr := RepoReport{Config: cfg, Old: pinFor(cfg, oldSnap), New: pinFor(cfg, newSnap)}

	if cfg.Local {
		analyzeLocal(ctx, g, &rr, oldSnap.SHA, newSnap.SHA)
	} else {
		analyzeRemote(ctx, gh, &rr)
	}

	for _, c := range rr.Commits {
		if keywordPattern.MatchString(c.Title) {
			rr.KeywordHits = append(rr.KeywordHits, c)
		}
	}
	rr.Notes = append(rr.Notes, dedupeDivergenceNotes(cfg, rr.Old, rr.New)...)
	return rr
}

// dedupeDivergenceNotes emits divergence notes for both refs, collapsing to a
// single "(both refs)" note when the divergence is identical at each end.
func dedupeDivergenceNotes(cfg RepoConfig, oldPin, newPin repoPin) []string {
	oldNotes := divergenceNotes(cfg, oldPin, "old")
	newNotes := divergenceNotes(cfg, newPin, "new")
	if len(oldNotes) == 1 && len(newNotes) == 1 {
		oldBody := strings.TrimPrefix(oldNotes[0], "divergent pins (old ref): ")
		newBody := strings.TrimPrefix(newNotes[0], "divergent pins (new ref): ")
		if oldBody == newBody {
			return []string{"divergent pins (both refs): " + oldBody}
		}
	}
	return append(oldNotes, newNotes...)
}

// analyzeLocal produces the commit changelog for the core repo from the local
// checkout, applying path filters.
func analyzeLocal(ctx context.Context, g gitRunner, rr *RepoReport, oldSHA, newSHA string) {
	rr.Old.PrimarySHA, rr.New.PrimarySHA = oldSHA, newSHA
	if oldSHA == newSHA {
		rr.Status = "identical"
		return
	}
	switch {
	case g.IsAncestor(ctx, oldSHA, newSHA):
		rr.Status = "ahead"
	case g.IsAncestor(ctx, newSHA, oldSHA):
		rr.Status = "behind" // new ref is strictly older: genuine rollback
	default:
		// Release branches diverge from each other by design; unlike an
		// external pin, this is normal for the core repo. LogRange still
		// yields exactly "what's in new that wasn't in old".
		rr.Status = "ahead"
		rr.Notes = append(rr.Notes, "refs are not in direct ancestry (e.g. different release lines); listing commits reachable from new but not old")
	}
	all, err := g.LogRange(ctx, oldSHA, newSHA)
	if err != nil {
		rr.Err = fmt.Sprintf("git log failed: %v", err)
		return
	}
	rr.TotalInRange = len(all)
	for _, c := range all {
		if len(rr.Config.IncludePaths) > 0 || len(rr.Config.ExcludePaths) > 0 {
			files, err := g.CommitFiles(ctx, c.SHA)
			if err != nil {
				rr.Notes = append(rr.Notes, fmt.Sprintf("could not list files of %s: %v", shortSHA(c.SHA), err))
				continue
			}
			if !pathMatch(files, rr.Config.IncludePaths, rr.Config.ExcludePaths) {
				continue
			}
		}
		title, pr := parseTitle(c.Title)
		author := c.AuthorName
		if login := noreplyLogin(c.AuthorEmail); login != "" {
			author = "@" + login
		}
		rr.Commits = append(rr.Commits, CommitEntry{SHA: c.SHA, Title: title, PR: pr, Author: author})
	}
}

// noreplyLogin extracts a GitHub login from a noreply email like
// "12345+octocat@users.noreply.github.com".
func noreplyLogin(email string) string {
	const suffix = "@users.noreply.github.com"
	if !strings.HasSuffix(email, suffix) {
		return ""
	}
	local := strings.TrimSuffix(email, suffix)
	if i := strings.LastIndex(local, "+"); i >= 0 {
		return local[i+1:]
	}
	return local
}

// pathMatch applies include/exclude prefix filters to a commit's file list.
func pathMatch(files, includes, excludes []string) bool {
	kept := files[:0:0]
	for _, f := range files {
		excluded := false
		for _, ex := range excludes {
			if strings.HasPrefix(f, ex) {
				excluded = true
				break
			}
		}
		if !excluded {
			kept = append(kept, f)
		}
	}
	if len(includes) == 0 {
		return len(kept) > 0
	}
	for _, f := range kept {
		for _, in := range includes {
			if strings.HasPrefix(f, in) {
				return true
			}
		}
	}
	return false
}

// analyzeRemote produces the commit changelog for an external repo via the
// GitHub compare API.
func analyzeRemote(ctx context.Context, gh *ghClient, rr *RepoReport) {
	oldSHA, newSHA := rr.Old.PrimarySHA, rr.New.PrimarySHA
	switch {
	case rr.Old.PrimaryVersion == "" && rr.New.PrimaryVersion == "":
		rr.Err = "no pin found at either ref"
		return
	case oldSHA == "" || newSHA == "":
		if rr.Old.PrimaryVersion == rr.New.PrimaryVersion {
			rr.Status = "identical"
			return
		}
		rr.Err = fmt.Sprintf("could not extract commit SHA from version(s) %q / %q",
			rr.Old.PrimaryVersion, rr.New.PrimaryVersion)
		return
	case oldSHA == newSHA:
		rr.Status = "identical"
		return
	}
	if gh == nil {
		rr.Err = "no GitHub client configured"
		return
	}
	res, err := gh.Compare(ctx, rr.Config.Owner, rr.Config.Name, oldSHA, newSHA)
	if err != nil {
		rr.Err = err.Error()
		return
	}
	rr.Status = res.Status
	// The compare API returns commits oldest-first; render newest-first.
	for i, j := 0, len(res.Commits)-1; i < j; i, j = i+1, j-1 {
		res.Commits[i], res.Commits[j] = res.Commits[j], res.Commits[i]
	}
	rr.Commits = res.Commits
	rr.TotalInRange = res.TotalCommits
	if rr.TotalInRange == 0 {
		rr.TotalInRange = res.AheadBy
	}
	rr.Truncated = res.TotalCommits > len(res.Commits)
}

// divergenceNotes reports when a repo's various pins at one ref point at
// different commits (multi-module divergence).
func divergenceNotes(cfg RepoConfig, pin repoPin, side string) []string {
	seen := map[string]string{} // sha -> label
	var order []string
	add := func(label, version string) {
		sha := VersionSHA(version)
		if sha == "" {
			return
		}
		if _, ok := seen[sha]; !ok {
			seen[sha] = label
			order = append(order, sha)
		}
	}
	for _, m := range cfg.GoModules {
		if v, ok := pin.ModuleVersions[m]; ok {
			add(strings.TrimPrefix(m, "github.com/smartcontractkit/"), v)
		}
	}
	for _, k := range cfg.PluginKeys {
		if p, ok := pin.PluginRefs[k]; ok {
			add("plugin:"+k, p.GitRef)
		}
	}
	if len(order) <= 1 {
		return nil
	}
	var parts []string
	for _, sha := range order {
		parts = append(parts, fmt.Sprintf("%s at `%s`", seen[sha], shortSHA(sha)))
	}
	return []string{fmt.Sprintf("divergent pins (%s ref): %s", side, strings.Join(parts, "; "))}
}

// repoFlags computes the top-level audit flags for one repo.
func repoFlags(rr RepoReport, oldSnap, newSnap DepSnapshot) []string {
	var flags []string
	name := rr.Config.Name

	// Plugin gitRef changes (spec output #4).
	for _, k := range rr.Config.PluginKeys {
		oldP, oldOK := oldSnap.Plugins[k]
		newP, newOK := newSnap.Plugins[k]
		switch {
		case oldOK && !newOK:
			flags = append(flags, fmt.Sprintf("**%s**: plugin `%s` REMOVED from plugins.public.yaml (was `%s`)", name, k, oldP.GitRef))
		case !oldOK && newOK:
			flags = append(flags, fmt.Sprintf("**%s**: plugin `%s` ADDED to plugins.public.yaml (`%s`)", name, k, newP.GitRef))
		case oldOK && oldP.GitRef != newP.GitRef:
			flags = append(flags, fmt.Sprintf("**%s**: plugin `%s` gitRef changed `%s` → `%s`", name, k, oldP.GitRef, newP.GitRef))
		}
	}

	// go.mod module added/removed.
	for _, m := range rr.Config.GoModules {
		oldV, oldOK := oldSnap.Modules[m]
		newV, newOK := newSnap.Modules[m]
		switch {
		case oldOK && !newOK:
			flags = append(flags, fmt.Sprintf("**%s**: go.mod module `%s` REMOVED (was `%s`)", name, shortModule(m), oldV))
		case !oldOK && newOK:
			flags = append(flags, fmt.Sprintf("**%s**: go.mod module `%s` ADDED (`%s`)", name, shortModule(m), newV))
		}
	}

	// Dual-source drift: plugin moduleURI matches a tracked go.mod module and
	// their SHAs disagree at the same ref.
	for _, k := range rr.Config.PluginKeys {
		for side, snap := range map[string]DepSnapshot{"old": oldSnap, "new": newSnap} {
			p, ok := snap.Plugins[k]
			if !ok {
				continue
			}
			if v, ok := snap.Modules[p.ModuleURI]; ok {
				if ps, ms := VersionSHA(p.GitRef), VersionSHA(v); ps != "" && ms != "" && ps != ms {
					flags = append(flags, fmt.Sprintf("**%s**: DRIFT at %s ref — plugin `%s` pins `%s` but go.mod `%s` pins `%s`",
						name, side, k, shortSHA(ps), shortModule(p.ModuleURI), shortSHA(ms)))
				}
			}
		}
	}

	// Rollback / divergence.
	switch rr.Status {
	case "behind":
		flags = append(flags, fmt.Sprintf("**%s**: ROLLBACK — new pin `%s` is BEHIND old pin `%s`", name, shortSHA(rr.New.PrimarySHA), shortSHA(rr.Old.PrimarySHA)))
	case "diverged":
		flags = append(flags, fmt.Sprintf("**%s**: DIVERGED — old and new pins share no direct ancestry (`%s` / `%s`)", name, shortSHA(rr.Old.PrimarySHA), shortSHA(rr.New.PrimarySHA)))
	}

	// Keyword callouts.
	for _, c := range rr.KeywordHits {
		flags = append(flags, fmt.Sprintf("**%s**: keyword match — %s ([`%s`](https://github.com/%s/%s/commit/%s))",
			name, c.Title, shortSHA(c.SHA), rr.Config.Owner, rr.Config.Name, c.SHA))
	}
	return flags
}

func shortModule(m string) string {
	return strings.TrimPrefix(m, "github.com/smartcontractkit/")
}
