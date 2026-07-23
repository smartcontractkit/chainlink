package onchain

// AllowAllIPsCIDR is the gateway whitelist used for dev/stage. Griddle nodes sit
// behind cluster networking, so the gateway itself is not the trust boundary.
const AllowAllIPsCIDR = "0.0.0.0/0"

// FallbackDeployerAddress is the well-known Anvil account 0 address, used only as a
// last-resort AdminAddress for a node when no discovered EVM address is available.
const FallbackDeployerAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
