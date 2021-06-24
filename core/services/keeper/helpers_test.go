package keeper

import "github.com/ethereum/go-ethereum"

func (rs *RegistrySynchronizer) ExportedFullSync() {
	rs.fullSync()
}

func (rs *RegistrySynchronizer) ExportedProcessLogs() {
	rs.processLogs()
}

func (executer *UpkeepExecuter) ExportedConstructCheckUpkeepCallMsg(upkeep UpkeepRegistration) (ethereum.CallMsg, error) {
	return executer.constructCheckUpkeepCallMsg(upkeep)
}
