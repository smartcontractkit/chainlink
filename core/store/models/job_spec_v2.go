package models

type JobSpecV2 struct {
	ID                            int32 `gorm: "primary_key"`
	OffchainreportingOracleSpecID int32
	OffchainreportingOracleSpec   *OffchainReportingOracleSpec
}

func (js JobSpecV2) TableName() string { return "jobs" }
