// Copyright (C) 2017 ScyllaDB
// Use of this source code is governed by a ALv2-style
// license that can be found in the LICENSE file.
// Imported from: https://github.com/scylladb/go-reflectx

package postgres

import (
	"strings"
	"testing"
)

var camelTosnakeTests = []struct {
	N string
	V string
}{
	{"a", "a"},
	{"snake", "snake"},
	{"A", "a"},
	{"ID", "id"},
	{"MOTD", "motd"},
	{"Snake", "snake"},
	{"SnakeTest", "snake_test"},
	{"APIResponse", "api_response"},
	{"SnakeID", "snake_id"},
	{"Snake_Id", "snake_id"},
	{"Snake_ID", "snake_id"},
	{"SnakeIDGoogle", "snake_id_google"},
	{"LinuxMOTD", "linux_motd"},
	{"OMGWTFBBQ", "omgwtfbbq"},
	{"omg_wtf_bbq", "omg_wtf_bbq"},
	{"woof_woof", "woof_woof"},
	{"_woof_woof", "_woof_woof"},
	{"woof_woof_", "woof_woof_"},
	{"WOOF", "woof"},
	{"Woof", "woof"},
	{"woof", "woof"},
	{"woof0_woof1", "woof0_woof1"},
	{"_woof0_woof1_2", "_woof0_woof1_2"},
	{"woof0_WOOF1_2", "woof0_woof1_2"},
	{"WOOF0", "woof0"},
	{"Woof1", "woof1"},
	{"woof2", "woof2"},
	{"woofWoof", "woof_woof"},
	{"woofWOOF", "woof_woof"},
	{"woof_WOOF", "woof_woof"},
	{"Woof_WOOF", "woof_woof"},
	{"WOOFWoofWoofWOOFWoofWoof", "woof_woof_woof_woof_woof_woof"},
	{"WOOF_Woof_woof_WOOF_Woof_woof", "woof_woof_woof_woof_woof_woof"},
	{"Woof_W", "woof_w"},
	{"Woof_w", "woof_w"},
	{"WoofW", "woof_w"},
	{"Woof_W_", "woof_w_"},
	{"Woof_w_", "woof_w_"},
	{"WoofW_", "woof_w_"},
	{"WOOF_", "woof_"},
	{"W_Woof", "w_woof"},
	{"w_Woof", "w_woof"},
	{"WWoof", "w_woof"},
	{"_W_Woof", "_w_woof"},
	{"_w_Woof", "_w_woof"},
	{"_WWoof", "_w_woof"},
	{"_WOOF", "_woof"},
	{"_woof", "_woof"},
}

func TestCamelToSnakeASCII(t *testing.T) {
	for _, test := range camelTosnakeTests {
		if actual := CamelToSnakeASCII(test.N); actual != test.V {
			t.Error("V", test.V, "got", actual, test)
		}
	}
}

func BenchmarkCamelToSnakeASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CamelToSnakeASCII(camelTosnakeTests[b.N%len(camelTosnakeTests)].N)
	}
}

func BenchmarkToLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.ToLower(camelTosnakeTests[b.N%len(camelTosnakeTests)].N)
	}
}
