package blockchain

// PolygonMultinodeClient represents a multi-node, EVM compatible client for the Klaytn network
type PolygonMultinodeClient struct {
	*EthereumMultinodeClient
}

// PolygonClient represents a single node, EVM compatible client for the Polygon network
type PolygonClient struct {
	*EthereumClient
}
