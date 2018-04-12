package logger

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

var levelColors = map[string]func(...interface{}) string{
	"default": color.New(color.FgWhite).SprintFunc(),
	"debug":   color.New(color.FgGreen).SprintFunc(),
	"info":    color.New(color.FgWhite).SprintFunc(),
	"warn":    color.New(color.FgYellow).SprintFunc(),
	"error":   color.New(color.FgRed).SprintFunc(),
	"panic":   color.New(color.FgRed).SprintFunc(),
	"fatal":   color.New(color.FgRed).SprintFunc(),
}

var blue = color.New(color.FgBlue).SprintFunc()

type PrettyConsole struct {
	zap.Sink
}

func (pc PrettyConsole) Write(b []byte) (int, error) {
	var js models.JSON
	err := json.Unmarshal(b, &js)
	if err != nil {
		return 0, err
	}

	headline := generateHeadline(js)
	details := generateDetails(js)
	return pc.Sink.Write([]byte(fmt.Sprintln(headline, details)))
}

func generateHeadline(js models.JSON) string {
	sec, dec := math.Modf(js.Get("ts").Float())
	headline := []interface{}{
		utils.ISO8601UTC(time.Unix(int64(sec), int64(dec*(1e9)))),
		" ",
		coloredLevel(js.Get("level")),
		js.Get("msg"),
		" ",
		blue(js.Get("caller")),
	}
	return fmt.Sprint(headline...)
}

// blacklist of key value pairs to show in details. This does not
// exclude it from being present in other logger sinks, like .jsonl files.
var blacklist = map[string]bool{"hash": true}

func generateDetails(js models.JSON) string {
	entries := js.Map()
	delete(entries, "level")
	delete(entries, "ts")
	delete(entries, "msg")
	delete(entries, "caller")

	var details string
	for k, v := range entries {
		if blacklist[k] || len(v.String()) == 0 {
			continue
		}
		details += fmt.Sprintf("%s=%v ", k, v)
	}

	if len(details) > 0 {
		details = fmt.Sprint("\n", details)
	}
	return details
}

func coloredLevel(level gjson.Result) string {
	color, ok := levelColors[level.String()]
	if !ok {
		color = levelColors["default"]
	}
	return color(fmt.Sprintf("%-8s", fmt.Sprint("[", strings.ToUpper(level.String()), "]")))
}
