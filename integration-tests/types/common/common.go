package common

type PgOpts struct {
	User     string
	Password string
	DbName   string
	Networks []string
	Port     string
}
