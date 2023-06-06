package models

type MercuryCredentials struct {
	URL      string
	Username string
	Password string
}

func (mc *MercuryCredentials) Validate() bool {
	return mc.URL != "" && mc.Username != "" && mc.Password != ""
}
