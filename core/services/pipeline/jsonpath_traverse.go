package pipeline

import (
	"math/big"
	"strings"

	"github.com/pkg/errors"
)

// traverseJSONPath walks decoded JSON (map/slice tree) along path the same way
// JSONParseTask does. When lax is false, a missing segment returns ErrKeypathNotFound.
func traverseJSONPath(decoded any, path []string, lax bool) (any, error) {
	for _, part := range path {
		switch d := decoded.(type) {
		case map[string]any:
			var exists bool
			decoded, exists = d[part]
			if !exists && lax {
				decoded = nil
				break
			} else if !exists {
				return nil, errJSONPathNotFound(path)
			}

		case []any:
			next, laxMiss, err := jsonPathArrayStep(d, part, path, lax)
			if err != nil {
				return nil, err
			}
			if laxMiss {
				decoded = nil
				break
			}
			decoded = next

		default:
			return nil, errJSONPathNotFound(path)
		}
	}
	return decoded, nil
}

func errJSONPathNotFound(path []string) error {
	return errors.Wrapf(ErrKeypathNotFound, `could not resolve path ["%v"]`, strings.Join(path, `","`))
}

// jsonPathArrayStep resolves one path segment against a JSON array. laxMiss is true when lax
// mode treats the segment as missing (decoded becomes nil). err is set for invalid index syntax
// or strict resolution failure.
func jsonPathArrayStep(d []any, part string, path []string, lax bool) (next any, laxMiss bool, err error) {
	bi, ok := big.NewInt(0).SetString(part, 10)
	if !ok {
		return nil, false, errors.Wrapf(ErrKeypathNotFound, "JSONParse task error: %v is not a valid array index", part)
	}
	if !bi.IsInt64() {
		if lax {
			return nil, true, nil
		}
		return nil, false, errJSONPathNotFound(path)
	}

	index := int(bi.Int64())
	if index < 0 {
		index = len(d) + index
	}
	if index < 0 || index >= len(d) {
		if lax {
			return nil, true, nil
		}
		return nil, false, errJSONPathNotFound(path)
	}
	return d[index], false, nil
}
