package models

type MercuryCredentials struct {
	LegacyURL string
	URL       string
	Username  string
	Password  string
}

func (mc *MercuryCredentials) Validate() bool {
	return mc.LegacyURL != "" && mc.URL != "" && mc.Username != "" && mc.Password != ""
}
