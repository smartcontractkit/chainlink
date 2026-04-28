package executable

import (
	"encoding/hex"
	"fmt"
	"strings"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

const (
	maxPayloadDiffBytes       = 64 * 1024
	maxPayloadDiffMatrixCells = 4_000_000
)

func diffPayloads(previousPayload, currentPayload []byte) (string, bool) {
	previous := formatPayloadForDiff(previousPayload)
	current := formatPayloadForDiff(currentPayload)
	if previous == current {
		return "", false
	}

	diff, truncated := lineDiff(previous, current)
	if len(diff) <= maxPayloadDiffBytes {
		return diff, truncated
	}
	return diff[:maxPayloadDiffBytes] + "\n... payload diff truncated ...", true
}

func formatPayloadForDiff(payload []byte) string {
	var req capabilitiespb.CapabilityRequest
	if err := proto.Unmarshal(payload, &req); err != nil {
		return fmt.Sprintf("unable to decode CapabilityRequest payload: %v\nraw_payload_hex:\n%s", err, hex.Dump(payload))
	}

	jsonPayload, jsonErr := protojson.MarshalOptions{
		UseProtoNames: true,
		Multiline:     true,
		Indent:        "  ",
	}.Marshal(&req)
	if jsonErr == nil {
		return string(jsonPayload)
	}

	textPayload, textErr := prototext.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}.Marshal(&req)
	if textErr == nil {
		return string(textPayload)
	}

	return fmt.Sprintf("unable to format CapabilityRequest payload as JSON (%v) or text (%v)\nraw_payload_hex:\n%s",
		jsonErr, textErr, hex.Dump(payload))
}

func lineDiff(previous, current string) (string, bool) {
	previousLines := splitDiffLines(previous)
	currentLines := splitDiffLines(current)
	if len(previousLines)*len(currentLines) > maxPayloadDiffMatrixCells {
		return oversizedDiff(previousLines, currentLines), true
	}

	lcs := make([][]int, len(previousLines)+1)
	for i := range lcs {
		lcs[i] = make([]int, len(currentLines)+1)
	}
	for i := len(previousLines) - 1; i >= 0; i-- {
		for j := len(currentLines) - 1; j >= 0; j-- {
			if previousLines[i] == currentLines[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else if lcs[i+1][j] >= lcs[i][j+1] {
				lcs[i][j] = lcs[i+1][j]
			} else {
				lcs[i][j] = lcs[i][j+1]
			}
		}
	}

	var b strings.Builder
	b.WriteString("--- previous_payload\n+++ current_payload\n")
	for i, j := 0, 0; i < len(previousLines) || j < len(currentLines); {
		switch {
		case i < len(previousLines) && j < len(currentLines) && previousLines[i] == currentLines[j]:
			i++
			j++
		case j >= len(currentLines) || (i < len(previousLines) && lcs[i+1][j] >= lcs[i][j+1]):
			writeDiffLine(&b, '-', previousLines[i])
			i++
		default:
			writeDiffLine(&b, '+', currentLines[j])
			j++
		}
	}
	return b.String(), false
}

func oversizedDiff(previousLines, currentLines []string) string {
	return fmt.Sprintf("--- previous_payload\n+++ current_payload\npayload diff omitted: formatted payloads are too large to diff safely (previousLines=%d, currentLines=%d)\n",
		len(previousLines), len(currentLines))
}

func splitDiffLines(s string) []string {
	if s == "" {
		return nil
	}
	lines := strings.SplitAfter(s, "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func writeDiffLine(b *strings.Builder, prefix byte, line string) {
	b.WriteByte(prefix)
	b.WriteString(line)
	if !strings.HasSuffix(line, "\n") {
		b.WriteByte('\n')
	}
}
