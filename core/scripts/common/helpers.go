package common

import (
	"flag"
	"fmt"
)

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ParseArgs(flagSet *flag.FlagSet, args []string, requiredArgs ...string) {
	PanicErr(flagSet.Parse(args))
	seen := map[string]bool{}
	argValues := map[string]string{}
	flagSet.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
		argValues[f.Name] = f.Value.String()
	})
	for _, req := range requiredArgs {
		if !seen[req] {
			PanicErr(fmt.Errorf("missing required -%s argument/flag", req))
		} else if req == "sub-id" && argValues[req] == "0" {
			PanicErr(fmt.Errorf("missing required -%s argument/flag", req))
		}
	}
}
