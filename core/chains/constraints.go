package chains

import (
	"database/sql"
	"database/sql/driver"
)

// ID types represent unique identifiers within a particular chain type. Using string is recommended.
type ID any

// Config types should have fields for chain configuration, and implement sql.Scanner and driver.Valuer for persistence in JSON format.
type Config interface {
	sql.Scanner
	driver.Valuer
}

// Node types should be a struct including these default fields:
//  ID        int32
//  Name      string
//  CreatedAt time.Time
//  UpdatedAt time.Time
type Node any
