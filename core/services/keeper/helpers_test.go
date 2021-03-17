package keeper

func (rs *RegistrySynchronizer) ExportedSyncRegistry() {
	rs.syncRegistry()
}
