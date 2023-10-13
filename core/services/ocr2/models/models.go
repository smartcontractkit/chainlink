package models

type MercuryCredentials struct {
	LegacyURL string
	URL       string
	Username  string
	Password  string
}

func (c *MercuryCredentials) GetUsername() string {
	return c.Username
}

func (c *MercuryCredentials) GetPassword() string {
	return c.Password
}

func (c *MercuryCredentials) GetURL() string {
	return c.URL
}

func (c *MercuryCredentials) GetLegacyURL() string {
	return c.LegacyURL
}
