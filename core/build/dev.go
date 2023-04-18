//go:build dev

package build

// Dev is a build tag that allow use of insecure config family.
// it's set to false by default, building the node with -tags=dev allows it to be overridden
// and enable the use of insecure configs.
const Dev = true
