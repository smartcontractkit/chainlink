// Package launcher provides the functionality to launch and manage OCR instances.
// For system-level documentation and diagrams, please refer to the README.md
// in this package directory.
//
// The CCIP launcher, at a high level, consumes updates from the Capabilities Registry,
// in particular, DON additions, updates, and removals, and depending on the changes,
// launches (in the case of additions and updates) or shuts down (in the case of removals)
// CCIP OCR instances.
//
// It achieves this by diffing the current state of the registry with the previous state,
// and then launching or shutting down instances as necessary. See the launcher's tick()
// method for the main logic.
//
// Diffing logic is contained within diff.go, and the main logic for launching and shutting
// down instances is contained within launcher.go.
//
// Active/candidate deployment support is provided by the ccipDeployment struct in deployment.go.
package launcher
