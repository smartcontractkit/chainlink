package node_config

type WebServer struct {
	*WebServerTLS       `toml:"TLS"`
	*WebServerRateLimit `toml:"RateLimit"`
	AllowOrigins        string `toml:"AllowOrigins"`
	HTTPPort            int    `toml:"HTTPPort"`
	SecureCookies       bool   `toml:"SecureCookies"`
	SessionTimeout      string `toml:"SessionTimeout"`
}

type WebServerTLS struct {
	HTTPSPort int `toml:"HTTPSPort"`
}

type WebServerRateLimit struct {
	Authenticated   int `toml:"Authenticated"`
	Unauthenticated int `toml:"Unauthenticated"`
}
