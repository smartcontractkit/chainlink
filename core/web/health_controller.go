package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type HealthController struct {
	App chainlink.Application
}

const (
	HealthStatusPassing = "passing"
	HealthStatusFailing = "failing"
)

// NOTE: We only implement the k8s readiness check, *not* the liveness check. Liveness checks are only recommended in cases
// where the app doesn't crash itself on panic, and if implemented incorrectly can cause cascading failures.
// See the following for more information:
// - https://srcco.de/posts/kubernetes-liveness-probes-are-dangerous.html
func (hc *HealthController) Readyz(c *gin.Context) {
	status := http.StatusOK

	checker := hc.App.GetHealthChecker()

	ready, errors := checker.IsReady()

	if !ready {
		status = http.StatusServiceUnavailable
	}

	c.Status(status)

	if _, ok := c.GetQuery("full"); !ok {
		return
	}

	checks := make([]presenters.Check, 0, len(errors))

	for name, err := range errors {
		status := HealthStatusPassing
		var output string

		if err != nil {
			status = HealthStatusFailing
			output = err.Error()
		}

		checks = append(checks, presenters.Check{
			JAID:   presenters.NewJAID(name),
			Name:   name,
			Status: status,
			Output: output,
		})
	}

	// return a json description of all the checks
	jsonAPIResponse(c, checks, "checks")
}

func (hc *HealthController) Health(c *gin.Context) {
	_, failing := c.GetQuery("failing")

	status := http.StatusOK

	checker := hc.App.GetHealthChecker()

	healthy, errors := checker.IsHealthy()

	if !healthy {
		status = http.StatusMultiStatus
	}

	c.Status(status)

	checks := make([]presenters.Check, 0, len(errors))
	for name, err := range errors {
		status := HealthStatusPassing
		var output string

		if err != nil {
			status = HealthStatusFailing
			output = err.Error()
		} else if failing {
			continue // omit from returned data
		}

		checks = append(checks, presenters.Check{
			JAID:   presenters.NewJAID(name),
			Name:   name,
			Status: status,
			Output: output,
		})
	}

	switch c.NegotiateFormat(gin.MIMEJSON, gin.MIMEHTML, gin.MIMEPlain) {
	case gin.MIMEJSON:
		break // default

	case gin.MIMEHTML:
		if err := newCheckTree(checks).WriteHTMLTo(c.Writer); err != nil {
			hc.App.GetLogger().Errorw("Failed to write HTML health report", "err", err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return

	case gin.MIMEPlain:
		if err := writeTextTo(c.Writer, checks); err != nil {
			hc.App.GetLogger().Errorw("Failed to write plaintext health report", "err", err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	slices.SortFunc(checks, presenters.CmpCheckName)
	jsonAPIResponseWithStatus(c, checks, "checks", status)
}

func writeTextTo(w io.Writer, checks []presenters.Check) error {
	slices.SortFunc(checks, presenters.CmpCheckName)
	for _, ch := range checks {
		status := "?  "
		switch ch.Status {
		case HealthStatusPassing:
			status = "ok "
		case HealthStatusFailing:
			status = "!  "
		}
		if _, err := fmt.Fprintf(w, "%s%s\n", status, ch.Name); err != nil {
			return err
		}
		if ch.Output != "" {
			if _, err := fmt.Fprintf(newLinePrefixWriter(w, "\t"), "\t%s", ch.Output); err != nil {
				return err
			}
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
	}
	return nil
}

type checkNode struct {
	Name   string // full
	Status string
	Output string

	Subs checkTree
}

type checkTree map[string]checkNode

func newCheckTree(checks []presenters.Check) checkTree {
	slices.SortFunc(checks, presenters.CmpCheckName)
	root := make(checkTree)
	for _, c := range checks {
		parts := strings.Split(c.Name, ".")
		node := root
		for _, short := range parts[:len(parts)-1] {
			n, ok := node[short]
			if !ok {
				n = checkNode{Subs: make(checkTree)}
				node[short] = n
			}
			node = n.Subs
		}
		p := parts[len(parts)-1]
		node[p] = checkNode{
			Name:   c.Name,
			Status: c.Status,
			Output: c.Output,
			Subs:   make(checkTree),
		}
	}
	return root
}

func (t checkTree) WriteHTMLTo(w io.Writer) error {
	if _, err := io.WriteString(w, `<style>
    details {
        margin: 0.0em 0.0em 0.0em 0.4em;
        padding: 0.3em 0.0em 0.0em 0.4em;
    }
    pre {
        margin-left:1em;
        margin-top: 0;
    }
    summary {
        padding-bottom: 0.4em;
    }
    details {
        border: thin solid black;
        border-bottom-color: rgba(0,0,0,0);
        border-right-color: rgba(0,0,0,0);
    }
    .passing:after {
        color: blue;
        content: " - (Passing)";
        font-size:small;
        text-transform: uppercase;
    }
    .failing:after {
        color: red;
        content: " - (Failing)";
        font-weight: bold;
        font-size:small;
        text-transform: uppercase;
    }
    summary.noexpand::marker {
        color: rgba(100,101,10,0);
    }
</style>`); err != nil {
		return err
	}
	return t.writeHTMLTo(newLinePrefixWriter(w, ""))
}

func (t checkTree) writeHTMLTo(w *linePrefixWriter) error {
	keys := maps.Keys(t)
	slices.Sort(keys)
	for _, short := range keys {
		node := t[short]
		if _, err := io.WriteString(w, `
<details open>`); err != nil {
			return err
		}
		var expand string
		if node.Output == "" && len(node.Subs) == 0 {
			expand = ` class="noexpand"`
		}
		if _, err := fmt.Fprintf(w, `
    <summary title="%s"%s><span class="%s">%s</span></summary>`, node.Name, expand, node.Status, short); err != nil {
			return err
		}
		if node.Output != "" {
			if _, err := w.WriteRawLinef("    <pre>%s</pre>", node.Output); err != nil {
				return err
			}
		}
		if len(node.Subs) > 0 {
			if err := node.Subs.writeHTMLTo(w.new("    ")); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, "\n</details>"); err != nil {
			return err
		}
	}
	return nil
}

type linePrefixWriter struct {
	w       io.Writer
	prefix  string
	prefixB []byte
}

func newLinePrefixWriter(w io.Writer, prefix string) *linePrefixWriter {
	prefix = "\n" + prefix
	return &linePrefixWriter{w: w, prefix: prefix, prefixB: []byte(prefix)}
}

func (w *linePrefixWriter) new(prefix string) *linePrefixWriter {
	prefix = w.prefix + prefix
	return &linePrefixWriter{w: w.w, prefix: prefix, prefixB: []byte(prefix)}
}

func (w *linePrefixWriter) Write(b []byte) (int, error) {
	return w.w.Write(bytes.ReplaceAll(b, []byte("\n"), w.prefixB))
}

func (w *linePrefixWriter) WriteString(s string) (n int, err error) {
	return io.WriteString(w.w, strings.ReplaceAll(s, "\n", w.prefix))
}

// WriteRawLinef writes a new newline with prefix, followed by s without modification.
func (w *linePrefixWriter) WriteRawLinef(s string, args ...any) (n int, err error) {
	return fmt.Fprintf(w.w, w.prefix+s, args...)
}
