package v2

import (
	"os"
	"strings"
)

var (
	CLDev = strings.ToLower(os.Getenv("CL_DEV")) == "true"
)
