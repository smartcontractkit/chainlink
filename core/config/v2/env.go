package v2

import (
	"os"
	"strings"
)

var (
	CLConfig = os.Getenv("CL_CONFIG")
	CLDev    = "true" == strings.ToLower(os.Getenv("CL_DEV"))
)
