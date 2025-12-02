package presenters

type VerifyDKGResultResource struct {
	JAID
	SHA string `json:"sha"`
}

// GetName implements the api2go EntityNamer interface
func (VerifyDKGResultResource) GetName() string {
	return "verifyDKGResult"
}

func NewVerifyDKGResultResource(sha string) VerifyDKGResultResource {
	return VerifyDKGResultResource{
		JAID: NewJAID(sha),
		SHA:  sha,
	}
}

type ExportDKGResultResource struct {
	JAID
	HexDKGResultPackage string `json:"hexDKGResultPackage"`
	SHA                 string `json:"sha"`
}

// GetName implements the api2go EntityNamer interface
func (ExportDKGResultResource) GetName() string {
	return "exportDKGResult"
}

func NewExportDKGResultResource(hexResult, sha string) ExportDKGResultResource {
	return ExportDKGResultResource{
		JAID:                NewJAID(sha),
		HexDKGResultPackage: hexResult,
		SHA:                 sha,
	}
}
