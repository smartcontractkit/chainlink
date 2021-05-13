package keeper

func (rs *RegistrySynchronizer) ExportedFullSync() {
	rs.fullSync()
}

func (rs *RegistrySynchronizer) ExportedProcessLogs() {
	rs.processLogs()
}
