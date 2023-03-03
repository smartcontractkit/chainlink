// +build go1.12

package memsize

import "reflect"

func iterateMap(m reflect.Value, fn func(k, v reflect.Value)) {
	it := m.MapRange()
	for it.Next() {
		fn(it.Key(), it.Value())
	}
}
