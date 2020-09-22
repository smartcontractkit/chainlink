package models

type OffchainreportingOracleSpec struct {
	ID int `gorm:"primary_key"`
}

func (OffchainreportingOracleSpec) TableName() string { return "offchainreporting_oracle_specs" }
