package external_adapter

// {
//     "data": {
//         "token": "eth",
//         "reserves": "platform",
//         "supplyChains": [
//             "5009297550715157269"
//         ],
//         "supplyChainBlocks": [
//             0
//         ]
//     }
// }

// EARequest represents the request structure sent to the secure mint external adapter.
type EARequest struct {
	Token             string   `json:"token"`
	Reserves          string   `json:"reserves"`
	SupplyChains      []string `json:"supplyChain,omitempty"`
	SupplyChainBlocks []uint64 `json:"supplyChainBlocks,omitempty"`
}

//	"mintables": {
//	    "5009297550715157269": {
//	        "mintable": "0",
//	        "block": 0
//	    }
//	},
//
//	"reserveInfo": {
//	    "reserveAmount": "10332550000000000000000",
//	    "timestamp": 1749483841486
//	},
//
//	"latestRelevantBlocks": {
//	    "5009297550715157269": 22667990
//	},

// EAResponse represents the response structure from the secure mint external adapter.
type EAResponse struct {
	Mintables            map[string]MintableInfo `json:"mintables"`
	ReserveInfo          ReserveInfo             `json:"reserveInfo"`
	LatestRelevantBlocks map[string]uint64       `json:"latestRelevantBlocks"`
}

type MintableInfo struct {
	Mintable string `json:"mintable"`
	Block    uint64 `json:"block"`
}

type ReserveInfo struct {
	ReserveAmount string `json:"reserveAmount"`
	Timestamp     int64  `json:"timestamp"`
}
