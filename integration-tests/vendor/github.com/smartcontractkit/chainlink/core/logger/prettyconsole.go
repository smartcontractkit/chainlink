package logger

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

var levelColors = map[string]func(...interface{}) string{
	"default": color.New(color.FgWhite).SprintFunc(),
	"debug":   color.New(color.FgGreen).SprintFunc(),
	"info":    color.New(color.FgWhite).SprintFunc(),
	"warn":    color.New(color.FgYellow).SprintFunc(),
	"error":   color.New(color.FgRed).SprintFunc(),
	"panic":   color.New(color.FgHiRed).SprintFunc(),
	"crit":    color.New(color.FgHiRed).SprintFunc(),
	"fatal":   color.New(color.FgHiRed).SprintFunc(),
}

var blue = color.New(color.FgBlue).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

// PrettyConsole wraps a Sink (Writer, Syncer, Closer), usually stdout, and
// formats the incoming json bytes with colors and white space for readability
// before passing on to the underlying Writer in Sink.
type PrettyConsole struct {
	zap.Sink
}

// Write reformats the incoming json bytes with colors, newlines and whitespace
// for better readability in console.
func (pc PrettyConsole) Write(b []byte) (int, error) {
	if !gjson.ValidBytes(b) {
		return 0, fmt.Errorf("unable to parse json for pretty console: %s", string(b))
	}
	js := gjson.ParseBytes(b)
	headline := generateHeadline(js)
	details := generateDetails(js)
	return pc.Sink.Write([]byte(fmt.Sprintln(headline, details)))
}

// Close is overridden to prevent accidental closure of stderr/stdout
func (pc PrettyConsole) Close() error {
	switch pc.Sink {
	case os.Stderr, os.Stdout:
		// Never close Stderr/Stdout because this will break any future go runtime logging from panics etc
		return nil
	default:
		return pc.Sink.Close()
	}
}

func generateHeadline(js gjson.Result) string {
	ts := js.Get("ts")
	var tsStr string
	if f := ts.Float(); f > 1 {
		sec, dec := math.Modf(f)
		tsStr = iso8601UTC(time.Unix(int64(sec), int64(dec*(1e9))))
	} else {
		// assume already formatted
		tsStr = ts.Str
	}
	headline := []interface{}{
		tsStr,
		" ",
		coloredLevel(js.Get("level")),
		fmt.Sprintf("%-50s", js.Get("msg")),
		" ",
		fmt.Sprintf("%-32s", blue(js.Get("caller"))),
	}
	return fmt.Sprint(headline...)
}

// detailsBlacklist of keys to show in details. This does not
// exclude it from being present in other logger sinks, like .jsonl files.
var detailsBlacklist = map[string]bool{
	"level":  true,
	"ts":     true,
	"msg":    true,
	"caller": true,
	"hash":   true,
}

func generateDetails(js gjson.Result) string {
	data := js.Map()
	keys := []string{}

	for k := range data {
		if detailsBlacklist[k] || len(data[k].String()) == 0 {
			continue
		}
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var details strings.Builder

	for _, v := range keys {
		details.WriteString(fmt.Sprintf("%s=%v ", green(v), data[v]))
	}

	return details.String()
}

func coloredLevel(level gjson.Result) string {
	color, ok := levelColors[level.String()]
	if !ok {
		color = levelColors["default"]
	}
	return color(fmt.Sprintf("%-8s", fmt.Sprint("[", strings.ToUpper(level.String()), "]")))
}

// iso8601UTC formats given time to ISO8601.
func iso8601UTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func prettyConsoleSink(s zap.Sink) func(*url.URL) (zap.Sink, error) {
	return func(*url.URL) (zap.Sink, error) {
		return PrettyConsole{s}, nil
	}
}
