package vrfv2

// BlockhashStoreJobSpec defines the BHS job for VRF v2 (coordinator address used for both V2 and V2Plus fields).
type BlockhashStoreJobSpec struct {
	Name                     string
	ExternalJobID            string
	CoordinatorV2Address     string
	CoordinatorV2PlusAddress string
	BlockhashStoreAddress    string
	FromAddresses            []string
	EVMChainID               string
	WaitBlocks               int
	LookbackBlocks           int
	PollPeriod               string
	RunTimeout               string
}

func (b *BlockhashStoreJobSpec) Type() string { return "blockhashstore" }

func (b *BlockhashStoreJobSpec) String() (string, error) {
	tmpl := `type = "blockhashstore"
schemaVersion = 1
name = "{{.Name}}"
externalJobID = "{{.ExternalJobID}}"
evmChainID = "{{.EVMChainID}}"
coordinatorV2Address     = "{{.CoordinatorV2Address}}"
coordinatorV2PlusAddress = "{{.CoordinatorV2PlusAddress}}"
blockhashStoreAddress    = "{{.BlockhashStoreAddress}}"
waitBlocks               = {{.WaitBlocks}}
lookbackBlocks           = {{.LookbackBlocks}}
pollPeriod               = "{{.PollPeriod}}"
runTimeout               = "{{.RunTimeout}}"
fromAddresses            = [{{range .FromAddresses}}"{{.}}",{{end}}]
`
	return marshallTemplate(b, "BlockhashStore Job", tmpl)
}
