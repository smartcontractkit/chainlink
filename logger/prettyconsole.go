package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/tidwall/gjson"
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

type PrettyConsole struct {
	io *os.File
}

func (PrettyConsole) Sync() error  { return nil }
func (PrettyConsole) Close() error { return nil }

func (pc PrettyConsole) Write(b []byte) (int, error) {
	var js models.JSON
	err := json.Unmarshal(b, &js)
	if err != nil {
		log.Panic(err)
	}

	output := []interface{}{
		coloredLevel(js.Get("level")),
		js.Get("msg"),
	}
	return fmt.Println(output...)
}

func coloredLevel(level gjson.Result) string {
	color, ok := levelColors[level.String()]
	if !ok {
		color = levelColors["default"]
	}
	return color(fmt.Sprintf("%-10s", fmt.Sprint("[", strings.ToUpper(level.String()), "]")))
}
