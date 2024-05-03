package main

import (
	_ "embed"
	"os"
	"text/template"
)

//go:embed merge_runner.go.tmpl
var mergeGoTmpl string

func main() {
	for i := 2; i <= 10; i++ {
		t, err := template.New("merge.go.tmpl").Funcs(map[string]any{
			"nis": func(i int) string {
				is := make([]byte, i)
				for j := 0; j < i; j++ {
					is[j] = 'I'
				}
				return string(is)
			},
			"rangeI": func(i int) []int {
				vals := make([]int, i)
				for j := 0; j < i; j++ {
					vals[j] = j + 1
				}
				return vals
			},
			"rangeIm2": func(i int) []int {
				vals := make([]int, i-2)
				for j := 2; j < i; j++ {
					vals[j-2] = j
				}
				return vals
			},
		}).Parse(mergeGoTmpl)

		output, err := os.Create("../merge_runner_gen.go")
		if err != nil {
			panic(err)
		}

		if err = t.Execute(output, 11); err != nil {
			panic(err)
		}
	}
}
