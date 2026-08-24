package engine

import (
	"fmt"
	"strconv"
	"strings"
)

// refLine renders a report header line for one ref, noting when an image tag
// was normalized to a git tag.
func refLine(snap DepSnapshot) string {
	if snap.ResolvedRef != "" {
		return fmt.Sprintf("`%s` (image tag → git tag `%s`) (`%s`)", snap.Ref, snap.ResolvedRef, snap.SHA)
	}
	return fmt.Sprintf("`%s` (`%s`)", snap.Ref, snap.SHA)
}

// displayName returns the human-facing name of a tracked repo.
func (r RepoReport) displayName() string {
	if r.Config.Local && len(r.Config.IncludePaths) > 0 {
		return fmt.Sprintf("%s (%s)", r.Config.Name, strings.Join(r.Config.IncludePaths, ", "))
	}
	return r.Config.Name
}

// commitURL returns the GitHub commit URL for an entry.
func (r RepoReport) commitURL(sha string) string {
	return fmt.Sprintf("https://github.com/%s/%s/commit/%s", r.Config.Owner, r.Config.Name, sha)
}

// prURL returns the GitHub PR URL for an entry.
func (r RepoReport) prURL(pr int) string {
	return fmt.Sprintf("https://github.com/%s/%s/pull/%d", r.Config.Owner, r.Config.Name, pr)
}

// formatEntry renders one commit line:
// "Title ([#1234](url)) ([`abc123`](url)) by @author"
func (r RepoReport) formatEntry(c CommitEntry) string {
	var b strings.Builder
	b.WriteString("- ")
	if c.Title == "" {
		b.WriteString("(no title)")
	} else {
		b.WriteString(c.Title)
	}
	if c.PR > 0 {
		fmt.Fprintf(&b, " ([#%d](%s))", c.PR, r.prURL(c.PR))
	}
	fmt.Fprintf(&b, " ([`%s`](%s))", shortSHA(c.SHA), r.commitURL(c.SHA))
	if c.Author != "" {
		fmt.Fprintf(&b, " by %s", c.Author)
	}
	return b.String()
}

// commitCountLine summarizes the commit count for a repo section header.
func (r RepoReport) commitCountLine() string {
	filtered := len(r.Config.IncludePaths) > 0 || len(r.Config.ExcludePaths) > 0
	switch {
	case r.Err != "":
		return "⚠️ compare failed"
	case r.Status == "identical":
		return "no changes"
	case r.Status == "behind":
		return "⚠️ rolled back (see flags)"
	case r.Status == "diverged":
		return "⚠️ diverged history (see flags)"
	case filtered && r.TotalInRange > len(r.Commits):
		return fmt.Sprintf("%s touching tracked paths (%d total in range)", pluralize(len(r.Commits)), r.TotalInRange)
	default:
		return pluralize(len(r.Commits))
	}
}

// pluralize renders "1 commit" / "N commits".
func pluralize(n int) string {
	if n == 1 {
		return "1 commit"
	}
	return fmt.Sprintf("%d commits", n)
}

// RenderMarkdown produces the full audit document.
func RenderMarkdown(rep *Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s Release Changelog\n\n", rep.ProductName)
	fmt.Fprintf(&b, "- **Old**: %s\n", refLine(rep.Old))
	fmt.Fprintf(&b, "- **New**: %s\n", refLine(rep.New))
	fmt.Fprintf(&b, "- **Generated**: %s\n", rep.Generated.Format("2006-01-02 15:04 UTC"))
	fmt.Fprintf(&b, "- **Core changelog**: [CHANGELOG.md at %s](https://github.com/smartcontractkit/chainlink/blob/%s/CHANGELOG.md)\n\n",
		rep.New.Ref, rep.New.SHA)

	// Flags
	b.WriteString("## ⚠️ Flags\n\n")
	if len(rep.Flags) == 0 {
		b.WriteString("No risk flags for this range. ✅\n\n")
	} else {
		for _, f := range rep.Flags {
			fmt.Fprintf(&b, "- %s\n", f)
		}
		b.WriteString("\n")
	}

	// go.mod changes
	b.WriteString("## go.mod changes (CCIP modules)\n\n")
	gomodAny := false
	for _, rr := range rep.Repos {
		if len(rr.Config.GoModules) == 0 {
			continue
		}
		gomodAny = true
		fmt.Fprintf(&b, "### %s\n\n", rr.Config.Name)
		for _, m := range rr.Config.GoModules {
			oldV, oldOK := rep.Old.Modules[m]
			newV, newOK := rep.New.Modules[m]
			fmt.Fprintf(&b, "- `%s`: %s\n", shortModule(m), versionTransition(oldV, oldOK, newV, newOK))
		}
		b.WriteString("\n")
	}
	if !gomodAny {
		b.WriteString("No tracked go.mod modules.\n\n")
	}

	// plugins.public.yaml changes
	b.WriteString("## plugins.public.yaml changes (CCIP plugins)\n\n")
	pluginsAny := false
	for _, rr := range rep.Repos {
		for _, k := range rr.Config.PluginKeys {
			pluginsAny = true
			oldP, oldOK := rep.Old.Plugins[k]
			newP, newOK := rep.New.Plugins[k]
			var oldV, newV string
			if oldOK {
				oldV = oldP.GitRef
			}
			if newOK {
				newV = newP.GitRef
			}
			fmt.Fprintf(&b, "- **%s** (`%s`): %s\n", k, rr.Config.Name, versionTransition(oldV, oldOK, newV, newOK))
		}
	}
	if !pluginsAny {
		b.WriteString("No tracked plugins.\n")
	}
	b.WriteString("\n")

	// Commit changelogs
	b.WriteString("## Commit changelogs\n\n")
	for _, rr := range rep.Repos {
		fmt.Fprintf(&b, "### %s — %s\n\n", rr.displayName(), rr.commitCountLine())
		if rr.Err != "" {
			fmt.Fprintf(&b, "⚠️ %s\n\n", rr.Err)
			continue
		}
		for _, note := range rr.Notes {
			fmt.Fprintf(&b, "> %s\n", note)
		}
		if len(rr.Notes) > 0 {
			b.WriteString("\n")
		}
		if rr.Truncated {
			fmt.Fprintf(&b, "> ⚠️ commit list truncated by GitHub API (showing %d of %d)\n\n", len(rr.Commits), rr.TotalInRange)
		}
		if rr.Status != "identical" {
			for _, c := range rr.Commits {
				b.WriteString(rr.formatEntry(c))
				b.WriteString("\n")
			}
			if len(rr.Commits) > 0 {
				b.WriteString("\n")
			}
		}
	}
	return b.String()
}

// RenderSlackSummary produces the compact Slack message (mrkdwn format).
func RenderSlackSummary(rep *Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "*%s Release Changelog* `%s` → `%s`\n", rep.ProductName, rep.Old.Ref, rep.New.Ref)
	fmt.Fprintf(&b, "<https://github.com/smartcontractkit/chainlink/compare/%s...%s|core compare> · ", shortSHA(rep.Old.SHA), shortSHA(rep.New.SHA))
	fmt.Fprintf(&b, "<https://github.com/smartcontractkit/chainlink/blob/%s/CHANGELOG.md|core changelog>\n\n", rep.New.SHA)

	if len(rep.Flags) == 0 {
		b.WriteString(":white_check_mark: *No risk flags for this range.*\n\n")
	} else {
		fmt.Fprintf(&b, ":warning: *Flags (%d)*\n", len(rep.Flags))
		for _, f := range rep.Flags {
			// Slack mrkdwn uses *bold*, not **bold**.
			fmt.Fprintf(&b, "• %s\n", strings.ReplaceAll(f, "**", "*"))
		}
		b.WriteString("\n")
	}

	b.WriteString("*Commits per repo*\n")
	var parts []string
	for _, rr := range rep.Repos {
		switch {
		case rr.Err != "":
			parts = append(parts, rr.Config.Name+": ⚠️ error")
		case rr.Status == "identical":
			parts = append(parts, rr.Config.Name+": no change")
		default:
			parts = append(parts, rr.Config.Name+": "+strconv.Itoa(len(rr.Commits)))
		}
	}
	b.WriteString(strings.Join(parts, " · "))
	b.WriteString("\n\nFull report attached below :point_down:")
	return b.String()
}

// versionTransition renders an old → new version transition.
func versionTransition(oldV string, oldOK bool, newV string, newOK bool) string {
	switch {
	case !oldOK && !newOK:
		return "not present at either ref"
	case !oldOK:
		return fmt.Sprintf("_(added)_ → `%s`", newV)
	case !newOK:
		return fmt.Sprintf("`%s` → _(removed)_", oldV)
	case oldV == newV:
		return fmt.Sprintf("no change (`%s`)", newV)
	default:
		return fmt.Sprintf("`%s` → `%s`", oldV, newV)
	}
}
