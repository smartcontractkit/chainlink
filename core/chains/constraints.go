package chains

// ID types represent unique identifiers within a particular chain type. Using string is recommended.
type ID any

// Config types should have fields for chain configuration, and normally implement sql.Scanner and driver.Valuer, but
// that is not enforced here since legacy types used mixed pointer and value receivers.
type Config interface {
	// sql.Scanner
	// driver.Valuer
}

// Node types should be a struct including these default fields:
//  ID        int32
//  Name      string
//  CreatedAt time.Time
//  UpdatedAt time.Time
type Node any
