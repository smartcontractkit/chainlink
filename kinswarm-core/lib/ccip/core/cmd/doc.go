// Package cmd is the front-end interface for the application
// as a command-line utility.
//
// # KeyStoreAuthenticator
//
// KeyStoreAuthenticator prompts the user for their password, which
// is used to unlock their keystore file to interact with the
// Ethereum blockchain. Since multiple keystore files can exist
// at the configured directory, the KeyStoreAuthenticator will try the
// password on all keystore files present.
//
// # Shell
//
// Shell is how the application is invoked from the command
// line. When you run the binary, for example `./chainlink n`,
// Shell.RunNode is called to start the Chainlink core.
// Similarly, running `./chainlink j` returns information on
// all jobs in the node, and `./chainlink s` with another
// argument as a JobID gives information specific to that job.
//
// # Renderer
//
// Renderer helps format and display data (based on the kind
// of data it is) to the command line.
package cmd
