package main

import (
	"github.com/smartcontractkit/wasp"
)

func main() {
	d, err := wasp.NewDashboard(nil, nil)
	if err != nil {
		panic(err)
	}
	if _, err := d.Deploy(); err != nil {
		panic(err)
	}
}
