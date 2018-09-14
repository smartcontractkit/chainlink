package migrations

func ExportedRegisterMigration(migration migration) {
	registerMigration(migration)
}

func ExportedAvailableMigrationTimestamps() []string {
	return availableMigrationTimestamps()
}
