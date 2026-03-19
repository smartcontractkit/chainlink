package vrfv2plus

// BlockhashForwarderJobSpec defines the BHF (blockhashmapper) job for a Chainlink node.
type BlockhashForwarderJobSpec struct {
	Name                       string
	ExternalJobID              string
	ForwardingAllowed          bool
	CoordinatorV2Address       string
	CoordinatorV2PlusAddress   string
	BlockhashStoreAddress      string
	BatchBlockhashStoreAddress string
	FromAddresses              []string
	EVMChainID                 string
	WaitBlocks                 int
	LookbackBlocks             int
	PollPeriod                 string
	RunTimeout                 string
}

func (b *BlockhashForwarderJobSpec) Type() string { return "blockheaderfeeder" }

func (b *BlockhashForwarderJobSpec) String() (string, error) {
	tmpl := `type = "blockheaderfeeder"
schemaVersion                 = 1
name                          = "{{.Name}}"
forwardingAllowed             = {{.ForwardingAllowed}}
coordinatorV2Address          = "{{.CoordinatorV2Address}}"
coordinatorV2PlusAddress      = "{{.CoordinatorV2PlusAddress}}"
blockhashStoreAddress	      = "{{.BlockhashStoreAddress}}"
batchBlockhashStoreAddress	  = "{{.BatchBlockhashStoreAddress}}"
fromAddresses                 = [{{range .FromAddresses}}"{{.}}",{{end}}]
evmChainID                    = "{{.EVMChainID}}"
externalJobID                 = "{{.ExternalJobID}}"
waitBlocks                    = {{.WaitBlocks}}
lookbackBlocks                = {{.LookbackBlocks}}
pollPeriod                    = "{{.PollPeriod}}"
runTimeout                    = "{{.RunTimeout}}"
`
	return marshallTemplate(b, "BlockheaderFeeder Job", tmpl)
}
