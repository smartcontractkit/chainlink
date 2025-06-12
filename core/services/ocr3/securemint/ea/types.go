package ea

// Request represents the request structure sent to the secure mint external adapter.
// Example (sent in the 'data' field):
//
//	{
//	    "data": {
//	        "token": "eth",
//	        "reserves": "platform",
//	        "supplyChains": [
//	            "5009297550715157269"
//	        ],
//	        "supplyChainBlocks": [
//	            0
//	        ]
//	    }
//	}
type Request struct {
	Token             string   `json:"token"`
	Reserves          string   `json:"reserves"`
	SupplyChains      []string `json:"supplyChains,omitempty"`
	SupplyChainBlocks []uint64 `json:"supplyChainBlocks,omitempty"`
}

// Response represents the response structure from the secure mint external adapter.
// Example:
//
//		{
//	    "mintables": {
//	        "5009297550715157269": {
//	            "mintable": "5",
//	            "block": 22667990
//	        }
//	    },
//	    "reserveInfo": {
//	        "reserveAmount": "10332550000000000000000",
//	        "timestamp": 1749483841486
//	    },
//	    "latestBlocks": {
//	        "5009297550715157269": 22667990
//	    }
//	}
type Response struct {
	Mintables    map[string]MintableInfo `json:"mintables"`
	ReserveInfo  ReserveInfo             `json:"reserveInfo"`
	LatestBlocks map[string]uint64       `json:"latestBlocks"`
}

type MintableInfo struct {
	Mintable string `json:"mintable"`
	Block    uint64 `json:"block"`
}

type ReserveInfo struct {
	ReserveAmount string `json:"reserveAmount"`
	Timestamp     int64  `json:"timestamp"`
}
