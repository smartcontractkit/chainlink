// Package store is used to keep application events in sync between
// the database on the node and the blockchain.
//
// Config
//
// Config contains the local configuration options that the application
// will adhere to.
//
// CallerSubscriberClient
//
// This makes use of Go-Ethereum's functions to interact with the blockchain.
// The underlying functions can be viewed here:
//  go-ethereum/rpc/client.go
//
// KeyStore
//
// KeyStore also utilizes Go-Ethereum's functions to store encrypted keys
// on the local file system.
// The underlying functions can be viewed here:
//  go-ethereum/accounts/keystore/keystore.go
//
// Store
//
// The Store is the persistence layer for the application. It saves the
// the application state and most interaction with the node needs to occur
// through the store.
//
// Tx Manager
//
// The transaction manager is used to synchronize interactions on the
// Ethereum blockchain with the application and database.
package store
