// Package modulemap maps short runtime names to go module import paths.
package modulemap

import "fmt"

// runtimeToModule maps .tool-versions plugin names to go install module paths.
var runtimeToModule = map[string]string{
	"mockery": "github.com/vektra/mockery/v2",
}

// ModulePath returns the go module path for a short runtime name.
func ModulePath(runtime string) (string, error) {
	mod, ok := runtimeToModule[runtime]
	if !ok {
		return "", fmt.Errorf("no go module mapping for: %s", runtime)
	}
	return mod, nil
}

// Modules returns all mapped module paths.
func Modules() []string {
	out := make([]string, 0, len(runtimeToModule))
	for _, mod := range runtimeToModule {
		out = append(out, mod)
	}
	return out
}
