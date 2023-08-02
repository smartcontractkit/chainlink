package node_config

type Log struct {
	Level       string `toml:"Level"`
	JSONConsole bool   `toml:"JSONConsole"`
}
