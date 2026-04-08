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
				return nil, errors.Wrapf(ErrKeypathNotFound, `could not resolve path ["%v"]`, strings.Join(path, `","`))
			}

		case []any:
			bigindex, ok := big.NewInt(0).SetString(part, 10)
			if !ok {
				return nil, errors.Wrapf(ErrKeypathNotFound, "JSONParse task error: %v is not a valid array index", part)
			} else if !bigindex.IsInt64() {
				if lax {
					decoded = nil
					break
				}
				return nil, errors.Wrapf(ErrKeypathNotFound, `could not resolve path ["%v"]`, strings.Join(path, `","`))
			}
			index := int(bigindex.Int64())
			if index < 0 {
				index = len(d) + index
			}

			exists := index >= 0 && index < len(d)
			if !exists && lax {
				decoded = nil
				break
			} else if !exists {
				return nil, errors.Wrapf(ErrKeypathNotFound, `could not resolve path ["%v"]`, strings.Join(path, `","`))
			}
			decoded = d[index]

		default:
			return nil, errors.Wrapf(ErrKeypathNotFound, `could not resolve path ["%v"]`, strings.Join(path, `","`))
		}
	}
	return decoded, nil
}
