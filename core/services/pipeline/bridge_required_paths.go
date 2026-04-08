package pipeline

import (
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// bridgeDataRootVarRegexp matches a data field that is only a single $(dotID)
// reference with no nested keypath (dotID: [a-zA-Z0-9_]+).
var bridgeDataRootVarRegexp = regexp.MustCompile(`^\$\(\s*([a-zA-Z0-9_]+)\s*\)$`)

// RequiredJSONPathsFromBridge returns JSON path segments (split the same way as
// jsonparse) that strict downstream jsonparse tasks require when they read this
// bridge task's string output directly.
//
// Limitations: does not follow merge/median/etc.; skips lax jsonparse, dynamic
// path/data containing "$(", data="$(bridge.field)", or jsonparse whose data
// input is not this bridge (by lowest output index or explicit $(bridgeID)).
func RequiredJSONPathsFromBridge(bridge Task) [][]string {
	if bridge == nil {
		return nil
	}
	bt, ok := bridge.(*BridgeTask)
	if !ok {
		return nil
	}
	bridgeID := bt.ID()
	bridgeDot := bt.DotID()

	seen := make(map[string]struct{})
	var out [][]string

	for _, d := range bt.GetDescendantTasks() {
		jp, ok := d.(*JSONParseTask)
		if !ok {
			continue
		}
		if jsonParseLaxStatic(jp) {
			continue
		}
		if strings.Contains(jp.Path, "$(") {
			continue
		}
		dataTrim := strings.TrimSpace(jp.Data)
		if strings.Contains(dataTrim, "$(") {
			m := bridgeDataRootVarRegexp.FindStringSubmatch(dataTrim)
			if len(m) != 2 || m[1] != bridgeDot {
				continue
			}
		}
		if !jsonParseDataReferencesBridge(jp, bridgeID, bridgeDot) {
			continue
		}
		segs, ok := jsonParseStaticPathSegments(jp)
		if !ok {
			continue
		}
		key := strings.Join(segs, "\x00")
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, segs)
	}
	return out
}

func jsonParseLaxStatic(j *JSONParseTask) bool {
	trimmed := strings.TrimSpace(j.Lax)
	if trimmed == "" {
		return false
	}
	b, err := strconv.ParseBool(trimmed)
	return err == nil && b
}

func jsonParseDataReferencesBridge(j *JSONParseTask, bridgeID int, bridgeDot string) bool {
	dataTrim := strings.TrimSpace(j.Data)
	if dataTrim == "" {
		src, ok := jsonParseLowestIndexPropagatingInput(j)
		return ok && src.ID() == bridgeID
	}
	m := bridgeDataRootVarRegexp.FindStringSubmatch(dataTrim)
	return len(m) == 2 && m[1] == bridgeDot
}

func jsonParseLowestIndexPropagatingInput(j *JSONParseTask) (Task, bool) {
	var (
		found        bool
		minOutputIdx int32 = math.MaxInt32
		src          Task
	)
	for _, dep := range j.Inputs() {
		if !dep.PropagateResult {
			continue
		}
		idx := dep.InputTask.OutputIndex()
		if !found || idx < minOutputIdx {
			found = true
			minOutputIdx = idx
			src = dep.InputTask
		}
	}
	return src, found
}

func jsonParseStaticPathSegments(j *JSONParseTask) ([]string, bool) {
	if strings.TrimSpace(j.Path) == "" {
		return nil, false
	}
	sep := strings.TrimSpace(j.Separator)
	if sep == "" {
		sep = ","
	}
	parts := strings.Split(j.Path, sep)
	if len(parts) == 0 {
		return nil, false
	}
	return parts, true
}

// jsonDecodeValidateRequiredPaths checks each path in the raw JSON body using
// streaming lookup (no full value unmarshal). Missing keys, JSON null, JSON numbers
// that are not valid for shopspring/decimal (including exponent/size limits), and the
// string "NaN" are treated as failure so bridge cache fallback can apply.
//
// For each present value, parseRequiredValidValue must succeed: decimals and floats
// via shopspring/decimal, integer literals (decimal or hex with optional sign) via
// math/big so huge 0x… values validate and scientific notation is not mistaken for
// hexadecimal.
func jsonDecodeValidateRequiredPaths(body []byte, paths [][]string) error {
	if len(paths) == 0 {
		return nil
	}
	for _, path := range paths {
		if len(path) == 0 {
			continue
		}

		value, dataType, _, err := jsonparser.Get(body, path...)
		if err != nil {
			return errors.Wrapf(err, "required path %q", strings.Join(path, ","))
		}

		if dataType == jsonparser.Null {
			return errors.Errorf("required path %q is null", strings.Join(path, ","))
		}

		if err := parseRequiredValidValue(value); err != nil {
			return errors.Wrapf(err, "required path %q", strings.Join(path, ","))
		}
	}
	return nil
}

// parseRequiredValidValue accepts values that can be parsed by shopspring/decimal or math/big.
func parseRequiredValidValue(value []byte) error {
	strValue := string(value)
	if _, err := decimal.NewFromString(strValue); err == nil {
		return nil
	}
	if _, ok := new(big.Int).SetString(strValue, 0); ok {
		return nil
	}
	return errors.Errorf("invalid value: %s", string(value))
}
