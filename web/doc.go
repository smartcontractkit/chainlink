// Package web handles receiving and supplying information
// within the node.
//
// Router
//
// Router defines the valid paths for the node and responds
// to requests.
//
// JobsController
//
// JobsController allows for the creation of Jobs to be added
// to the node, and shows the current jobs which have already
// been added.
//
// JobRunsController
//
// JobRunsController allows for the creation of JobRuns within
// a given Job on the node.
//
// BridgeTypesController
//
// BridgeTypesController allows for the creation of BridgeTypes
// on the node. BridgeTypes are the external adapters which add
// functionality not available in the core, from outside the node.
package web
