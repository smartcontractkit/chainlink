package s4

type PluginConfig struct {
	ProductName             string
	NSnapshotShards         uint
	MaxObservationEntries   uint
	MaxReportEntries        uint
	MaxDeleteExpiredEntries uint
}
