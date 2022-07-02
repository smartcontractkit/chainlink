package chainlink

import "net/url"

func (l *legacyGeneralConfig) DatabaseURL() url.URL {
	//TODO https://app.shortcut.com/chainlinklabs/story/33624/add-secrets-toml
	panic("implement me")
}

func (l *legacyGeneralConfig) ExplorerAccessKey() string {
	//TODO https://app.shortcut.com/chainlinklabs/story/33624/add-secrets-toml
	panic("implement me")
}

func (l *legacyGeneralConfig) ExplorerSecret() string {
	//TODO https://app.shortcut.com/chainlinklabs/story/33624/add-secrets-toml
	panic("implement me")
}

func (l *legacyGeneralConfig) SessionSecret() ([]byte, error) {
	//TODO https://app.shortcut.com/chainlinklabs/story/33624/add-secrets-toml
	panic("implement me")
}
