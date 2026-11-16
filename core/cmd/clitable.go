package cmd

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

const (
	cliTableColumn = "║"
	cliTableRow    = "═"
	cliTableCenter = "╬"
)

var tableWidthCondition = runewidth.NewCondition()

var cliTableANSI = func() *regexp.Regexp {
	escape := "\x1b"
	bell := "\x07"
	terminator := "(" + regexp.QuoteMeta(escape+"\\") + "|" + regexp.QuoteMeta(bell) + ")"
	csi := regexp.QuoteMeta(escape+"[") + "[\x30-\x3f]*[\x20-\x2f]*[\x40-\x7e]"
	osc := regexp.QuoteMeta(escape+"]") + ".*?" + terminator
	return regexp.MustCompile("(" + csi + "|" + osc + ")")
}()

type cliTable struct {
	writer         io.Writer
	headers        [][]string
	rows           [][][]string
	widthCondition *runewidth.Condition
}

func newCLITable(writer io.Writer, headers []string) *cliTable {
	condition := *tableWidthCondition
	return &cliTable{
		writer:         writer,
		headers:        prepareCLITableCells(headers, true),
		widthCondition: &condition,
	}
}

func (t *cliTable) Append(row []string) error {
	t.rows = append(t.rows, prepareCLITableCells(row, false))
	return nil
}

func prepareCLITableCells(cells []string, header bool) [][]string {
	result := make([][]string, 0)
	for i, cell := range cells {
		cell = strings.TrimSpace(cell)
		if header {
			cell = cliTableTitle(strings.Join(cliTableSplitCamelCase(cell), " "))
		}
		lines := strings.Split(cell, "\n")
		for len(result) < len(lines) {
			result = append(result, make([]string, len(cells)))
		}
		for lineIndex := range result {
			if lineIndex < len(lines) {
				result[lineIndex][i] = lines[lineIndex]
			}
		}
	}
	return result
}

func cliTableSplitCamelCase(src string) (entries []string) {
	if !utf8.ValidString(src) {
		return []string{src}
	}
	entries = []string{}
	var runes [][]rune
	lastClass := 0
	for _, r := range src {
		class := 4
		switch {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	for _, group := range runes {
		value := string(group)
		if len(group) > 0 && strings.TrimSpace(value) != "" && value != "_" {
			entries = append(entries, value)
		}
	}
	return entries
}

func cliTableTitle(name string) string {
	originalLength := len(name)
	runes := []rune(name)
	for i, r := range runes {
		switch r {
		case '_':
			runes[i] = ' '
		case '.':
			if (i != 0 && !cliTableNumericOrSpace(runes[i-1])) || (i != len(runes)-1 && !cliTableNumericOrSpace(runes[i+1])) {
				runes[i] = ' '
			}
		}
	}
	name = strings.TrimSpace(string(runes))
	if len(name) == 0 && originalLength > 0 {
		name = " "
	}
	return strings.ToUpper(name)
}

func cliTableNumericOrSpace(r rune) bool {
	return ('0' <= r && r <= '9') || r == ' '
}

func (t *cliTable) width(value string) int {
	plain := cliTableANSI.ReplaceAllLiteralString(value, "")
	return t.widthCondition.StringWidth(plain)
}

func (t *cliTable) updateWidths(lines [][]string, widths []int) {
	for _, line := range lines {
		for column, cell := range line {
			contentWidth := 0
			for _, visualLine := range strings.Split(cell, "\n") {
				lineWidth := t.width(strings.TrimSpace(visualLine))
				if lineWidth > contentWidth {
					contentWidth = lineWidth
				}
			}
			totalWidth := contentWidth + 2
			if totalWidth < 1 {
				totalWidth = 1
			}
			if totalWidth > widths[column] {
				widths[column] = totalWidth
			}
		}
	}
}

func (t *cliTable) columnCount() int {
	count := 0
	if len(t.headers) > 0 && len(t.headers[0]) > count {
		count = len(t.headers[0])
	}
	for _, row := range t.rows {
		if len(row) > 0 && len(row[0]) > count {
			count = len(row[0])
		}
	}
	return count
}

func (t *cliTable) columnWidths(count int) []int {
	widths := make([]int, count)
	t.updateWidths(t.headers, widths)
	for _, row := range t.rows {
		t.updateWidths(row, widths)
	}
	for i, width := range widths {
		if width == 0 {
			widths[i] = 1
		}
	}
	return widths
}

func padCLITableLine(line []string, columns int) []string {
	if len(line) >= columns {
		return line
	}
	padded := make([]string, columns)
	copy(padded, line)
	return padded
}

func (t *cliTable) padRight(value string, width int) string {
	if extra := width - t.width(value); extra > 0 {
		return value + strings.Repeat(" ", extra)
	}
	return value
}

func (t *cliTable) truncate(value string, maxWidth int) string {
	if maxWidth < 0 {
		return ""
	}
	displayWidth := t.width(value)
	if displayWidth == 0 {
		if maxWidth == 0 && len(value) > 0 {
			return ""
		}
		return value
	}
	if maxWidth == 0 {
		return ""
	}
	if displayWidth <= maxWidth {
		return value
	}

	var content bytes.Buffer
	var ansiSequence bytes.Buffer
	inANSI := false
	currentWidth := 0
	for _, r := range value {
		if r == '\x1b' {
			inANSI = true
			ansiSequence.Reset()
			ansiSequence.WriteRune(r)
			continue
		}
		if inANSI {
			ansiSequence.WriteRune(r)
			sequence := ansiSequence.Bytes()
			terminated := false
			if len(sequence) >= 2 {
				switch sequence[1] {
				case '[':
					terminated = len(sequence) >= 3 && r >= 0x40 && r <= 0x7e
				case ']':
					terminated = r == '\x07' || (len(sequence) > 1 && sequence[len(sequence)-2] == '\x1b' && r == '\\')
				}
			}
			if terminated {
				inANSI = false
				content.Write(ansiSequence.Bytes())
				ansiSequence.Reset()
			}
			continue
		}

		runeWidth := t.widthCondition.RuneWidth(r)
		if maxWidth == 0 {
			break
		}
		if currentWidth+runeWidth > maxWidth {
			break
		}
		content.WriteRune(r)
		currentWidth += runeWidth
	}
	return content.String()
}

func (t *cliTable) formatCell(content string, width int, header bool) string {
	if width <= 0 {
		return ""
	}
	contentWidth := t.width(content)
	availableContentWidth := width - 2
	if availableContentWidth < 0 {
		availableContentWidth = 0
	}
	if contentWidth > availableContentWidth {
		content = t.truncate(content, availableContentWidth)
		contentWidth = t.width(content)
	}
	totalPadding := width - contentWidth
	if totalPadding < 0 {
		totalPadding = 0
	}

	var result string
	if header {
		leftPadding := (totalPadding-2)/2 + 1
		rightPadding := totalPadding - leftPadding
		result = strings.Repeat(" ", leftPadding) + content + strings.Repeat(" ", rightPadding)
	} else {
		rightPadding := totalPadding - 1
		result = " " + content
		if rightPadding > 0 {
			result += strings.Repeat(" ", rightPadding)
		}
	}

	finalWidth := t.width(result)
	if finalWidth > width {
		result = t.truncate(result, width)
	} else if finalWidth < width {
		result = t.padRight(result, width)
	}
	return result
}

func (t *cliTable) renderContentLine(line []string, widths []int, header bool) string {
	line = padCLITableLine(line, len(widths))
	var output strings.Builder
	output.WriteString(cliTableColumn)
	for column, width := range widths {
		output.WriteString(t.formatCell(line[column], width, header))
		if column < len(widths)-1 {
			output.WriteString(cliTableColumn)
		}
	}
	output.WriteString(cliTableColumn)
	output.WriteByte('\n')
	return output.String()
}

func (t *cliTable) renderHorizontalRule(widths []int) string {
	if len(widths) == 0 {
		return ""
	}
	columnWidth := t.width(cliTableColumn)
	segmentWidth := t.width(cliTableRow)
	var output strings.Builder
	for column, width := range widths {
		adjustedWidth := width
		if column == 0 {
			adjustedWidth += columnWidth // Empty left corner
		}
		if column == len(widths)-1 {
			adjustedWidth += columnWidth // Empty right corner
		}

		repeat := adjustedWidth / segmentWidth
		if repeat < 1 && adjustedWidth > 0 {
			repeat = 1
		}
		segment := strings.Repeat(cliTableRow, repeat)
		actualWidth := t.width(segment)
		if actualWidth > adjustedWidth {
			segment = t.truncate(segment, adjustedWidth)
		} else if actualWidth < adjustedWidth {
			remainingWidth := adjustedWidth - actualWidth
			for i := 0; i < remainingWidth/segmentWidth; i++ {
				segment += cliTableRow
			}
			if t.width(segment) < adjustedWidth {
				segment = t.padRight(segment, adjustedWidth)
			}
		}
		output.WriteString(segment)
		if column < len(widths)-1 {
			// Empty specialized junctions fall back to the center symbol
			output.WriteString(cliTableCenter)
		}
	}
	output.WriteByte('\n')
	return output.String()
}

func (t *cliTable) render() error {
	columns := t.columnCount()
	if columns == 0 {
		return nil
	}
	widths := t.columnWidths(columns)
	write := func(value string) {
		_, _ = t.writer.Write([]byte(value))
	}

	if len(t.headers) > 0 {
		write(t.renderHorizontalRule(widths))
		for _, line := range t.headers {
			write(t.renderContentLine(line, widths, true))
		}
		if len(t.rows) > 0 {
			write(t.renderHorizontalRule(widths))
		}
	} else if len(t.rows) > 0 {
		write(t.renderHorizontalRule(widths))
	}

	for _, row := range t.rows {
		for lineIndex, line := range row {
			if lineIndex > 0 {
				blank := true
				for _, cell := range line {
					if strings.TrimSpace(cell) != "" {
						blank = false
						break
					}
				}
				if blank {
					continue
				}
			}
			write(t.renderContentLine(line, widths, false))
		}
	}
	write(t.renderHorizontalRule(widths))
	return nil
}
