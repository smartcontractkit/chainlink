package protobuf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type TagPrefix int

// Possible tag options.
const (
	TagNone TagPrefix = iota
	TagOptional
	TagRequired
)

func ParseTag(field reflect.StructField) (id int, opt TagPrefix, name string) {
	tag := field.Tag.Get("protobuf")
	if tag == "" {
		return
	}
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		if part == "opt" {
			opt = TagOptional
		} else if part == "req" {
			opt = TagRequired
		} else {
			i, err := strconv.Atoi(part)
			if err != nil {
				name = part
			} else {
				id = int(i)
			}
		}
	}
	return
}

// ProtoField contains cached reflected metadata for struct fields.
type ProtoField struct {
	ID     int64
	Prefix TagPrefix
	Name   string // If non-empty, tag-defined field name.
	Index  []int
	Field  reflect.StructField
}

func (p *ProtoField) Required() bool {
	return p.Prefix == TagRequired || p.Field.Type.Kind() != reflect.Ptr
}

var cache = map[reflect.Type][]*ProtoField{}
var cacheLock sync.Mutex

func ProtoFields(t reflect.Type) []*ProtoField {
	cacheLock.Lock()
	idx, ok := cache[t]
	cacheLock.Unlock()
	if ok {
		return idx
	}
	id := 0
	idx = innerFieldIndexes(&id, t)
	seen := map[int64]struct{}{}
	for _, i := range idx {
		if _, ok := seen[i.ID]; ok {
			panic(fmt.Sprintf("protobuf ID %d reused in %s.%s", i.ID, t.PkgPath(), t.Name()))
		}
		seen[i.ID] = struct{}{}
	}
	cacheLock.Lock()
	defer cacheLock.Unlock()
	cache[t] = idx
	return idx
}

func innerFieldIndexes(id *int, v reflect.Type) []*ProtoField {
	if v.Kind() == reflect.Ptr {
		return innerFieldIndexes(id, v.Elem())
	}
	out := []*ProtoField{}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		*id++
		tid, prefix, name := ParseTag(f)
		if tid != 0 {
			*id = tid
		}
		if f.Anonymous {
			*id--
			for _, inner := range innerFieldIndexes(id, f.Type) {
				inner.Index = append([]int{i}, inner.Index...)
				out = append(out, inner)
			}
		} else {
			out = append(out, &ProtoField{
				ID:     int64(*id),
				Prefix: prefix,
				Name:   name,
				Index:  []int{i},
				Field:  f,
			})
		}
	}
	return out

}
