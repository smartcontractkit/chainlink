package node_config

type Database struct {
	MaxIdleConns     int  `toml:"MaxIdleConns"`
	MaxOpenConns     int  `toml:"MaxOpenConns"`
	MigrateOnStartup bool `toml:"MigrateOnStartup"`
}
