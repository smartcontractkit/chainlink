package config

const DefaultSchedule = "*/30 * * * * *"

type Config struct {
	Schedule string `json:"schedule"`
}
