package s4

type PluginConfig struct {
	DONID                   string
	ProductName             string
	NSnapshotShards         uint
	MaxObservationEntries   uint
	MaxReportEntries        uint
	MaxDeleteExpiredEntries uint
}
