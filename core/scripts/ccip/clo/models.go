package clo

type GenericGraphQLResponseBody[T any] struct {
	Data   T `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var listChains = `query ListChains {
	ccip {
		chains {
			id
			network {name, chainID}
			contracts {name, address}
		}
	}
}
`

type ListChainsResponseData struct {
	Ccip struct {
		Chains []ChainC `json:"chains"`
	} `json:"ccip"`
}

type ChainC struct {
	ID        string      `json:"id"`
	Network   NetworkC    `json:"network"`
	Contracts []ContractC `json:"contracts"`
}

type NetworkC struct {
	Name    string `json:"name"`
	ChainID string `json:"chainID"`
}

type ContractC struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

var listLanes = `query ListLanes {
	ccip {
		lanes {
			id
			legA {
				id
				source {
					chain {
						network {
							name
							chainID
						}
					}
					contracts {
						name
						address
						deployedAt
					}
				}
				destination {
					chain {
						network {
							name
							chainID
						}
					}
					contracts {
						name
						address
						deployedAt
					}
				}
			}
			legB {
				id
				source {
					chain {
						network {
							name
							chainID
						}
					}
					contracts {
						name
						address
						deployedAt
					}
				}
				destination {
					chain {
						network {
							name
							chainID
						}
					}
					contracts {
						name
						address
						deployedAt
					}
				}
			}
		}
	}
}
`

type ListLanesResponse struct {
	Ccip struct {
		Lanes []LaneL `json:"lanes"`
	} `json:"ccip"`
}

type LaneL struct {
	ID   string `json:"id"`
	LegA struct {
		LegL
	} `json:"legA"`
	LegB struct {
		LegL
	} `json:"legB"`
}

type LegL struct {
	ID          string       `json:"id"`
	Source      SourceL      `json:"source"`
	Destination DestinationL `json:"destination"`
}

type SourceL struct {
	Chain     ChainL      `json:"chain"`
	Contracts []ContractL `json:"contracts"`
}

type DestinationL struct {
	Chain     ChainL      `json:"chain"`
	Contracts []ContractL `json:"contracts"`
}

type ChainL struct {
	Network NetworkL `json:"network"`
}

type NetworkL struct {
	Name    string `json:"name"`
	ChainID string `json:"chainID"`
}

type ContractL struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	DeployedAt string `json:"deployedAt"`
}
