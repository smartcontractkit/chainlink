// +build !dev 

package web

import (
	"embed"
)

// Go's new embed feature doesn't allow us to embed things outside of the current module.
// To get around this, we need to make sure that the assets we want to embed are available
// inside this module. To achieve this, we direct webpack to output all of the compiled assets
// in this module's folder under the "assets" directory.

//go:embed "assets"
var uiEmbedFs embed.FS

// assetFs is the singleton file system instance that is used to serve the static
// assets for the operator UI.
var assetFs = NewEmbedFileSystem(uiEmbedFs, "assets")
