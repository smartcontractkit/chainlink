package migrations

func ExportedRegisterMigration(migration migration) {
	registerMigration(migration)
}

func ExportedClearRegisteredMigrations() {
	migrationMutex.Lock()
	availableMigrations = make(map[string]migration)
	migrationMutex.Unlock()
}

func ExportedAvailableMigrationTimestamps() []string {
	return availableMigrationTimestamps()
}
