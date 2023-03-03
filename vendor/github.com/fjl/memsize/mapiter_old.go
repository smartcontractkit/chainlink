// +build !go1.12

package memsize

import "reflect"

func iterateMap(m reflect.Value, fn func(k, v reflect.Value)) {
	for _, k := range m.MapKeys() {
		fn(k, m.MapIndex(k))
	}
}
