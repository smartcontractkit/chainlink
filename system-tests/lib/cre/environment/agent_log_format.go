package environment

import (
	"encoding/json"
	"fmt"
	"strings"
)

func prettifyAgentLogLine(line string) string {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return ""
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(trimmed), &payload); err != nil {
		return trimmed
	}

	message, _ := payload["message"].(string)
	if message == "" {
		return trimmed
	}

	level, _ := payload["level"].(string)
	if level == "" {
		level = "info"
	}

	cmd, _ := payload["Cmd"].(string)
	if cmd != "" {
		return fmt.Sprintf("[%s] %s (cmd=%s)", level, message, cmd)
	}

	return fmt.Sprintf("[%s] %s", level, message)
}
