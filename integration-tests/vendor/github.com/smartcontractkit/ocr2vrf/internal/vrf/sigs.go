package vrf

import (
	"bytes"
	"sort"
)

type sig []byte

type sigs []sig

var _ sort.Interface = sigs{}

func (s sigs) Len() int           { return len(s) }
func (s sigs) Less(i, j int) bool { return bytes.Compare(s[i], s[j]) < 0 }
func (s sigs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
